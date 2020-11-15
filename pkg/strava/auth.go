package strava

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog/log"
)

const (
	authKey = "_stravaAuthTokens"
)

// AuthService is the API for auth endpoints
type AuthService service

// Refresh returns a new access token
func (s *AuthService) Refresh(ctx context.Context) (*AuthTokens, error) {
	token, err := s.client.provider.RefreshToken(s.client.refreshToken)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &AuthTokens{
		UpdatedAt:    now,
		ExpiresAt:    token.Expiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		CreatedAt:    now,
	}, nil
}

// AuthRequired is gin middleware for ensuring the user is authenticated
func AuthRequired(provider goth.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		authTokensJSON := session.Get(authKey)
		if authTokensJSON == nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization keys found"})
			return
		}

		authTokens := &AuthTokens{}
		err := json.Unmarshal(authTokensJSON.([]byte), authTokens)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session value"})
			return
		}

		now := time.Now()
		if authTokens.ExpiresAt.After(now) {
			// we have a valid token
			c.Set(authKey, authTokens)
			return
		}

		log.Warn().Time("expiresAt", authTokens.ExpiresAt).Time("now", now).Msg("authrequired")

		token, err := provider.RefreshToken(authTokens.RefreshToken)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		authTokens = &AuthTokens{
			UpdatedAt:    time.Now(),
			ExpiresAt:    token.Expiry,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			CreatedAt:    authTokens.CreatedAt,
			AthleteID:    authTokens.AthleteID,
		}

		err = saveTokens(c, authTokens)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.Set(authKey, authTokens)
	}
}

// AuthHandler starts the OAuth session
func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

// AuthCallbackHandler supports the callback from Strava
func AuthCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		log.Info().
			Str("athleteID", user.UserID).
			Time("expiresAt", user.ExpiresAt).
			Send()

		athleteID, err := strconv.Atoi(user.UserID)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		authTokens := &AuthTokens{
			AthleteID:    athleteID,
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
			ExpiresAt:    user.ExpiresAt,
		}

		err = saveTokens(c, authTokens)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, authTokens)
	}
}

// saveTokens saves the auth tokens to the database and session
func saveTokens(c *gin.Context, authToken *AuthTokens) error {
	session := sessions.Default(c)
	accessTokenJSON, err := json.Marshal(authToken)
	if err != nil {
		return err
	}

	session.Set(authKey, accessTokenJSON)
	err = session.Save()
	if err != nil {
		return err
	}
	return nil
}
