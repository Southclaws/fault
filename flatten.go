package fault

import (
	"errors"
)

type Chain struct {
	Root   error
	Errors []Step
}

type Step struct {
	Location string
	Message  string
}

// Flatten attempts to derive more useful structured information from an error
// chain. If the input is a fault error, the output will contain an easy to use
// error chain list with location information and individual error messages.
func Flatten(err error) *Chain {
	if err == nil {
		return nil
	}

	var f Chain
	for err != nil {
		var next error

		switch unwrapped := err.(type) {
		case *container:
			step := Step{}
			step.Location = unwrapped.location

			// NOTE: Because fault containers do not have messages, they only
			// exist to contain other errors that actually contain information,
			// peek the next error in the chain in order to get its message and
			// add it to the current step before appending it to the chain.
			next = errors.Unwrap(err)
			if _, ok := next.(*container); !ok {
				step.Message = next.Error()
			}

			f.Errors = append([]Step{step}, f.Errors...)
		}

		next = errors.Unwrap(err)

		// If the next error in the chain is nil, that means `err` is the last
		// error in the chain. This error is the root or "external" error.
		if next == nil {
			f.Root = err
		}

		err = next
	}

	return &f
}
