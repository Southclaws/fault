package fault_test

import (
	"errors"
	"testing"

	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/fault"
)

func TestWithMessageNativeErrors(t *testing.T) {
	fn := func() error {
		return errors.New("root cause")
	}

	rootcause := fn()

	level1 := fault.WithMessage(rootcause, "level 1")
	level2 := fault.WithMessage(level1, "level 2")
	level3 := fault.WithMessage(level2, "level 3")

	gotString := level3.Error()
	gotContext := fault.Context(level3)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotString)
	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Message)

	assert.Len(t, gotContext.Trace, 3)
}

func TestWithMessagePkgErrors(t *testing.T) {
	fn := func() error {
		return pkg_errors.New("root cause")
	}

	rootcause := fn()

	level1 := fault.WithMessage(rootcause, "level 1")
	level2 := fault.WithMessage(level1, "level 2")
	level3 := fault.WithMessage(level2, "level 3")

	gotString := level3.Error()
	gotContext := fault.Context(level3)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotString)
	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Message)
	assert.Len(t, gotContext.Trace, 4)
}

func TestWithMessageNil(t *testing.T) {
	assert.Nil(t, fault.WithMessage(nil, ""))
}
