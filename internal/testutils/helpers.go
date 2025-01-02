package testutils

import (
	"reflect"
)

// extractModelFieldsByTag creates a map of model fields keyed by their tfsdk tags.
func extractModelFieldsByTag(modelType reflect.Type) map[string]string {
	modelFields := make(map[string]string)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tfsdkTag := field.Tag.Get("tfsdk")
		if tfsdkTag != "" {
			modelFields[tfsdkTag] = field.Name
		}
	}
	return modelFields
}
