package web

import (
	"encoding/json"
	"net/http"

	"github.com/bzimmer/gravl/pkg"
)

// VersionHandler .
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(map[string]string{
			"build_version":   pkg.BuildVersion,
			"build_timestamp": pkg.BuildTime,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
