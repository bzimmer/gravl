package cmd

import (
	"encoding/json"
	"io"
	"os"
)

func newEncoder(writer io.Writer, compact bool) *json.Encoder {
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

func newDecoder(reader io.Reader) *json.Decoder {
	if reader == nil {
		reader = os.Stdin
	}
	return json.NewDecoder(reader)
}
