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

// SubscriptionAcknowledgement describes the details of webhook subscription
type SubscriptionAcknowledgement struct {
	ID int `json:"id"`
}

// Subscription describes the details of webhook subscription
type Subscription struct {
	ID            int       `json:"id"`
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
	SubscriptionID int               `json:"subscription_id"`
	EventTime      int               `json:"event_time"`
	Updates        map[string]string `json:"updates"`
}

// WebhookSubscriber .
type WebhookSubscriber interface {

	// SubscriptionRequest receives a callback during the subscription request flow
	SubscriptionRequest(challenge, verify string) error

	// MessageReceived is called every time a message is received from Strava
	MessageReceived(*WebhookMessage) error
}

// Subscribe to a webhook
func (s *WebhookService) Subscribe(ctx context.Context, callbackURL, verifyToken string) (*SubscriptionAcknowledgement, error) {
	uri := "push_subscriptions"
	req, err := s.client.newWebhookRequest(ctx, http.MethodPost, uri,
		map[string]string{"callback_url": callbackURL, "verify_token": verifyToken})
	if err != nil {
		return nil, err
	}
	ack := &SubscriptionAcknowledgement{}
	err = s.client.Do(req, ack)
	if err != nil {
		return nil, err
	}
	return ack, err
}

// Unsubscribe to a webhook
func (s *WebhookService) Unsubscribe(ctx context.Context, subscriptionID int) error {
	uri := fmt.Sprintf("push_subscriptions/%d", subscriptionID)
	req, err := s.client.newWebhookRequest(ctx, http.MethodDelete, uri, make(map[string]string))
	if err != nil {
		return err
	}
	return s.client.Do(req, nil)
}

// Subscriptions returns a list of subscriptions
func (s *WebhookService) Subscriptions(ctx context.Context) (*[]Subscription, error) {
	uri := fmt.Sprintf("push_subscriptions?client_id=%s&client_secret=%s",
		s.client.config.ClientID, s.client.config.ClientSecret)
	req, err := s.client.newWebhookRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	subs := &[]Subscription{}
	err = s.client.Do(req, subs)
	if err != nil {
		return nil, err
	}
	return subs, err
}

// WebhookSubscriptionHandler handles subscription requests from Strava
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

// WebhookEventHandler receives the webhook callbacks from Strava
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
