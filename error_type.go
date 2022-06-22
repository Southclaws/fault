package fault

import (
	"encoding/json"
)

// Error implements the Go error type and supports metadata that can easily be
// logged or sent as a response to clients.
type Error struct {
	underlying error
	private    Metadata
	public     Metadata
}

// Metadata provides a way to annotate errors with additional information that
// can easily be logged in a structured way.
type Metadata map[string]interface{}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.underlying.Error()
}

func (e *Error) Cause() error  { return e.underlying }
func (e *Error) Unwrap() error { return e.underlying }

// MarshalJSON implements the json.Marshaler interface.
func (e *Error) MarshalJSON() ([]byte, error) {
	serialised := e.public
	serialised["error"] = e.underlying.Error()
	return json.Marshal(serialised)
}

// unwrapper is from the standard library.
type unwrapper interface {
	Unwrap() error
}

// causer is for the pkg/friendsofgo errors package
type causer interface {
	Cause() error
}

// interface assertions
var (
	_ error          = &Error{}
	_ json.Marshaler = &Error{}
	_ unwrapper      = &Error{}
	_ causer         = &Error{}
)
