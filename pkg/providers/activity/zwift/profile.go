package zwift

import (
	"context"
	"fmt"
	"net/http"
)

// ProfileService is the API for profile endpoints
type ProfileService service

const Me = "me"

// Profile returns the profile for the id
// The `profileID` can be empty or "me" to return the profile for the authenticated user
func (s *ProfileService) Profile(ctx context.Context, profileID string) (*Profile, error) {
	if profileID == "" {
		profileID = Me
	}
	uri := fmt.Sprintf("api/profiles/%s", profileID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	var profile *Profile
	if err = s.client.do(req, &profile); err != nil {
		return nil, err
	}
	return profile, err
}
