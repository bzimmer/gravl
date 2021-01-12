package strava_test

// type TokenSource struct{}

// func (s TokenSource) Token() (*oauth2.Token, error) {
// 	return nil, nil
// }

// type Config struct {
// 	oauth2.Config
// }

// func (c *Config) TokenSource(ctx context.Context) oauth2.TokenSource {
// 	return TokenSource{}
// }

// func Test_Refresh(t *testing.T) {
// 	t.Parallel()
// 	a := assert.New(t)

// 	client, err := newClient(http.StatusOK, "refresh.json")
// 	a.NoError(err)

// 	// cfg := Config{}
// 	// err = strava.WithConfig(cfg)(client)
// 	// a.NoError(err)

// 	ctx := context.Background()
// 	tokens, err := client.Auth.Refresh(ctx)
// 	a.NoError(err, "failed to refresh")
// 	a.NotNil(tokens)
// }
