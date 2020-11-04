package strava_test

import (
	"context"
	"net/http"
	"testing"

	au "github.com/markbates/goth/providers/strava"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/bzimmer/gravl/pkg/strava"
)

func Test_Refresh(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	provider := au.New("foo", "bar", "", "read")
	provider.HTTPClient = &http.Client{
		Transport: &common.TestDataTransport{
			Status:      http.StatusOK,
			Filename:    "refresh.json",
			ContentType: "application/json",
		},
	}
	client, err := newClient(http.StatusOK, "")
	strava.WithProvider(provider)(client)

	ctx := context.Background()
	tokens, err := client.Auth.Refresh(ctx)
	a.NoError(err, "failed to refresh")
	a.NotNil(tokens)
	a.Equal("andthisbetherefreshtoken", tokens.RefreshToken)
	a.Equal("andthisbetheaccesstoken", tokens.AccessToken)
}
