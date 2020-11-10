package strava_test

import (
	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/transport"
)

func newClient(status int, filename string) (*strava.Client, error) {
	return newClienter(status, filename, nil, nil)
}

func newClienter(status int, filename string, requester transport.Requester, responder transport.Responder) (*strava.Client, error) {
	return strava.NewClient(
		strava.WithTransport(&transport.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
			Requester:   requester,
			Responder:   responder}),
		strava.WithAPICredentials("fooKey", "barToken"))
}
