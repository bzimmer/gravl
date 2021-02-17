package internal

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
)

const id = int64(823722321)

// FileBasedTestSuite for file based stores
type FileBasedTestSuite struct {
	suite.Suite
	Persistent bool
	Opener     func(path string) (store.Store, error)
}

// tempfile creates new tempfiles
// The caller is responsible for deleting the file
func (s *FileBasedTestSuite) tempfile() (*os.File, func()) {
	f, err := ioutil.TempFile("", "FileBasedTestSuite")
	if err != nil {
		s.T().FailNow()
	}
	return f, func() {
		if s.T().Failed() {
			s.T().Logf("test failed; not removing temp file: %s", f.Name())
			return
		}
		_ = os.Remove(f.Name())
	}
}

// TestRetrieveMissingActivity tests retrieving an activity which does not exist
func (s *FileBasedTestSuite) TestRetrieveMissingActivity() {
	a := s.Assert()

	f, remove := s.tempfile()
	defer remove()
	db, err := s.Opener(f.Name())
	a.NoError(err)
	a.NotNil(db)

	ctx := context.Background()
	act, err := db.Activity(ctx, 299289299288234)
	a.Equal(store.ErrNotFound, err)
	a.Nil(act)
	a.NoError(db.Close())
}

// TestCancel tests reading from the `Done` channel of a context
func (s *FileBasedTestSuite) TestCancel() {
	a := s.Assert()
	ctx := context.Background()

	f, remove := s.tempfile()
	defer remove()
	db, err := s.Opener(f.Name())
	a.NoError(err)
	a.NotNil(db)

	err = db.Save(ctx, &strava.Activity{ID: id, Name: "foobar"})
	a.NoError(err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	cancel() // force the store to deal with the cancel

	acts := db.Activities(ctx)
	for {
		select {
		case res, ok := <-acts:
			if !ok {
				return
			}
			// ensure the cancel was handled; activities might have
			// been produced in the meantime (`select` is random) but
			// ignore them and wait for the canceled context error
			if res.Err != nil {
				a.Equal(context.Canceled, res.Err)
				a.Nil(res.Activity)
				return
			}
		case <-time.After(time.Millisecond * 500):
			a.FailNow("should have handled the cancel by now")
		}
	}
}

func (s *FileBasedTestSuite) TestAddRemove() {
	a := s.Assert()
	ctx := context.Background()

	f, remove := s.tempfile()
	defer remove()
	db, err := s.Opener(f.Name())
	a.NoError(err)
	a.NotNil(db)

	err = db.Save(ctx, &strava.Activity{ID: id, Name: "foobar"})
	a.NoError(err)
	act, err := db.Activity(ctx, id)
	a.NoError(err)
	ok, err := db.Exists(ctx, id)
	a.NoError(err)
	a.True(ok)
	err = db.Remove(ctx, act)
	a.NoError(err)
	ok, err = db.Exists(ctx, id)
	a.NoError(err)
	a.False(ok)
	a.NoError(db.Close())
}

func (s *FileBasedTestSuite) TestLifecycle() {
	if !s.Persistent {
		s.T().Skip("skipping lifecycle because store is not persistent")
		return
	}
	a := s.Assert()
	ctx := context.Background()

	f, remove := s.tempfile()
	defer remove()
	db, err := s.Opener(f.Name())
	a.NoError(err)
	a.NotNil(db)

	err = db.Save(ctx, &strava.Activity{ID: id, Name: "foobar"})
	a.NoError(err)
	err = db.Close()
	a.NoError(err)

	db, err = s.Opener(f.Name())
	a.NoError(err)
	acts := db.Activities(ctx)
	res := <-acts
	a.NoError(res.Err)
	a.Equal(id, res.Activity.ID)
	a.Equal("foobar", res.Activity.Name)

	act, err := db.Activity(ctx, id)
	a.NoError(err)
	a.Equal(id, act.ID)
	a.Equal("foobar", act.Name)

	ok, err := db.Exists(ctx, id)
	a.NoError(err)
	a.True(ok)
	ok, err = db.Exists(ctx, 99887766)
	a.NoError(err)
	a.False(ok)
	a.NoError(db.Close())
}
