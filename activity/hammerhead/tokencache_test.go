package hammerhead //nolint:testpackage // exercises unexported token-cache internals

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestLoadCachedTokenMissing(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	fs := afero.NewMemMapFs()
	a.Nil(loadCachedToken(fs, ""))
}

func TestLoadCachedTokenInvalidJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	fs := afero.NewMemMapFs()
	path, err := tokenCachePath("")
	a.NoError(err)
	a.NoError(fs.MkdirAll("/x", 0o700))
	a.NoError(afero.WriteFile(fs, path, []byte("not json"), 0o600))
	a.Nil(loadCachedToken(fs, ""))
}

func TestLoadCachedTokenMissingRefreshToken(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	fs := afero.NewMemMapFs()
	path, err := tokenCachePath("")
	a.NoError(err)
	a.NoError(afero.WriteFile(fs, path, []byte(`{"access_token":"foo"}`), 0o600))
	a.Nil(loadCachedToken(fs, ""))
}

func TestSaveAndLoadCachedTokenRoundTrip(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	fs := afero.NewMemMapFs()

	token := &oauth2.Token{AccessToken: "access", RefreshToken: "refresh"}
	a.NoError(saveCachedToken(fs, token, ""))

	loaded := loadCachedToken(fs, "")
	a.NotNil(loaded)
	a.Equal("access", loaded.AccessToken)
	a.Equal("refresh", loaded.RefreshToken)
}

func TestSaveCachedTokenReadOnly(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
	token := &oauth2.Token{AccessToken: "access", RefreshToken: "refresh"}
	a.Error(saveCachedToken(fs, token, ""))
}

func TestTokenCachePathOverride(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	path, err := tokenCachePath("/custom/cache")
	a.NoError(err)
	a.Equal("/custom/cache/hammerhead-token.json", path)
}

func TestSaveAndLoadCachedTokenCustomDir(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	fs := afero.NewMemMapFs()
	dir := "/custom/cache"

	token := &oauth2.Token{AccessToken: "custom-access", RefreshToken: "custom-refresh"}
	a.NoError(saveCachedToken(fs, token, dir))

	loaded := loadCachedToken(fs, dir)
	a.NotNil(loaded)
	a.Equal("custom-access", loaded.AccessToken)
	a.Equal("custom-refresh", loaded.RefreshToken)

	// default path should be empty
	a.Nil(loadCachedToken(fs, ""))
}
