package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSafeFormatRFC3339(t *testing.T) {
	ctx := context.Background()

	// Test case: Valid time pointer
	t.Run("valid time pointer", func(t *testing.T) {
		testTime := time.Date(2024, 11, 15, 10, 0, 0, 0, time.UTC)
		formatted := SafeFormatRFC3339(ctx, &testTime, time.RFC3339)
		assert.Equal(t, "2024-11-15T10:00:00Z", formatted, "Expected formatted time to match RFC3339 format")
	})

	// Test case: Nil time pointer
	t.Run("nil time pointer", func(t *testing.T) {
		formatted := SafeFormatRFC3339(ctx, nil, time.RFC3339)
		assert.Equal(t, "", formatted, "Expected empty string for nil time pointer")
	})
}

func TestValidateRFC3339(t *testing.T) {
	ctx := context.Background()

	// Test case: Valid RFC3339 timestamp
	t.Run("valid RFC3339 timestamp", func(t *testing.T) {
		err := ValidateRFC3339(ctx, "2024-11-15T10:00:00Z")
		assert.NoError(t, err, "Expected no error for valid RFC3339 timestamp")
	})

	// Test case: Invalid RFC3339 timestamp
	t.Run("invalid RFC3339 timestamp", func(t *testing.T) {
		err := ValidateRFC3339(ctx, "15-11-2024 10:00:00")
		assert.Error(t, err, "Expected error for invalid RFC3339 timestamp")
		assert.Contains(t, err.Error(), "received timestamp does not match RFC3339 format", "Expected specific error message")
	})

	// Test case: Empty string
	t.Run("empty string", func(t *testing.T) {
		err := ValidateRFC3339(ctx, "")
		assert.Error(t, err, "Expected error for empty string")
		assert.Contains(t, err.Error(), "received timestamp does not match RFC3339 format", "Expected specific error message")
	})
}
