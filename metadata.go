package fault

import "errors"

func Metadata(err error) map[string]any {
	m := make(map[string]any)

	for err != nil {
		if f, ok := err.(interface {
			Key() string
			Value() any
		}); ok && f.Key() != "" {
			m[f.Key()] = f.Value()
		}

		err = errors.Unwrap(err)
	}

	return m
}
