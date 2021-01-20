package analysis

import "context"

// Context of the analysis pass
type Context struct {
	context.Context

	// Units to observe when performing analysis and returning results
	Units Units
}

// WithContext creates a new context using the parent with units
func WithContext(parent context.Context, units Units) *Context {
	return &Context{Context: parent, Units: units}
}
