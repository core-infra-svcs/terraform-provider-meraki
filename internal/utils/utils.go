package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"strconv"
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

func ListInt64TypeToInt32Array(data types.List) ([]int32, diag.Diagnostics) {
	var diags diag.Diagnostics
	var int32Array []int32
	for _, i := range data.Elements() {
		groupIdInt, err := strconv.Atoi(i.String())
		if err != nil {
			diags.AddError("Failed to extract list of int64 to int32", err.Error())
			return int32Array, diags
		}
		int32Array = append(int32Array, int32(groupIdInt))
	}

	return int32Array, diags

}
