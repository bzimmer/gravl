package bunt_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/bunt"
	"github.com/bzimmer/gravl/pkg/store/internal"
)

func TestBuntStore(t *testing.T) {
	suite.Run(t, &internal.FileBasedTestSuite{
		Opener: func(path string) (store.Store, error) {
			return bunt.Open(path)
		},
	})
}
