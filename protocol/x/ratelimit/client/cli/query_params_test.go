//go:build all || integration_test

package cli_test

import (
	"strconv"
	"testing"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryParams(t *testing.T) {
	// TODO(CORE-823): implement query for `x/ratelimit`
}
