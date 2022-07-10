package fault

import "encoding/json"

// MarshalJSON implements the json.Marshaler interface for the Fault error type.
// This allows you to easily serialise these errors for responses or logging.
func (e *fault) MarshalJSON() ([]byte, error) {
	return json.Marshal(Context(e))
}
