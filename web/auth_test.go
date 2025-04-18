package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/web"
)

func TestAuthHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mux := http.NewServeMux()
	svr := httptest.NewServer(mux)
	defer svr.Close()

	cfg := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL: svr.URL + "/auth",
		},
	}
	mux.HandleFunc("/auth", web.AuthHandler(cfg, "foo-state-bar"))

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/auth", http.NoBody)
	a.NoError(err)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusFound, w.Code)

	res := w.Result()
	defer res.Body.Close()
	loc := res.Header.Get("Location")
	u, err := url.Parse(loc)
	a.NoError(err)
	a.Contains(loc, svr.URL+"/auth")
	a.Equal("foo-state-bar", u.Query().Get("state"))
}

func TestAuthCallbackHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []struct {
		name  string
		state string
		code  string
		res   int
		json  bool
	}{
		{
			name:  "success",
			state: "foo-state-bar",
			code:  "foo-code-bar",
			res:   http.StatusOK,
		},
		{
			name: "no state",
			code: "foo-code-bar",
			res:  http.StatusBadRequest,
		},
		{
			name:  "no code",
			state: "foo-state-bar",
			res:   http.StatusBadRequest,
		},
		{
			name:  "bad json",
			json:  true,
			state: "foo-state-bar",
			code:  "foo-code-bar",
			res:   http.StatusInternalServerError,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux := http.NewServeMux()
			svr := httptest.NewServer(mux)
			defer svr.Close()

			cfg := &oauth2.Config{
				Endpoint: oauth2.Endpoint{
					TokenURL: svr.URL + "/token",
				},
				RedirectURL: svr.URL + "/callback",
			}
			mux.HandleFunc("/callback", web.AuthCallbackHandler(cfg, "foo-state-bar"))
			mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
				var data []byte
				if tt.json {
					data = []byte("garbage")
				} else {
					data = []byte(`{
					"access_token":"99881100332255",
					"token_type":"bearer",
					"expires_in":3600,
					"refresh_token":"ThisIsGood"}`)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Content-Length", strconv.Itoa(len(data)))
				n, err := w.Write(data)
				a.NoError(err)
				a.Equal(n, len(data))
				w.WriteHeader(http.StatusOK)
			})

			form := url.Values{}
			form.Set("state", tt.state)
			form.Set("code", tt.code)
			data := form.Encode()
			body := strings.NewReader(data)

			ctx := t.Context()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/callback", body)
			a.NoError(err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", strconv.Itoa(len(data)))

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			a.Equal(tt.res, w.Code)
			if tt.res != http.StatusOK {
				return
			}

			res := w.Result()
			defer res.Body.Close()

			token := new(oauth2.Token)
			dec := json.NewDecoder(res.Body)
			a.NoError(dec.Decode(&token))
			a.Equal("99881100332255", token.AccessToken)
		})
	}
}
