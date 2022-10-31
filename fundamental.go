package fault

import "fmt"

// New creates a new basic fault error.
func New(message string) error {
	return &fundamental{
		msg:      message,
		location: getLocation(),
	}
}

// Newf includes formatting specifiers.
func Newf(message string, va ...any) error {
	return &fundamental{
		msg:      fmt.Sprintf(message, va...),
		location: getLocation(),
	}
}

type fundamental struct {
	msg      string
	location string
}

func (f *fundamental) Error() string {
	return f.msg
}
