package common

import (
	"encoding/json"
	"io"
	"os"
)

// NewEncoder .
func NewEncoder(writer io.Writer, compact bool) *json.Encoder {
	if writer == nil {
		writer = os.Stdout
	}
	encoder := json.NewEncoder(writer)
	if !compact {
		encoder.SetIndent("", " ")
	}
	encoder.SetEscapeHTML(false)
	return encoder
}

// NewDecoder .
func NewDecoder(reader io.Reader) *json.Decoder {
	if reader == nil {
		reader = os.Stdin
	}
	return json.NewDecoder(reader)
}
