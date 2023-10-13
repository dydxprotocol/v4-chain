package json

import (
	"encoding/json"
)

// IsValidJSON checks if a JSON string is well-formed.
func IsValidJSON(str string) error {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(str), &js)
	if err != nil {
		return err
	}
	return nil
}
