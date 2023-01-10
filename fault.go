// Package fault provides an extensible yet ergonomic mechanism for wrapping
// errors. It implements this as a kind of middleware style pattern by providing
// a simple option-style interface that can be passed to a call to `fault.Wrap`.
//
// See the GitHub repository for full documentation and examples.
package fault

import (
	"fmt"
	"runtime"
	"strings"
)

// Wrapper describes a kind of middleware that packages can satisfy in order to
// decorate errors with additional, domain-specific context.
type Wrapper func(err error) error

// Wrap wraps an error with all of the wrappers provided.
func Wrap(err error, w ...Wrapper) error {
	if err == nil {
		return nil
	}

	for _, fn := range w {
		err = fn(err)
	}

	c := &container{
		cause:    err,
		location: getLocation(),
	}

	return c
}

type container struct {
	cause    error
	location string
}

// Error behaves like most error wrapping libraries, it gives you all the error
// messages conjoined with ": ". This is useful only for internal error reports,
// never show this to an end-user or include it in responses as it may reveal
// internal technical information about your application stack.
func (f *container) Error() string {
	errs := []string{}
	chain := Flatten(f)

	// reverse iterate since the chain is in caller order
	for i := len(chain) - 1; i >= 0; i-- {
		message := chain[i].Message
		if message != "" && !isInternalString(message) {
			errs = append(errs, chain[i].Message)
		}
	}

	message := strings.Join(errs, ": ")
	if message == "" {
		message = "(no error message provided)"
	}

	return message
}

func (f *container) Unwrap() error { return f.cause }

func (f *container) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			u := Flatten(f)
			for _, v := range u {
				if v.Message != "" {
					fmt.Fprintf(s, "%s\n", v.Message)
				}
				if v.Location != "" {
					fmt.Fprintf(s, "\t%s\n", v.Location)
				}
			}
			return
		}

		fallthrough

	case 's':
		fmt.Fprint(s, f.Error())
	}
}

func getLocation() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}

// isInternalString returns true for messages like <fctx> which are placeholders
func isInternalString(s string) bool {
	return strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">")
}
