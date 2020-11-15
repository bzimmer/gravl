package noaa

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Options(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	c, err := NewClient()
	a.NoError(err)
	a.NotNil(c)

	a.Equal(0*time.Second, c.client.Timeout)
	err = WithTimeout(15 * time.Millisecond)(c)
	a.NoError(err)
	a.Equal(15*time.Millisecond, c.client.Timeout)

	req, err := c.newAPIRequest(http.MethodGet, "/foo")
	a.NoError(err)
	a.NotNil(req)
	a.Equal("application/geo+json", req.Header.Get("Accept"))

	m := &http.Client{Timeout: 100 * time.Millisecond}
	a.NotEqual(m, c.client)
	err = WithHTTPClient(m)(c)
	a.NoError(err)
	a.Equal(m, c.client)
}
