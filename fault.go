package fault

import (
	"errors"
)

func Sentinel(text string) error {
	return &fault{
		underlying: errors.New(text),
		location:   getLocation(),
	}
}
