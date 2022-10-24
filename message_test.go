package fault_test

import (
	"testing"

	"github.com/Southclaws/fault"
	"github.com/stretchr/testify/assert"
)

func TestFault_Msg(t *testing.T) {
	err := fault.New("root cause")

	err = fault.Wrap(err, fault.Msg("oh no"))
	err = fault.Wrap(err)
	err = fault.Wrap(err, fault.Msg("second"))
	err = fault.Wrap(err)

	assert.Equal(t, "second: oh no: root cause", err.Error())
}
