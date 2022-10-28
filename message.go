package fault

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg }

func (w *withMessage) Unwrap() error { return w.cause }

// Msg provides a way to decorate a wrapped error with an additional message.
// It's similar to the `Wrap(error, string)` APIs common in the errors ecosystem
// and is implemented as just a simple string exposed as the `.Error()` value.
func Msg(s string) func(error) error {
	return func(err error) error {
		return &withMessage{err, s}
	}
}
