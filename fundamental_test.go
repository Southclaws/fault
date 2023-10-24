package fault

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testNewExampleWraper struct{ e error }

func (e *testNewExampleWraper) Error() string { return "example wrapper" }
func (e *testNewExampleWraper) Unwrap() error { return e.e }

func testNewExampleWrap() func(error) error {
	return func(err error) error {
		return &testNewExampleWraper{err}
	}
}

func TestNew(t *testing.T) {
	r := require.New(t)
	var err error

	err = New("TestNew example")
	r.Error(err)

	err = New("TestNew example", testNewExampleWrap())
	r.Error(err)
	_, casts := err.(*testNewExampleWraper)
	r.True(casts)
}

func TestNewf(t *testing.T) {
	r := require.New(t)

	err := Newf("TestNew example %s", "one")
	r.Error(err)
	r.Equal("TestNew example one", err.Error())
}
