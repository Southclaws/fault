package fault

type wrapper interface {
	Wrap(err error) error
}

// Wrap wraps an error with all of the wrappers provided.
// TODO: stack traces and all that good stuff.
func Wrap(err error, w ...wrapper) error {
	for _, w := range w {
		err = w.Wrap(err)
	}
	return err
}
