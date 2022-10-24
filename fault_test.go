package fault_test

import (
	"testing"

	"github.com/Southclaws/fault"
	"github.com/kr/pretty"
)

func Test_fault_Format(t *testing.T) {
	err := fault.New("root cause")
	err = fault.Wrap(err, fault.Msg("wrapped"))
	err = fault.Wrap(err, fault.Msg("wrapped"))

	f := err.(interface{ Stack() fault.Stack })

	pretty.Printf("%+v\n", f.Stack().String())
}
