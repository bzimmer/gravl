package dynamo

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
)

const (
	tableName         = "activities"
	attributeID       = "ID"
	attributeActivity = "act"
)

type ddb struct {
	svc *dynamodb.Client
}

// Open a dynamoDB database
func Open(ctx context.Context, cfg aws.Config) (store.Store, error) {
	return &ddb{svc: dynamodb.NewFromConfig(cfg)}, nil
}

// Close the source
func (b *ddb) Close() error {
	return nil
}

// Activities returns a channel of activities and errors for an athlete
func (b *ddb) Activities(ctx context.Context) <-chan *strava.ActivityResult {
	acts := make(chan *strava.ActivityResult)
	go func() {
		defer close(acts)
		var n float64
		input := &dynamodb.ScanInput{TableName: aws.String(tableName)}
		for {
			output, err := b.svc.Scan(ctx, input)
			if err != nil {
				acts <- &strava.ActivityResult{Err: err}
				return
			}
			if len(output.LastEvaluatedKey) == 0 {
				return
			}
			input.ExclusiveStartKey = output.LastEvaluatedKey
			for _, item := range output.Items {
				if math.Mod(n, 100) == 0 {
					log.Info().Float64("n", n).Str("db", "dynamodb").Msg("activities")
				}
				switch x := item[attributeActivity].(type) {
				case *types.AttributeValueMemberB:
					var act *strava.Activity
					if err := json.Unmarshal(x.Value, &act); err != nil {
						acts <- &strava.ActivityResult{Err: err}
						return
					}
					acts <- &strava.ActivityResult{Activity: act}
				}
				n++
			}
		}
	}()
	return acts
}

// Activity returns a fully populated Activity
func (b *ddb) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	output, err := b.get(ctx, activityID)
	if err != nil {
		return nil, err
	}
	switch x := output.Item[attributeActivity].(type) {
	case *types.AttributeValueMemberB:
		var act *strava.Activity
		if err := json.Unmarshal(x.Value, &act); err != nil {
			return nil, err
		}
		return act, nil
	default:
		return nil, fmt.Errorf("unexpected type %T", x)
	}
}

// Exists returns true if the activity exists, false otherwise
func (b *ddb) Exists(ctx context.Context, activityID int64) (bool, error) {
	output, err := b.get(ctx, activityID)
	if err != nil {
		return false, err
	}
	return output.Item != nil, nil
}

// Save the activities to the source
func (b *ddb) Save(ctx context.Context, acts ...*strava.Activity) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      make(map[string]types.AttributeValue),
	}
	for _, act := range acts {
		data, err := json.Marshal(act)
		if err != nil {
			return err
		}
		input.Item[attributeID] = &types.AttributeValueMemberN{Value: strconv.FormatInt(act.ID, 10)}
		input.Item[attributeActivity] = &types.AttributeValueMemberB{Value: data}
		_, err = b.svc.PutItem(ctx, input)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove the activities from the source
func (b *ddb) Remove(ctx context.Context, acts ...*strava.Activity) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       make(map[string]types.AttributeValue),
	}
	for _, act := range acts {
		input.Key[attributeID] = &types.AttributeValueMemberN{Value: strconv.FormatInt(act.ID, 10)}
		_, err := b.svc.DeleteItem(ctx, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *ddb) get(ctx context.Context, activityID int64) (*dynamodb.GetItemOutput, error) {
	return b.svc.GetItem(ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				attributeID: &types.AttributeValueMemberN{Value: strconv.FormatInt(activityID, 10)},
			},
		},
	)
}
