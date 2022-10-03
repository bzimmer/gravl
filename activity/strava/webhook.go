package strava

import (
	"context"
	"strconv"

	"github.com/bzimmer/activity/strava"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl"
)

func list(c *cli.Context, f func(sub *strava.WebhookSubscription) error) error {
	client := gravl.Runtime(c).Strava
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
		gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "webhook", c.Command.Name}, 1)
		log.Info().Time("created", sub.CreatedAt).Str("url", sub.CallbackURL).Int64("id", sub.ID).Msg("webhook")
	}
	return nil
}

func whlist(c *cli.Context) error {
	active := false
	if err := list(c, func(sub *strava.WebhookSubscription) error {
		active = true
		return gravl.Runtime(c).Encoder.Encode(sub)
	}); err != nil {
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
		if err := list(c, func(sub *strava.WebhookSubscription) error {
			args = append(args, strconv.FormatInt(sub.ID, 10))
			return nil
		}); err != nil {
			return err
		}
	}
	if len(args) == 0 {
		log.Info().Msg("no active subscriptions")
		return nil
	}
	return entityWithArgs(c, func(ctx context.Context, client *strava.Client, id int64) (any, error) {
		log.Info().Int64("id", id).Msg("unsubscribing")
		return id, client.Webhook.Unsubscribe(ctx, id)
	}, args)
}

func whsubscribe(c *cli.Context) error {
	client := gravl.Runtime(c).Strava
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()
	ack, err := client.Webhook.Subscribe(ctx, c.String("url"), c.String("verify"))
	if err != nil {
		return err
	}
	gravl.Runtime(c).Metrics.IncrCounter([]string{Provider, "webhook", c.Command.Name}, 1)
	return gravl.Runtime(c).Encoder.Encode(ack)
}

func whlistCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List all active webhook subscriptions",
		Action: whlist,
	}
}

func whunsubscribeCommand() *cli.Command {
	return &cli.Command{
		Name:    "unsubscribe",
		Aliases: []string{"delete", "remove"},
		Usage:   "Unsubscribe an active webhook subscription (or all if specified)",
		Action:  whunsubscribe,
	}
}

func whsubscribeCommand() *cli.Command {
	return &cli.Command{
		Name:  "subscribe",
		Usage: "Subscribe for webhook notifications",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Usage: "Address where webhook events will be sent (max length 255 characters"},
			&cli.StringFlag{
				Name:  "verify",
				Usage: "String chosen by the application owner for client security"},
		},
		Action: whsubscribe,
	}
}

func webhookCommand() *cli.Command {
	return &cli.Command{
		Name:  "webhook",
		Usage: "Manage webhook subscriptions",
		Subcommands: []*cli.Command{
			whlistCommand(),
			whsubscribeCommand(),
			whunsubscribeCommand(),
		},
	}
}
