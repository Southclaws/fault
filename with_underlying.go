package fault

func Wrap(err error, text string) error {
	return &fault{
		underlying: err,
		msg:        text,
		location:   getLocation(),
	}
}
