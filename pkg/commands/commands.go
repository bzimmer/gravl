package commands

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/bzimmer/gravl/pkg/eval"
	"github.com/bzimmer/gravl/pkg/eval/antonmedv"
)

// Filterer returns a filterer for the expression
func Filterer(q string) (eval.Filterer, error) {
	return antonmedv.Filterer(q)
}

// Mapper returns a mapper for the expression
func Mapper(q string) (eval.Mapper, error) {
	return antonmedv.Mapper(q)
}

// Evaluator returns an evaluator for the expression
func Evaluator(q string) (eval.Evaluator, error) {
	return antonmedv.Evaluator(q)
}

// Token produces a random token of length `n`
func Token(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
