package strava

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/net/ngrok"
	"github.com/bzimmer/gravl/pkg/web"
)

const webhookPath = "/webhook"

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

func tunnel(c *cli.Context) (*ngrok.Tunnel, error) {
	ng, err := ngrok.NewClient(ngrok.WithHTTPTracing(c.Bool("http-tracing")))
	if err != nil {
		return nil, err
	}
	tns, err := ng.Tunnels.Tunnels(c.Context)
	if err != nil {
		return nil, err
	}

	port := c.String("port")
	for _, tn := range tns {
		if tn.Proto != "https" {
			continue
		}
		u, err := url.Parse(tn.Config.Address)
		if err != nil {
			return nil, err
		}
		if u.Port() == port {
			return tn, nil
		}
	}
	return nil, fmt.Errorf("no tunnel for port %s", port)
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
	mux.Handle(webhookPath, handle(strava.NewWebhookHandler(sub)))
	address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
	log.Info().Str("address", address).Str("verify", sub.verify).Msg("serving ...")
	return http.ListenAndServe(address, mux)
}

func whdaemon(c *cli.Context) error {
	// url:
	//   if url: nothing
	//   if "" or ngrok: query for local ngrok endpoints and ensure a tunnel exists for the port
	// verify:
	//   if verify: nothing
	//   if "": generate token
	// launch the server
	// remove all active subscriptions
	// start the new one
	uri := c.String("url")
	if uri == "" || uri == "ngrok" {
		tn, err := tunnel(c)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				log.Warn().Msg("Make sure ngrok is running.")
			}
			return err
		}
		uri = strings.TrimSuffix(tn.PublicURL, "/")
		uri = fmt.Sprintf("%s%s", uri, webhookPath)
	}
	verify := c.String("verify")
	if verify == "" {
		t, err := commands.Token(16)
		if err != nil {
			return err
		}
		verify = t
	}
	fmt.Printf("gravl strava webhook serve --verify '%s'\n", verify)
	fmt.Printf("gravl strava webhook subscribe --url '%s' --verify '%s'\n", uri, verify)
	return nil
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

var whdaemonCommand = &cli.Command{
	Name:    "daemon",
	Aliases: []string{""},
	Usage:   "Start the webhook listener and a new subscription",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Value: 9003,
			Usage: "Port on which to listen"},
		&cli.StringFlag{
			Name:  "url",
			Value: "",
			Usage: "Address where webhook events will be sent (max length 255 characters"},
		verifyFlag,
	},
	Action: whdaemon,
}

var webhookCommand = &cli.Command{
	Name:  "webhook",
	Usage: "Manage webhook subscriptions",
	Subcommands: []*cli.Command{
		whdaemonCommand,
		whlistCommand,
		whserveCommand,
		whsubscribeCommand,
		whunsubscribeCommand,
	},
}
