package fault

import (
	"errors"
	"fmt"
	"runtime"
)

func Trace(err error) []string {
	minimalStack := []string{}

	for err != nil {
		if f, ok := err.(interface {
			Location() string
		}); ok {
			if loc := f.Location(); loc != "" {
				minimalStack = append(minimalStack, loc)
			}
		}

		err = errors.Unwrap(err)
	}

	return minimalStack
}

func getLocation() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
	}

	return fmt.Sprintf("%s:%d", file, line)
}
