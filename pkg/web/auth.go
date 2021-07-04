package web

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

// AuthHandler redirects to the oauth provider's credential acceptance page
func AuthHandler(c *oauth2.Config, state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := c.AuthCodeURL(state)
		http.Redirect(w, r, u, http.StatusFound)
	}
}

// AuthCallback receives the callback from the oauth provider with the credentials
func AuthCallbackHandler(c *oauth2.Config, state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s := r.Form.Get("state")
		if s != state {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}

		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := c.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(*token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
