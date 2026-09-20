package main

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeFaultError struct {
	code int
}

func (f *fakeFaultError) Error() string       { return "fake fault" }
func (f *fakeFaultError) HTTPStatusCode() int { return f.code }

func TestExitCode(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  error
		code int
	}{
		{name: "nil error", err: nil, code: 1},
		{name: "plain error", err: errors.New("boom"), code: 1},
		{name: "non-rate-limit fault", err: &fakeFaultError{code: http.StatusInternalServerError}, code: 1},
		{name: "rate limited fault", err: &fakeFaultError{code: http.StatusTooManyRequests}, code: exitTempFail},
		{
			name: "wrapped rate limited fault",
			err:  fmt.Errorf("wrap: %w", &fakeFaultError{code: http.StatusTooManyRequests}),
			code: exitTempFail,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, exitCode(tt.err))
		})
	}
}
