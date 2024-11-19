package administered

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meraki/dashboard-api-go/client"
	"time"
)

// mapAPIResponseToState converts the API response into the Terraform state model.
func mapAPIResponseToState(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) (dataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data dataSourceModel

	tflog.Trace(ctx, "[identities_me] Starting to marshal API response to Terraform model")

	// Map top-level fields
	// Assign a static ID
	data.Id = types.StringValue("identities_me")
	data.Name = mapName(ctx, inlineResp)
	data.Email = mapEmail(ctx, inlineResp)
	data.LastUsedDashboardAt = mapLastUsedDashboardAt(ctx, inlineResp)

	// Map Authentication fields
	data.Authentication, diags = mapAuthentication(ctx, inlineResp.Authentication, diags)

	tflog.Trace(ctx, "[identities_me] Successfully marshaled API response to Terraform model")
	return data, diags
}

// mapName maps the Name field from the API response to the Terraform state.
func mapName(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) types.String {
	tflog.Trace(ctx, "[identities_me] Mapping Name field")
	if inlineResp.HasName() {
		return types.StringValue(inlineResp.GetName())
	}
	return types.StringNull()
}

// mapEmail maps the Email field from the API response to the Terraform state.
func mapEmail(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) types.String {
	tflog.Trace(ctx, "[identities_me] Mapping Email field")
	if inlineResp.HasEmail() {
		return types.StringValue(inlineResp.GetEmail())
	}
	return types.StringNull()
}

// mapLastUsedDashboardAt maps and validates the LastUsedDashboardAt field from the API response.
func mapLastUsedDashboardAt(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) types.String {
	tflog.Trace(ctx, "[identities_me] Mapping LastUsedDashboardAt field")
	if !inlineResp.HasLastUsedDashboardAt() {
		return types.StringNull()
	}

	lastUsedFormatted := utils.SafeFormatRFC3339(ctx, inlineResp.LastUsedDashboardAt, time.RFC3339)
	if err := utils.ValidateRFC3339(ctx, lastUsedFormatted); err != nil {
		tflog.Error(ctx, "[identities_me] Validation failed for LastUsedDashboardAt", map[string]interface{}{
			"value": lastUsedFormatted,
			"error": err.Error(),
		})
		return types.StringNull()
	}

	tflog.Trace(ctx, "[identities_me] LastUsedDashboardAt is valid RFC3339 format", map[string]interface{}{
		"value": lastUsedFormatted,
	})
	return types.StringValue(lastUsedFormatted)
}

// mapAuthentication maps the Authentication object from the API response.
func mapAuthentication(ctx context.Context, auth *client.GetAdministeredIdentitiesMe200ResponseAuthentication, diags diag.Diagnostics) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Starting to map Authentication fields")

	if auth == nil {
		return types.ObjectNull(authenticationType), diags
	}

	apiValue, apiDiags := mapApi(ctx, auth.Api)
	diags = append(diags, apiDiags...)

	samlValue, samlDiags := mapSaml(ctx, auth.Saml)
	diags = append(diags, samlDiags...)

	twoFactorValue, twoFactorDiags := mapTwoFactor(ctx, auth.TwoFactor)
	diags = append(diags, twoFactorDiags...)

	authAttrValues := map[string]attr.Value{
		"mode":       types.StringValue(auth.GetMode()),
		"api":        apiValue,
		"saml":       samlValue,
		"two_factor": twoFactorValue,
	}

	authValue, authDiags := types.ObjectValue(authenticationType, authAttrValues)
	diags = append(diags, authDiags...)

	tflog.Trace(ctx, "[identities_me] Successfully mapped Authentication fields")
	return authValue, diags
}

// mapApi maps the API object from the Authentication section of the API response.
func mapApi(ctx context.Context, api *client.GetAdministeredIdentitiesMe200ResponseAuthenticationApi) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "[identities_me] Mapping API field")

	if api == nil || !api.HasKey() {
		tflog.Trace(ctx, "[identities_me] API Key field is not present, setting to null")
		return types.ObjectNull(apiType), diags
	}

	keyValue, keyDiags := mapKey(ctx, api.Key)
	diags = append(diags, keyDiags...)

	apiAttrValue := map[string]attr.Value{
		"key": keyValue,
	}

	apiValue, apiDiags := types.ObjectValue(apiType, apiAttrValue)
	diags = append(diags, apiDiags...)

	return apiValue, diags
}

// mapKey maps the Key object from the API section of the API response.
func mapKey(ctx context.Context, key *client.GetAdministeredIdentitiesMe200ResponseAuthenticationApiKey) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Mapping Key field in API")

	if key == nil {
		tflog.Trace(ctx, "[identities_me] Key field is not present, setting to null")
		return types.ObjectNull(keyType), nil
	}

	keyAttrValue := map[string]attr.Value{
		"created": types.BoolValue(key.GetCreated()),
	}

	keyValue, diags := types.ObjectValue(keyType, keyAttrValue)
	return keyValue, diags
}

// mapSaml maps the SAML object from the Authentication section of the API response.
func mapSaml(ctx context.Context, saml *client.GetAdministeredIdentitiesMe200ResponseAuthenticationSaml) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Mapping SAML field")

	if saml == nil {
		tflog.Trace(ctx, "[identities_me] SAML field is not present, setting to null")
		return types.ObjectNull(samlType), nil
	}

	samlAttrValue := map[string]attr.Value{
		"enabled": types.BoolValue(saml.GetEnabled()),
	}

	samlValue, diags := types.ObjectValue(samlType, samlAttrValue)
	return samlValue, diags
}

// mapTwoFactor maps the TwoFactor object from the Authentication section of the API response.
func mapTwoFactor(ctx context.Context, twoFactor *client.GetAdministeredIdentitiesMe200ResponseAuthenticationTwoFactor) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Mapping TwoFactor field")

	if twoFactor == nil {
		tflog.Trace(ctx, "[identities_me] TwoFactor field is not present, setting to null")
		return types.ObjectNull(twoFactorType), nil
	}

	twoFactorAttrValue := map[string]attr.Value{
		"enabled": types.BoolValue(twoFactor.GetEnabled()),
	}

	twoFactorValue, diags := types.ObjectValue(twoFactorType, twoFactorAttrValue)
	return twoFactorValue, diags
}
