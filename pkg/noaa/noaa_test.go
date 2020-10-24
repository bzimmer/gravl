package noaa

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestDataTransport struct {
	status   int
	filename string
}

func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var data []byte
	if t.filename != "" {
		data, _ = ioutil.ReadFile("testdata/" + t.filename)
	} else {
		data = make([]byte, 0)
	}

	header := make(http.Header)
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode:    t.status,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}, nil
}

func newClient(status int, filename string) (*Client, error) {
	return NewClient(
		WithTransport(&TestDataTransport{
			status:   status,
			filename: filename,
		}))
}

func Test_Options(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	c, err := newClient(http.StatusOK, "")
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
