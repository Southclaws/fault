package fault

import "fmt"

func Wrap(err error, text string) error {
	return &fault{
		underlying: err,
		msg:        fmt.Sprintf("%s: %s", text, err.Error()),
		location:   getLocation(),
	}
}
