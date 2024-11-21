package cellular

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// mapModelToApiPayload maps the Terraform resource model to the API request payload.
func mapModelToApiPayload(model *resourceModel) (*openApiClient.UpdateDeviceCellularSimsRequest, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	payload := openApiClient.NewUpdateDeviceCellularSimsRequest()

	// Map Sims
	simsPayload, simsDiags := mapSimsToApiPayload(model.Sims)
	diagnostics.Append(simsDiags...)
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	payload.SetSims(simsPayload)

	// Map SimFailOver
	failOverPayload, failOverDiags := mapSimFailOverToApiPayload(model.SimFailOver)
	diagnostics.Append(failOverDiags...)
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	if failOverPayload != nil {
		payload.SimFailover = failOverPayload
	}

	// Map SimOrdering (handles ordering of SIMs by slot)
	if !model.SimOrdering.IsNull() && !model.SimOrdering.IsUnknown() {
		// Convert Terraform Set to []string
		simOrdering := asStringSlice(model.SimOrdering)

		// Convert []string to []UpdateDeviceCellularSimsRequestSimsInner
		orderedSims := make([]openApiClient.UpdateDeviceCellularSimsRequestSimsInner, len(simOrdering))
		for i, slot := range simOrdering {
			orderedSims[i] = openApiClient.UpdateDeviceCellularSimsRequestSimsInner{
				Slot: &slot,
			}
		}
		payload.SetSims(orderedSims)
	}

	return payload, diagnostics
}

// mapApiResponseToModel maps the API response payload back to the Terraform resource model.
func mapApiResponseToModel(apiResponse map[string]interface{}, model *resourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	// Map Sims
	if sims, ok := apiResponse["sims"].([]interface{}); ok {
		model.Sims, diagnostics = mapSimsToTerraformModel(sims)
		if diagnostics.HasError() {
			return diagnostics
		}
	} else {
		model.Sims = types.ListNull(types.ObjectType{AttrTypes: ResourceModelSimAttrTypes()})
	}

	// Map SimFailOver
	if simFailOver, ok := apiResponse["simFailOver"].(map[string]interface{}); ok {
		model.SimFailOver = mapSimFailOverToTerraformModel(simFailOver)
	} else {
		model.SimFailOver = types.ObjectNull(ResourceModelSimFailOverAttrTypes())
	}

	// Map SimOrdering
	model.SimOrdering, diagnostics = mapSimOrderingToTerraformModel(apiResponse["simOrdering"])
	return diagnostics
}

func mapSimOrderingToTerraformModel(data interface{}) (types.Set, diag.Diagnostics) {
	if simOrdering, ok := data.([]interface{}); ok {
		simOrderingStrings := make([]string, len(simOrdering))
		for i, slot := range simOrdering {
			if slotStr, ok := slot.(string); ok {
				simOrderingStrings[i] = slotStr
			}
		}
		return asSetOfStrings(simOrderingStrings), nil
	}
	return types.SetNull(types.StringType), nil
}

func mapSimFailOverToTerraformModel(apiFailOver map[string]interface{}) types.Object {
	return types.ObjectValueMust(ResourceModelSimFailOverAttrTypes(), map[string]attr.Value{
		"enabled": types.BoolValue(apiFailOver["enabled"].(bool)),
		"timeout": types.Int64Value(int64(apiFailOver["timeout"].(float64))),
	})
}

func mapSimsToApiPayload(sims types.List) ([]openApiClient.UpdateDeviceCellularSimsRequestSimsInner, diag.Diagnostics) {
	if sims.IsNull() || sims.IsUnknown() {
		return nil, nil
	}

	var simsData []ResourceModelSim
	diags := sims.ElementsAs(context.Background(), &simsData, false)
	if diags.HasError() {
		return nil, diags
	}

	var apiSims []openApiClient.UpdateDeviceCellularSimsRequestSimsInner
	for _, sim := range simsData {
		slot := sim.Slot.ValueString()
		isPrimary := sim.IsPrimary.ValueBool()

		apiSim := openApiClient.UpdateDeviceCellularSimsRequestSimsInner{
			Slot:      &slot,
			IsPrimary: &isPrimary,
		}

		// Map APNs
		apnsPayload, apnsDiags := mapApnsToApiPayload(sim.Apns)
		diags.Append(apnsDiags...)
		if diags.HasError() {
			return nil, diags
		}
		apiSim.SetApns(apnsPayload)

		apiSims = append(apiSims, apiSim)
	}

	return apiSims, diags
}

func mapSimsToTerraformModel(apiSims []interface{}) (types.List, diag.Diagnostics) {
	if len(apiSims) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: ResourceModelSimAttrTypes()}), nil
	}

	var sims []attr.Value
	diagnostics := diag.Diagnostics{}
	for _, apiSim := range apiSims {
		if simData, ok := apiSim.(map[string]interface{}); ok {
			// Map APNs
			apns, apnsDiags := mapApnsToTerraformModel(simData["apns"].([]openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner))
			diagnostics.Append(apnsDiags...)
			if diagnostics.HasError() {
				return types.ListNull(types.ObjectType{AttrTypes: ResourceModelSimAttrTypes()}), diagnostics
			}

			// Create SIM object
			sim := types.ObjectValueMust(ResourceModelSimAttrTypes(), map[string]attr.Value{
				"slot":       types.StringValue(simData["slot"].(string)),
				"is_primary": types.BoolValue(simData["isPrimary"].(bool)),
				"apns":       apns,
			})
			sims = append(sims, sim)
		}
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: ResourceModelSimAttrTypes()}, sims), diagnostics
}

