package fault

import "fmt"

// New creates a new basic fault error.
func New(message string) error {
	stack := callers(3)
	return &fault{
		root:   true,
		msg:    message,
		stack:  stack,
		global: stack.isGlobal(),
	}
}

type wrapper func(err error) error

// Wrap wraps an error with all of the wrappers provided.
func Wrap(err error, w ...wrapper) error {
	if err == nil {
		return nil
	}

	// callers(4) skips runtime.Callers, stack.callers, this method, and Wrap(f)
	stack := callers(4)
	// caller(3) skips stack.caller, this method, and Wrap(f)
	// caller(skip) has a slightly different meaning which is why it's not 4 as above
	frame := caller(3)
	switch e := err.(type) {
	case *fault:
		if e.root {
			if e.global {
				// create a new root error for global values to make sure nothing interferes with the stack
				err = &fault{
					root:   true,
					global: e.global,
					stack:  stack,
				}
			} else {
				// insert the frame into the stack
				e.stack.insertPC(*stack)
			}
		} else {
			if root, ok := Cause(err).(*fault); ok {
				root.stack.insertPC(*stack)
			}
		}
	default:
		// return a new root error that wraps the external error
		return &fault{
			root:  true,
			msg:   e.Error(),
			cause: e,
			stack: stack,
		}
	}

	// run all the decorators after the stack info is figured out.
	for _, fn := range w {
		err = fn(err)
	}

	return &fault{
		root:  false,
		cause: err,
		stack: callers(3),
		frame: frame,
	}
}

// Cause returns the root cause of the error, which is defined as the first error in the chain. The original
// error is returned if it does not implement `Unwrap() error` and nil is returned if the error is nil.
func Cause(err error) error {
	for {
		uerr := Unwrap(err)
		if uerr == nil {
			return err
		}
		err = uerr
	}
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method
// returning error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// StackFrames returns the trace of a root error in the form of a program counter slice.
// This method is currently called by an external error tracing library (Sentry).
func (e *fault) StackFrames() []uintptr {
	return *e.stack
}

// StackFrames returns the trace of an error in the form of a program counter slice.
// Use this method if you want to pass the eris stack trace to some other error tracing library.
func StackFrames(err error) []uintptr {
	for err != nil {
		switch err := err.(type) {
		case *fault:
			return err.StackFrames()
		default:
			return []uintptr{}
		}
	}
	return []uintptr{}
}

type fault struct {
	root   bool   // is this error the first in the chain?
	global bool   // is this a globally declared sentinel error?
	msg    string // root error message
	cause  error  // if this wraps another error
	stack  *stack // the full stack trace of a root error
	frame  *frame // the stack frame of a wrapped error
}

func (f *fault) Error() string {
	if f.root {
		return f.msg
	} else {
		return f.cause.Error()
	}
}

type Fault struct {
	Message  string
	Root     error
	External error
	Stack    Stack
	Chain    []Breadcrumb
}

type Breadcrumb struct {
	Message string
	Frame   StackFrame
}

func Get(err error) *Fault {
	if err == nil {
		return nil
	}

	var f Fault
	for err != nil {
		switch err := err.(type) {
		case *fault:
			if err.root {
				f.Message = err.msg
				f.Stack = err.stack.get()
			} else {
				link := Breadcrumb{
					Message: err.msg,
					Frame:   err.frame.get(),
				}
				f.Chain = append([]Breadcrumb{link}, f.Chain...)
			}

		default:
			f.External = err
			return &f
		}

		err = Unwrap(err)
	}

	return &f
}

func (f *fault) Format(s fmt.State, verb rune) {
	u := Get(f)

	s.Write([]byte(u.Message + "\n"))

	for _, v := range u.Chain {
		if v.Message != "" {
			s.Write([]byte(fmt.Sprintf("%s\n", v.Message)))
		}
		s.Write([]byte(fmt.Sprintf("\t%s\n", v.Frame.String())))
	}
}

func (f *fault) Unwrap() error { return f.cause }

func (f *fault) Stack() Stack {
	return f.stack.get()
}
