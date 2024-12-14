package testutils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"reflect"
	"strings"
	"testing"
)

// ValidateResourceSchemaModelConsistency validates schema-model alignment for a resource.
func ValidateResourceSchemaModelConsistency(t *testing.T, schemaAttributes map[string]schema.Attribute, modelStruct interface{}) {
	if err := ResourceSchemaModelConsistency(schemaAttributes, modelStruct, "Resource"); err != nil {
		t.Error(err)
	}
}

// ResourceSchemaModelConsistency performs schema-model validation for resources.
func ResourceSchemaModelConsistency(schemaAttributes map[string]schema.Attribute, modelStruct interface{}, entityType string) error {
	modelType := reflect.TypeOf(modelStruct)

	// Ensure modelStruct is a struct
	if modelType.Kind() != reflect.Struct {
		return fmt.Errorf("%s modelStruct must be a struct; got %s", entityType, modelType.Kind())
	}

	// Map model fields by their tfsdk tags for efficient lookups
	modelFields := extractModelFieldsByTag(modelType)

	// Perform validation checks
	var errors []string
	errors = append(errors, validateResourceSchemaAttributes(schemaAttributes, modelFields, entityType)...)
	errors = append(errors, validateResourceModelFields(schemaAttributes, modelFields, entityType)...)

	// Aggregate errors and return
	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

// validateSchemaAttributes checks if all schema attributes exist in the model fields.
func validateResourceSchemaAttributes(schemaAttributes map[string]schema.Attribute, modelFields map[string]string, entityType string) []string {
	var errors []string
	for attrKey := range schemaAttributes {
		if _, found := modelFields[attrKey]; !found {
			errors = append(errors, fmt.Sprintf("%s schema attribute %q does not match any field in the model struct", entityType, attrKey))
		}
	}
	return errors
}

// validateModelFields checks if all model fields exist in the schema attributes.
func validateResourceModelFields(schemaAttributes map[string]schema.Attribute, modelFields map[string]string, entityType string) []string {
	var errors []string
	for tfsdkTag, fieldName := range modelFields {
		if _, found := schemaAttributes[tfsdkTag]; !found {
			errors = append(errors, fmt.Sprintf("%s model field %q (tfsdk tag: %q) does not match any schema attribute", entityType, fieldName, tfsdkTag))
		}
	}
	return errors
}
