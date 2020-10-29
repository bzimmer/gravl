package common_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/common"
)

func Test_VerboseTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := []http.RoundTripper{
		&common.VerboseTransport{
			Transport: &common.TestDataTransport{
				Filename:    "transport.txt",
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		},
		&common.VerboseTransport{
			Event: log.Debug(),
			Transport: &common.TestDataTransport{
				Filename:    "transport.txt",
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		},
	}
	for _, transport := range tests {
		client := http.Client{
			Transport: transport,
		}
		res, err := client.Get("http://example.com")
		a.NoError(err)
		a.NotNil(res)

		body, err := ioutil.ReadAll(res.Body)
		a.NoError(err)
		a.NotNil(res)
		a.Equal(
			`The mountains are calling & I must go & I will work on while I can, studying incessantly.`,
			strings.Trim(string(body), "\n"),
		)
	}
}

func Test_TestDataTransport(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	tests := [][]string{
		{"transport.txt", `The mountains are calling & I must go & I will work on while I can, studying incessantly.`},
		{"", ""}}

	for _, test := range tests {
		client := http.Client{
			Transport: &common.TestDataTransport{
				Filename:    test[0],
				Status:      http.StatusOK,
				ContentType: "text/plain",
			},
		}
		res, err := client.Get("http://example.com")
		a.NoError(err)
		a.NotNil(res)

		body, err := ioutil.ReadAll(res.Body)
		a.NoError(err)
		a.NotNil(res)
		a.Equal(test[1], strings.Trim(string(body), "\n"))
	}

	client := http.Client{
		Transport: &common.TestDataTransport{
			Filename:    "~garbage~",
			Status:      http.StatusOK,
			ContentType: "text/plain",
		},
	}
	res, err := client.Get("http://example.com")
	a.Error(err)
	a.Nil(res)
}
