package errmsg

import (
	"errors"
	"strings"
)

// Problem describes an error message that is intended for end-user consumption.
// It's named "problem" because this word feels like a non-technical alternative
// to the word "error" which is 1. overused and 2. typically used mainly in tech
// to describe internal tech issues. Using a different word helps differentiate.
type Problem = string

type withMessage struct {
	underlying error
	internal   string
	external   string
}

// Wrap wraps an error with an internal and an external message. The internal
// message is just the same as any error wrapping library but the external error
// message is intended for display to the end-user. It's recommended to use full
// punctuation and grammar and end the message with a period.
func Wrap(err error, internal, external string) error {
	if err == nil {
		panic("nil error passed to Wrap")
	}

	return &withMessage{
		err, internal, external,
	}
}

func (e *withMessage) Error() string   { return e.internal }
func (e *withMessage) Unwrap() error   { return e.underlying }
func (e *withMessage) Message() string { return e.external }

// GetProblem returns a space-joined string of all end-user intended messages in
// the error chain. This message can then be displayed/sent to end users.
func GetProblem(err error) Problem {
	return Problem(strings.Join(GetProblems(err), " "))
}

// GetProblems returns all end-user intended messages in the input error chain.
func GetProblems(err error) []Problem {
	p := []Problem{}

	for err != nil {
		var wm *withMessage
		if errors.As(err, &wm) {
			p = append(p, wm.external)
		}

		err = errors.Unwrap(err)
	}

	return p
}
