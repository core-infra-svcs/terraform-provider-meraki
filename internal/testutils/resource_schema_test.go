package testutils_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"testing"
)

type MockResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

var MockResourceSchema = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The unique identifier for the resource.",
		Computed:            true,
	},
	"name": schema.StringAttribute{
		MarkdownDescription: "The name of the resource.",
		Optional:            true,
	},
	"enabled": schema.BoolAttribute{
		MarkdownDescription: "Indicates whether the resource is enabled.",
		Optional:            true,
	},
}

func TestResourceSchemaModelConsistency_Matching(t *testing.T) {
	err := testutils.ResourceSchemaModelConsistency(MockResourceSchema, MockResourceModel{}, "Resource")
	if err != nil {
		t.Errorf("Unexpected error for matching schema and model: %s", err)
	}
}

func TestResourceSchemaModelConsistency_SchemaExtraAttribute(t *testing.T) {
	extraSchema := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the resource.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the resource.",
			Optional:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the resource is enabled.",
			Optional:            true,
		},
		"extra_field": schema.StringAttribute{
			MarkdownDescription: "An extra field not in the model.",
			Optional:            true,
		},
	}

	err := testutils.ResourceSchemaModelConsistency(extraSchema, MockResourceModel{}, "Resource")
	if err == nil {
		t.Errorf("Expected error due to schema attribute 'extra_field' not matching model, but got none")
		return
	}

	expectedError := `Resource schema attribute "extra_field" does not match any field in the model struct`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}

func TestResourceSchemaModelConsistency_ModelExtraField(t *testing.T) {
	type ExtraFieldModel struct {
		ID         types.String `tfsdk:"id"`
		Name       types.String `tfsdk:"name"`
		Enabled    types.Bool   `tfsdk:"enabled"`
		ExtraField types.String `tfsdk:"extra_field"`
	}

	err := testutils.ResourceSchemaModelConsistency(MockResourceSchema, ExtraFieldModel{}, "Resource")
	if err == nil {
		t.Errorf("Expected error due to model field 'extra_field' not matching schema, but got none")
		return
	}

	expectedError := `Resource model field "ExtraField" (tfsdk tag: "extra_field") does not match any schema attribute`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}

func TestResourceSchemaModelConsistency_NonStructModel(t *testing.T) {
	nonStructModel := "invalid"

	err := testutils.ResourceSchemaModelConsistency(MockResourceSchema, nonStructModel, "Resource")
	if err == nil {
		t.Errorf("Expected error due to non-struct model, but got none")
		return
	}

	expectedError := `Resource modelStruct must be a struct; got string`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}

func TestResourceSchemaModelConsistency_EmptySchema(t *testing.T) {
	emptySchema := map[string]schema.Attribute{}

	err := testutils.ResourceSchemaModelConsistency(emptySchema, MockResourceModel{}, "Resource")
	if err == nil {
		t.Errorf("Expected error due to model fields not matching empty schema, but got none")
		return
	}

	expectedErrors := []string{
		`Resource model field "ID" (tfsdk tag: "id") does not match any schema attribute`,
		`Resource model field "Name" (tfsdk tag: "name") does not match any schema attribute`,
		`Resource model field "Enabled" (tfsdk tag: "enabled") does not match any schema attribute`,
	}

	for _, expectedError := range expectedErrors {
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("Expected error message not found: %q\nGot:\n%s", expectedError, err.Error())
		}
	}
}

func TestResourceSchemaModelConsistency_EmptyModel(t *testing.T) {
	type EmptyModel struct{}

	err := testutils.ResourceSchemaModelConsistency(MockResourceSchema, EmptyModel{}, "Resource")
	if err == nil {
		t.Errorf("Expected error due to schema attributes not matching empty model, but got none")
		return
	}

	expectedErrors := []string{
		`Resource schema attribute "id" does not match any field in the model struct`,
		`Resource schema attribute "name" does not match any field in the model struct`,
		`Resource schema attribute "enabled" does not match any field in the model struct`,
	}

	for _, expectedError := range expectedErrors {
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("Expected error message not found: %q\nGot:\n%s", expectedError, err.Error())
		}
	}
}
