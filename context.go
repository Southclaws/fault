package fault

import (
	"errors"
)

// POI is a point of interest in the code. Each one of these represents a
// location where fault has been used to wrap an error with additional context.
type POI struct {
	Message  string `json:"message,omitempty"`
	Location string `json:"location,omitempty"`
	Key      string `json:"key,omitempty"`
	Value    any    `json:"value,omitempty"`
}

// Context provides something similar to a stack trace, but much more minimal,
// focused and simple.
//
// A stack trace contains every single line of code in the call stack (including
// Go runtime internals) as well as a whole load of information about the super
// low level details of the state of the process at the time of the error. This
// information may be useful in some contexts but in a lot of cases in real
// world applications where business logic and data flow are among the primary
// focuses, this highly granular style of assembly-level stack trace can create
// noise and make finding the source of what often turns out to be application
// level rather than binary level errors take longer than is necessary.
//
// Instead, Context provides a set of signposts, decorated with the available
// contextual information in a simple list. The Fault library provides a way to
// annotate errors as you ascend the call stack and this unrolls all that info
// in the error chain out to a simple list.
//
func Context(err error) []POI {
	m := []POI{}

	for err != nil {
		step := POI{
			Message: err.Error(),
		}

		if f, ok := err.(withLocation); ok {
			step.Location = f.Location()
		}

		if f, ok := err.(withValue); ok {
			step.Key = f.Key()
			step.Value = f.Value()
		}

		m = append(m, step)

		err = errors.Unwrap(err)
	}

	return m
}

type withLocation interface {
	Location() string
}

type withValue interface {
	Key() string
	Value() any
}
