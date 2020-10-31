package strava_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// func Test_WebhookSubscribe(t *testing.T) {
// 	t.Parallel()
// 	a := assert.New(t)
// 	client, err := NewTestClient(func(req *http.Request) (*http.Response, error) {
// 		a.Equal(http.MethodPost, req.Method)
// 		a.Equal("clientID", req.FormValue("client_id"))
// 		a.Equal("clientSecret", req.FormValue("client_secret"))
// 		a.Equal("https://example.com/wh/callback", req.FormValue("callback_url"))
// 		a.Equal("verifyToken123", req.FormValue("verify_token"))
// 		return &http.Response{
// 			StatusCode: http.StatusOK,
// 			Body: ioutil.NopCloser(bytes.NewBufferString(`{
// 				"application_id":31898,
// 				"callback_url":"https://4cbd9ddb9748.ngrok.io/strava/webhook",
// 				"created_at":"2020-10-06T01:27:31Z",
// 				"id":165928,
// 				"resource_state":2,
// 				"updated_at":"2020-10-06T01:27:31Z"}`)),
// 			Header: make(http.Header),
// 		}, nil
// 	})
// 	err = WithWebhookCredentials("clientID", "clientSecret")(client)
// 	a.NoError(err)
// 	ctx := context.Background()
// 	msg, err := client.Webhook.Subscribe(ctx, "https://example.com/wh/callback", "verifyToken123")
// 	a.NoError(err)
// 	a.NotNil(msg)
// }

func Test_WebhookUnsubscribe(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusNoContent, "")
	ctx := context.Background()
	err = client.Webhook.Unsubscribe(ctx, 882722)
	a.NoError(err)
}

func Test_WebhookSubscriptions(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "subscriptions.json")
	ctx := context.Background()
	msgs, err := client.Webhook.Subscriptions(ctx)
	a.NoError(err)
	a.NotNil(msgs)
	a.Equal(1, len(*msgs))
}

type TestSubscriber struct {
	verify    string
	challenge string
	msg       *strava.WebhookMessage
	fail      bool
}

func (t *TestSubscriber) SubscriptionRequest(challenge string, verify string) error {
	t.verify = verify
	t.challenge = challenge
	if t.fail {
		return errors.New("failed")
	}
	return nil
}

func (t *TestSubscriber) MessageReceived(msg *strava.WebhookMessage) error {
	t.msg = msg
	if t.fail {
		return errors.New("failed")
	}
	return nil
}

func setupTestRouter() (*TestSubscriber, *gin.Engine) {
	sub := &TestSubscriber{fail: false}
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/webhook", strava.WebhookSubscriptionHandler(sub))
	r.POST("/webhook", strava.WebhookEventHandler(sub))
	return sub, r
}

func Test_WebhookEventHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub, router := setupTestRouter()

	reader := strings.NewReader(`
	{
		"aspect_type": "update",
		"event_time": 1516126040,
		"object_id": 1360128428,
		"object_type": "activity",
		"owner_id": 18637089,
		"subscription_id": 120475,
		"updates": {
			"title": "Messy",
			"type": "Bike"

		}
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/webhook", reader)
	router.ServeHTTP(w, req)

	a.Equal(200, w.Code)
	a.NotNil(sub)
	a.Equal(18637089, sub.msg.OwnerID)
	a.Equal("Bike", sub.msg.Updates["type"])

	reader = strings.NewReader(`
	{
		"aspect_type": "update",
		"event_time": 1516126040,
		"object_id": 1360128428,
		"object_type": "activity",
		"owner_id": 18637089,
		"subscription_id": 120475,
		"updates": {
			"title": "Messy",
			"type": "Bike"

		}
	}`)

	sub.fail = true
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/webhook", reader)
	router.ServeHTTP(w, req)

	a.Equal(500, w.Code)
}

func Test_WebhookSubscriptionHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub, router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/webhook?hub.verify_token=bar&hub.challenge=baz", nil)
	router.ServeHTTP(w, req)

	a.Equal(200, w.Code)
	a.NotNil(sub)
	a.Equal("bar", sub.verify)
	a.Equal("baz", sub.challenge)

	sub.fail = true
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/webhook?hub.verify_token=bar&hub.challenge=baz", nil)
	router.ServeHTTP(w, req)

	a.Equal(500, w.Code)
}
