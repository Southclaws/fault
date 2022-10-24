package fault

// New creates a new basic fault error.
func New(message string) error {
	return &fault{
		root:  true,
		msg:   message,
		stack: callers(3),
	}
}

type wrapper func(err error) error

// Wrap wraps an error with all of the wrappers provided.
func Wrap(err error, w ...wrapper) error {
	for _, fn := range w {
		err = fn(err)
	}
	return &fault{
		root:  false,
		cause: err,
		stack: callers(3),
	}
}

type fault struct {
	root  bool   // is this error the first in the chain?
	msg   string // root error message
	cause error  // nil if root == true, otherwise, the wrapped error
	stack *stack // stack pc
}

func (f *fault) Error() string {
	if f.root {
		return f.msg
	} else {
		return f.cause.Error()
	}
}
func (f *fault) Unwrap() error { return f.cause }

func (f *fault) Stack() Stack {
	return f.stack.get()
}
