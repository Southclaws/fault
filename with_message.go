package fault

import "fmt"

// WithMessage is just like all the other error wrapping libraries. Stores a
// string message alongside an error in the chain.
func WithMessage(parent error, message string) error {
	if parent == nil {
		return nil
	}

	return &fault{
		underlying: parent,
		msg:        fmt.Sprintf("%s: %s", message, parent.Error()),
		location:   getLocation(),
	}
}
