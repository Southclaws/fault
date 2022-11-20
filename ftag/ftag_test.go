package ftag

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapWithKind(t *testing.T) {
	err := Wrap(errors.New("a problem"), NotFound)
	out := Get(err)

	assert.Equal(t, NotFound, out)
}

func TestWrapWithKindChanging(t *testing.T) {
	err := Wrap(errors.New("a problem"), Internal)
	err = Wrap(err, Internal)
	err = Wrap(err, Internal)
	err = Wrap(err, InvalidArgument)
	err = Wrap(err, InvalidArgument)
	err = Wrap(err, NotFound)
	out := Get(err)

	assert.Equal(t, NotFound, out, "Should always pick the most recent kind from an error chain.")
}

func TestWrapNil(t *testing.T) {
	assert.Panics(t, func() {
		Wrap(nil, NotFound)
	})
}
