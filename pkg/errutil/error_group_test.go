package errutil_test

import (
	"errors"
	"testing"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
	"github.com/stretchr/testify/assert"
)

func TestErrorGroup_AddAndErrors(t *testing.T) {
	eg := errutil.NewErrorGroup()
	assert.Empty(t, eg.Errors())

	eg.Add(nil)
	assert.Empty(t, eg.Errors())

	eg.Add(errors.New("e1"))
	eg.Add(errors.New("e2"))
	errs := eg.Errors()
	assert.Len(t, errs, 2)
	assert.Equal(t, "e1", errs[0].Error())
	assert.Equal(t, "e2", errs[1].Error())
}

func TestErrorGroup_ErrorString(t *testing.T) {
	eg := errutil.NewErrorGroup(errutil.WithSeparator(" | "))
	eg.Add(errors.New("one"))
	eg.Add(errors.New("two"))
	assert.Equal(t, "one | two", eg.Error())
}
