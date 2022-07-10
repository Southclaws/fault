package fault

// fault implements the Go error type and supports metadata that can easily be
// logged or sent as a response to clients.
type fault struct {
	// the wrapped error value, either a standard library primitive or any other
	// error type from the ecosystem of error libraries.
	underlying error

	// a simple message annotating this error in the error chain.
	msg string

	// location context of this particular error context so we don't need to
	// store a full stack trace of mostly useless info.
	location string

	// a key-value pair much like context.valueCtx for storing any metadata.
	key   string
	value any
}

// New creates an error with a simple minimal location string.
func New(text string) error {
	return &fault{
		msg:      text,
		location: getLocation(),
	}
}

// Implements all the interfaces for compatibility with the errors ecosystem.

func (e *fault) Error() string {
	// if e.underlying != nil {
	// 	return e.msg + ": " + e.underlying.Error()
	// }
	return e.msg
}

func (e *fault) Location() string { return e.location }
func (e *fault) Value() any       { return e.value }
func (e *fault) Key() string      { return e.key }
func (e *fault) Cause() error     { return e.underlying }
func (e *fault) Unwrap() error    { return e.underlying }
func (e *fault) String() string   { return e.Error() }
