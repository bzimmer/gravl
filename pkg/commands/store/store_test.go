package store_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

const N = 1122

func tempfile(t *testing.T, pattern string) *os.File {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("file (%s): %s", pattern, f.Name())
	return f
}

func TestStoreIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	a := assert.New(t)

	var err error
	c := internal.Gravl("-c", "store", "export", "-i", fmt.Sprintf("fake,n=%d", N))
	<-c.Start()
	a.True(c.Success())

	f := tempfile(t, "fake_input")
	defer os.Remove(f.Name())

	n := len(c.Status().Stdout)
	a.Equal(N, n)
	_, err = f.WriteString(c.Stdout())
	a.NoError(err)
	err = f.Close()
	a.NoError(err)

	activityfile := fmt.Sprintf("file,path=%s", f.Name())
	c = internal.Gravl("-c", "store", "export", "-i", activityfile)
	<-c.Start()
	a.True(c.Success())
	t.Logf("%s %v", c.Name, c.Args)

	f = tempfile(t, "file_export")
	_, err = f.WriteString(c.Stdout())
	a.NoError(err)
	err = f.Close()
	a.NoError(err)

	p := len(c.Status().Stdout)
	a.Equal(n, p)

	f = tempfile(t, "bunt_input")
	buntfile := fmt.Sprintf("bunt,path=%s", f.Name())
	c = internal.Gravl("store", "update", "-i", activityfile, "-o", buntfile)
	<-c.Start()
	a.True(c.Success())

	err = f.Close()
	a.NoError(err)

	tests := []struct {
		n          int
		err, array bool
		name       string
		args       []string
	}{
		{name: "export all activities as attribute", n: N,
			args: []string{"-c", "store", "export", "-i", buntfile, "-B", ".ID, .Name, .StartDateLocal"}},
		{name: "export all activities", n: N, args: []string{"-c", "store", "export", "-i", buntfile}},
		{name: "remove requires a filter", err: true, args: []string{"-c", "store", "remove", "-i", buntfile}},
		{name: "delete all", n: N, array: true, args: []string{"-c", "store", "remove", "-f", "true", "-i", buntfile}},
		{name: "delete but nothing exists", n: 0, array: true, args: []string{"-c", "store", "remove", "-f", "true", "-i", buntfile}},
		{name: "export all activities", n: 0, args: []string{"-c", "store", "export", "-i", buntfile}},
	}
	for _, tt := range tests {
		tt := tt
		c = internal.Gravl(tt.args...)
		<-c.Start()
		switch tt.err {
		case true:
			a.False(c.Success())
		case false:
			a.True(c.Success())
			var x int
			gjson.ForEachLine(c.Stdout(), func(res gjson.Result) bool {
				if tt.array {
					a.True(res.IsArray())
					x += len(res.Array())
				} else {
					x++
				}
				return true
			})
			a.Equal(tt.n, x)
		}
	}
}
