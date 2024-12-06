package networksSettings

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// FromGetNetworkSettings200Response transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *resourceModel) FromGetNetworkSettings200Response(ctx context.Context, data *resourceModel, networkSettings200Response *openApiClient.GetNetworkSettings200Response) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksSettingsResourceModel FromGetNetworkSettings200Response")
	if networkSettings200Response == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API networkSettings200Response for IPv6")}
	}

	// Id
	m.Id = data.NetworkId

	// NetworkId
	m.NetworkId = data.NetworkId

	// LocalStatusPageEnabled
	if networkSettings200Response.HasLocalStatusPageEnabled() {
		m.LocalStatusPageEnabled = jsontypes.BoolValue(networkSettings200Response.GetLocalStatusPageEnabled())
	} else {
		m.LocalStatusPageEnabled = jsontypes.BoolNull()
	}

	// RemoteStatusPageEnabled
	if networkSettings200Response.HasRemoteStatusPageEnabled() {
		m.RemoteStatusPageEnabled = jsontypes.BoolValue(networkSettings200Response.GetRemoteStatusPageEnabled())
	} else {
		m.RemoteStatusPageEnabled = jsontypes.BoolNull()
	}

	// SecurePortEnabled
	if networkSettings200Response.SecurePort.HasEnabled() {
		m.SecurePortEnabled = jsontypes.BoolValue(networkSettings200Response.SecurePort.GetEnabled())
	} else {
		m.SecurePortEnabled = jsontypes.BoolNull()
	}

	// FipsEnabled
	if networkSettings200Response.Fips.GetEnabled() {
		m.FipsEnabled = jsontypes.BoolValue(networkSettings200Response.Fips.GetEnabled())
	} else {
		m.FipsEnabled = jsontypes.BoolValue(false)
	}

	// NamedVlans
	if networkSettings200Response.NamedVlans.GetEnabled() {
		m.NamedVlansEnabled = jsontypes.BoolValue(networkSettings200Response.NamedVlans.GetEnabled())
	} else {
		m.NamedVlansEnabled = jsontypes.BoolNull()
	}

	// LocalStatusPage
	var localStatusPage NetworksSettingsResourceModelLocalStatusPage
	localStatusPageDiags := localStatusPage.FromGetNetworkSettings200Response(ctx, data, networkSettings200Response)
	if localStatusPageDiags.HasError() {
		tflog.Error(ctx, "[fail] NetworksSettingsResourceModel localStatusPage Variable")
		return localStatusPageDiags
	}

	/*
		localStatusPage := NetworksSettingsResourceModelLocalStatusPage{
					Authentication: authenticationObject,
				}

				localStatusPageObject, localStatusPageObjectDiags := types.ObjectValueFrom(ctx, NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAttrTypes(), localStatusPage)
				if localStatusPageObjectDiags.HasError() {
					tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPage localStatusPage object")
					return localStatusPageObjectDiags
				}
	*/

	localStatusPageObject, localStatusPageObjectDiags := types.ObjectValueFrom(ctx, NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAttrTypes(), localStatusPage)
	if localStatusPageObjectDiags.HasError() {
		tflog.Error(ctx, "[fail] NetworksSettingsResourceModel localStatusPageObject")
		return localStatusPageObjectDiags
	}

	m.LocalStatusPage = localStatusPageObject

	tflog.Info(ctx, "[end] NetworksSettingsResourceModel FromGetNetworkSettings200Response")

	return nil
}

// FromGetNetworkSettings200Response transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *NetworksSettingsResourceModelLocalStatusPage) FromGetNetworkSettings200Response(ctx context.Context, data *resourceModel, networkSettings200Response *openApiClient.GetNetworkSettings200Response) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksSettingsResourceModelLocalStatusPage FromGetNetworkSettings200Response")

	if networkSettings200Response == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API networkSettings200Response for IPv6")}
	}

	if networkSettings200Response.HasLocalStatusPage() {

		var authentication NetworksSettingsResourceModelLocalStatusPageAuthentication
		authenticationDiags := authentication.FromGetNetworkSettings200Response(ctx, data, networkSettings200Response)
		if authenticationDiags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPage authentication variable")
			return authenticationDiags
		}

		authenticationObject, authenticationObjectDiags := types.ObjectValueFrom(ctx, NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes(), authentication)
		if authenticationObjectDiags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthentication authentication object")
			return authenticationObjectDiags
		}

		m.Authentication = authenticationObject

	} else {

		authenticationObject := types.ObjectNull(NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes())

		m.Authentication = authenticationObject

	}

	tflog.Info(ctx, "[end] NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthentication FromGetNetworkSettings200Response")

	return nil
}

