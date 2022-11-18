package tests

import (
	"errors"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
)

var (
	errSentinelStdlib = errors.New("stdlib sentinel error")
	errSentinelFault  = fault.New("fault sentinel error")
)

func rootCause(kind int) error {
	if kind == 0 {
		return nil
	} else if kind == 1 {
		return errSentinelStdlib
	} else if kind == 2 {
		return errSentinelFault
	} else if kind == 3 {
		return errors.New("stdlib root cause error")
	} else if kind == 4 {
		return fault.New("fault root cause error")
	} else if kind == 5 {
		return fmt.Errorf("errorf wrapped: %w", errSentinelStdlib)
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

func errorCallerMid(kind int) error {
	err := errorCallerDeep(kind)
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to call function"))
	}

	return nil
}

func errorCaller(kind int) error {
	err := errorCallerMid(kind)
	if err != nil {
		return fault.Wrap(err)
	}

	return nil
}
