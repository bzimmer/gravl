package cyclinganalytics

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RidesService .
type RidesService service

type RideOptions struct {
	Streams []string
	Curves  struct {
		AveragePower   bool
		EffectivePower bool
	}
}

func (r *RideOptions) values() *url.Values {
	v := &url.Values{}
	if r.Streams != nil {
		v.Set("streams", strings.Join(r.Streams, ","))
	}
	if r.Curves.AveragePower && r.Curves.EffectivePower {
		v.Set("curves", "true")
	} else {
		v.Set("power_curve", fmt.Sprintf("%t", r.Curves.AveragePower))
		v.Set("epower_curve", fmt.Sprintf("%t", r.Curves.EffectivePower))
	}
	return v
}

func (s *RidesService) Ride(ctx context.Context, rideID int64, opts RideOptions) (*Ride, error) {
	uri := fmt.Sprintf("ride/%d", rideID)

	params := opts.values()
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, params)
	if err != nil {
		return nil, err
	}
	ride := &Ride{}
	err = s.client.do(req, ride)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

func (s *RidesService) Rides(ctx context.Context, userID UserID) ([]*Ride, error) {
	uri := "me/rides"
	if userID == Me {
		uri = fmt.Sprintf("%d/rides", userID)
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Rides []*Ride `json:"rides"`
	}
	res := &Response{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res.Rides, nil
}

// // ValidStream returns true if the strean name is valid
// func ValidStream(stream string) bool { // nolint
// 	// https://www.cyclinganalytics.com/developer/api#/ride/ride_id
// 	switch stream {
// 	case "cadence":
// 		return true
// 	case "distance":
// 		// The sequence of distance values for this stream, in kilometers [float]
// 		return true
// 	case "elevation":
// 		// The sequence of elevation values for this stream, in meters [float]
// 		return true
// 	case "gears":
// 		return true
// 	case "gradient":
// 		// The sequence of grade values for this stream, as percents of a grade [float]
// 		return true
// 	case "heart_rate_variability":
// 		return true
// 	case "heartrate":
// 		// The sequence of heart rate values for this stream, in beats per minute [integer]
// 		return true
// 	case "latitude":
// 		return true
// 	case "longitude":
// 		return true
// 	case "lrbalance":
// 		return true
// 	case "pedal_smoothness":
// 		return true
// 	case "platform_center_offset":
// 		return true
// 	case "power":
// 		return true
// 	case "power_direction":
// 		return true
// 	case "power_phase":
// 		return true
// 	case "respiration_rate":
// 		return true
// 	case "smo2":
// 		return true
// 	case "speed":
// 		// The sequence of speed values for this stream, in meters per second [float]
// 		return true
// 	case "temperature":
// 		// The sequence of temperature values for this stream, in celsius degrees [float]
// 		return true
// 	case "thb":
// 		return true
// 	case "torque_effectiveness":
// 		return true
// 	}
// 	return false
// }
