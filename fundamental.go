package fault

// New creates a new basic fault error.
func New(message string) error {
	return &fundamental{
		msg:      message,
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
