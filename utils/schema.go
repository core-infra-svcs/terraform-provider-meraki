package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewStringDefault returns a struct that implements the defaults.String interface for a resource Schema
func NewStringDefault(defaultValue string) defaults.String {
	return &stringDefault{defaultValue: defaultValue}
}

type stringDefault struct {
	defaultValue string
}

func (d *stringDefault) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value is '%s'", d.defaultValue)
}

func (d *stringDefault) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value is `%s`", d.defaultValue)
}

func (d *stringDefault) DefaultString(ctx context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringValue(d.defaultValue)
}

// NewInt64Default returns a struct that implements the defaults.Int64 interface for a resource Schema
func NewInt64Default(defaultValue int64) defaults.Int64 {
	return &int64Default{defaultValue: defaultValue}
}

type int64Default struct {
	defaultValue int64
}

func (d *int64Default) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value is '%v'", d.defaultValue)
}

func (d *int64Default) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value is `%s`", d.defaultValue)
}

func (d *int64Default) DefaultInt64(ctx context.Context, req defaults.Int64Request, resp *defaults.Int64Response) {
	resp.PlanValue = types.Int64Value(d.defaultValue)
}

// NewBoolDefault returns a struct that implements the defaults.Bool interface for a resource Schema
func NewBoolDefault(defaultValue bool) defaults.Bool {
	return &boolDefault{defaultValue: defaultValue}
}

type boolDefault struct {
	defaultValue bool
}

func (d *boolDefault) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value is '%s'", d.defaultValue)
}

func (d *boolDefault) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value is `%s`", d.defaultValue)
}

func (d *boolDefault) DefaultBool(ctx context.Context, req defaults.BoolRequest, resp *defaults.BoolResponse) {
	resp.PlanValue = types.BoolValue(d.defaultValue)
}

// NewFloat64Default returns a struct that implements the defaults.Float64 interface for a resource Schema
func NewFloat64Default(defaultValue float64) defaults.Float64 {
	return &float64Default{defaultValue: defaultValue}
}

type float64Default struct {
	defaultValue float64
}

func (d *float64Default) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value is '%v'", d.defaultValue)
}

func (d *float64Default) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value is `%s`", d.defaultValue)
}

func (d *float64Default) DefaultFloat64(ctx context.Context, req defaults.Float64Request, resp *defaults.Float64Response) {
	resp.PlanValue = types.Float64Value(d.defaultValue)
}

// RequiresReplaceIfSensitive Custom plan modifier to require replacement if the sensitive attribute changes
type RequiresReplaceIfSensitive struct{}

func (r RequiresReplaceIfSensitive) Description(context.Context) string {
	return "Requires replacement if sensitive attribute changes."
}

func (r RequiresReplaceIfSensitive) MarkdownDescription(context.Context) string {
	return "Requires replacement if sensitive attribute changes."
}

func (r RequiresReplaceIfSensitive) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueString() != req.StateValue.ValueString() {
		resp.RequiresReplace = true
	}
}
