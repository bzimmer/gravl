package hammerhead

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"golang.org/x/oauth2"
)

// tokenCachePath returns the file used to persist refreshed Hammerhead tokens
// across separate CLI invocations. Confirmed by testing (2026-08-10): Hammerhead
// rotates the refresh token on every use, so relying solely on a static
// HAMMERHEAD_REFRESH_TOKEN env var works exactly once; the cache carries the
// rotated token forward to the next invocation.
//
// dir overrides the directory; when empty the OS user config directory is used
// (e.g. ~/.config/gravl on Linux).
func tokenCachePath(dir string) (string, error) {
	if dir == "" {
		d, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(d, "gravl")
	}
	return filepath.Join(dir, Provider+"-token.json"), nil
}

func loadCachedToken(fs afero.Fs, dir string) *oauth2.Token {
	path, err := tokenCachePath(dir)
	if err != nil {
		return nil
	}
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil
	}
	var token oauth2.Token
	if err = json.Unmarshal(data, &token); err != nil || token.RefreshToken == "" {
		return nil
	}
	return &token
}

func saveCachedToken(fs afero.Fs, token *oauth2.Token, dir string) error {
	path, err := tokenCachePath(dir)
	if err != nil {
		return err
	}
	if err = fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(token) //nolint:gosec // oauth token cached on disk for the authenticated user
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, path, data, 0o600)
}
