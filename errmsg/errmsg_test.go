package errmsg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOne(t *testing.T) {
	err := Wrap(errors.New("a problem"), "shit happened", "Shit happened.")
	out := GetProblem(err)

	assert.Equal(t, "Shit happened.", out)
}

func TestWithNone(t *testing.T) {
	err := errors.New("a problem")
	out := GetProblem(err)

	assert.Equal(t, "", out)
}

func TestWithMany(t *testing.T) {
	err := errors.New("the original problem")

	err = Wrap(err, "layer 1", "The post was not found.")
	err = Wrap(err, "layer 2", "Unable to reply to post.")
	err = Wrap(err, "layer 3", "Your reply draft has been saved however we could not publish it.")
	out := GetProblem(err)

	assert.Equal(t, "Your reply draft has been saved however we could not publish it. Unable to reply to post. The post was not found.", out)
}
