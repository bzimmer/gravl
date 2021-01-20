package analysis

import "context"

type Context struct {
	context.Context

	// Units of the resulting Activities
	Units Units
}

// WithContext creates a new context using the parent with units
func WithContext(parent context.Context, units Units) *Context {
	return &Context{Context: parent, Units: units}
}
