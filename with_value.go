package fault

func WithValue(parent error, message, key string, val any) error {
	if parent == nil {
		panic("cannot create error from nil parent")
	}

	return &fault{
		msg:        message,
		underlying: parent,
		key:        key,
		value:      val,
		location:   getLocation(),
	}
}
