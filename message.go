package fault

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg }

func (w *withMessage) Unwrap() error { return w.cause }

func Msg(s string) func(error) error {
	return func(err error) error {
		return &withMessage{err, s}
	}
}
