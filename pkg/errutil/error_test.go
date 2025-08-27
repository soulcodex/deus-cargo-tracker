package errutil_test

import (
	"testing"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sentinel struct{}

func (s sentinel) Error() string { return "sentinel" }

func TestBaseError_Basics(t *testing.T) {
	be := errutil.NewError("oops")
	assert.Equal(t, "oops", be.Error())
	require.NoError(t, be.Unwrap())
	assert.Equal(t, errutil.SeverityInfo, be.Severity())
	assert.True(t, be.Metadata().IsEmpty())
}

func TestBaseError_WrapAndUnwrap(t *testing.T) {
	be := errutil.NewError("outer")
	inner := sentinel{}
	_ = be.Wrap(inner)
	assert.Equal(t, "outer: sentinel", be.Error())
	require.Equal(t, inner, be.Unwrap())
}

func TestWithSeverityOption(t *testing.T) {
	be := errutil.NewError("m", errutil.WithSeverity(errutil.SeverityWarning))
	assert.Equal(t, errutil.SeverityWarning, be.Severity())

	// invalid severity is ignored
	be2 := errutil.NewError("m2", errutil.WithSeverity(errutil.ErrorSeverity(999)))
	assert.Equal(t, errutil.SeverityInfo, be2.Severity())
}

func TestMetadataOptions(t *testing.T) {
	// WithMetadata
	md := errutil.NewErrorMetadata()
	md.Set("k", "v")
	be := errutil.NewError("m", errutil.WithMetadata(md))
	assert.Equal(t, "v", be.Metadata().Get("k"))

	// WithMetadataKeyValue
	be2 := errutil.NewError("m2", errutil.WithMetadataKeyValue("a", 123))
	assert.Equal(t, 123, be2.Metadata().Get("a"))
}
