package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"math"
)

// Int32Pointer Convert int64 to *int32 with error handling
func Int32Pointer(input int64) (*int32, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check for overflow before conversion
	if input <= int64(math.MaxInt32) && input >= int64(math.MinInt32) {
		result := int32(input)
		return &result, diags
	} else {
		diags.AddError("Value Conversion Error", fmt.Sprintf("input %d is out of int32 range", input))
		return nil, diags
	}
}

// Float32Pointer Convert float64 to *float32 with error handling
func Float32Pointer(input float64) (*float32, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check for overflow before conversion
	if input <= math.MaxFloat32 && input >= -math.MaxFloat32 {
		result := float32(input)
		return &result, diags
	} else {
		diags.AddError("Value Conversion Error", fmt.Sprintf("input %f is out of float32 range", input))
		return nil, diags
	}
}

// Float64Pointer Convert float32 to *float64 with error handling
func Float64Pointer(input float32) (*float64, diag.Diagnostics) {
	var diags diag.Diagnostics

	// No overflow check needed since float64 has a larger range than float32
	result := float64(input)
	return &result, diags
}
