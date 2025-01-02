package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
	"time"
)

// SafeFormatRFC3339 safely formats a *time.Time to a string using the provided layout.
// Returns an empty string if the time pointer is nil.
func SafeFormatRFC3339(ctx context.Context, t *time.Time, layout string) string {
	if t == nil {
		tflog.Trace(ctx, "safeFormatTime: Received nil time pointer, returning empty string")
		return ""
	}

	formattedTime := t.Format(layout)
	tflog.Trace(ctx, "safeFormatTime: Successfully formatted time", map[string]interface{}{
		"formattedTime": formattedTime,
		"layout":        layout,
	})
	return formattedTime
}

// ValidateRFC3339 validates whether a given string conforms to the RFC3339 format.
// Returns an error if the format is invalid.
func ValidateRFC3339(value string) error {
	ctx := context.Background()
	tflog.Trace(ctx, "validateRFC3339: Validating RFC3339 format", map[string]interface{}{
		"value": value,
	})

	re := regexp.MustCompile(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`)
	if !re.MatchString(value) {
		err := fmt.Errorf("received timestamp does not match RFC3339 format: %s", value)
		tflog.Trace(ctx, "validateRFC3339: Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	tflog.Trace(ctx, "validateRFC3339: Validation succeeded", map[string]interface{}{
		"value": value,
	})
	return nil
}
