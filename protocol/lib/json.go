package lib

import (
	"encoding/json"
	"fmt"
)

// MaybeGetStructJsonString returns the json representation of a struct, or a formatted string using
// %+v if the json conversion encounters an error.
func MaybeGetJsonString(i interface{}) string {
	jsonData, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%+v", i)
	}

	return string(jsonData)
}
