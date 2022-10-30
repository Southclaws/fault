package fdesc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOne(t *testing.T) {
	err := Wrap(errors.New("a problem"), "shit happened", "Shit happened.")
	out := GetIssue(err)

	assert.Equal(t, "Shit happened.", out)
}

func TestWithNone(t *testing.T) {
	err := errors.New("a problem")
	out := GetIssue(err)

	assert.Equal(t, "", out)
}

func TestWithMany(t *testing.T) {
	err := errors.New("the original problem")

	err = Wrap(err, "layer 1", "The post was not found.")
	err = Wrap(err, "layer 2", "Unable to reply to post.")
	err = Wrap(err, "layer 3", "Your reply draft has been saved however we could not publish it.")
	out := GetIssue(err)

	assert.Equal(t, "Your reply draft has been saved however we could not publish it. Unable to reply to post. The post was not found.", out)
}

func TestWithManySlice(t *testing.T) {
	err := errors.New("the original problem")

	err = Wrap(err, "layer 1", "The post was not found.")
	err = Wrap(err, "layer 2", "Unable to reply to post.")
	err = Wrap(err, "layer 3", "Your reply draft has been saved however we could not publish it.")
	out := GetIssues(err)

	assert.Len(t, out, 3)
	assert.Equal(t, []string{"Your reply draft has been saved however we could not publish it.", "Unable to reply to post.", "The post was not found."}, out)
}

func TestNil(t *testing.T) {
	assert.Panics(t, func() {
		Wrap(nil, "oh no", ":(")
	})
}
