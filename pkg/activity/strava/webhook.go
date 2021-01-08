package strava

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//  API documented at https://developers.strava.com/docs/webhooks/

// WebhookService is the API for webhook endpoints
type WebhookService service

// WebhookAcknowledgement is the ack from Strava a webhook subscription has been received
type WebhookAcknowledgement struct {
	ID int64 `json:"id"`
}

// WebhookSubscription describes the details of webhook subscription
type WebhookSubscription struct {
	ID            int64     `json:"id"`
	ResourceState int       `json:"resource_state"`
	ApplicationID int       `json:"application_id"`
	CallbackURL   string    `json:"callback_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// WebhookMessage is the incoming webhook message
type WebhookMessage struct {
	ObjectType     string            `json:"object_type"`
	ObjectID       int               `json:"object_id"`
	AspectType     string            `json:"aspect_type"`
	OwnerID        int               `json:"owner_id"`
	SubscriptionID int64             `json:"subscription_id"`
	EventTime      int               `json:"event_time"`
	Updates        map[string]string `json:"updates"`
}

// WebhookSubscriber provides callbacks on webhook messages
type WebhookSubscriber interface {

	// SubscriptionRequest receives a callback during the subscription request flow
	SubscriptionRequest(challenge, verify string) error

	// MessageReceived is called every time a message is received from Strava
	MessageReceived(*WebhookMessage) error
}

// Subscribe to a webhook
func (s *WebhookService) Subscribe(ctx context.Context, callbackURL, verifyToken string) (*WebhookAcknowledgement, error) {
	uri := "push_subscriptions"
	req, err := s.client.newWebhookRequest(ctx, http.MethodPost, uri,
		map[string]string{"callback_url": callbackURL, "verify_token": verifyToken})
	if err != nil {
		return nil, err
	}
	ack := &WebhookAcknowledgement{}
	err = s.client.do(req, ack)
	if err != nil {
		return nil, err
	}
	return ack, err
}

// Unsubscribe to a webhook
func (s *WebhookService) Unsubscribe(ctx context.Context, subscriptionID int64) error {
	uri := fmt.Sprintf("push_subscriptions/%d", subscriptionID)
	// the empty body for credentials to be included in the request
	req, err := s.client.newWebhookRequest(ctx, http.MethodDelete, uri, map[string]string{})
	if err != nil {
		return err
	}
	return s.client.do(req, nil)
}

// List active webhook subscriptions
func (s *WebhookService) List(ctx context.Context) ([]*WebhookSubscription, error) {
	uri := fmt.Sprintf("push_subscriptions?client_id=%s&client_secret=%s",
		s.client.config.ClientID, s.client.config.ClientSecret)
	req, err := s.client.newWebhookRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	var subs []*WebhookSubscription
	err = s.client.do(req, &subs)
	if err != nil {
		return nil, err
	}
	return subs, err
}

// WebhookSubscriptionHandler handles subscription requests from Strava (GET)
func WebhookSubscriptionHandler(subscriber WebhookSubscriber) func(c *gin.Context) {
	return func(c *gin.Context) {
		verify, _ := c.GetQuery("hub.verify_token")
		challenge, _ := c.GetQuery("hub.challenge")
		if subscriber != nil {
			// if err is not nil the verification failed
			err := subscriber.SubscriptionRequest(challenge, verify)
			if err != nil {
				c.Abort()
				_ = c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"hub.challenge": challenge})
	}
}

// WebhookEventHandler receives the webhook callbacks from Strava (POST)
func WebhookEventHandler(subscriber WebhookSubscriber) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := &WebhookMessage{}
		err := c.BindJSON(m)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "ok"})
			return
		}
		if subscriber != nil {
			err := subscriber.MessageReceived(m)
			if err != nil {
				c.Abort()
				_ = c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
