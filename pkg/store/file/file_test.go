package file_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/file"
	"github.com/bzimmer/gravl/pkg/store/internal"
)

func TestFileStore(t *testing.T) {
	suite.Run(t, &internal.FileBasedTestSuite{
		Opener: func(path string) (store.Store, error) {
			return file.Open(path, file.Flush(true))
		},
	})
}
