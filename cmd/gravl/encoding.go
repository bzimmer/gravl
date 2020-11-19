package gravl

import (
	"encoding/json"
	"io"
	"os"
)

func newEncoder(writer io.Writer, compact bool) *json.Encoder {
	if writer == nil {
		writer = os.Stdout
	}
	m := json.NewEncoder(writer)
	if !compact {
		m.SetIndent("", " ")
	}
	m.SetEscapeHTML(false)
	return m
}

func newDecoder(reader io.Reader) *json.Decoder {
	if reader == nil {
		reader = os.Stdin
	}
	return json.NewDecoder(reader)
}
