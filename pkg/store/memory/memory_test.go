package memory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/internal"
	"github.com/bzimmer/gravl/pkg/store/memory"
)

type provider struct{}

func (p *provider) Activities() ([]*strava.Activity, error) {
	return []*strava.Activity{
		{ID: 1, Name: "foo"},
		{ID: 2, Name: "bar"},
		{ID: 3, Name: "baz"},
	}, nil
}

func (p *provider) Close(map[int64]*strava.Activity) error {
	return nil
}

func TestFileStore(t *testing.T) {
	suite.Run(t, &internal.FileBasedTestSuite{
		Opener: func(path string) (store.Store, error) {
			return memory.Open(&provider{})
		},
	})
}
