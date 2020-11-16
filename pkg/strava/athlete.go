package strava

import (
	"context"
	"fmt"
	"net/http"
)

// AthleteService is the API for athlete endpoints
type AthleteService service

// Athlete returns the currently authenticated athlete
func (s *AthleteService) Athlete(ctx context.Context) (*Athlete, error) {
	uri := "athlete"
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	ath := &Athlete{}
	err = s.client.Do(req, ath)
	if err != nil {
		return nil, err
	}
	return ath, err
}

// Stats returns the activity stats of an athlete. Only includes data
// from activities set to Everyone visibilty.
func (s *AthleteService) Stats(ctx context.Context, id int) (*Stats, error) {
	uri := fmt.Sprintf("athletes/%d/stats", id)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	sts := &Stats{}
	err = s.client.Do(req, sts)
	if err != nil {
		return nil, err
	}
	return sts, err
}
