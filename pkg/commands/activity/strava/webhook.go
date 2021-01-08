package strava

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/web"
)

type sub struct {
	verify string
}

func (s *sub) SubscriptionRequest(challenge, verify string) error {
	return nil
}
func (s *sub) MessageReceived(m *strava.WebhookMessage) error {
	return encoding.Encode(m)
}

func subscriber(c *cli.Context) (*sub, error) {
	t := c.String("verify")
	if t == "" {
		var err error
		t, err = commands.Token(16)
		if err != nil {
			return nil, err
		}
	}
	return &sub{verify: t}, nil
}

func list(c *cli.Context, f func(sub *strava.WebhookSubscription) error) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	subs, err := client.Webhook.List(ctx)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		if err = f(sub); err != nil {
			return err
		}
	}
	return nil
}

func whlist(c *cli.Context) error {
	return list(c, func(sub *strava.WebhookSubscription) error {
		return encoding.Encode(sub)
	})
}

func whunsubscribe(c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) == 0 {
		err := list(c, func(sub *strava.WebhookSubscription) error {
			args = append(args, fmt.Sprintf("%d", sub.ID))
			return nil
		})
		if err != nil {
			return err
		}
	}
	return entityWithArgs(c, func(ctx context.Context, client *strava.Client, id int64) (interface{}, error) {
		log.Info().Int64("id", id).Msg("unsubscribing")
		return id, client.Webhook.Unsubscribe(ctx, id)
	}, args)
}

func whsubscribe(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	ack, err := client.Webhook.Subscribe(ctx, c.String("url"), c.String("verify"))
	if err != nil {
		return err
	}
	return encoding.Encode(ack)
}

func whserve(c *cli.Context) error {
	r := gin.New()
	r.Use(gin.Recovery(), web.LogMiddleware())
	s, err := subscriber(c)
	if err != nil {
		return err
	}
	r.GET("/wh", strava.WebhookSubscriptionHandler(s))
	r.POST("/wh", strava.WebhookEventHandler(s))

	address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
	log.Info().Str("address", address).Str("verify", s.verify).Msg("serving ...")
	return r.Run(address)
}

var verifyFlag = &cli.StringFlag{
	Name:  "verify",
	Value: "",
	Usage: "String chosen by the application owner for client security"}

var listCommand = &cli.Command{
	Name:   "list",
	Usage:  "List all active webhook subscriptions",
	Action: whlist,
}

var unsubscribeCommand = &cli.Command{
	Name:    "unsubscribe",
	Aliases: []string{"delete"},
	Usage:   "Unsubscribe an active webhook subscription (or all if specified)",
	Action:  whunsubscribe,
}

var subscribeCommand = &cli.Command{
	Name:  "subscribe",
	Usage: "Subscribe for webhook notications",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Value: "",
			Usage: "Address where webhook events will be sent (max length 255 characters"},
		verifyFlag,
	},
	Action: whsubscribe,
}

var serveCommand = &cli.Command{
	Name:  "serve",
	Usage: "Run the webhook listener",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Value: 9003,
			Usage: "Port on which to listen"},
		verifyFlag,
	},
	Action: whserve,
}

var webhookCommand = &cli.Command{
	Name:  "webhook",
	Usage: "Manage webhook subscriptions",
	Subcommands: []*cli.Command{
		listCommand,
		serveCommand,
		subscribeCommand,
		unsubscribeCommand,
	},
}
