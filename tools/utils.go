package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strconv"
)

// MapStringValue Extracts a string from an interface and returns a Terraform type
func MapStringValue(m map[string]interface{}, key string, diags *diag.Diagnostics) types.String {
	var result types.String

	if v := m[key]; v != nil {
		result = types.StringValue(v.(string))
	} else {
		diags.AddWarning(
			"String extraction error",
			fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, v))
		result = types.StringNull()
	}

	return result
}

// MapBoolValue Extracts a boolean from an interface and returns a Terraform type
func MapBoolValue(m map[string]interface{}, key string, diags *diag.Diagnostics) types.Bool {
	var result types.Bool
	if v := m[key]; v != nil {

		if _, ok := v.(string); ok {

			b, _ := strconv.ParseBool(v.(string))
			result = types.BoolValue(b)
		} else {
			diags.AddWarning(
				"Bool extraction error",
				fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, v))
			result = types.BoolValue(v.(bool))
		}

	} else {
		result = types.BoolNull()
	}
	return result
}

// MerakiBoolType -
type MerakiBoolType struct {
	basetypes.BoolType
}

func (mst MerakiBoolType) Equal(o attr.Type) bool {
	_, ok := o.(MerakiBoolType)
	return ok
}

func (mst MerakiBoolType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mst.BoolType.ValueFromTerraform(ctx, in)

	return MerakiBool{BoolValue: val.(types.Bool)}, err
}

type MerakiBool struct {
	basetypes.BoolValue
}

func (m MerakiBool) Type(_ context.Context) attr.Type {
	return MerakiBoolType{}
}

func (m MerakiBool) Equal(value attr.Value) bool {
	if v, ok := value.(basetypes.BoolValue); ok {
		return m.BoolValue.Equal(v)
	}

	v, ok := value.(MerakiBool)
	if !ok {
		return false
	}

	return m.BoolValue.Equal(v.BoolValue)
}

func (m *MerakiBool) UnmarshalJSON(bytes []byte) error {
	m.BoolValue = types.BoolNull()

	var b *bool
	if err := json.Unmarshal(bytes, &b); err != nil {
		return err
	}

	if b != nil {
		m.BoolValue = types.BoolValue(*b)
	}
	return nil
}

// MerakiStringType -
type MerakiStringType struct {
	basetypes.StringType
}

func (mst MerakiStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mst.StringType.ValueFromTerraform(ctx, in)

	return MerakiString{StringValue: val.(types.String)}, err
}

func (mst MerakiStringType) Equal(o attr.Type) bool {
	_, ok := o.(MerakiStringType)
	return ok
}

type MerakiString struct {
	basetypes.StringValue
}

func (m *MerakiString) UnmarshalJSON(bytes []byte) error {
	m.StringValue = types.StringNull()

	var b *string
	if err := json.Unmarshal(bytes, &b); err != nil {
		return err
	}

	if b != nil {
		m.StringValue = types.StringValue(*b)
	}
	return nil
}

func (m MerakiString) Type(_ context.Context) attr.Type {
	return MerakiStringType{}
}

func (m MerakiString) Equal(value attr.Value) bool {
	if v, ok := value.(basetypes.StringValue); ok {
		return m.StringValue.Equal(v)
	}

	v, ok := value.(MerakiString)
	if !ok {
		return false
	}

	return m.StringValue.Equal(v.StringValue)
}
