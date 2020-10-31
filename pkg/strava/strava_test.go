package strava_test

import (
	"github.com/bzimmer/gravl/pkg/common"
	"github.com/bzimmer/gravl/pkg/strava"
)

func newClient(status int, filename string) (*strava.Client, error) {
	return newClienter(status, filename, nil, nil)
}

func newClienter(status int, filename string, requester common.Requester, responder common.Responder) (*strava.Client, error) {
	return strava.NewClient(
		strava.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
			Requester:   requester,
			Responder:   responder}),
		strava.WithAPICredentials("fooKey", "barToken"))
}
