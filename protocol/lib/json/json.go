package json

import (
	"encoding/json"
)

// IsValidJSON checks if a JSON string is well-formed.
func IsValidJSON(str string) error {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js)
}
