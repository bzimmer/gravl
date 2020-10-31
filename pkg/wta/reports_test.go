package wta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_query(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	q := query("foobar")
	a.NotNil(q)
	a.Equal("author=foobar&b_size=100&filter=Search&hiketypes%3Alist=day-hike&hiketypes%3Alist=multi-night-backpack&hiketypes%3Alist=overnight&hiketypes%3Alist=snowshoe-xc-ski&month=all&subregion=all", q.RawQuery)
}
