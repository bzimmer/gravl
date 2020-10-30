package wx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/common/wx"
)

func Test_WindBearing(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	b, err := wx.WindBearing("W")
	a.NoError(err)
	a.Equal(270.0, b)
	b, err = wx.WindBearing("bar")
	a.Zero(b)
	a.Error(err)
	b, err = wx.WindBearing("")
	a.Zero(b)
	a.Error(err)
}
