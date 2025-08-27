package errutil

import (
	"strings"
	"sync"
)

// ErrorGroup aggregates multiple errors into one.
type ErrorGroup struct {
	mu     sync.RWMutex
	errors []error
	sep    string
}

// ErrorGroupOption customizes ErrorGroup behavior.
type ErrorGroupOption func(*ErrorGroup)

// WithSeparator sets the string used between error messages.
func WithSeparator(sep string) ErrorGroupOption {
	return func(g *ErrorGroup) {
		g.sep = sep
	}
}

// NewErrorGroup constructs an empty ErrorGroup.
func NewErrorGroup(opts ...ErrorGroupOption) *ErrorGroup {
	g := &ErrorGroup{
		mu:     sync.RWMutex{},
		errors: []error{},
		sep:    "; ",
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Add appends a non-nil error to the group.
func (g *ErrorGroup) Add(err error) {
	if err == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.errors = append(g.errors, err)
}

// Errors returns a snapshot of all collected errors.
func (g *ErrorGroup) Errors() []error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	snap := make([]error, len(g.errors))
	copy(snap, g.errors)
	return snap
}

// Error implements error by joining sub-errors with the separator.
func (g *ErrorGroup) Error() string {
	slice := g.Errors()
	parts := make([]string, len(slice))
	for i, e := range slice {
		parts[i] = e.Error()
	}
	return strings.Join(parts, g.sep)
}
