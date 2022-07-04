package fault

import (
	"encoding/json"
)

const SerialisedErrorKey = "error"

// MarshalJSON implements the json.Marshaler interface for the Fault error type.
func (e *fault) MarshalJSON() ([]byte, error) {
	return json.Marshal(Context(e))
}
