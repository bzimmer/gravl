package strava

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fault(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	f := func() error {
		return &Fault{Message: "foo"}
	}
	err := f()
	a.Error(err)
	a.Equal("foo", err.Error())
}
