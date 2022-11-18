package tests

import (
	"errors"
	"testing"

	"github.com/Southclaws/fault/ftag"
	"github.com/stretchr/testify/assert"
)

func TestWrapWithKind(t *testing.T) {
	err := ftag.Wrap(errors.New("a problem"), ftag.NotFound)
	out := ftag.Get(err)

	assert.Equal(t, ftag.NotFound, out)
}

func TestWrapWithKindChanging(t *testing.T) {
	err := ftag.Wrap(errors.New("a problem"), ftag.Internal)
	err = ftag.Wrap(err, ftag.Internal)
	err = ftag.Wrap(err, ftag.Internal)
	err = ftag.Wrap(err, ftag.InvalidArgument)
	err = ftag.Wrap(err, ftag.InvalidArgument)
	err = ftag.Wrap(err, ftag.NotFound)
	out := ftag.Get(err)

	assert.Equal(t, ftag.NotFound, out, "Should always pick the most recent kind from an error chain.")
}

func TestWrapNil(t *testing.T) {
	assert.Panics(t, func() {
		ftag.Wrap(nil, ftag.NotFound)
	})
}
