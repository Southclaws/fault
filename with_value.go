package fault

import "fmt"

// WithValue wraps an error along with some additional metadata expressed as a
// key-value pair. This allows errors to contain structured information for
// later retrieval for use in structured logging systems, HTTP and gRPC response
// mechanisms and general diagnostic tooling.
func WithValue(parent error, message, key, val string) error {
	if parent == nil {
		return nil
	}

	return &fault{
		msg:        fmt.Sprintf("%s: %s", message, parent.Error()),
		underlying: parent,
		key:        key,
		value:      val,
		location:   getLocation(),
	}
}