// FromGetNetworkSettings200Response transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *NetworksSettingsResourceModelLocalStatusPageAuthentication) FromGetNetworkSettings200Response(ctx context.Context, data *resourceModel, networkSettings200Response *openApiClient.GetNetworkSettings200Response) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksSettingsResourceModelLocalStatusPageAuthentication FromGetNetworkSettings200Response")

	if networkSettings200Response == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API networkSettings200Response for IPv6")}
	}

	// Authentication
	if networkSettings200Response.LocalStatusPage.HasAuthentication() {

		if networkSettings200Response.LocalStatusPage.Authentication.HasEnabled() {

			// Enabled
			if networkSettings200Response.LocalStatusPage.Authentication.HasEnabled() {
				m.Enabled = jsontypes.BoolValue(networkSettings200Response.LocalStatusPage.Authentication.GetEnabled())
			}

			// Username
			if networkSettings200Response.LocalStatusPage.Authentication.HasUsername() {
				m.Username = jsontypes.StringValue(networkSettings200Response.LocalStatusPage.Authentication.GetUsername())
			}
		}

	} else {
		m.Enabled = jsontypes.BoolNull()
		m.Username = jsontypes.StringNull()
		m.Password = jsontypes.StringNull()
	}

	if !data.LocalStatusPage.IsNull() {
		// Check Terraform Plan for Password Value
		var LocalStatusPagePlanData NetworksSettingsResourceModelLocalStatusPage

		// Extract LocalStatusPage into Plan Data Variable
		diags := data.LocalStatusPage.As(ctx, &LocalStatusPagePlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPageAuthentication localStatusPage variable")
			return diags
		}

		// Further Extract Object Plan Data into LocalStatusPageAuthentication Variable
		var LocalStatusPageAuthenticationPlanData NetworksSettingsResourceModelLocalStatusPageAuthentication

		diags = LocalStatusPagePlanData.Authentication.As(ctx, &LocalStatusPageAuthenticationPlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPageAuthentication localStatusPageAuthentication variable")
			return diags
		}

		// Extract Password if in Plan
		if m.Password.IsNull() {
			if LocalStatusPageAuthenticationPlanData.Password.ValueString() != "" {
				m.Password = jsontypes.StringValue(LocalStatusPageAuthenticationPlanData.Password.ValueString())
			}
		}

		// Extract Username if in Plan
		if m.Username.IsNull() {

			if !LocalStatusPageAuthenticationPlanData.Username.IsUnknown() {
				m.Username = LocalStatusPageAuthenticationPlanData.Username
			}
		}
	}

	tflog.Info(ctx, "[end] NetworksSettingsResourceModelLocalStatusPageAuthentication FromGetNetworkSettings200Response")

	return nil
}

func createUpdateHttpReqPayload(ctx context.Context, data *resourceModel) (openApiClient.UpdateNetworkSettingsRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] createUpdateHttpReqPayload Function Call")

	payload := *openApiClient.NewUpdateNetworkSettingsRequest()

	tflog.Info(ctx, "[start] LocalStatusPageEnabled")
	if !data.LocalStatusPageEnabled.IsUnknown() && !data.LocalStatusPageEnabled.IsNull() {

		payload.SetLocalStatusPageEnabled(data.LocalStatusPageEnabled.ValueBool())

	}

	tflog.Info(ctx, "[start] RemoteStatusPageEnabled")
	if !data.RemoteStatusPageEnabled.IsUnknown() && !data.RemoteStatusPageEnabled.IsNull() {
		payload.SetRemoteStatusPageEnabled(data.RemoteStatusPageEnabled.ValueBool())
	}

	tflog.Info(ctx, "[start] LocalStatusPage")
	if !data.LocalStatusPage.IsUnknown() && !data.LocalStatusPage.IsNull() {

		var localStatusPage openApiClient.UpdateNetworkSettingsRequestLocalStatusPage

		var localStatusPagePlanData NetworksSettingsResourceModelLocalStatusPage
		diags := data.LocalStatusPage.As(ctx, &localStatusPagePlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"LocalStatusPagePlanDataReq",
				fmt.Sprintf("%s", diags.Errors()),
			)
			return payload, resp
		}

		var authenticationPlanData NetworksSettingsResourceModelLocalStatusPageAuthentication
		diags = localStatusPagePlanData.Authentication.As(ctx, &authenticationPlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"LocalStatusPageAuthenticationPlanDataReq",
				fmt.Sprintf("%s", diags.Errors()),
			)
			return payload, resp
		}

		if localStatusPage.HasAuthentication() {
			var authentication openApiClient.UpdateNetworkSettingsRequestLocalStatusPageAuthentication

			authentication.SetEnabled(authenticationPlanData.Enabled.ValueBool())
			authentication.SetPassword(authenticationPlanData.Password.ValueString())
			localStatusPage.SetAuthentication(authentication)
		}

		payload.SetLocalStatusPage(localStatusPage)

	}

	tflog.Info(ctx, "[end] LocalStatusPage")

	if !data.SecurePortEnabled.IsUnknown() && !data.SecurePortEnabled.IsNull() {
		var securePort openApiClient.GetNetworkSettings200ResponseSecurePort
		securePort.SetEnabled(data.SecurePortEnabled.ValueBool())
		payload.SetSecurePort(securePort)
	}

	//NamedVlans
	if !data.NamedVlansEnabled.IsUnknown() && !data.NamedVlansEnabled.IsNull() {
		var namedVlans openApiClient.UpdateNetworkSettingsRequestNamedVlans
		namedVlans.SetEnabled(data.NamedVlansEnabled.ValueBool())
		payload.SetNamedVlans(namedVlans)
	}

	tflog.Info(ctx, "[end] createUpdateHttpReqPayload Function Call")

	return payload, nil
}

/*
type NetworksSettingsResourceModelClientPrivacy struct {
	ExpireDataOlderThan jsontypes.Int64 `tfsdk:"expireDataOlderThan"`
	ExpireDataBefore    string          `tfsdk:"expireDataBefore"`
}
*/

type NetworksSettingsResourceModelLocalStatusPageAuthentication struct {
	Enabled  jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Username jsontypes.String `tfsdk:"username" json:"username"`
	Password jsontypes.String `tfsdk:"password" json:"password"`
}

func NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":  jsontypes.BoolType,
		"username": jsontypes.StringType,
		"password": jsontypes.StringType,
	}
}

type NetworksSettingsResourceModelLocalStatusPage struct {
	Authentication types.Object `tfsdk:"authentication" json:"authentication"`
}

func NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"authentication": types.ObjectType{AttrTypes: NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes()},
	}

}
