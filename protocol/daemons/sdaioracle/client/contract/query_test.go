package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryDaiConversionRate(t *testing.T) {
	// Test with uninitialized client
	assert.Panics(t, func() {
		_, _ = QueryDaiConversionRate(nil)
	}, "Expected panic with uninitialized client")
}
