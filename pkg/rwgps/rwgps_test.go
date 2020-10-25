package rwgps_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bzimmer/wta/pkg/common"
	rw "github.com/bzimmer/wta/pkg/rwgps"
)

func newClient(status int, filename string) (*rw.Client, error) {
	return rw.NewClient(
		rw.WithTransport(&common.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
			RequestFunc: func(req *http.Request) error {
				var body map[string]interface{}
				decoder := json.NewDecoder(req.Body)
				err := decoder.Decode(&body)
				if err != nil {
					return err
				}
				if body["auth_token"] != "barToken" {
					return fmt.Errorf("expected authToke == 'barToken', not '%s'", body["auth_token"])
				}
				if body["apikey"] != "fooKey" {
					return fmt.Errorf("expected authToke == 'fooKey', not '%s'", body["apikey"])
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
