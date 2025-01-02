package utils

import (
	"context"
	"fmt"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConvertResourceSchemaToDataSourceSchema converts a resource schema to a data source schema while carrying over all settings.
func ConvertResourceSchemaToDataSourceSchema(resourceAttrs map[string]schema.Attribute) map[string]datasourceSchema.Attribute {
	dataSourceAttrs := make(map[string]datasourceSchema.Attribute)
	for key, attr := range resourceAttrs {
		dataSourceAttrs[key] = convertResourceAttrToDataSourceAttr(attr)
	}
	return dataSourceAttrs
}

// Helper function to convert a single resource schema attribute to a data source schema attribute.
func convertResourceAttrToDataSourceAttr(attr schema.Attribute) datasourceSchema.Attribute {
	switch a := attr.(type) {
	case schema.StringAttribute:
		return datasourceSchema.StringAttribute{
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	case schema.BoolAttribute:
		return datasourceSchema.BoolAttribute{
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	case schema.Int64Attribute:
		return datasourceSchema.Int64Attribute{
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	case schema.ListAttribute:
		return datasourceSchema.ListAttribute{
			ElementType:         a.ElementType,
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	case schema.SingleNestedAttribute:
		return datasourceSchema.SingleNestedAttribute{
			Attributes:          ConvertResourceSchemaToDataSourceSchema(a.Attributes),
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	case schema.MapAttribute:
		return datasourceSchema.MapAttribute{
			ElementType:         a.ElementType,
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	case schema.SetAttribute:
		return datasourceSchema.SetAttribute{
			ElementType:         a.ElementType,
			CustomType:          a.CustomType,
			MarkdownDescription: a.MarkdownDescription,
			Description:         a.Description,
			DeprecationMessage:  a.DeprecationMessage,
			Validators:          a.Validators,
			Computed:            true,
			Sensitive:           a.Sensitive,
		}
	default:
		return nil // Handle unsupported attribute types gracefully
	}
}

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
	return fmt.Sprintf("Default value is `%v`", d.defaultValue)
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
	return fmt.Sprintf("Default value is '%v'", d.defaultValue)
}

func (d *boolDefault) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value is `%v`", d.defaultValue)
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
	return fmt.Sprintf("Default value is `%v`", d.defaultValue)
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
