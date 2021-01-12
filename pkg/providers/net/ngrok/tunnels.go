package ngrok

import (
	"context"
	"net/http"
)

// TunnelsService returns all active tunnels for the ngrok daemon running on localhost
type TunnelsService service

// Tunnels returns established tunnels
func (s *TunnelsService) Tunnels(ctx context.Context) ([]*Tunnel, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "tunnels")
	if err != nil {
		return nil, err
	}

	type Response struct {
		Tunnels []*Tunnel `json:"tunnels"`
		URI     string    `json:"uri"`
	}

	var tunnels Response
	err = s.client.do(req, &tunnels)
	if err != nil {
		return nil, err
	}
	return tunnels.Tunnels, err
}
