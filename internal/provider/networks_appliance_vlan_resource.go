package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"strconv"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksApplianceVLANsResource{}
	_ resource.ResourceWithConfigure   = &NetworksApplianceVLANsResource{}
	_ resource.ResourceWithImportState = &NetworksApplianceVLANsResource{}
)

func NewNetworksApplianceVLANsResource() resource.Resource {
	return &NetworksApplianceVLANsResource{}
}

// NetworksApplianceVLANsResource defines the resource implementation.
type NetworksApplianceVLANsResource struct {
	client *openApiClient.APIClient
}

type NetworksApplianceVLANsResourceModel struct {
	Id                     jsontypes.String `tfsdk:"id" json:"-"`
	NetworkId              jsontypes.String `tfsdk:"network_id" json:"networkId"`
	VlanId                 jsontypes.Int64  `tfsdk:"vlan_id" json:"id"`
	InterfaceId            jsontypes.String `tfsdk:"interface_id" json:"interfaceId,omitempty"`
	Name                   jsontypes.String `tfsdk:"name" json:"name"`
	Subnet                 jsontypes.String `tfsdk:"subnet" json:"subnet"`
	ApplianceIp            jsontypes.String `tfsdk:"appliance_ip" json:"applianceIp"`
	GroupPolicyId          jsontypes.String `tfsdk:"group_policy_id" json:"groupPolicyId"`
	TemplateVlanType       jsontypes.String `tfsdk:"template_vlan_type" json:"templateVlanType"`
	Cidr                   jsontypes.String `tfsdk:"cidr" json:"cidr"`
	Mask                   jsontypes.Int64  `tfsdk:"mask" json:"mask"`
	DhcpRelayServerIps     types.List       `tfsdk:"dhcp_relay_server_ips" json:"dhcpRelayServerIps"`
	DhcpHandling           jsontypes.String `tfsdk:"dhcp_handling" json:"dhcpHandling"`
	DhcpLeaseTime          jsontypes.String `tfsdk:"dhcp_lease_time" json:"dhcpLeaseTime"`
	DhcpBootOptionsEnabled jsontypes.Bool   `tfsdk:"dhcp_boot_options_enabled" json:"dhcpBootOptionsEnabled"`
	DhcpBootNextServer     jsontypes.String `tfsdk:"dhcp_boot_next_server" json:"dhcpBootNextServer"`
	DhcpBootFilename       jsontypes.String `tfsdk:"dhcp_boot_filename" json:"dhcpBootFilename"`
	FixedIpAssignments     types.Map        `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges       types.List       `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
	DnsNameservers         jsontypes.String `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	DhcpOptions            types.List       `tfsdk:"dhcp_options" json:"dhcpOptions"`
	VpnNatSubnet           jsontypes.String `tfsdk:"vpn_nat_subnet" json:"vpnNatSubnet"`
	MandatoryDhcp          types.Object     `tfsdk:"mandatory_dhcp" json:"MandatoryDhcp"`
	IPv6                   types.Object     `tfsdk:"ipv6" json:"ipv6"`
}

type NetworksApplianceVLANsResourceModelIpNameMapping struct {
	Ip   jsontypes.String `tfsdk:"ip" json:"ip"`
	Name jsontypes.String `tfsdk:"name" json:"name"`
}

type NetworksApplianceVLANsResourceModelReservedIpRange struct {
	Start   jsontypes.String `tfsdk:"start" json:"start"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
}

