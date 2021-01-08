package strava

import (
	"context"
	"fmt"
	"net/http"

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
	active := false
	err := list(c, func(sub *strava.WebhookSubscription) error {
		active = true
		return encoding.Encode(sub)
	})
	if err != nil {
		return err
	}
	if !active {
		log.Info().Msg("no active subscriptions")
	}
	return nil
}

func whunsubscribe(c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) == 0 {
		log.Info().Msg("querying active subscriptions")
		err := list(c, func(sub *strava.WebhookSubscription) error {
			args = append(args, fmt.Sprintf("%d", sub.ID))
			return nil
		})
		if err != nil {
			return err
		}
	}
	if len(args) == 0 {
		log.Info().Msg("no active subscriptions")
		return nil
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
	sub, err := subscriber(c)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	handle := web.NewLogHandler(log.Logger)
	mux.Handle("/webhook", handle(strava.NewWebhookHandler(sub)))
	address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
	log.Info().Str("address", address).Str("verify", sub.verify).Msg("serving ...")
	return http.ListenAndServe(address, mux)
}

var verifyFlag = &cli.StringFlag{
	Name:  "verify",
	Value: "",
	Usage: "String chosen by the application owner for client security"}

var whlistCommand = &cli.Command{
	Name:   "list",
	Usage:  "List all active webhook subscriptions",
	Action: whlist,
}

var whunsubscribeCommand = &cli.Command{
	Name:    "unsubscribe",
	Aliases: []string{"delete", "remove"},
	Usage:   "Unsubscribe an active webhook subscription (or all if specified)",
	Action:  whunsubscribe,
}

var whsubscribeCommand = &cli.Command{
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

var whserveCommand = &cli.Command{
	Name:    "serve",
	Aliases: []string{"listen"},
	Usage:   "Run the webhook listener",
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
		whlistCommand,
		whserveCommand,
		whsubscribeCommand,
		whunsubscribeCommand,
	},
}
