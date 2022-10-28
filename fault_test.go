package fault_test

import (
	"errors"
	"testing"

	"github.com/Southclaws/fault"
	"github.com/stretchr/testify/assert"
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
		return fault.Wrap(err, fault.Msg("failed to call function"))
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

func TestFullCallStack(t *testing.T) {
	t.Run("sentinel_stdlib", func(t *testing.T) {
		a := assert.New(t)
		err := errorCaller(1)
		chain := fault.Flatten(err)

		a.ErrorContains(err, "failed to call function: stdlib sentinel error")
		a.ErrorContains(chain.Root, "stdlib sentinel error")
		a.Len(chain.Errors, 3)

		e0 := chain.Errors[0]
		a.Equal("stdlib sentinel error", e0.Message)
		a.Contains(e0.Location, "fault/fault_test.go:33")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:42")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:51")
	})

	t.Run("sentinel_fault", func(t *testing.T) {
		a := assert.New(t)
		err := errorCaller(2)
		chain := fault.Flatten(err)

		a.ErrorContains(err, "failed to call function: fault sentinel error")
		a.ErrorContains(chain.Root, "fault sentinel error")
		a.Len(chain.Errors, 3)

		e0 := chain.Errors[0]
		a.Equal("fault sentinel error", e0.Message)
		a.Contains(e0.Location, "fault/fault_test.go:33")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:42")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:51")
	})

	t.Run("inline_stdlib", func(t *testing.T) {
		a := assert.New(t)
		err := errorCaller(3)
		chain := fault.Flatten(err)

		a.ErrorContains(err, "failed to call function: stdlib root cause error")
		a.ErrorContains(chain.Root, "stdlib root cause error")
		a.Len(chain.Errors, 3)

		e0 := chain.Errors[0]
		a.Equal("stdlib root cause error", e0.Message)
		a.Contains(e0.Location, "fault/fault_test.go:33")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:42")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:51")
	})

	t.Run("inline_fault", func(t *testing.T) {
		a := assert.New(t)
		err := errorCaller(4)
		chain := fault.Flatten(err)

		a.ErrorContains(err, "failed to call function: fault root cause error")
		a.ErrorContains(chain.Root, "fault root cause error")
		a.Len(chain.Errors, 3)

		e0 := chain.Errors[0]
		a.Equal("fault root cause error", e0.Message)
		a.Contains(e0.Location, "fault/fault_test.go:33")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:42")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:51")
	})
}
