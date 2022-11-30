package fault

import (
	"errors"
)

// Chain represents an unwound error chain. Each step is a useful error. Errors
// without messages (such as those that only contain other errors) are omitted.
type Chain []Step

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
func Flatten(err error) Chain {
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

			f = append([]Step{{
				Location: lastLocation,
				Message:  message,
			}}, f...)

			lastLocation = ""
			lastMessage = message
		}
	}

	return f
}
