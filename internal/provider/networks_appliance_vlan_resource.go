package provider

import (
	"context"
	"fmt"
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
	Id                     types.String `tfsdk:"id" json:"-"`
	NetworkId              types.String `tfsdk:"network_id" json:"networkId"`
	VlanId                 types.Int64  `tfsdk:"vlan_id" json:"id"`
	InterfaceId            types.String `tfsdk:"interface_id" json:"interfaceId,omitempty"`
	Name                   types.String `tfsdk:"name" json:"name"`
	Subnet                 types.String `tfsdk:"subnet" json:"subnet"`
	ApplianceIp            types.String `tfsdk:"appliance_ip" json:"applianceIp"`
	GroupPolicyId          types.String `tfsdk:"group_policy_id" json:"groupPolicyId"`
	TemplateVlanType       types.String `tfsdk:"template_vlan_type" json:"templateVlanType"`
	Cidr                   types.String `tfsdk:"cidr" json:"cidr"`
	Mask                   types.Int64  `tfsdk:"mask" json:"mask"`
	DhcpRelayServerIps     types.List   `tfsdk:"dhcp_relay_server_ips" json:"dhcpRelayServerIps"`
	DhcpHandling           types.String `tfsdk:"dhcp_handling" json:"dhcpHandling"`
	DhcpLeaseTime          types.String `tfsdk:"dhcp_lease_time" json:"dhcpLeaseTime"`
	DhcpBootOptionsEnabled types.Bool   `tfsdk:"dhcp_boot_options_enabled" json:"dhcpBootOptionsEnabled"`
	DhcpBootNextServer     types.String `tfsdk:"dhcp_boot_next_server" json:"dhcpBootNextServer"`
	DhcpBootFilename       types.String `tfsdk:"dhcp_boot_filename" json:"dhcpBootFilename"`
	FixedIpAssignments     types.Map    `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges       types.List   `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
	DnsNameservers         types.String `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	DhcpOptions            types.List   `tfsdk:"dhcp_options" json:"dhcpOptions"`
	VpnNatSubnet           types.String `tfsdk:"vpn_nat_subnet" json:"vpnNatSubnet"`
	MandatoryDhcp          types.Object `tfsdk:"mandatory_dhcp" json:"MandatoryDhcp"`
	IPv6                   types.Object `tfsdk:"ipv6" json:"ipv6"`
}

type NetworksApplianceVLANsResourceModelIpNameMapping struct {
	Ip   types.String `tfsdk:"ip" json:"ip"`
	Name types.String `tfsdk:"name" json:"name"`
}

type NetworksApplianceVLANsResourceModelReservedIpRange struct {
	Start   types.String `tfsdk:"start" json:"start"`
	End     types.String `tfsdk:"end" json:"end"`
	Comment types.String `tfsdk:"comment" json:"comment"`
}

type NetworksApplianceVLANsResourceModelDhcpOption struct {
	Code  types.String `tfsdk:"code" json:"code"`
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

type NetworksApplianceVLANsResourceModelFixedIpAssignment struct {
	Ip   types.String `tfsdk:"ip"`
	Name types.String `tfsdk:"name"`
}

type NetworksApplianceVLANsResourceModelIpv6 struct {
	Enabled           types.Bool `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments types.List `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

type NetworksApplianceVLANsResourceModelIpv6PrefixAssignment struct {
	Autonomous         types.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       types.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 types.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             types.Object `tfsdk:"origin" json:"origin"`
}

type NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin struct {
	Type       types.String `tfsdk:"type" json:"type"`
	Interfaces types.List   `tfsdk:"interfaces" json:"interfaces"`
}

type NetworksApplianceVLANsResourceModelMandatoryDhcp struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
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

		var ipv6 NetworksApplianceVLANsResourceModelIpv6
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
			var prefixAssignmentMap map[string]NetworksApplianceVLANsResourceModelIpv6PrefixAssignment

			// Use ElementsAs to convert the elements
			if prefixAssignmentMapDiags := ipv6.PrefixAssignments.ElementsAs(ctx, &prefixAssignmentMap, false); diags.HasError() {
				resp.AddError(
					"Create Payload Failure", fmt.Sprintf("%v", prefixAssignmentMapDiags),
				)
			}

			for _, prefixAssignment := range prefixAssignmentMap {
				var originInterfaces []string

				// Extract the Origin object
				var origin NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
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

	// Id
	if data.Id.IsUnknown() {
		data.Id = types.StringValue("example-id")
	}

	// VlanId
	if data.VlanId.IsUnknown() {
		// check api response for vlanId
		if response.HasId() {

			// API returns string, openAPI spec defines int
			idStr := response.GetId()

			// Check if the string is empty
			if idStr == "" {

				// TODO Handle the case where the ID string is empty (fail or warn?)
				resp.AddWarning("CreatePayloadResponse VlanId Error", "Received empty VlanId from response")
				data.VlanId = types.Int64Null()

			} else {
				// Convert string to int
				vlanId, err := strconv.Atoi(idStr)
				if err != nil {
					// Handle the error if conversion fails
					resp.AddError("CreatePayloadResponse VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s': %v", idStr, err))
				}

				// set new vlanId
				data.VlanId = types.Int64Value(int64(vlanId))
			}
		} else {
			data.VlanId = types.Int64Null()
		}
	}

	// InterfaceId
	if response.HasInterfaceId() {
		data.InterfaceId = types.StringValue(response.GetInterfaceId())
	} else {
		data.InterfaceId = types.StringNull()
	}

	// Name
	if response.HasName() {
		data.Name = types.StringValue(response.GetName())
	} else {
		data.Name = types.StringNull()
	}

	// Subnet
	if response.HasSubnet() {
		data.Subnet = types.StringValue(response.GetSubnet())
	} else {
		data.Subnet = types.StringNull()
	}

	// ApplianceIp
	if response.HasApplianceIp() {
		data.ApplianceIp = types.StringValue(response.GetApplianceIp())
	} else {
		data.ApplianceIp = types.StringNull()
	}

	// GroupPolicyId
	if response.HasGroupPolicyId() {
		data.GroupPolicyId = types.StringValue(response.GetGroupPolicyId())
	} else {
		data.GroupPolicyId = types.StringNull()
	}

	// TemplateVlanType
	if response.HasTemplateVlanType() {
		data.TemplateVlanType = types.StringValue(response.GetTemplateVlanType())
	} else {
		data.TemplateVlanType = types.StringNull()
	}

	// Cidr
	if response.HasCidr() {
		data.Cidr = types.StringValue(response.GetCidr())
	} else {
		data.Cidr = types.StringNull()
	}

	// Mask
	if response.HasMask() {
		data.Mask = types.Int64Value(int64(response.GetMask()))
	} else {
		data.Mask = types.Int64Null()
	}

	// DhcpRelayServerIps
	if data.DhcpRelayServerIps.IsUnknown() {
		data.DhcpRelayServerIps = basetypes.NewListNull(types.StringType)
	}

	// DhcpHandling
	if data.DhcpHandling.IsUnknown() {
		data.DhcpHandling = types.StringNull()
	}

	// DhcpLeaseTime
	if data.DhcpLeaseTime.IsUnknown() {
		data.DhcpLeaseTime = types.StringNull()
	}

	// DhcpBootOptionsEnabled
	if data.DhcpBootOptionsEnabled.IsUnknown() {
		data.DhcpBootOptionsEnabled = types.BoolNull()
	}

	// DhcpBootNextServer
	if data.DhcpBootNextServer.IsUnknown() {
		data.DhcpBootNextServer = types.StringNull()
	}

	// DhcpBootFilename
	if data.DhcpBootFilename.IsUnknown() {
		data.DhcpBootFilename = types.StringNull()
	}

	// FixedIpAssignments
	if data.FixedIpAssignments.IsUnknown() {

		var fixedIpAssignments NetworksApplianceVLANsResourceModelFixedIpAssignment

		profileObjectValue, fixedIpAssignmentsDiags := types.MapValueFrom(ctx, types.ObjectType{}, fixedIpAssignments)
		if fixedIpAssignmentsDiags.HasError() {
			resp.Append(fixedIpAssignmentsDiags...)
		}

		data.FixedIpAssignments = profileObjectValue

	}

	// ReservedIpRanges
	if data.ReservedIpRanges.IsUnknown() {

		data.ReservedIpRanges = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"comment": types.StringType,
				"end":     types.StringType,
				"start":   types.StringType,
			},
		})
	}

	// DnsNameservers
	if data.DnsNameservers.IsUnknown() {
		data.DnsNameservers = types.StringNull()
	}

	// DhcpOptions
	if data.DhcpOptions.IsUnknown() {
		data.DhcpOptions = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"code":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
		})
	}

	// VpnNatSubnet
	if data.VpnNatSubnet.IsUnknown() {
		data.VpnNatSubnet = types.StringNull()
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANsResourceModelMandatoryDhcp{}

		// Enabled
		if response.MandatoryDhcp.HasEnabled() {
			mandatoryDhcp.Enabled = types.BoolValue(response.MandatoryDhcp.GetEnabled())
		}

		mandatoryDhcpAttributes := map[string]attr.Type{
			"enabled": types.BoolType,
		}

		objectVal, diags := types.ObjectValueFrom(ctx, mandatoryDhcpAttributes, mandatoryDhcp)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.MandatoryDhcp = objectVal
	} else {
		data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
			"enabled": types.BoolType,
		})
	}

	// IPv6
	if response.HasIpv6() {
		ipv6Response := NetworksApplianceVLANsResourceModelIpv6{}

		// Set the 'enabled' attribute
		if response.Ipv6.HasEnabled() {
			ipv6Response.Enabled = types.BoolValue(response.Ipv6.GetEnabled())
		}

		// Define the attribute types for the 'origin' object in each prefix assignment
		originTypes := map[string]attr.Type{
			"type":       types.StringType,
			"interfaces": types.ListType{ElemType: types.StringType},
		}

		// Define the attribute types for each prefix assignment
		prefixAssignmentAttrTypes := map[string]attr.Type{
			"autonomous":           types.BoolType,
			"static_prefix":        types.StringType,
			"static_appliance_ip6": types.StringType,
			"origin":               types.ObjectType{AttrTypes: originTypes},
		}

		// Handling PrefixAssignments
		if response.Ipv6.HasPrefixAssignments() {
			var prefixAssignmentsValues []attr.Value

			for _, prefixAssignmentResponse := range response.Ipv6.PrefixAssignments {
				// Initialize originValue
				var originValue attr.Value

				// Handling the 'origin' object if it exists
				if prefixAssignmentResponse.HasOrigin() {
					var interfacesListValues []attr.Value
					for _, interfaceValue := range prefixAssignmentResponse.Origin.GetInterfaces() {
						interfacesListValues = append(interfacesListValues, types.StringValue(interfaceValue))
					}

					// Construct the interfaces list
					interfacesList, _ := basetypes.NewListValue(types.StringType, interfacesListValues)

					// Create the 'origin' object
					originObject := map[string]attr.Value{
						"type":       types.StringValue(prefixAssignmentResponse.Origin.GetType()),
						"interfaces": interfacesList,
					}
					originValue, _ = types.ObjectValueFrom(ctx, originTypes, originObject)
				}

				// Construct the prefix assignment object
				prefixAssignmentObject := map[string]attr.Value{
					"autonomous":           types.BoolValue(prefixAssignmentResponse.GetAutonomous()),
					"static_prefix":        types.StringValue(prefixAssignmentResponse.GetStaticPrefix()),
					"static_appliance_ip6": types.StringValue(prefixAssignmentResponse.GetStaticApplianceIp6()),
					"origin":               originValue,
				}

				// Convert the map to a Terraform object value
				prefixAssignmentValue, _ := types.ObjectValueFrom(ctx, prefixAssignmentAttrTypes, prefixAssignmentObject)
				prefixAssignmentsValues = append(prefixAssignmentsValues, prefixAssignmentValue)
			}

			// Set the 'prefixAssignments' attribute using NewListValue
			ipv6Response.PrefixAssignments, _ = basetypes.NewListValue(
				types.ObjectType{AttrTypes: prefixAssignmentAttrTypes},
				prefixAssignmentsValues,
			)
		}

		// Convert the ipv6Response to a Terraform object value
		ipv6Value, _ := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"enabled":            types.BoolType,
			"prefix_assignments": types.ListType{ElemType: types.ObjectType{AttrTypes: prefixAssignmentAttrTypes}},
		}, map[string]attr.Value{
			"enabled":            ipv6Response.Enabled,
			"prefix_assignments": ipv6Response.PrefixAssignments,
		})

		// Assuming 'data' is your resource data structure
		data.IPv6.As(ctx, ipv6Value, basetypes.ObjectAsOptions{})
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

		var ipv6 NetworksApplianceVLANsResourceModelIpv6
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
			var prefixAssignmentMap map[string]NetworksApplianceVLANsResourceModelIpv6PrefixAssignment

			// Use ElementsAs to convert the elements
			if prefixAssignmentMapDiags := ipv6.PrefixAssignments.ElementsAs(ctx, &prefixAssignmentMap, false); diags.HasError() {
				resp.AddError(
					"Create Payload Failure", fmt.Sprintf("%v", prefixAssignmentMapDiags),
				)
			}

			for _, prefixAssignment := range prefixAssignmentMap {
				var originInterfaces []string

				// Extract the Origin object
				var origin NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
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

	// Id
	if data.Id.IsUnknown() {
		data.Id = types.StringValue("example-id")
	}

	// VlanId
	if data.VlanId.IsUnknown() {
		// check api response for vlanId
		if response.HasId() {

			// API returns string, openAPI spec defines int
			idStr := response.GetId()

			// Check if the string is empty
			if idStr == "" {

				// TODO Handle the case where the ID string is empty (fail or warn?)
				resp.AddWarning("CreatePayloadResponse VlanId Error", "Received empty VlanId from response")
				data.VlanId = types.Int64Null()

			} else {
				// Convert string to int
				vlanId, err := strconv.Atoi(idStr)
				if err != nil {
					// Handle the error if conversion fails
					resp.AddError("CreatePayloadResponse VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s': %v", idStr, err))
				}

				// set new vlanId
				data.VlanId = types.Int64Value(int64(vlanId))
			}
		} else {
			data.VlanId = types.Int64Null()
		}
	}

	// InterfaceId
	if response.HasInterfaceId() {
		data.InterfaceId = types.StringValue(response.GetInterfaceId())
	} else {
		data.InterfaceId = types.StringNull()
	}

	// Name
	if response.HasName() {
		data.Name = types.StringValue(response.GetName())
	} else {
		data.Name = types.StringNull()
	}

	// Subnet
	if response.HasSubnet() {
		data.Subnet = types.StringValue(response.GetSubnet())
	} else {
		data.Subnet = types.StringNull()
	}

	// ApplianceIp
	if response.HasApplianceIp() {
		data.ApplianceIp = types.StringValue(response.GetApplianceIp())
	} else {
		data.ApplianceIp = types.StringNull()
	}

	// GroupPolicyId
	if response.HasGroupPolicyId() {
		data.GroupPolicyId = types.StringValue(response.GetGroupPolicyId())
	} else {
		data.GroupPolicyId = types.StringNull()
	}

	// TemplateVlanType
	if response.HasTemplateVlanType() {
		data.TemplateVlanType = types.StringValue(response.GetTemplateVlanType())
	} else {
		data.TemplateVlanType = types.StringNull()
	}

	// Cidr
	if response.HasCidr() {
		data.Cidr = types.StringValue(response.GetCidr())
	} else {
		data.Cidr = types.StringNull()
	}

	// Mask
	if response.HasMask() {
		data.Mask = types.Int64Value(int64(response.GetMask()))
	} else {
		data.Mask = types.Int64Null()
	}

	// DhcpRelayServerIps
	if data.DhcpRelayServerIps.IsUnknown() {

		if response.HasDhcpRelayServerIps() {
			var dhcpRelayServerIps []attr.Value
			for _, dhcpRelayServerIp := range response.GetDhcpRelayServerIps() {
				dhcpRelayServerIps = append(dhcpRelayServerIps, types.StringValue(dhcpRelayServerIp))

			}
			dhcpRelayServerIpsRespData, dhcpRelayServerIpsDiags := basetypes.NewListValue(types.StringType, dhcpRelayServerIps)
			if dhcpRelayServerIpsDiags.HasError() {
				resp.Append(dhcpRelayServerIpsDiags...)
			}

			data.DhcpRelayServerIps = dhcpRelayServerIpsRespData
		}

	} else {
		data.DhcpRelayServerIps = basetypes.NewListNull(types.StringType)
	}

	// DhcpHandling
	if data.DhcpHandling.IsUnknown() {
		data.DhcpHandling = types.StringValue(response.GetDhcpHandling())

	} else {
		data.DhcpHandling = types.StringNull()
	}

	// DhcpLeaseTime
	if data.DhcpLeaseTime.IsUnknown() {
		data.DhcpLeaseTime = types.StringValue(response.GetDhcpLeaseTime())

	} else {
		data.DhcpLeaseTime = types.StringNull()
	}

	// DhcpBootOptionsEnabled
	if data.DhcpBootOptionsEnabled.IsUnknown() {
		data.DhcpBootOptionsEnabled = types.BoolValue(response.GetDhcpBootOptionsEnabled())
	} else {
		data.DhcpBootOptionsEnabled = types.BoolNull()
	}

	// DhcpBootNextServer
	if data.DhcpBootNextServer.IsUnknown() {
		data.DhcpBootNextServer = types.StringValue(response.GetDhcpBootNextServer())
	} else {
		data.DhcpBootNextServer = types.StringNull()
	}

	// DhcpBootFilename
	if data.DhcpBootFilename.IsUnknown() {
		data.DhcpBootFilename = types.StringValue(response.GetDhcpBootFilename())

	} else {
		data.DhcpBootFilename = types.StringNull()
	}

	// FixedIpAssignments
	if data.FixedIpAssignments.IsUnknown() {
		if response.HasFixedIpAssignments() {
			fixedIpAssignmentsMap := make(map[string]attr.Value)

			fixedIpAssignmentAttrTypes := map[string]attr.Type{
				"ip":   types.StringType,
				"name": types.StringType,
			}

			for macAddress, assignmentInterface := range response.GetFixedIpAssignments() {
				// Check if the value is indeed a map with expected fields
				assignmentMap, ok := assignmentInterface.(map[string]interface{})
				if !ok {
					resp.AddError("failed fixedIpAssignmentMap", "assignmentMap not ok")
					continue
				}

				// Extract IP and Name from the map, asserting their types
				ip, ipOk := assignmentMap["ip"].(string)
				if !ipOk {
					resp.AddError("failed fixedIpAssignmentMap", "ip not ok")
					continue
				}

				name, nameOk := assignmentMap["name"].(string)
				if !nameOk {
					resp.AddError("failed fixedIpAssignmentMap", "name not ok")
					continue
				}

				fixedIpAssignmentObject := map[string]attr.Value{
					"ip":   types.StringValue(ip),
					"name": types.StringValue(name),
				}

				fixedIpAssignmentValue, fixedIpAssignmentsDiags := types.ObjectValueFrom(ctx, fixedIpAssignmentAttrTypes, fixedIpAssignmentObject)
				if fixedIpAssignmentsDiags.HasError() {
					resp.Append(fixedIpAssignmentsDiags...)
					continue
				}

				fixedIpAssignmentsMap[macAddress] = fixedIpAssignmentValue
			}

			fixedIpAssignmentsValue, fixedIpAssignmentsDiags := types.MapValueFrom(ctx, types.MapType{ElemType: types.ObjectType{AttrTypes: fixedIpAssignmentAttrTypes}}, fixedIpAssignmentsMap)
			if fixedIpAssignmentsDiags.HasError() {
				resp.Append(fixedIpAssignmentsDiags...)
			} else {
				data.FixedIpAssignments = fixedIpAssignmentsValue
			}
		}
	}

	// ReservedIpRanges
	if data.ReservedIpRanges.IsUnknown() {

		data.ReservedIpRanges = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"comment": types.StringType,
				"end":     types.StringType,
				"start":   types.StringType,
			},
		})

	} else {

		reservedIpRangeAttrTypes := map[string]attr.Type{
			"comment": types.StringType,
			"end":     types.StringType,
			"start":   types.StringType,
		}

		var reservedIpRangesList []attr.Value

		for _, reservedIpRange := range response.GetReservedIpRanges() {
			reservedIpRangeObject := map[string]attr.Value{
				"comment": types.StringValue(reservedIpRange.GetComment()),
				"end":     types.StringValue(reservedIpRange.GetEnd()),
				"start":   types.StringValue(reservedIpRange.GetStart()),
			}

			reservedIpRangeValue, reservedIpRangeDiags := types.ObjectValueFrom(ctx, reservedIpRangeAttrTypes, reservedIpRangeObject)
			if reservedIpRangeDiags.HasError() {
				resp.Append(reservedIpRangeDiags...)
				continue
			}

			reservedIpRangesList = append(reservedIpRangesList, reservedIpRangeValue)
		}

		// Creating a ListValue from the list of objects.
		reservedIpRangesValue, reservedIpRangesDiags := basetypes.NewListValue(types.ObjectType{AttrTypes: reservedIpRangeAttrTypes}, reservedIpRangesList)
		if reservedIpRangesDiags.HasError() {
			resp.Append(reservedIpRangesDiags...)
		}

		data.ReservedIpRanges = reservedIpRangesValue

	}

	// DnsNameservers
	if data.DnsNameservers.IsUnknown() {
		data.DnsNameservers = types.StringValue(response.GetDnsNameservers())
	} else {
		data.DnsNameservers = types.StringNull()
	}

	// DhcpOptions
	if data.DhcpOptions.IsUnknown() {

		data.DhcpOptions = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"code":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
		})

	} else {

		// Define the structure of each object in the list
		dhcpOptionsAttrTypes := map[string]attr.Type{
			"code":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		}

		// Create a slice to hold the object values
		var interfaceAttrValues []attr.Value

		for _, dhcpOptionsRange := range response.GetDhcpOptions() {
			dhcpOptionsObject := map[string]attr.Value{
				"code":  types.StringValue(dhcpOptionsRange.GetCode()),
				"type":  types.StringValue(dhcpOptionsRange.GetType()),
				"value": types.StringValue(dhcpOptionsRange.GetValue()),
			}

			dhcpOptionValue, dhcpOptionsDiags := types.ObjectValueFrom(ctx, dhcpOptionsAttrTypes, dhcpOptionsObject)
			if dhcpOptionsDiags.HasError() {
				resp.Append(dhcpOptionsDiags...)
				continue
			}

			// Add the constructed object to the slice
			interfaceAttrValues = append(interfaceAttrValues, dhcpOptionValue)
		}

		// Correctly create a ListValue from the slice of ObjectValue
		dhcpOptionsValue, reservedIpRangesDiags := basetypes.NewListValue(types.ObjectType{AttrTypes: dhcpOptionsAttrTypes}, interfaceAttrValues)
		if reservedIpRangesDiags.HasError() {
			resp.Append(reservedIpRangesDiags...)
		}

		data.DhcpOptions = dhcpOptionsValue

	}

	// VpnNatSubnet
	if data.VpnNatSubnet.IsUnknown() {
		data.VpnNatSubnet = types.StringValue(response.GetVpnNatSubnet())

	} else {
		data.VpnNatSubnet = types.StringNull()
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANsResourceModelMandatoryDhcp{}

		// Enabled
		if response.MandatoryDhcp.HasEnabled() {
			mandatoryDhcp.Enabled = types.BoolValue(response.MandatoryDhcp.GetEnabled())
		}

		mandatoryDhcpAttributes := map[string]attr.Type{
			"enabled": types.BoolType,
		}

		objectVal, diags := types.ObjectValueFrom(ctx, mandatoryDhcpAttributes, mandatoryDhcp)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.MandatoryDhcp = objectVal
	} else {
		data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
			"enabled": types.BoolType,
		})
	}

	// TODO: (copy/paste from CreateRespPayload when it's fixed)
	// IPv6
	if response.HasIpv6() {
		ipv6Response := NetworksApplianceVLANsResourceModelIpv6{}

		// Set the 'enabled' attribute
		if response.Ipv6.HasEnabled() {
			ipv6Response.Enabled = types.BoolValue(response.Ipv6.GetEnabled())
		}

		// Define the attribute types for the 'origin' object in each prefix assignment
		originTypes := map[string]attr.Type{
			"type":       types.StringType,
			"interfaces": types.ListType{ElemType: types.StringType},
		}

		// Define the attribute types for each prefix assignment
		prefixAssignmentAttrTypes := map[string]attr.Type{
			"autonomous":           types.BoolType,
			"static_prefix":        types.StringType,
			"static_appliance_ip6": types.StringType,
			"origin":               types.ObjectType{AttrTypes: originTypes},
		}

		// Handling PrefixAssignments
		if response.Ipv6.HasPrefixAssignments() {
			var prefixAssignmentsValues []attr.Value

			for _, prefixAssignmentResponse := range response.Ipv6.PrefixAssignments {
				// Initialize originValue
				var originValue attr.Value

				// Handling the 'origin' object if it exists
				if prefixAssignmentResponse.HasOrigin() {
					var interfacesListValues []attr.Value
					for _, interfaceValue := range prefixAssignmentResponse.Origin.GetInterfaces() {
						interfacesListValues = append(interfacesListValues, types.StringValue(interfaceValue))
					}

					// Construct the interfaces list
					interfacesList, _ := basetypes.NewListValue(types.StringType, interfacesListValues)

					// Create the 'origin' object
					originObject := map[string]attr.Value{
						"type":       types.StringValue(prefixAssignmentResponse.Origin.GetType()),
						"interfaces": interfacesList,
					}
					originValue, _ = types.ObjectValueFrom(ctx, originTypes, originObject)
				}

				// Construct the prefix assignment object
				prefixAssignmentObject := map[string]attr.Value{
					"autonomous":           types.BoolValue(prefixAssignmentResponse.GetAutonomous()),
					"static_prefix":        types.StringValue(prefixAssignmentResponse.GetStaticPrefix()),
					"static_appliance_ip6": types.StringValue(prefixAssignmentResponse.GetStaticApplianceIp6()),
					"origin":               originValue,
				}

				// Convert the map to a Terraform object value
				prefixAssignmentValue, _ := types.ObjectValueFrom(ctx, prefixAssignmentAttrTypes, prefixAssignmentObject)
				prefixAssignmentsValues = append(prefixAssignmentsValues, prefixAssignmentValue)
			}

			// Set the 'prefixAssignments' attribute using NewListValue
			ipv6Response.PrefixAssignments, _ = basetypes.NewListValue(
				types.ObjectType{AttrTypes: prefixAssignmentAttrTypes},
				prefixAssignmentsValues,
			)
		}

		// Convert the ipv6Response to a Terraform object value
		ipv6Value, _ := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"enabled":            types.BoolType,
			"prefix_assignments": types.ListType{ElemType: types.ObjectType{AttrTypes: prefixAssignmentAttrTypes}},
		}, map[string]attr.Value{
			"enabled":            ipv6Response.Enabled,
			"prefix_assignments": ipv6Response.PrefixAssignments,
		})

		// Assuming 'data' is your resource data structure
		data.IPv6.As(ctx, ipv6Value, basetypes.ObjectAsOptions{})
	}

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
				Computed: true,
			},
			"vlan_id": schema.Int64Attribute{
				Required: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"interface_id": schema.StringAttribute{
				MarkdownDescription: "The Interface ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new VLAN",
				Optional:            true,
				Computed:            true,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The subnet of the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"appliance_ip": schema.StringAttribute{
				MarkdownDescription: "The local IP of the appliance on the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"group_policy_id": schema.StringAttribute{
				MarkdownDescription: " desired group policy to apply to the VLAN",
				Optional:            true,
			},
			"vpn_nat_subnet": schema.StringAttribute{
				MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
				Optional:            true,
			},
			"dhcp_handling": schema.StringAttribute{
				MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_relay_server_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
				Optional:    true,
			},
			"dhcp_lease_time": schema.StringAttribute{
				MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_options_enabled": schema.BoolAttribute{
				MarkdownDescription: "Use DHCP boot options specified in other properties",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_next_server": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option to direct boot clients to the server to load the boot file from",
				Optional:            true,
			},
			"dhcp_boot_filename": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option for boot filename ",
				Optional:            true,
			},
			"fixed_ip_assignments": schema.MapNestedAttribute{
				Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Optional:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Optional:            true,
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
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "The last IP in the reserved range",
							Optional:            true,
							Computed:            true,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "A text comment for the reserved range",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"dns_nameservers": schema.StringAttribute{
				MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
				Optional:            true,
				Computed:            true,
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
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type for the DHCP option. One of: 'text', 'ip', 'hex' or 'integer'",
							Optional:            true,
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value for the DHCP option",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"template_vlan_type": schema.StringAttribute{
				MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
				Optional:            true,
				Computed:            true,
			},
			"mask": schema.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
			},
			"ipv6": schema.SingleNestedAttribute{
				Description: "IPv6 configuration on the VLAN",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
					},
					"prefix_assignments": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Prefix assignments on the VLAN",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"autonomous": schema.BoolAttribute{
									MarkdownDescription: "Auto assign a /64 prefix from the origin to the VLAN",
									Optional:            true,
									Computed:            true,
								},
								"static_prefix": schema.StringAttribute{
									MarkdownDescription: "Manual configuration of a /64 prefix on the VLAN",
									Optional:            true,
									Computed:            true,
								},
								"static_appliance_ip6": schema.StringAttribute{
									MarkdownDescription: "Manual configuration of the IPv6 Appliance IP",
									Optional:            true,
									Computed:            true,
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
										},
										"interfaces": schema.ListAttribute{
											ElementType: types.StringType,
											Description: "Interfaces associated with the prefix",
											Optional:    true,
											Computed:    true,
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
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
						Optional:            true,
						Computed:            true,
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
