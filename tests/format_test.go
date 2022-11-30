package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatStdlibSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(1)

	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: stdlib sentinel error", fmt.Sprintf("%v", err))
	a.Regexp(`stdlib sentinel error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatFaultSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(2)

	a.Equal("failed to call function: fault sentinel error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: fault sentinel error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: fault sentinel error", fmt.Sprintf("%v", err))
	a.Regexp(`fault sentinel error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(3)

	a.Equal("failed to call function: stdlib root cause error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: stdlib root cause error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: stdlib root cause error", fmt.Sprintf("%v", err))
	a.Regexp(`stdlib root cause error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatFaultInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(4)

	a.Equal("failed to call function: fault root cause error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: fault root cause error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: fault root cause error", fmt.Sprintf("%v", err))
	a.Regexp(`fault root cause error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibErrorfWrappedError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(5)

	a.Equal("failed to call function: errorf wrapped: stdlib sentinel error: stdlib sentinel error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: errorf wrapped: stdlib sentinel error: stdlib sentinel error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: errorf wrapped: stdlib sentinel error: stdlib sentinel error", fmt.Sprintf("%v", err))
	a.Regexp(`stdlib sentinel error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibErrorfWrappedExternalError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(6)

	a.Equal("failed to call function: errorf wrapped external: external error wrapped with errorf: stdlib external error: external error wrapped with errorf: stdlib external error: stdlib external error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: errorf wrapped external: external error wrapped with errorf: stdlib external error: external error wrapped with errorf: stdlib external error: stdlib external error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: errorf wrapped external: external error wrapped with errorf: stdlib external error: external error wrapped with errorf: stdlib external error: stdlib external error", fmt.Sprintf("%v", err))
	a.Regexp(`errorf wrapped external: external error wrapped with errorf: stdlib external error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibErrorfWrappedExternallyWrappedError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(7)

	a.Equal("failed to call function: errorf wrapped external: external error wrapped with pkg/errors: github.com/pkg/errors external error: external error wrapped with pkg/errors: github.com/pkg/errors external error: github.com/pkg/errors external error", fmt.Sprintf("%s", err.Error()))
	a.Equal("failed to call function: errorf wrapped external: external error wrapped with pkg/errors: github.com/pkg/errors external error: external error wrapped with pkg/errors: github.com/pkg/errors external error: github.com/pkg/errors external error", fmt.Sprintf("%s", err))
	a.Equal("failed to call function: errorf wrapped external: external error wrapped with pkg/errors: github.com/pkg/errors external error: external error wrapped with pkg/errors: github.com/pkg/errors external error: github.com/pkg/errors external error", fmt.Sprintf("%v", err))
	a.Regexp(`errorf wrapped external: external error wrapped with pkg/errors: github.com/pkg/errors external error
\s+.+fault/tests/test_callers.go:29
failed to call function
\s+.+fault/tests/test_callers.go:20
`, fmt.Sprintf("%+v", err))
}

func TestFormatStdlibErrorfWrappedExternallyWrappedErrorBLank(t *testing.T) {
	a := assert.New(t)

	err := errorProducerFromRootCause(8)

	a.Equal("external error wrapped with pkg/errors: github.com/pkg/errors external error: github.com/pkg/errors external error", fmt.Sprintf("%s", err.Error()))
	a.Equal("external error wrapped with pkg/errors: github.com/pkg/errors external error: github.com/pkg/errors external error", fmt.Sprintf("%s", err))
	a.Equal("external error wrapped with pkg/errors: github.com/pkg/errors external error: github.com/pkg/errors external error", fmt.Sprintf("%v", err))

	a.Regexp(
		`external error wrapped with pkg/errors: github.com/pkg/errors external error
\s+.+fault/tests/test_callers.go:29
`, fmt.Sprintf("%+v", err))
}
