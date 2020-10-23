package rwgps

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gj "github.com/paulmach/go.geojson"
)

// TripsService .
type TripsService service

// Route .
func (s *TripsService) Route(ctx context.Context, routeID int64) (*gj.FeatureCollection, error) {
	uri := fmt.Sprintf("routes/%d.json", routeID)
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	res := &TripResponse{}
	err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return newFeatureCollection("route", res.Route)
}

// Trip .
func (s *TripsService) Trip(ctx context.Context, tripID int64) (*gj.FeatureCollection, error) {
	uri := fmt.Sprintf("trips/%d.json", tripID)
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	res := &TripResponse{}
	err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return newFeatureCollection("trip", res.Trip)
}

func newFeatureCollection(activity string, trip *Trip) (*gj.FeatureCollection, error) {
	if trip == nil {
		return gj.NewFeatureCollection(), nil
	}

	coords := make([][]float64, len(trip.TrackPoints))
	feature := gj.NewFeature(gj.NewLineStringGeometry(coords))

	feature.ID = trip.ID
	feature.Properties["type"] = activity
	feature.Properties["name"] = trip.Name
	feature.Properties["user_id"] = trip.UserID
	feature.Properties["description"] = trip.Description
	feature.Properties["departed_at"] = trip.DepartedAt.Format(time.RFC3339)

	for i, tp := range trip.TrackPoints {
		coords[i] = []float64{tp.Longitude, tp.Latitude, tp.Elevation}
	}

	fc := gj.NewFeatureCollection()
	fc.AddFeature(feature)
	return fc, nil
}
