package web

import (
	"encoding/json"
	"net/http"

	"github.com/bzimmer/gravl/version"
)

func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(map[string]string{
			"build_version":   version.BuildVersion,
			"build_timestamp": version.BuildTime,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
