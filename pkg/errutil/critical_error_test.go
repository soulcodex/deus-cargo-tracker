package errutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

func TestNewCriticalError(t *testing.T) {
	ce := errutil.NewCriticalError("fatal")
	assert.Equal(t, "fatal", ce.Error())
	assert.Equal(t, errutil.SeverityFatal, ce.Severity())
	assert.Equal(t, "critical", ce.Metadata().Get("error.type"))
	assert.True(t, errutil.IsCriticalError(ce))

	var err interface{} = ce
	assert.True(t, errutil.IsCriticalError(err.(error)))
}

func TestNewCriticalErrorWithMetadata(t *testing.T) {
	md := errutil.NewErrorMetadata()
	md.Set("foo", "bar")
	ce := errutil.NewCriticalErrorWithMetadata("msg", md)
	assert.Equal(t, "bar", ce.Metadata().Get("foo"))
}
