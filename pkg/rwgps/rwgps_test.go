package rwgps_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bzimmer/gravl/pkg/common"
	rw "github.com/bzimmer/gravl/pkg/rwgps"
)

var (
	tests = map[string]string{
		"apikey":     "fooKey",
		"version":    "2",
		"auth_token": "barToken",
	}
)

func newClient(status int, filename string) (*rw.Client, error) {
	return rw.NewClient(
		rw.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
			Validator: func(req *http.Request) error {
				var body map[string]interface{}
				decoder := json.NewDecoder(req.Body)
				err := decoder.Decode(&body)
				if err != nil {
					return err
				}
				// confirm the body has the expected key:value pairs
				for key, value := range tests {
					v := body[key]
					if v != value {
						return fmt.Errorf("expected %s == '%v', not '%v'", key, value, v)
					}
				}
				return nil
			},
		}),
		rw.WithAPIKey("fooKey"),
		rw.WithAuthToken("barToken"),
		rw.WithAPIVersion(2),
		rw.WithAccept("application/json"),
	)
}
