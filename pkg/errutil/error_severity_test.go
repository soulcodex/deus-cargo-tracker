package errutil_test

import (
	"testing"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
	"github.com/stretchr/testify/assert"
)

func TestSeverity_StringAndValid(t *testing.T) {
	assert.Equal(t, "info", errutil.SeverityInfo.String())
	assert.Equal(t, "warning", errutil.SeverityWarning.String())
	assert.Equal(t, "fatal", errutil.SeverityFatal.String())
	assert.Equal(t, "unknown", errutil.ErrorSeverity(42).String())
	assert.True(t, errutil.SeverityFatal.IsValid())
	assert.False(t, errutil.ErrorSeverity(-1).IsValid())
}
