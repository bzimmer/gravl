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

const (
	tripType  = "trip"
	routeType = "route"
)

// Trip .
func (s *TripsService) Trip(ctx context.Context, tripID int64) (*gj.FeatureCollection, error) {
	return s.trip(ctx, tripType, fmt.Sprintf("trips/%d.json", tripID))
}

// Route .
func (s *TripsService) Route(ctx context.Context, routeID int64) (*gj.FeatureCollection, error) {
	return s.trip(ctx, routeType, fmt.Sprintf("routes/%d.json", routeID))
}

func (s *TripsService) trip(ctx context.Context, activity, uri string) (*gj.FeatureCollection, error) {
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	res := &TripResponse{}
	err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	switch activity {
	case tripType:
		return newFeatureCollection(activity, res.Trip)
	case routeType:
		return newFeatureCollection(activity, res.Route)
	}
	return nil, fmt.Errorf("unknown activity type {%s}", activity)
}

func newFeatureCollection(activity string, trip *Trip) (*gj.FeatureCollection, error) {
	fc := gj.NewFeatureCollection()

	if trip == nil {
		return fc, nil
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

	fc.AddFeature(feature)
	return fc, nil
}
