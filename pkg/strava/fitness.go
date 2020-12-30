package strava

import (
	"context"
	"fmt"
	"net/http"
)

// FitnessService is the API for fitness endpoints
type FitnessService service

// TrainingLoad returns the training load for an athlete
func (s *FitnessService) TrainingLoad(ctx context.Context, userID int) ([]*TrainingLoad, error) {
	uri := fmt.Sprintf("fitness/%d", userID)
	req, err := s.client.newWebRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	type Response struct {
		Data      []*TrainingLoad `json:"data"`
		Reference interface{}     `json:"reference"`
	}
	var res []*Response
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("expected one result, found %d", len(res))
	}
	return res[0].Data, nil
}
