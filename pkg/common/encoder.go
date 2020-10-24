package common

import (
	"encoding/json"
	"os"
)

// NewEncoder .
func NewEncoder(compact bool) *json.Encoder {
	encoder := json.NewEncoder(os.Stdout)
	if !compact {
		encoder.SetIndent("", " ")
	}
	encoder.SetEscapeHTML(false)
	return encoder
}

// NewDecoder .
func NewDecoder() *json.Decoder {
	return json.NewDecoder(os.Stdin)
}
