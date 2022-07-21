package fault_test

import (
	"encoding/json"
	"testing"

	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/fault"
)

func TestMarshalJSON(t *testing.T) {
	fn := func() error {
		return pkg_errors.New("root cause")
	}

	rootcause := fn()

	level1 := fault.WithValue(rootcause, "level 1", "context_at_level_1", "x")
	level2 := fault.WithValue(level1, "level 2", "context_at_level_2", "y")
	level3 := fault.WithValue(level2, "level 3", "context_at_level_3", "z")

	b, err := json.MarshalIndent(level3, "", "  ")
	assert.NoError(t, err)
	assert.Contains(t, string(b), "marshal_json_test.go:22")
	assert.Contains(t, string(b), "marshal_json_test.go:21")
	assert.Contains(t, string(b), "marshal_json_test.go:20")
}
