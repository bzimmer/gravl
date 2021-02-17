package fake_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/fake"
	"github.com/bzimmer/gravl/pkg/store/internal"
)

func TestFakeStore(t *testing.T) {
	tests := []struct {
		fuzz bool
		acts int
		name string
	}{
		{name: "fuzzing enabled", acts: 108, fuzz: true},
		{name: "fuzzing disabled", acts: 108, fuzz: false},
	}
	for _, tt := range tests {
		w := tt
		suite.Run(t, &internal.FileBasedTestSuite{
			Persistent: false,
			Opener: func(path string) (store.Store, error) {
				return fake.Open(w.acts, w.fuzz)
			},
		})
	}
}
