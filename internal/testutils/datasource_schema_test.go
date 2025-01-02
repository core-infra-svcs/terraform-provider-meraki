package testutils_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"testing"
)

type MockDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

var MockDataSourceSchema = map[string]schema.Attribute{
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

func TestDataSourceSchemaModelConsistency_Matching(t *testing.T) {
	// Positive test: The schema and model match perfectly
	err := testutils.DataSourceSchemaModelConsistency(MockDataSourceSchema, MockDataSourceModel{}, "Data Source")
	if err != nil {
		t.Errorf("Unexpected error for matching schema and model: %s", err)
	}
}

func TestDataSourceSchemaModelConsistency_SchemaExtraAttribute(t *testing.T) {
	// Schema has an extra attribute not in the model
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

	err := testutils.DataSourceSchemaModelConsistency(extraSchema, MockDataSourceModel{}, "Data Source")
	if err == nil {
		t.Errorf("Expected error due to schema attribute 'extra_field' not matching model, but got none")
		return
	}

	// Check if the error message is as expected
	expectedError := `Data Source schema attribute "extra_field" does not match any field in the model struct`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}

func TestDataSourceSchemaModelConsistency_ModelExtraField(t *testing.T) {
	// Model has an extra field not in the schema
	type ExtraFieldModel struct {
		ID         types.String `tfsdk:"id"`
		Name       types.String `tfsdk:"name"`
		Enabled    types.Bool   `tfsdk:"enabled"`
		ExtraField types.String `tfsdk:"extra_field"`
	}

	err := testutils.DataSourceSchemaModelConsistency(MockDataSourceSchema, ExtraFieldModel{}, "Data Source")
	if err == nil {
		t.Errorf("Expected error due to model field 'extra_field' not matching schema, but got none")
		return
	}

	// Check if the error message is as expected
	expectedError := `Data Source model field "ExtraField" (tfsdk tag: "extra_field") does not match any schema attribute`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}

func TestDataSourceSchemaModelConsistency_NonStructModel(t *testing.T) {
	// Non-struct model should result in an error
	nonStructModel := "invalid"

	err := testutils.DataSourceSchemaModelConsistency(MockDataSourceSchema, nonStructModel, "Data Source")
	if err == nil {
		t.Errorf("Expected error due to non-struct model, but got none")
		return
	}

	expectedError := `Data Source modelStruct must be a struct; got string`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}

// Test case for empty schema
func TestDataSourceSchemaModelConsistency_EmptySchema(t *testing.T) {
	emptySchema := map[string]schema.Attribute{}

	err := testutils.DataSourceSchemaModelConsistency(emptySchema, MockDataSourceModel{}, "Data Source")
	if err == nil {
		t.Errorf("Expected error due to model fields not matching empty schema, but got none")
		return
	}

	expectedErrors := []string{
		`Data Source model field "ID" (tfsdk tag: "id") does not match any schema attribute`,
		`Data Source model field "Name" (tfsdk tag: "name") does not match any schema attribute`,
		`Data Source model field "Enabled" (tfsdk tag: "enabled") does not match any schema attribute`,
	}

	for _, expectedError := range expectedErrors {
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("Expected error message not found: %q\nGot:\n%s", expectedError, err.Error())
		}
	}
}
func TestDataSourceSchemaModelConsistency_EmptyModel(t *testing.T) {
	// Empty model with a valid schema should not throw any errors
	type EmptyModel struct{}

	err := testutils.DataSourceSchemaModelConsistency(MockDataSourceSchema, EmptyModel{}, "Data Source")
	if err == nil {
		t.Errorf("Expected error due to schema attributes not matching empty model, but got none")
		return
	}

	expectedError := `Data Source schema attribute "id" does not match any field in the model struct`
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message not found. Got:\n%s", err.Error())
	}
}
