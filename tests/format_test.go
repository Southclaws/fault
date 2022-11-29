package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatStdlibSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(1)
	full := err.Error()
	format := fmt.Sprintf("%+v", err)

	a.Equal("failed to call function: stdlib sentinel error", full)
	a.Equal(`stdlib sentinel error
	d:/Work/makeroom/fault/tests/test_callers.go:29
failed to call function
	d:/Work/makeroom/fault/tests/test_callers.go:20
	d:/Work/makeroom/fault/tests/test_callers.go:11
`, format)
}

func TestFormatFaultSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(2)
	full := err.Error()
	format := fmt.Sprintf("%+v", err)

	a.Equal("failed to call function: fault sentinel error", full)
	a.Equal(`fault sentinel error
	d:/Work/makeroom/fault/tests/test_callers.go:29
failed to call function
	d:/Work/makeroom/fault/tests/test_callers.go:20
	d:/Work/makeroom/fault/tests/test_callers.go:11
`, format)
}

func TestFormatStdlibInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(3)
	full := err.Error()
	format := fmt.Sprintf("%+v", err)

	a.Equal("failed to call function: stdlib root cause error", full)
	a.Equal(`stdlib root cause error
	d:/Work/makeroom/fault/tests/test_callers.go:29
failed to call function
	d:/Work/makeroom/fault/tests/test_callers.go:20
	d:/Work/makeroom/fault/tests/test_callers.go:11
`, format)
}

func TestFormatFaultInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(4)
	full := err.Error()
	format := fmt.Sprintf("%+v", err)

	a.Equal("failed to call function: fault root cause error", full)
	a.Equal(`fault root cause error
	d:/Work/makeroom/fault/tests/test_callers.go:29
failed to call function
	d:/Work/makeroom/fault/tests/test_callers.go:20
	d:/Work/makeroom/fault/tests/test_callers.go:11
`, format)
}

func TestFormatStdlibErrorfWrappedError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(5)
	full := err.Error()
	format := fmt.Sprintf("%+v", err)

	a.Equal("failed to call function: errorf wrapped: stdlib sentinel error", full)
	a.Equal(`errorf wrapped: stdlib sentinel error
	d:/Work/makeroom/fault/tests/test_callers.go:29
failed to call function
	d:/Work/makeroom/fault/tests/test_callers.go:20
	d:/Work/makeroom/fault/tests/test_callers.go:11
`, format)
}

func TestFormatStdlibErrorfWrappedExternalError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(6)
	full := err.Error()
	format := fmt.Sprintf("%+v", err)

	a.Equal("failed to call function: errorf wrapped external: external error wrapped with errorf: stdlib external error", full)
	a.ErrorContains(err, "external error wrapped with errorf: stdlib external error")
	a.Equal(`errorf wrapped external: external error wrapped with errorf: stdlib external error
	d:/Work/makeroom/fault/tests/test_callers.go:29
failed to call function
	d:/Work/makeroom/fault/tests/test_callers.go:20
	d:/Work/makeroom/fault/tests/test_callers.go:11
`, format)
}
