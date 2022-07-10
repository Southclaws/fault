package fault_test

import (
	"errors"
	"testing"

	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/fault"
)

func TestWithValueNativeErrors(t *testing.T) {
	fn := func() error {
		return errors.New("root cause")
	}

	rootcause := fn()

	level1 := fault.WithValue(rootcause, "level 1", "context_at_level_1", "x")
	level2 := fault.WithValue(level1, "level 2", "context_at_level_2", "y")
	level3 := fault.WithValue(level2, "level 3", "context_at_level_3", "z")

	gotString := level3.Error()
	gotContext := fault.Context(level3)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotString)
	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Message)

	assert.Len(t, gotContext.Trace, 3)
	assert.Equal(t, []fault.Location{
		{
			Message:  "level 3: level 2: level 1: root cause",
			Location: "/Users/southclaws/Work/fault/with_value_test.go:22",
		},
		{
			Message:  "level 2: level 1: root cause",
			Location: "/Users/southclaws/Work/fault/with_value_test.go:21",
		},
		{
			Message:  "level 1: root cause",
			Location: "/Users/southclaws/Work/fault/with_value_test.go:20",
		},
	}, gotContext.Trace)

	assert.Equal(t, map[string]any{
		"context_at_level_1": "x",
		"context_at_level_2": "y",
		"context_at_level_3": "z",
	}, gotContext.Values)
}

func TestWithValuePkgErrors(t *testing.T) {
	fn := func() error {
		return pkg_errors.New("root cause")
	}

	rootcause := fn()

	level1 := fault.WithValue(rootcause, "level 1", "context_at_level_1", "x")
	level2 := fault.WithValue(level1, "level 2", "context_at_level_2", "y")
	level3 := fault.WithValue(level2, "level 3", "context_at_level_3", "z")

	gotString := level3.Error()
	gotContext := fault.Context(level3)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotString)
	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Message)

	assert.Len(t, gotContext.Trace, 4)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Trace[0].Message)
	assert.Contains(t, gotContext.Trace[0].Location, "with_value_test")

	assert.Equal(t, "level 2: level 1: root cause", gotContext.Trace[1].Message)
	assert.Contains(t, gotContext.Trace[1].Location, "with_value_test")

	assert.Equal(t, "level 1: root cause", gotContext.Trace[2].Message)
	assert.Contains(t, gotContext.Trace[2].Location, "with_value_test")

	assert.Equal(t, "root cause", gotContext.Trace[3].Message)
	assert.Contains(t, gotContext.Trace[0].Location, "with_value_test.go")

	assert.Equal(t, map[string]any{
		"context_at_level_1": "x",
		"context_at_level_2": "y",
		"context_at_level_3": "z",
	}, gotContext.Values)
}

func TestWithValuePkgErrorsMixed(t *testing.T) {
	fn := func() error {
		return pkg_errors.New("root cause")
	}

	rootcause := fn()

	level1 := fault.WithValue(rootcause, "level 1", "context_at_level_1", "x")
	level2 := pkg_errors.Wrap(level1, "level 2")
	level3 := fault.WithValue(level2, "level 3", "context_at_level_3", "z")

	gotString := level3.Error()
	gotContext := fault.Context(level3)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotString)
	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Message)

	assert.Len(t, gotContext.Trace, 4)

	assert.Equal(t, "level 3: level 2: level 1: root cause", gotContext.Trace[0].Message)
	assert.Contains(t, gotContext.Trace[0].Location, "with_value_test")

	assert.Equal(t, "level 2: level 1: root cause", gotContext.Trace[1].Message)
	assert.Contains(t, gotContext.Trace[1].Location, "with_value_test")

	assert.Equal(t, "level 1: root cause", gotContext.Trace[2].Message)
	assert.Contains(t, gotContext.Trace[2].Location, "with_value_test")

	assert.Equal(t, "root cause", gotContext.Trace[3].Message)
	assert.Contains(t, gotContext.Trace[0].Location, "with_value_test.go")

	assert.Equal(t, map[string]any{
		"context_at_level_1": "x",
		"context_at_level_3": "z",
	}, gotContext.Values)
}

func TestWithValueNil(t *testing.T) {
	assert.Nil(t, fault.WithValue(nil, "", "", ""))
}
