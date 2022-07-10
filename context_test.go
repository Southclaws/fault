package fault_test

import (
	"errors"
	"testing"

	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/fault"
)

var (
	ErrNativeSentinel    = errors.New("native sentinel error")
	ErrPkgErrorsSentinel = pkg_errors.New("pkg/errors sentinel")
	ErrSentinel          = fault.New("fault sentinel error")
)

func TestNativeError(t *testing.T) {
	want := fault.ErrorInfo{
		Message: "oops",
		Values:  map[string]any{},
		Trace:   []fault.Location{},
	}
	got := fault.Context(errors.New("oops"))
	assert.Equal(t, want, got)
}

func TestPkgErrorsWrapped(t *testing.T) {
	want := fault.ErrorInfo{
		Message: "some context: oops",
		Values:  map[string]any{},
		Trace: []fault.Location{
			{
				Message:  "some context: oops",
				Location: "context_test.go",
			},
		},
	}
	got := fault.Context(pkg_errors.Wrap(errors.New("oops"), "some context"))
	assert.Equal(t, want, got)
}

func TestNativeSentinel(t *testing.T) {
	want := fault.ErrorInfo{
		Message: "native sentinel error",
		Values:  map[string]any{},
		Trace:   []fault.Location{},
	}
	got := fault.Context(ErrNativeSentinel)
	assert.Equal(t, want, got)
}

func TestNativeSentinelWrapped(t *testing.T) {
	want := fault.ErrorInfo{
		Message: "some context: native sentinel error",
		Values:  map[string]any{},
		Trace: []fault.Location{
			{
				Message:  "some context: native sentinel error",
				Location: "context_test.go",
			},
		},
	}
	got := fault.Context(pkg_errors.Wrap(ErrNativeSentinel, "some context"))
	assert.Equal(t, want, got)
}

func TestFaultSentinel(t *testing.T) {
	want := fault.ErrorInfo{
		Message: "fault sentinel error",
		Values:  map[string]any{},
		Trace:   []fault.Location{},
	}
	got := fault.Context(ErrSentinel)
	assert.Equal(t, want.Message, got.Message)
	assert.Equal(t, want.Values, got.Values)
	assert.Len(t, got.Trace, 1)
	assert.Equal(t, "fault sentinel error", got.Trace[0].Message)
	assert.Contains(t, got.Trace[0].Location, "fault/context_test.go")
}

func TestFaultSentinelWrapped(t *testing.T) {
	want := fault.ErrorInfo{
		Message: "some context: fault sentinel error",
		Values:  map[string]any{},
		Trace:   []fault.Location{},
	}
	got := fault.Context(pkg_errors.Wrap(ErrSentinel, "some context"))
	assert.Equal(t, want.Message, got.Message)
	assert.Equal(t, want.Values, got.Values)
	assert.Len(t, got.Trace, 2)
	assert.Equal(t, "some context: fault sentinel error", got.Trace[0].Message)
	assert.Equal(t, "fault sentinel error", got.Trace[1].Message)
	assert.Contains(t, got.Trace[0].Location, "context_test.go")
	assert.Contains(t, got.Trace[1].Location, "context_test.go")
}
