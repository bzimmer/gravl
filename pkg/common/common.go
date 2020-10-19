package common

import "net/http"

// RoundTripperFunc wraps a func to make it into a http.RoundTripper. Similar to http.HandleFunc.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// // WithVerboseLogging .
// func WithVerboseLogging(debug bool) func(*Client) error {
// 	return func(client *Client) error {
// 		if !debug {
// 			return nil
// 		}
// 		transport := client.client.Transport
// 		if transport == nil {
// 			transport = http.DefaultTransport
// 		}
// 		client.client.Transport = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
// 			dump, _ := httputil.DumpRequestOut(req, true)
// 			log.Debug().Str("req", string(dump)).Msg("sending")
// 			res, err := transport.RoundTrip(req)
// 			dump, _ = httputil.DumpResponse(res, true)
// 			log.Debug().Str("res", string(dump)).Msg("received")
// 			return res, err
// 		})
// 		return nil
// 	}
// }
