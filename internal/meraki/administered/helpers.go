package administered

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meraki/dashboard-api-go/client"
	"time"
)

// Helper Function for Safe Time Formatting
func safeFormatTime(t *time.Time, layout string) string {
	if t == nil {
		return ""
	}
	return t.Format(layout)
}

func mapTopLevelFields(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) (types.String, types.String, types.String) {
	tflog.Trace(ctx, "Mapping top-level fields: Name, Email, and LastUsedDashboardAt")

	name := types.StringNull()
	if inlineResp.HasName() {
		name = types.StringValue(inlineResp.GetName())
	}

	email := types.StringNull()
	if inlineResp.HasEmail() {
		email = types.StringValue(inlineResp.GetEmail())
	}

	lastUsed := types.StringNull()
	if inlineResp.HasLastUsedDashboardAt() {
		lastUsed = types.StringValue(safeFormatTime(inlineResp.LastUsedDashboardAt, time.RFC3339))
	}

	return name, email, lastUsed
}

func mapAuthentication(ctx context.Context, auth *client.GetAdministeredIdentitiesMe200ResponseAuthentication) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Starting to map Authentication fields")

	apiValue, apiDiags := mapApi(ctx, auth.Api)
	diags = append(diags, apiDiags...)

	samlValue, samlDiags := mapSaml(ctx, auth.Saml)
	diags = append(diags, samlDiags...)

	twoFactorValue, twoFactorDiags := mapTwoFactor(ctx, auth.TwoFactor)
	diags = append(diags, twoFactorDiags...)

	authAttrValues := map[string]attr.Value{
		"mode":         types.StringValue(auth.GetMode()),
		"api":          apiValue,
		"saml_enabled": samlValue,
		"two_factor":   twoFactorValue,
		"api_key_created": types.BoolValue(
			auth.HasApi() && auth.Api.HasKey() && auth.Api.Key.GetCreated(),
		),
	}

	authValue, authDiags := types.ObjectValue(authenticationAttrType, authAttrValues)
	diags = append(diags, authDiags...)

	tflog.Trace(ctx, "Successfully mapped Authentication fields")
	return authValue, diags
}

func mapApi(ctx context.Context, api *client.GetAdministeredIdentitiesMe200ResponseAuthenticationApi) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tflog.Trace(ctx, "Mapping API field")

	if api == nil || !api.HasKey() {
		tflog.Trace(ctx, "API Key field is not present, setting to null")
		return types.ObjectNull(apiAttrType), diags
	}

	keyValue, keyDiags := mapKey(ctx, api.Key)
	diags = append(diags, keyDiags...)

	apiAttrValue := map[string]attr.Value{
		"key": keyValue,
	}

	apiValue, apiDiags := types.ObjectValue(apiAttrType, apiAttrValue)
	diags = append(diags, apiDiags...)

	return apiValue, diags
}

func mapKey(ctx context.Context, key *client.GetAdministeredIdentitiesMe200ResponseAuthenticationApiKey) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "Mapping Key field in API")

	if key == nil {
		tflog.Trace(ctx, "Key field is not present, setting to null")
		return types.ObjectNull(keyAttrType), nil
	}

	keyAttrValue := map[string]attr.Value{
		"created": types.BoolValue(key.GetCreated()),
	}

	keyValue, diags := types.ObjectValue(keyAttrType, keyAttrValue)
	return keyValue, diags
}

func mapSaml(ctx context.Context, saml *client.GetAdministeredIdentitiesMe200ResponseAuthenticationSaml) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "Mapping SAML field")

	if saml == nil {
		tflog.Trace(ctx, "SAML field is not present, setting to null")
		return types.ObjectNull(samlAttrType), nil
	}

	samlAttrValue := map[string]attr.Value{
		"enabled": types.BoolValue(saml.GetEnabled()),
	}

	samlValue, diags := types.ObjectValue(samlAttrType, samlAttrValue)
	return samlValue, diags
}

func mapTwoFactor(ctx context.Context, twoFactor *client.GetAdministeredIdentitiesMe200ResponseAuthenticationTwoFactor) (basetypes.ObjectValue, diag.Diagnostics) {
	tflog.Trace(ctx, "Mapping TwoFactor field")

	if twoFactor == nil {
		tflog.Trace(ctx, "TwoFactor field is not present, setting to null")
		return types.ObjectNull(twoFactorAttrType), nil
	}

	twoFactorAttrValue := map[string]attr.Value{
		"enabled": types.BoolValue(twoFactor.GetEnabled()),
	}

	twoFactorValue, diags := types.ObjectValue(twoFactorAttrType, twoFactorAttrValue)
	return twoFactorValue, diags
}

func marshalState(ctx context.Context, inlineResp *client.GetAdministeredIdentitiesMe200Response) (identitiesMeAttrModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data identitiesMeAttrModel

	tflog.Trace(ctx, "Starting to marshal API response to Terraform model")

	// Map top-level fields
	data.Name, data.Email, data.LastUsedDashboardAt = mapTopLevelFields(ctx, inlineResp)

	// Map Authentication fields
	if inlineResp.HasAuthentication() {
		tflog.Trace(ctx, "Mapping Authentication fields")
		auth, authDiags := mapAuthentication(ctx, inlineResp.Authentication)
		diags = append(diags, authDiags...)
		if diags.HasError() {
			tflog.Trace(ctx, "Error encountered while mapping Authentication fields")
			return data, diags
		}
		data.Authentication = auth
	} else {
		tflog.Trace(ctx, "Authentication field is not present, setting to null")
		data.Authentication = types.ObjectNull(authenticationAttrType)
	}

	tflog.Trace(ctx, "Successfully marshaled API response to Terraform model")
	return data, diags
}