type NetworksApplianceVLANsResourceModelDhcpOption struct {
	Code  jsontypes.String `tfsdk:"code" json:"code"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
	Value jsontypes.String `tfsdk:"value" json:"value"`
}

type NetworksApplianceVLANsResourceModelIpv6Configuration struct {
	Enabled           jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments types.List     `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

type NetworksApplianceVLANsResourceModelFixedIpAssignment struct {
	Ip   jsontypes.String `tfsdk:"ip"`
	Name jsontypes.String `tfsdk:"name"`
}

type NetworksApplianceVLANsResourceModelPrefixAssignment struct {
	Autonomous         jsontypes.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       jsontypes.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 jsontypes.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             types.Object     `tfsdk:"origin" json:"origin"`
}

type NetworksApplianceVLANsResourceModelOrigin struct {
	Type       jsontypes.String `tfsdk:"type" json:"type"`
	Interfaces types.List       `tfsdk:"interfaces" json:"interfaces"`
}

type NetworksApplianceVLANsResourceModelMandatoryDhcp struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

func CreatePayloadRequest(ctx context.Context, data *NetworksApplianceVLANsResourceModel) (*openApiClient.CreateNetworkApplianceVlanRequest, diag.Diagnostics) {

	resp := diag.Diagnostics{}

	// Initialize the payload
	payload := openApiClient.NewCreateNetworkApplianceVlanRequest(data.Id.ValueString(), data.Name.ValueString())

	// Id
	if !data.VlanId.IsUnknown() && !data.VlanId.IsNull() {

		// API returns this as string, openAPI spec has set as Integer
		payload.SetId(fmt.Sprintf("%v", data.VlanId.ValueInt64()))
	}

	// Name
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	// Subnet
	if !data.Subnet.IsUnknown() && !data.Subnet.IsNull() {
		payload.SetSubnet(data.Subnet.ValueString())
	}

	// ApplianceIp
	if !data.ApplianceIp.IsUnknown() && !data.ApplianceIp.IsNull() {
		payload.SetApplianceIp(data.ApplianceIp.ValueString())
	}

	// GroupPolicyId
	if !data.GroupPolicyId.IsUnknown() && !data.GroupPolicyId.IsNull() {
		payload.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	}

	// TemplateVlanType
	if !data.TemplateVlanType.IsUnknown() && !data.TemplateVlanType.IsNull() {
		payload.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	}

	// Cidr
	if !data.Cidr.IsUnknown() && !data.Cidr.IsNull() {
		payload.SetCidr(data.Cidr.ValueString())
	}

	// Mask
	if !data.Mask.IsUnknown() && !data.Mask.IsNull() {
		payload.SetMask(int32(data.Mask.ValueInt64()))
	}

	// IPV6
	if !data.IPv6.IsUnknown() && !data.IPv6.IsNull() {
		ipv6Payload := openApiClient.NewUpdateNetworkApplianceSingleLanRequestIpv6()

		var ipv6 NetworksApplianceVLANsResourceModelIpv6Configuration
		diags := data.IPv6.As(ctx, &ipv6, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		// Enabled
		ipv6Payload.SetEnabled(ipv6.Enabled.ValueBool())

		// Handle Prefix Assignments
		var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
		if !ipv6.PrefixAssignments.IsUnknown() && !ipv6.PrefixAssignments.IsNull() {
			// Create a variable to hold the converted map elements
			var prefixAssignmentMap map[string]NetworksApplianceVLANsResourceModelPrefixAssignment

			// Use ElementsAs to convert the elements
			if prefixAssignmentMapDiags := ipv6.PrefixAssignments.ElementsAs(ctx, &prefixAssignmentMap, false); diags.HasError() {
				resp.AddError(
					"Create Payload Failure", fmt.Sprintf("%v", prefixAssignmentMapDiags),
				)
			}

			for _, prefixAssignment := range prefixAssignmentMap {
				var originInterfaces []string

				// Extract the Origin object
				var origin NetworksApplianceVLANsResourceModelOrigin
				if diags = prefixAssignment.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{}); diags.HasError() {
					resp.AddError(
						"Create Payload Failure", fmt.Sprintf("%v", diags),
					)
				}

				// Assuming origin.Interfaces is a list of strings
				if !origin.Interfaces.IsUnknown() && !origin.Interfaces.IsNull() {
					var interfaceList []types.String
					if diags = origin.Interfaces.ElementsAs(ctx, &interfaceList, true); diags.HasError() {
						resp.AddError(
							"Create Payload Failure", fmt.Sprintf("%v", diags),
						)
					}

					for _, iface := range interfaceList {
						if !iface.IsUnknown() && !iface.IsNull() {
							originInterfaces = append(originInterfaces, iface.ValueString())
						}
					}
				}

				prefixAssignments = append(prefixAssignments, openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
					Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
					StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
					StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
					Origin: &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin{
						Type:       origin.Type.ValueString(),
						Interfaces: originInterfaces,
					},
				})
			}
		}

		ipv6Payload.SetPrefixAssignments(prefixAssignments)
		payload.SetIpv6(*ipv6Payload)
	}

	// MandatoryDhcp
	if !data.MandatoryDhcp.IsUnknown() && !data.MandatoryDhcp.IsNull() {
		mandatoryDhcpPayload := openApiClient.NewGetNetworkApplianceVlans200ResponseInnerMandatoryDhcp()
		var mandatoryDhcp NetworksApplianceVLANsResourceModelMandatoryDhcp

		diags := data.MandatoryDhcp.As(ctx, &mandatoryDhcp, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		mandatoryDhcpPayload.SetEnabled(mandatoryDhcp.Enabled.ValueBool())

		payload.SetMandatoryDhcp(*mandatoryDhcpPayload)
	}

	return payload, nil
}

func CreatePayloadResponse(ctx context.Context, data *NetworksApplianceVLANsResourceModel, response *openApiClient.CreateNetworkApplianceVlan201Response) diag.Diagnostics {

	resp := diag.Diagnostics{}

	if response.HasId() {

		// API returns vlanId as a string instead of an Integer
		vlanId, err := strconv.Atoi(response.GetId())
		if err != nil {
			resp.AddError("VlanId Response", fmt.Sprintf("\n%v", err))
		}
		data.VlanId = jsontypes.Int64Value(int64(vlanId))
		data.Id = jsontypes.StringValue(fmt.Sprintf("%v", data.VlanId.ValueInt64()))
	}

	if response.HasInterfaceId() {
		data.InterfaceId = jsontypes.StringValue(response.GetInterfaceId())
	}

	if response.HasName() {
		data.Name = jsontypes.StringValue(response.GetName())
	}

	if response.HasSubnet() {
		data.Subnet = jsontypes.StringValue(response.GetSubnet())
	}

	if response.HasApplianceIp() {
		data.ApplianceIp = jsontypes.StringValue(response.GetApplianceIp())
	}

	if response.HasGroupPolicyId() {
		data.GroupPolicyId = jsontypes.StringValue(response.GetGroupPolicyId())
	}

	if response.HasTemplateVlanType() {
		data.TemplateVlanType = jsontypes.StringValue(response.GetTemplateVlanType())
	}

	if response.HasCidr() {
		data.Cidr = jsontypes.StringValue(response.GetCidr())
	}

	if response.HasMask() {
		data.Mask = jsontypes.Int64Value(int64(response.GetMask()))
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANsResourceModelMandatoryDhcp{}

		mandatoryDhcp.Enabled = jsontypes.BoolValue(response.MandatoryDhcp.GetEnabled())

		mandatoryDhcpAttributes := map[string]attr.Type{
			"enabled": jsontypes.BoolType,
		}

		objectVal, diags := types.ObjectValueFrom(ctx, mandatoryDhcpAttributes, mandatoryDhcp)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.MandatoryDhcp = objectVal
	}

	// IPv6
	if response.HasIpv6() {
		ipv6Response := NetworksApplianceVLANsResourceModelIpv6Configuration{}

		ipv6Response.Enabled = jsontypes.BoolValue(response.Ipv6.GetEnabled())

		// PrefixAssignments
		for _, prefixAssignmentResponse := range response.Ipv6.PrefixAssignments {
			prefixAssignment := NetworksApplianceVLANsResourceModelPrefixAssignment{}

			if prefixAssignmentResponse.HasAutonomous() {
				prefixAssignment.Autonomous = jsontypes.BoolValue(prefixAssignmentResponse.GetAutonomous())
			}

			if prefixAssignmentResponse.HasStaticPrefix() {
				prefixAssignment.StaticPrefix = jsontypes.StringValue(prefixAssignmentResponse.GetStaticPrefix())
			}

			if prefixAssignmentResponse.HasStaticApplianceIp6() {
				prefixAssignment.StaticApplianceIp6 = jsontypes.StringValue(prefixAssignmentResponse.GetStaticApplianceIp6())
			}

			// Origins
			if prefixAssignmentResponse.HasOrigin() {
				originModel := NetworksApplianceVLANsResourceModelOrigin{}
				originModel.Type = jsontypes.StringValue(prefixAssignmentResponse.Origin.GetType())

				// Prepare a slice to hold interface values as tftypes.Values
				interfaceTfValues := make([]tftypes.Value, len(prefixAssignmentResponse.Origin.Interfaces))

				// Convert each interface string to a tftypes.Value
				for i, interfaceResponse := range prefixAssignmentResponse.Origin.Interfaces {
					interfaceTfValues[i] = tftypes.NewValue(tftypes.String, interfaceResponse)
				}

				// Define the type for the list of interface values
				interfaceListType := tftypes.List{ElementType: tftypes.String}

				// Create a tftypes.Value representing a list of interfaces
				interfacesTfList := tftypes.NewValue(interfaceListType, interfaceTfValues)

				// Convert the tftypes.Value list to a Terraform ListValue
				interfaceListValue := basetypes.ListValue{}
				if diags := interfaceListValue.ElementsAs(ctx, interfacesTfList, false); diags.HasError() {
					resp.Append(diags...)
				}

				originModel.Interfaces = interfaceListValue

				// Define the attribute types for the origin object
				originAttributes := map[string]attr.Type{
					"type":       jsontypes.StringType,
					"interfaces": types.ListType{ElemType: jsontypes.StringType},
				}

				// Convert the origin model to a Terraform ObjectValue
				originObjectValue, originDiags := types.ObjectValueFrom(ctx, originAttributes, originModel)
				if originDiags.HasError() {
					resp.Append(originDiags...)
				}

				prefixAssignment.Origin = originObjectValue
			}

			PrefixAssignmentsMap := basetypes.ListValue{}

			PrefixAssignmentsMap.ElementsAs(ctx, prefixAssignment, false)

			ipv6Response.PrefixAssignments = PrefixAssignmentsMap
		}

		ipv6Attributes := map[string]attr.Type{
			"enabled":            jsontypes.BoolType,
			"prefix_assignments": types.ObjectType{},
		}

		ipv6Diags := diag.Diagnostics{}
		data.IPv6, ipv6Diags = types.ObjectValueFrom(ctx, ipv6Attributes, ipv6Response)
		if ipv6Diags.HasError() {
			resp.Append(ipv6Diags...)
		}
	}

	return resp
}

func UpdatePayloadRequest(ctx context.Context, data *NetworksApplianceVLANsResourceModel) (*openApiClient.UpdateNetworkApplianceVlanRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	payload := openApiClient.NewUpdateNetworkApplianceVlanRequest()

	// Name
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	// Subnet
	if !data.Subnet.IsUnknown() && !data.Subnet.IsNull() {
		payload.SetSubnet(data.Subnet.ValueString())
	}

	// ApplianceIp
	if !data.ApplianceIp.IsUnknown() && !data.ApplianceIp.IsNull() {
		payload.SetApplianceIp(data.ApplianceIp.ValueString())
	}

	// GroupPolicyId
	if !data.GroupPolicyId.IsUnknown() && !data.GroupPolicyId.IsNull() {
		payload.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	}

	// VpnNatSubnet
	if !data.VpnNatSubnet.IsUnknown() && !data.VpnNatSubnet.IsNull() {
		payload.SetVpnNatSubnet(data.VpnNatSubnet.ValueString())
	}

	// DhcpHandling
	if !data.DhcpHandling.IsUnknown() && !data.DhcpHandling.IsNull() {
		payload.SetDhcpHandling(data.DhcpHandling.ValueString())
	}

	// DhcpRelayServerIps
	if !data.DhcpRelayServerIps.IsUnknown() && !data.DhcpRelayServerIps.IsNull() {
		var dhcpRelayServerIps []string

		for _, dhcpRelayServerIp := range data.DhcpRelayServerIps.Elements() {
			dhcpRelayServerIps = append(dhcpRelayServerIps, dhcpRelayServerIp.String())
		}

		payload.SetDhcpRelayServerIps(dhcpRelayServerIps)
	}

	// DhcpLeaseTime
	if !data.DhcpLeaseTime.IsUnknown() && !data.DhcpLeaseTime.IsNull() {
		payload.SetDhcpLeaseTime(data.DhcpLeaseTime.ValueString())
	}

	// DhcpBootOptionsEnabled
	if !data.DhcpBootOptionsEnabled.IsUnknown() && !data.DhcpBootOptionsEnabled.IsNull() {
		payload.SetDhcpBootOptionsEnabled(data.DhcpBootOptionsEnabled.ValueBool())
	}

	// DhcpBootNextServer
	if !data.DhcpBootNextServer.IsUnknown() && !data.DhcpBootNextServer.IsNull() {
		payload.SetDhcpBootNextServer(data.DhcpBootNextServer.ValueString())
	}

	// DhcpBootFilename
	if !data.DhcpBootFilename.IsUnknown() && !data.DhcpBootFilename.IsNull() {
		payload.SetDhcpBootFilename(data.DhcpBootFilename.ValueString())
	}

	// FixedIpAssignments
	if !data.FixedIpAssignments.IsUnknown() && !data.FixedIpAssignments.IsNull() {

		var fixedIpAssignments map[string]interface{}
		fixedIpAssignmentsDiags := data.FixedIpAssignments.ElementsAs(ctx, &fixedIpAssignments, false)
		if fixedIpAssignmentsDiags.HasError() {
			resp.AddError(
				"Create Payload Failure, FixedIpAssignments", fmt.Sprintf("%v", fixedIpAssignmentsDiags),
			)
		}

		payload.SetFixedIpAssignments(fixedIpAssignments)
	}

	// ReservedIpRanges
	if !data.ReservedIpRanges.IsUnknown() && !data.ReservedIpRanges.IsNull() {
		var reservedIpRanges []openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner
		var reservedIpRangesData []NetworksApplianceVLANsResourceModelReservedIpRange

		reservedIpRangesDiags := data.ReservedIpRanges.ElementsAs(ctx, reservedIpRangesData, false)
		if reservedIpRangesDiags.HasError() {
			resp.AddError(
				"Create Payload Failure, ReservedIpRanges", fmt.Sprintf("%v", reservedIpRangesDiags),
			)
		}

		for _, reservedIpRangeData := range reservedIpRangesData {
			var reservedIpRangePayload openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner

			reservedIpRangePayload.SetStart(reservedIpRangeData.Start.ValueString())
			reservedIpRangePayload.SetEnd(reservedIpRangeData.End.ValueString())
			reservedIpRangePayload.SetComment(reservedIpRangeData.Comment.ValueString())

			reservedIpRanges = append(reservedIpRanges, reservedIpRangePayload)
		}

		payload.SetReservedIpRanges(reservedIpRanges)
	}

	// DnsNameservers
	if !data.DnsNameservers.IsUnknown() && !data.DnsNameservers.IsNull() {
		payload.SetDnsNameservers(data.DnsNameservers.ValueString())
	}

	// DhcpOptions
	if !data.DhcpOptions.IsUnknown() && !data.DhcpOptions.IsNull() {
		var dhcpOptions []openApiClient.GetNetworkApplianceVlans200ResponseInnerDhcpOptionsInner
		var dhcpOptionsData []NetworksApplianceVLANsResourceModelDhcpOption

		dhcpOptionsDiags := data.DhcpOptions.ElementsAs(ctx, dhcpOptionsData, false)
		if dhcpOptionsDiags.HasError() {
			resp.AddError(
				"Create Payload Failure, DhcpOptions", fmt.Sprintf("%v", dhcpOptionsDiags),
			)
		}

		for _, dhcpOption := range dhcpOptionsData {

			var dhcpOptionPayload openApiClient.GetNetworkApplianceVlans200ResponseInnerDhcpOptionsInner

			dhcpOptionPayload.SetCode(dhcpOption.Code.ValueString())
			dhcpOptionPayload.SetType(dhcpOption.Type.ValueString())
			dhcpOptionPayload.SetValue(dhcpOption.Value.ValueString())

			dhcpOptions = append(dhcpOptions, dhcpOptionPayload)

		}

		payload.SetDhcpOptions(dhcpOptions)
	}

	// TemplateVlanType
	if !data.TemplateVlanType.IsUnknown() && !data.TemplateVlanType.IsNull() {
		payload.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	}

	// Cidr
	if !data.Cidr.IsUnknown() && !data.Cidr.IsNull() {
		payload.SetCidr(data.Cidr.ValueString())
	}

	// Mask
	if !data.Mask.IsUnknown() && !data.Mask.IsNull() {
		payload.SetMask(int32(data.Mask.ValueInt64()))
	}

	// IPV6
	if !data.IPv6.IsUnknown() && !data.IPv6.IsNull() {
		ipv6Payload := openApiClient.NewUpdateNetworkApplianceSingleLanRequestIpv6()

		var ipv6 NetworksApplianceVLANsResourceModelIpv6Configuration
		diags := data.IPv6.As(ctx, &ipv6, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		// Enabled
		ipv6Payload.SetEnabled(ipv6.Enabled.ValueBool())

		// Handle Prefix Assignments
		var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
		if !ipv6.PrefixAssignments.IsUnknown() && !ipv6.PrefixAssignments.IsNull() {
			// Create a variable to hold the converted map elements
			var prefixAssignmentMap map[string]NetworksApplianceVLANsResourceModelPrefixAssignment

			// Use ElementsAs to convert the elements
			if prefixAssignmentMapDiags := ipv6.PrefixAssignments.ElementsAs(ctx, &prefixAssignmentMap, false); diags.HasError() {
				resp.AddError(
					"Create Payload Failure", fmt.Sprintf("%v", prefixAssignmentMapDiags),
				)
			}

			for _, prefixAssignment := range prefixAssignmentMap {
				var originInterfaces []string

				// Extract the Origin object
				var origin NetworksApplianceVLANsResourceModelOrigin
				if diags = prefixAssignment.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{}); diags.HasError() {
					resp.AddError(
						"Create Payload Failure", fmt.Sprintf("%v", diags),
					)
				}

				// Assuming origin.Interfaces is a list of strings
				if !origin.Interfaces.IsUnknown() && !origin.Interfaces.IsNull() {
					var interfaceList []types.String
					if diags = origin.Interfaces.ElementsAs(ctx, &interfaceList, true); diags.HasError() {
						resp.AddError(
							"Create Payload Failure", fmt.Sprintf("%v", diags),
						)
					}

					for _, iface := range interfaceList {
						if !iface.IsUnknown() && !iface.IsNull() {
							originInterfaces = append(originInterfaces, iface.ValueString())
						}
					}
				}

				prefixAssignments = append(prefixAssignments, openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
					Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
					StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
					StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
					Origin: &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin{
						Type:       origin.Type.ValueString(),
						Interfaces: originInterfaces,
					},
				})
			}
		}

		ipv6Payload.SetPrefixAssignments(prefixAssignments)
		payload.SetIpv6(*ipv6Payload)
	}

	// MandatoryDhcp
	if !data.MandatoryDhcp.IsUnknown() && !data.MandatoryDhcp.IsNull() {
		mandatoryDhcpPayload := openApiClient.NewGetNetworkApplianceVlans200ResponseInnerMandatoryDhcp()
		var mandatoryDhcp NetworksApplianceVLANsResourceModelMandatoryDhcp

		diags := data.MandatoryDhcp.As(ctx, &mandatoryDhcp, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		mandatoryDhcpPayload.SetEnabled(mandatoryDhcp.Enabled.ValueBool())

		payload.SetMandatoryDhcp(*mandatoryDhcpPayload)
	}

	return payload, nil
}

func UpdatePayloadResponse(ctx context.Context, data *NetworksApplianceVLANsResourceModel, response *openApiClient.GetNetworkApplianceVlans200ResponseInner) diag.Diagnostics {
	resp := diag.Diagnostics{}

	if response.Id != nil {
		data.Id = jsontypes.StringValue(*response.Id)
	}
	if response.InterfaceId != nil {
		data.InterfaceId = jsontypes.StringValue(*response.InterfaceId)
	}
	if response.Name != nil {
		data.Name = jsontypes.StringValue(*response.Name)
	}
	if response.Subnet != nil {
		data.Subnet = jsontypes.StringValue(*response.Subnet)
	}
	if response.ApplianceIp != nil {
		data.ApplianceIp = jsontypes.StringValue(*response.ApplianceIp)
	}
	if response.GroupPolicyId != nil {
		data.GroupPolicyId = jsontypes.StringValue(*response.GroupPolicyId)
	}
	if response.TemplateVlanType != nil {
		data.TemplateVlanType = jsontypes.StringValue(*response.TemplateVlanType)
	}
	if response.Cidr != nil {
		data.Cidr = jsontypes.StringValue(*response.Cidr)
	}
	if response.Mask != nil {
		data.Mask = jsontypes.Int64Value(int64(*response.Mask))
	}
	if response.DhcpHandling != nil {
		data.DhcpHandling = jsontypes.StringValue(*response.DhcpHandling)
	}
	if response.DhcpLeaseTime != nil {
		data.DhcpLeaseTime = jsontypes.StringValue(*response.DhcpLeaseTime)
	}

	// DHCP Relay Server IPs
	if len(response.DhcpRelayServerIps) > 0 {
		var dhcpRelayServerIpsList types.List

		var dhcpRelayServerIpsElems []attr.Value
		for _, ip := range response.DhcpRelayServerIps {
			dhcpRelayServerIpsElems = append(dhcpRelayServerIpsElems, jsontypes.StringValue(ip))
		}

		dhcpRelayServerIpsList, _ = basetypes.NewListValue(types.StringType, dhcpRelayServerIpsElems)
		data.DhcpRelayServerIps = dhcpRelayServerIpsList

		if response.DhcpBootOptionsEnabled != nil {
			data.DhcpBootOptionsEnabled = jsontypes.BoolValue(*response.DhcpBootOptionsEnabled)
		}
		if response.DhcpBootNextServer != nil {
			data.DhcpBootNextServer = jsontypes.StringValue(*response.DhcpBootNextServer)
		}
		if response.DhcpBootFilename != nil {
			data.DhcpBootFilename = jsontypes.StringValue(*response.DhcpBootFilename)
		}
		if response.VpnNatSubnet != nil {
			data.VpnNatSubnet = jsontypes.StringValue(*response.VpnNatSubnet)
		}
	}

	// NetworksApplianceVLANsResourceModelFixedIpAssignment (map[string]interface)
	if len(response.FixedIpAssignments) > 0 {

		fixedIpAssignmentsMap := make(map[string]attr.Value)

		for mac, assignment := range response.FixedIpAssignments {
			assignmentMap, ok := assignment.(map[string]interface{})
			if !ok {
				resp.AddError(
					"Failed to render response for fixedIpAssignments",
					fmt.Sprintf("mac: %s, assignment:%v", mac, assignment),
				)
				continue
			}

			ip, ipOk := assignmentMap["ip"].(string)
			name, nameOk := assignmentMap["name"].(string)
			if !ipOk || !nameOk {
				resp.AddError(
					"Failed to render ip/name for fixedIpAssignments",
					fmt.Sprintf("ip: %s, name:%v", ip, name),
				)
				continue
			}

			// Create a NetworksApplianceVLANsResourceModelFixedIpAssignment instance
			fixedIpAssignmentData := map[string]attr.Type{}

			fixedIpAssignmentAttr := NetworksApplianceVLANsResourceModelFixedIpAssignment{
				Ip:   jsontypes.StringValue(ip),
				Name: jsontypes.StringValue(name),
			}

			// Construct the types.Object for NetworksApplianceVLANsResourceModelFixedIpAssignment
			fixedIpAssignmentObj, diags := types.ObjectValueFrom(ctx, fixedIpAssignmentData, fixedIpAssignmentAttr)
			if diags.HasError() {
				resp.AddError(
					"Failed to create object for fixedIpAssignments",
					fmt.Sprintf("%v", diags),
				)
				continue
			}

			fixedIpAssignmentsMap[mac] = fixedIpAssignmentObj
		}

		// Construct the final types.Object to hold the map of FixedIpAssignments
		fixedIpAssignmentsObject, diags := types.MapValue(types.ObjectType{}, fixedIpAssignmentsMap)
		if diags.HasError() {
			resp.AddError(
				"Failed to create map for fixedIpAssignments",
				fmt.Sprintf("%v", diags),
			)
		}

		data.FixedIpAssignments = fixedIpAssignmentsObject
	}

	// Reserved IP Ranges
	var reservedIpRangesList []attr.Value

	// Define the attribute types for NetworksApplianceVLANsResourceModelReservedIpRange
	rangeAttrTypes := map[string]attr.Type{
		"start":   types.StringType,
		"end":     types.StringType,
		"comment": types.StringType,
	}

	for _, rangeItem := range response.ReservedIpRanges {

		rangeMap := make(map[string]attr.Value)
		if rangeItem.Start != nil {
			rangeMap["start"] = jsontypes.StringValue(*rangeItem.Start)
		}
		if rangeItem.End != nil {
			rangeMap["end"] = jsontypes.StringValue(*rangeItem.End)
		}
		if rangeItem.Comment != nil {
			rangeMap["comment"] = jsontypes.StringValue(*rangeItem.Comment)
		}

		// Convert rangeMap to types.Object
		rangeObject, diags := types.ObjectValueFrom(ctx, rangeAttrTypes, rangeMap)
		if diags.HasError() {
			resp.AddError(
				"Failed to create object for reservedIpRanges",
				fmt.Sprintf("%v", diags),
			)
			continue
		}

		reservedIpRangesList = append(reservedIpRangesList, rangeObject)
	}

	// Define the ListType for the reserved IP ranges
	listType := types.ListType{ElemType: types.ObjectType{AttrTypes: rangeAttrTypes}}

	// Convert the slice of attr.Value to a ListValue
	listValue := basetypes.ListValue{}
	listValue.ElementsAs(ctx, reservedIpRangesList, false)

	// Construct the types.List to hold the collection of reserved IP ranges
	reservedIpRanges, reservedIpRangesDiags := listType.ValueFromList(ctx, listValue)
	if reservedIpRangesDiags.HasError() {
		resp.AddError(
			"Failed to create list for reservedIpRanges",
			fmt.Sprintf("%v", reservedIpRangesDiags),
		)
	}

	data.ReservedIpRanges, reservedIpRangesDiags = reservedIpRanges.ToListValue(ctx)
	if reservedIpRangesDiags.HasError() {
		resp.AddError(
			"Failed to create list for reservedIpRanges",
			fmt.Sprintf("%v", reservedIpRangesDiags),
		)
	}

	// DHCP Options
	dhcpOptionAttrTypes := map[string]attr.Type{
		"code":  types.StringType,
		"type":  types.StringType,
		"value": types.StringType,
	}

	var dhcpOptionsList []attr.Value

	for _, option := range response.DhcpOptions {
		dhcpOptionMap := make(map[string]attr.Value)

		if option.Code != "" {
			dhcpOptionMap["code"] = jsontypes.StringValue(option.Code)
		}
		if option.Type != "" {
			dhcpOptionMap["type"] = jsontypes.StringValue(option.Type)
		}
		if option.Value != "" {
			dhcpOptionMap["value"] = jsontypes.StringValue(option.Value)
		}

		// Convert dhcpOptionMap to types.Object
		dhcpOptionObject, diags := types.ObjectValueFrom(ctx, dhcpOptionAttrTypes, dhcpOptionMap)
		if diags.HasError() {
			resp.AddError(
				"Failed to create list for dhcpOptionObject",
				fmt.Sprintf("%v", reservedIpRangesDiags),
			)
			continue
		}

		dhcpOptionsList = append(dhcpOptionsList, dhcpOptionObject)
	}

	// Define the ListType for the DHCP options
	listType = types.ListType{ElemType: types.ObjectType{AttrTypes: dhcpOptionAttrTypes}}

	// Convert the slice of attr.Value to a ListValue
	listValue = basetypes.ListValue{}
	listValue.ElementsAs(ctx, dhcpOptionsList, false)

	// Construct the types.List to hold the collection of DHCP options
	dhcpOptions, dhcpOptionsDiags := listType.ValueFromList(ctx, listValue)
	if dhcpOptionsDiags.HasError() {
		// Handle errors
	}

	DhcpOptionsData, DhcpOptionsDiags := dhcpOptions.ToListValue(ctx)
	if DhcpOptionsDiags.HasError() {
		resp.AddError(
			"Failed to create list for DhcpOptions",
			fmt.Sprintf("%v", DhcpOptionsDiags),
		)
	}

	data.DhcpOptions = DhcpOptionsData

	return resp
}

func (r *NetworksApplianceVLANsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *NetworksApplianceVLANsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage the VLANs for an MX network",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"vlan_id": schema.Int64Attribute{
				Required:   true,
				CustomType: jsontypes.Int64Type,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"interface_id": schema.StringAttribute{
				MarkdownDescription: "The Interface ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The subnet of the VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"appliance_ip": schema.StringAttribute{
				MarkdownDescription: "The local IP of the appliance on the VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"group_policy_id": schema.StringAttribute{
				MarkdownDescription: " desired group policy to apply to the VLAN",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"vpn_nat_subnet": schema.StringAttribute{
				MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_handling": schema.StringAttribute{
				MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_relay_server_ips": schema.ListAttribute{
				ElementType: jsontypes.StringType,
				Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
				Optional:    true,
			},
			"dhcp_lease_time": schema.StringAttribute{
				MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_options_enabled": schema.BoolAttribute{
				MarkdownDescription: "Use DHCP boot options specified in other properties",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"dhcp_boot_next_server": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option to direct boot clients to the server to load the boot file from",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_filename": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option for boot filename ",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"fixed_ip_assignments": schema.MapNestedAttribute{
				Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"reserved_ip_ranges": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DHCP reserved IP ranges on the VLAN",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							MarkdownDescription: "The first IP in the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "The last IP in the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "A text comment for the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"dns_nameservers": schema.StringAttribute{
				MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_options": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The list of DHCP options that will be included in DHCP responses. Each object in the list should have \"code\", \"type\", and \"value\" properties.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							MarkdownDescription: "The code for the DHCP option. This should be an integer between 2 and 254.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type for the DHCP option. One of: 'text', 'ip', 'hex' or 'integer'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value for the DHCP option",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"template_vlan_type": schema.StringAttribute{
				MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"mask": schema.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"ipv6": schema.SingleNestedAttribute{
				Description: "IPv6 configuration on the VLAN",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"prefix_assignments": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Prefix assignments on the VLAN",
						NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"autonomous": schema.BoolAttribute{
								MarkdownDescription: "Auto assign a /64 prefix from the origin to the VLAN",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
							"static_prefix": schema.StringAttribute{
								MarkdownDescription: "Manual configuration of a /64 prefix on the VLAN",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"static_appliance_ip6": schema.StringAttribute{
								MarkdownDescription: "Manual configuration of the IPv6 Appliance IP",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"origin": schema.SingleNestedAttribute{
								MarkdownDescription: "The origin of the prefix",
								Optional:            true,
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "Type of the origin",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"interfaces": schema.ListAttribute{
										ElementType: jsontypes.StringType,
										Description: "Interfaces associated with the prefix",
										Optional:    true,
									},
								},
							},
						}},
					},
				},
			},
			"mandatory_dhcp": schema.SingleNestedAttribute{
				Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
		},
	}
}

func (r *NetworksApplianceVLANsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NetworksApplianceVLANsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadReqDiags := CreatePayloadRequest(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlanRequest(*payload).Execute()

	// Meraki API seems to return http status code 201 as an error.
	if err != nil && httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"HTTP Client Create Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	payloadRespDiags := CreatePayloadResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceVLANsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns string, OpenAPI defines integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	payloadRespDiags := UpdatePayloadResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceVLANsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadReqDiags := UpdatePayloadRequest(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns this as string, openAPI spec has set as Integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).UpdateNetworkApplianceVlanRequest(*payload).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	payloadRespDiags := UpdatePayloadResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceVLANsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns a string, OpenAPI spec defines an integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	httpResp, err := r.client.ApplianceApi.DeleteNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).Execute()
	if err != nil && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"HTTP Client Delete Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksApplianceVLANsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, admin_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vlan_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
