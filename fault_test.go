package fault_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Southclaws/fault"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

type User struct{ name string }

var ErrBadID = fault.Sentinel("user ID invalid")

// Or, same as normal
// var ErrBadID = errors.New("user ID invalid")

func GetUser(ctx context.Context, id int) (*User, error) {
	if id < 0 {
		return nil, ErrBadID
	}

	user, err := databaseCall()
	if err != nil {
		return nil, fault.Wrap(err, "failed to call database")
	}

	if user == nil {
		return nil, fault.WithValue(err, "user not found", "user_id", "admin123")
	}

	return user, nil
}

func databaseCall() (*User, error) {
	return nil, errors.New("internal problem")
}

func TestStuff(t *testing.T) {
	ctx := context.WithValue(context.Background(), "userID", "admin123")

	_, err := GetUser(ctx, 124)

	b, _ := json.MarshalIndent(err, "", "  ")

	fmt.Println(string(b))

	pretty.Println(err)
}

func TestUnwrap(t *testing.T) {
	ctx := context.WithValue(context.Background(), "userID", "admin123")

	// fmt.Println(errors.Is(err, ErrBadID))

	_, err := GetUser(ctx, -2)

	err = fault.WithValue(err, "step 1", "userID", "admin123")

	err = fault.WithValue(errors.Wrap(err, "problem"), "step 2", "traceID", "0xdead")

	err = fault.WithValue(err, "step 3", "requestID", "69")

	// meta := fault.Context(err)
	// trace := fault.Trace(err)

	// pretty.Println(meta)
	// pretty.Println(trace)

	b, _ := json.MarshalIndent(err, "", "  ")
	fmt.Println(string(b))

	_, err = GetUser(ctx, 1)

	err = fault.WithValue(err, "step 1", "userID", "admin123")

	err = fault.WithValue(errors.Wrap(err, "failed to do xyz"), "step 2", "traceID", "0xdead")

	err = fault.WithValue(err, "step 3", "requestID", "69")

	b, _ = json.MarshalIndent(err, "", "  ")
	fmt.Println(string(b))
}
