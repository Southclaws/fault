package tests

import (
	"testing"

	"github.com/Southclaws/fault"
	"github.com/stretchr/testify/assert"
)

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
		a.Contains(e0.Location, "fault/fault_test.go:3")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:4")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:5")
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
		a.Contains(e0.Location, "fault/fault_test.go:3")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:4")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:5")
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
		a.Contains(e0.Location, "fault/fault_test.go:3")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:4")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:5")
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
		a.Contains(e0.Location, "fault/fault_test.go:3")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:4")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:5")
	})

	t.Run("errorf_wrapped", func(t *testing.T) {
		a := assert.New(t)
		err := errorCaller(5)
		chain := fault.Flatten(err)

		a.ErrorContains(err, "failed to call function: errorf wrapped: stdlib sentinel error")
		a.ErrorContains(chain.Root, "stdlib sentinel error")
		a.Len(chain.Errors, 3)

		e0 := chain.Errors[0]
		a.Equal("errorf wrapped: stdlib sentinel error", e0.Message)
		a.Contains(e0.Location, "fault/fault_test.go:3")

		e1 := chain.Errors[1]
		a.Equal("failed to call function", e1.Message)
		a.Contains(e1.Location, "fault/fault_test.go:4")

		e2 := chain.Errors[2]
		a.Equal("", e2.Message)
		a.Contains(e2.Location, "fault/fault_test.go:5")
	})
}
