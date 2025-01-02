package vpn

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

func NetworkApplianceVpnSiteToSiteVpnResourcePayload(ctx context.Context, data *resourceModel) (*openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequest, diag.Diagnostics) {

	payload := *openApiClient.NewUpdateNetworkApplianceVpnSiteToSiteVpnRequest(data.Mode.ValueString())

	// For mode value "none" hubs and subnets should be empty
	if data.Mode.ValueString() != "none" {

		// Hubs
		var hubs []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner

		if !data.Hubs.IsUnknown() && !data.Hubs.IsNull() {
			var hubsPayload []resourceModelHubs
			diags := data.Hubs.ElementsAs(ctx, &hubsPayload, false)
			if diags.HasError() {
				return nil, diags
			}

			for _, hubValue := range hubsPayload {
				var hubData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner

				hubData.SetHubId(hubValue.HubId.ValueString())
				hubData.SetUseDefaultRoute(hubValue.UseDefaultRoute.ValueBool())
				hubs = append(hubs, hubData)
			}

			payload.SetHubs(hubs)
		} else {
			payload.SetHubs(nil)
		}

		if !data.Subnets.IsUnknown() && !data.Subnets.IsNull() {
			// Subnets
			var subnets []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
			var subnetsPayload []resourceModelSubnets
			diags := data.Subnets.ElementsAs(ctx, &subnetsPayload, false)
			if diags.HasError() {
				return nil, diags
			}

			for _, subnetValue := range subnetsPayload {
				var subnetData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
				subnetData.SetLocalSubnet(subnetValue.LocalSubnet.ValueString())
				subnetData.SetUseVpn(subnetValue.UseVpn.ValueBool())
				subnets = append(subnets, subnetData)
			}

			payload.SetSubnets(subnets)
		}
	} else {
		payload.SetSubnets(nil)
	}

	data.Id = data.NetworkId

	return &payload, nil

}

func NetworksApplianceVpnSiteToSiteVpnResourceResponse(ctx context.Context, response *openApiClient.GetNetworkApplianceVpnSiteToSiteVpn200Response, data *resourceModel) (*resourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	data.Mode = jsontypes.StringValue(response.GetMode())

	// Hubs
	var hubs []resourceModelHubs
	for _, element := range response.GetHubs() {
		var hub resourceModelHubs
		hub.UseDefaultRoute = jsontypes.BoolValue(element.GetUseDefaultRoute())
		hub.HubId = jsontypes.StringValue(element.GetHubId())
		hubs = append(hubs, hub)

	}

	hubAttributes := map[string]attr.Type{
		"use_default_route": types.BoolType,
		"hub_id":            types.StringType,
	}

	hubSchema := types.ObjectType{
		AttrTypes: hubAttributes,
	}

	data.Hubs, diags = types.ListValueFrom(ctx, hubSchema, hubs)
	if diags.HasError() {
		return data, diags
	}

	// Subnets
	var subnets []resourceModelSubnets
	for _, element := range response.GetSubnets() {
		var subnet resourceModelSubnets
		subnet.UseVpn = jsontypes.BoolValue(element.GetUseVpn())
		subnet.LocalSubnet = jsontypes.StringValue(element.GetLocalSubnet())

		subnets = append(subnets, subnet)

	}

	subnetAttributes := map[string]attr.Type{
		"use_vpn":      types.BoolType,
		"local_subnet": types.StringType,
	}

	subnetSchema := types.ObjectType{
		AttrTypes: subnetAttributes,
	}

	data.Subnets, diags = types.ListValueFrom(ctx, subnetSchema, subnets)
	if diags.HasError() {
		return data, diags
	}
	return data, nil
}
