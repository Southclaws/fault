package fault

import (
	"errors"
)

// Chain represents an unwound error chain. It contains a reference to the root
// error (the error at the start of the chain, which does not provide any way
// to `errors.Unwrap` it further) as well as a slice of "Step" objects which are
// points where the error was wrapped.
type Chain struct {
	Root   error // TODO: deprecate
	Errors []Step
}

// Step represents a location where an error was wrapped by Fault. The location
// is present if the error being wrapped contained stack information and the
// message is present if the underlying error provided a message. Note that not
// all errors provide errors or locations. If both are missing, it's omitted.
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

	// first, flatten the call tree into an array so it's easier to work with.
	flat := []error{}
	for err != nil {
		flat = append(flat, err)
		err = errors.Unwrap(err)
	}

	lastLocation := ""
	lastMessage := ""

	var f Chain
	for i := 0; i < len(flat); i++ {
		err := flat[i]

		// lookahead if possible
		var next error
		if i+1 < len(flat) {
			next = flat[i+1]
		}

		switch unwrapped := err.(type) {
		// NOTE: Because fault containers do not have messages, they only
		// exist to contain other errors that actually contain information,
		// store the container's recorded location for usage with the next item.
		case *container:
			lastLocation = unwrapped.location

		default:
			message := err.Error()

			// de-duplicate nested error messages
			// TODO: improve this by destructuring/splitting nested duplicates.
			if lastMessage == message {
				continue
			}

			f.Errors = append([]Step{{
				Location: lastLocation,
				Message:  message,
			}}, f.Errors...)

			lastLocation = ""
			lastMessage = message
		}

		// If the next error in the chain is nil, that means `err` is the last
		// error in the chain. This error is the root or "external" error.
		if next == nil {
			f.Root = err
		}
	}

	return &f
}
