package fault

import (
	"errors"
	"fmt"
	"runtime"
)

// New creates a new basic fault error.
func New(message string) error {
	return &fault{
		root:     true,
		msg:      message,
		location: getLocation(),
	}
}

type wrapper func(err error) error

// Wrap wraps an error with all of the wrappers provided.
func Wrap(err error, w ...wrapper) error {
	if err == nil {
		return nil
	}

	for _, fn := range w {
		err = fn(err)
	}

	return &fault{
		root:     false,
		cause:    err,
		location: getLocation(),
	}
}

type fault struct {
	root     bool   // is this error the first in the chain?
	global   bool   // is this a globally declared sentinel error?
	msg      string // root error message
	cause    error  // if this wraps another error
	location Location
}

func (f *fault) Error() string {
	if f.root {
		return f.msg
	} else {
		return f.cause.Error()
	}
}

type Fault struct {
	Message string
	Root    error
	Trace   []Location
}

type Location string

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
			} else {
				f.Trace = append([]Location{err.location}, f.Trace...)
			}

		default:
			f.Root = err
			f.Message = err.Error()
			return &f
		}

		err = errors.Unwrap(err)
	}

	return &f
}

func (f *fault) Format(s fmt.State, verb rune) {
	u := Get(f)

	s.Write([]byte(u.Message + "\n"))

	for _, v := range u.Trace {
		// if v.Message != "" {
		// 	s.Write([]byte(fmt.Sprintf("%s\n", v.Message)))
		// }
		s.Write([]byte(fmt.Sprintf("\t%s\n", v)))
	}
}

func (f *fault) Unwrap() error { return f.cause }

func getLocation() Location {
	_, file, line, _ := runtime.Caller(2)
	return Location(fmt.Sprintf("%s:%d", file, line))
}
