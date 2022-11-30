package tests

import (
	"testing"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/stretchr/testify/assert"
)

func TestFlattenStdlibSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(1)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: stdlib sentinel error", full)
	a.Equal("stdlib sentinel error", root)
	a.Len(chain.Errors, 2)

	e0 := chain.Errors[0]
	a.Equal("stdlib sentinel error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")
}

func TestFlattenFaultSentinelError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(2)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: fault sentinel error", full)
	a.Equal("fault sentinel error", root)
	a.Len(chain.Errors, 2)

	e0 := chain.Errors[0]
	a.Equal("fault sentinel error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")
}

func TestFlattenStdlibInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(3)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: stdlib root cause error", full)
	a.Equal("stdlib root cause error", root)
	a.Len(chain.Errors, 2)

	e0 := chain.Errors[0]
	a.Equal("stdlib root cause error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")
}

func TestFlattenFaultInlineError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(4)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to call function: fault root cause error", full)
	a.Equal("fault root cause error", root)
	a.Len(chain.Errors, 2)

	e0 := chain.Errors[0]
	a.Equal("fault root cause error", e0.Message)
	a.Contains(e0.Location, "test_callers.go:29")

	e1 := chain.Errors[1]
	a.Equal("failed to call function", e1.Message)
	a.Contains(e1.Location, "test_callers.go:20")
}

func TestFlattenStdlibErrorfWrappedError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(5)
	chain := fault.Flatten(err)
	full := err.Error()

	a.Equal("failed to call function: errorf wrapped: stdlib sentinel error: stdlib sentinel error", full)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("stdlib sentinel error", e0.Message)
	a.Empty(e0.Location)

	e1 := chain.Errors[1]
	a.Equal("errorf wrapped: stdlib sentinel error", e1.Message)
	a.Contains(e1.Location, "test_callers.go:29")

	e2 := chain.Errors[2]
	a.Equal("failed to call function", e2.Message)
	a.Contains(e2.Location, "test_callers.go:20")
}

func TestFlattenStdlibErrorfWrappedExternalError(t *testing.T) {
	a := assert.New(t)

	err := errorCaller(6)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	// NOTE: the way that other libraries handle wrapped errors isn't great, the
	// message is often just a join of nested strings so it's not easy to split.
	a.Equal("failed to call function: errorf wrapped external: external error wrapped with errorf: stdlib external error: external error wrapped with errorf: stdlib external error: stdlib external error", full)
	a.Equal("stdlib external error", root)
	a.Len(chain.Errors, 4)

	e0 := chain.Errors[0]
	a.Equal("stdlib external error", e0.Message)
	a.Empty(e0.Location)

	e1 := chain.Errors[1]
	a.Equal("external error wrapped with errorf: stdlib external error", e1.Message)
	a.Empty(e1.Location)

	e2 := chain.Errors[2]
	a.Equal("errorf wrapped external: external error wrapped with errorf: stdlib external error", e2.Message)
	a.Contains(e2.Location, "test_callers.go:29")

	e3 := chain.Errors[3]
	a.Equal("failed to call function", e3.Message)
	a.Contains(e3.Location, "test_callers.go:20")
}

func TestFlattenStdlibErrorfWrappedExternallyWrappedError(t *testing.T) {
	a := assert.New(t)

	err := rootCause(8)
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("external error wrapped with pkg/errors: github.com/pkg/errors external error", full)
	a.Equal("github.com/pkg/errors external error", root)
	a.Len(chain.Errors, 2)

	e0 := chain.Errors[0]
	a.Equal("github.com/pkg/errors external error", e0.Message)
	a.Empty(e0.Location)

	e1 := chain.Errors[1]
	a.Equal("external error wrapped with pkg/errors: github.com/pkg/errors external error", e1.Message)
	a.Empty(e1.Location)
}

func TestFlattenStdlibErrorfWrappedExternallyWrappedErrorBrokenChain(t *testing.T) {
	a := assert.New(t)

	original := externalWrappedPostgresError()
	err := fault.Wrap(original, fmsg.With("failed to query"))
	chain := fault.Flatten(err)
	full := err.Error()
	root := chain.Root.Error()

	a.Equal("failed to query: external pg error: fatal: your sql was wrong bro (SQLSTATE 123): fatal: your sql was wrong bro (SQLSTATE 123)", full)
	a.Equal("fatal: your sql was wrong bro (SQLSTATE 123)", root)
	a.Len(chain.Errors, 3)

	e0 := chain.Errors[0]
	a.Equal("fatal: your sql was wrong bro (SQLSTATE 123)", e0.Message)

	e1 := chain.Errors[1]
	a.Equal("external pg error: fatal: your sql was wrong bro (SQLSTATE 123)", e1.Message)

	e2 := chain.Errors[2]
	a.Equal("failed to query", e2.Message)
}
