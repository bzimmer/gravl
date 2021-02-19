package activity

import (
	"io"
)

// File for uploading
type File struct {
	io.Reader `json:"-"`
	Name      string `json:"name"`
	Format    Format `json:"format"`
}

func (f *File) Close() error {
	if f.Reader == nil {
		return nil
	}
	if x, ok := f.Reader.(io.Closer); ok {
		return x.Close()
	}
	return nil
}
