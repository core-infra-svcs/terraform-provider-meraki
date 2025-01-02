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

// MarshalIdentitiesMeForRead maps the API response into the Terraform state model for the identities_me resource.
func MarshalIdentitiesMeForRead(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) (DataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data DataSourceModel

	tflog.Trace(ctx, "[identities_me] Starting to marshal API response to Terraform state")

	// Assign static ID
	data.Id = types.StringValue("identities_me")
	data.Name = MarshalNameForRead(ctx, inlineResp)
	data.Email = MarshalEmailForRead(ctx, inlineResp)
	data.LastUsedDashboardAt = MarshalLastUsedDashboardAtForRead(ctx, inlineResp)

	// Map Authentication fields
	data.Authentication, diags = MarshalAuthenticationForRead(ctx, inlineResp.Authentication, diags)

	tflog.Trace(ctx, "[identities_me] Successfully marshaled API response to Terraform state")
	return data, diags
}

// MarshalNameForRead maps the Name field from the API response to the Terraform state.
func MarshalNameForRead(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) types.String {
	tflog.Trace(ctx, "[identities_me] Mapping Name field")
	if inlineResp.HasName() {
		return types.StringValue(inlineResp.GetName())
	}
	return types.StringNull()
}

// MarshalEmailForRead maps the Email field from the API response to the Terraform state.
func MarshalEmailForRead(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) types.String {
	tflog.Trace(ctx, "[identities_me] Mapping Email field")
	if inlineResp.HasEmail() {
		return types.StringValue(inlineResp.GetEmail())
	}
	return types.StringNull()
}

// MarshalLastUsedDashboardAtForRead maps and validates the LastUsedDashboardAt field from the API response.
func MarshalLastUsedDashboardAtForRead(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) types.String {
	tflog.Trace(ctx, "[identities_me] Mapping LastUsedDashboardAt field")
	if !inlineResp.HasLastUsedDashboardAt() {
		return types.StringNull()
	}

	lastUsedFormatted := utils.SafeFormatRFC3339(ctx, inlineResp.LastUsedDashboardAt, time.RFC3339)
	if err := utils.ValidateRFC3339(lastUsedFormatted); err != nil {
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

// MarshalAuthenticationForRead maps the Authentication object from the API response to the Terraform state.
func MarshalAuthenticationForRead(ctx context.Context, auth *client.GetAdministeredIdentitiesMe200ResponseAuthentication, diags diag.Diagnostics) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Starting to map Authentication fields")

	if auth == nil {
		return types.ObjectNull(AuthenticationType), diags
	}

	apiValue, apiDiags := MarshalAuthenticationAPIForRead(ctx, auth.Api)
	diags = append(diags, apiDiags...)

	samlValue, samlDiags := MarshalAuthenticationSAMLForRead(ctx, auth.Saml)
	diags = append(diags, samlDiags...)

	twoFactorValue, twoFactorDiags := MarshalAuthenticationTwoFactorForRead(ctx, auth.TwoFactor)
	diags = append(diags, twoFactorDiags...)

	authAttrValues := map[string]attr.Value{
		"mode":       types.StringValue(auth.GetMode()),
		"api":        apiValue,
		"saml":       samlValue,
		"two_factor": twoFactorValue,
	}

	authValue, authDiags := types.ObjectValue(AuthenticationType, authAttrValues)
	diags = append(diags, authDiags...)

	tflog.Trace(ctx, "[identities_me] Successfully mapped Authentication fields")
	return authValue, diags
}

// MarshalAuthenticationAPIForRead maps the API object from the Authentication section of the API response to Terraform state.
func MarshalAuthenticationAPIForRead(ctx context.Context, api *client.GetAdministeredIdentitiesMe200ResponseAuthenticationApi) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "[identities_me] Mapping API field")

	if api == nil || !api.HasKey() {
		return types.ObjectNull(APIType), diags
	}

	keyValue, keyDiags := MarshalAuthenticationAPIKeyForRead(ctx, api.Key)
	diags = append(diags, keyDiags...)

	apiAttrValue := map[string]attr.Value{
		"key": keyValue,
	}

	apiValue, apiDiags := types.ObjectValue(APIType, apiAttrValue)
	diags = append(diags, apiDiags...)

	return apiValue, diags
}

// MarshalAuthenticationAPIKeyForRead maps the API Key object to Terraform state.
func MarshalAuthenticationAPIKeyForRead(ctx context.Context, key *client.GetAdministeredIdentitiesMe200ResponseAuthenticationApiKey) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Mapping Key field in API")

	if key == nil {
		return types.ObjectNull(KeyType), nil
	}

	keyAttrValue := map[string]attr.Value{
		"created": types.BoolValue(key.GetCreated()),
	}

	keyValue, diags := types.ObjectValue(KeyType, keyAttrValue)
	return keyValue, diags
}

// MarshalAuthenticationSAMLForRead maps the SAML object from the Authentication section of the API response to Terraform state.
func MarshalAuthenticationSAMLForRead(ctx context.Context, saml *client.GetAdministeredIdentitiesMe200ResponseAuthenticationSaml) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Mapping SAML field")

	if saml == nil {
		return types.ObjectNull(SAMLType), nil
	}

	samlAttrValue := map[string]attr.Value{
		"enabled": types.BoolValue(saml.GetEnabled()),
	}

	samlValue, diags := types.ObjectValue(SAMLType, samlAttrValue)
	return samlValue, diags
}

// MarshalAuthenticationTwoFactorForRead maps the TwoFactor object from the Authentication section of the API response to Terraform state.
func MarshalAuthenticationTwoFactorForRead(ctx context.Context, twoFactor *client.GetAdministeredIdentitiesMe200ResponseAuthenticationTwoFactor) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "[identities_me] Mapping TwoFactor field")

	if twoFactor == nil {
		return types.ObjectNull(TwoFactorType), nil
	}

	twoFactorAttrValue := map[string]attr.Value{
		"enabled": types.BoolValue(twoFactor.GetEnabled()),
	}

	twoFactorValue, diags := types.ObjectValue(TwoFactorType, twoFactorAttrValue)
	return twoFactorValue, diags
}
