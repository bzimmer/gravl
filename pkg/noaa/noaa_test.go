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
	WithTimeout(15 * time.Millisecond)(c)
	a.Equal(15*time.Millisecond, c.client.Timeout)

	a.Equal("application/geo+json", c.header.Get("Accept"))
	req, err := c.newAPIRequest(http.MethodGet, "/foo")
	a.NoError(err)
	a.NotNil(req)
	a.Equal("application/geo+json", req.Header.Get("Accept"))
	WithAccept("application/foobar")(c)
	req, err = c.newAPIRequest(http.MethodGet, "/bar")
	a.NoError(err)
	a.NotNil(req)
	a.Equal("application/foobar", req.Header.Get("Accept"))

	m := &http.Client{Timeout: 100 * time.Millisecond}
	a.NotEqual(m, c.client)
	WithHTTPClient(m)(c)
	a.Equal(m, c.client)
}