func mapApnsToApiPayload(apns types.List) ([]openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner, diag.Diagnostics) {
	if apns.IsNull() || apns.IsUnknown() {
		return nil, nil
	}

	var apnsData []ResourceModelApns
	diags := apns.ElementsAs(context.Background(), &apnsData, false)
	if diags.HasError() {
		return nil, diags
	}

	var apiApns []openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
	for _, apn := range apnsData {
		apiApn := openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner{
			Name:           apn.Name.ValueString(),
			AllowedIpTypes: asStringSlice(apn.AllowedIpTypes),
		}

		// Handle Authentication
		if !apn.Authentication.IsNull() && !apn.Authentication.IsUnknown() {
			auth := mapAuthenticationToApiPayload(apn.Authentication)
			if auth != nil {
				apiApn.SetAuthentication(*auth)
			}
		}

		apiApns = append(apiApns, apiApn)
	}

	return apiApns, nil
}

func mapApnsToTerraformModel(apiApns []openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner) (types.List, diag.Diagnostics) {
	if len(apiApns) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: ResourceModelApnsAttrTypes()}), nil
	}

	var apns []attr.Value
	for _, apiApn := range apiApns {
		var auth types.Object
		if apiApn.HasAuthentication() {
			auth = mapAuthenticationToTerraformModel(apiApn.Authentication)
		} else {
			auth = types.ObjectNull(ResourceModelAuthenticationAttrTypes())
		}

		apn := types.ObjectValueMust(ResourceModelApnsAttrTypes(), map[string]attr.Value{
			"name":             types.StringValue(apiApn.GetName()),
			"allowed_ip_types": asSetOfStrings(apiApn.GetAllowedIpTypes()),
			"authentication":   auth,
		})
		apns = append(apns, apn)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: ResourceModelApnsAttrTypes()}, apns), nil
}

// mapSimFailOverToApiPayload converts Terraform SimFailOver object to API payload.
func mapSimFailOverToApiPayload(simFailOver types.Object) (*openApiClient.UpdateDeviceCellularSimsRequestSimFailover, diag.Diagnostics) {
	if simFailOver.IsNull() || simFailOver.IsUnknown() {
		return nil, nil
	}

	var failOverData ResourceModelSimFailOver
	options := basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	}
	diags := simFailOver.As(context.Background(), &failOverData, options)
	if diags.HasError() {
		return nil, diags
	}

	// Extract values and create pointers
	enabled := failOverData.Enabled.ValueBool()
	timeout := int32(failOverData.Timeout.ValueInt64())

	return &openApiClient.UpdateDeviceCellularSimsRequestSimFailover{
		Enabled: &enabled,
		Timeout: &timeout,
	}, nil
}

// mapAuthenticationToApiPayload converts Terraform Authentication object to API payload.
func mapAuthenticationToApiPayload(auth types.Object) *openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInnerAuthentication {
	if auth.IsNull() || auth.IsUnknown() {
		return nil
	}

	var authData ResourceModelAuthentication
	options := basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	}
	_ = auth.As(context.Background(), &authData, options)

	// Extract values and create pointers
	authType := authData.Type.ValueString()
	username := authData.Username.ValueString()
	password := authData.Password.ValueString()

	return &openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInnerAuthentication{
		Type:     &authType,
		Username: &username,
		Password: &password,
	}
}

// mapAuthenticationToTerraformModel converts API Authentication response to Terraform model.
func mapAuthenticationToTerraformModel(apiAuth *openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInnerAuthentication) types.Object {
	if apiAuth == nil {
		return types.ObjectNull(ResourceModelAuthenticationAttrTypes())
	}

	return types.ObjectValueMust(ResourceModelAuthenticationAttrTypes(), map[string]attr.Value{
		"type":     types.StringValue(apiAuth.GetType()),
		"username": types.StringValue(apiAuth.GetUsername()),
		"password": types.StringValue(apiAuth.GetPassword()),
	})
}

// Utility Functions

// asStringSlice converts a Terraform Set of Strings to a slice of strings.
func asStringSlice(set types.Set) []string {
	var result []string
	if !set.IsNull() && !set.IsUnknown() {
		_ = set.ElementsAs(context.Background(), &result, false)
	}
	return result
}

// asSetOfStrings converts a slice of strings to a Terraform Set of Strings.
func asSetOfStrings(strings []string) types.Set {
	elements := make([]attr.Value, len(strings))
	for i, s := range strings {
		elements[i] = types.StringValue(s)
	}
	return types.SetValueMust(types.StringType, elements)
}
