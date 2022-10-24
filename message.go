package fault

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }

func Msg(s string) func(error) error {
	return func(err error) error {
		return &withMessage{err, s}
	}
}
