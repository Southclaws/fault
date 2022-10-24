package fault_test

import (
	"errors"
	"testing"

	"github.com/Southclaws/fault"
	"github.com/kr/pretty"
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
	} else {
		return fault.New("fault root cause error")
	}
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
		return fault.Wrap(err)
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

func Test_fault_Format(t *testing.T) {
	err := errorCaller(1)

	f := err.(interface{ Stack() fault.Stack })

	pretty.Printf("%+v\n", f.Stack())
}
