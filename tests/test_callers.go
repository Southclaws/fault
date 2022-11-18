package tests

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
)

func errorCaller(kind int) error {
	err := errorCallerMid(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func errorCallerMid(kind int) error {
	err := errorCallerDeep(kind)
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to call function"))
	}

	return nil
}

func errorCallerDeep(kind int) error {
	err := rootCause(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}
