package store_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
	t.Logf("%s %v", c.Name, c.Args)

	err = f.Close()
	a.NoError(err)

	c = internal.Gravl("-c", "store", "export", "-i", buntfile)
	<-c.Start()
	a.True(c.Success())
	p = len(c.Status().Stdout)
	a.Equal(n, p)
	a.Equal(N, p)
}
