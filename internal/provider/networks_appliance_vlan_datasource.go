package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"strconv"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksApplianceVLANsDatasource{}

func NewNetworksApplianceVLANsDatasource() datasource.DataSource {
	return &NetworksApplianceVLANsDatasource{}
}

// NetworksApplianceVLANsDatasource defines the resource implementation.
type NetworksApplianceVLANsDatasource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceVLANsDatasourceModel NetworksApplianceVLANResourceModel describes the resource data model.
type NetworksApplianceVLANsDatasourceModel struct {
	Id        types.String                            `tfsdk:"id" json:"-"`
	NetworkId types.String                            `tfsdk:"network_id" json:"network_id"`
	List      []NetworksApplianceVLANsDataSourceModel `tfsdk:"list"`
}

type NetworksApplianceVLANsDataSourceModel struct {
	Id                     types.String `tfsdk:"id" json:"-"`
	NetworkId              types.String `tfsdk:"network_id" json:"networkId"`
	VlanId                 types.Int64  `tfsdk:"vlan_id" json:"-"`
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

type NetworksApplianceVLANsDataSourceModelIpNameMapping struct {
	Ip   types.String `tfsdk:"ip" json:"ip"`
	Name types.String `tfsdk:"name" json:"name"`
}

type NetworksApplianceVLANsDataSourceModelReservedIpRange struct {
	Start   types.String `tfsdk:"start" json:"start"`
	End     types.String `tfsdk:"end" json:"end"`
	Comment types.String `tfsdk:"comment" json:"comment"`
}

type NetworksApplianceVLANsDataSourceFixedIpAssignmentTerraform struct {
	IP   types.String `tfsdk:"ip"`
	Name types.String `tfsdk:"name"`
}

type NetworksApplianceVLANsDataSourceFixedIpAssignment struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

func (n *NetworksApplianceVLANsDataSourceModelReservedIpRange) FromTerraformValue(ctx context.Context, val tftypes.Value) error {
	// Assuming val is a tftypes.Object with "start", "end", and "comment"
	var data map[string]tftypes.Value
	if !val.IsKnown() || val.IsNull() {
		return errors.New("comment is unknown or null")
	}

	conversionErr := val.As(&data)
	if conversionErr != nil {
		return conversionErr
	}

	var start, end, comment string
	if conversionErr = data["start"].As(&start); conversionErr != nil {
		return conversionErr
	}
	n.Start = basetypes.NewStringValue(start)

	if conversionErr = data["end"].As(&end); conversionErr != nil {
		return conversionErr
	}
	n.End = basetypes.NewStringValue(end)

	if conversionErr = data["comment"].As(&comment); conversionErr != nil {
		return conversionErr
	}
	n.Comment = basetypes.NewStringValue(comment)

	return nil
}

type NetworksApplianceVLANsDataSourceModelDhcpOption struct {
	Code  types.String `tfsdk:"code" json:"code"`
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

func (n *NetworksApplianceVLANsDataSourceModelDhcpOption) FromTerraformValue(ctx context.Context, val tftypes.Value) error {
	// Assuming val is a tftypes.Object with "code", "type", and "value"
	var data map[string]tftypes.Value
	if !val.IsKnown() || val.IsNull() {
		return errors.New("value is unknown or null")
	}

	conversionErr := val.As(&data)
	if conversionErr != nil {
		return conversionErr
	}

	var code, typ, value string
	if conversionErr = data["code"].As(&code); conversionErr != nil {
		return conversionErr
	}
	n.Code = basetypes.NewStringValue(code)

	if conversionErr = data["type"].As(&typ); conversionErr != nil {
		return conversionErr
	}
	n.Type = basetypes.NewStringValue(typ)

	if conversionErr = data["value"].As(&value); conversionErr != nil {
		return conversionErr
	}
	n.Value = basetypes.NewStringValue(value)

	return nil
}

type NetworksApplianceVLANsDataSourceModelFixedIpAssignment struct {
	Ip   types.String `tfsdk:"ip"`
	Name types.String `tfsdk:"name"`
}

// NetworksApplianceVLANsDataSourceModelIpv6 represents the IPv6 configuration for a VLAN resource model.
type NetworksApplianceVLANsDataSourceModelIpv6 struct {
	Enabled           types.Bool `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments types.List `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

func NetworksApplianceVLANsDataSourceModelIpv6AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":            types.BoolType,
		"prefix_assignments": types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentAttrTypes()}},
	}
}

// ToAPIPayload converts the Terraform resource data model into the API payload.
func (m *NetworksApplianceVLANsDataSourceModelIpv6) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6, diag.Diagnostics) {

	payload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6{}

	// Convert 'Enabled' field
	payload.Enabled = m.Enabled.ValueBoolPointer()

	var prefixAssignments []NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment

	// Convert 'PrefixAssignments' field
	err := m.PrefixAssignments.ElementsAs(ctx, &prefixAssignments, false)
	if err != nil {
		return nil, err.Errors()
	}

	for _, prefixAssignment := range prefixAssignments {

		var prefixAssignmentPayload openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner

		prefixAssignmentPayload.SetStaticPrefix(prefixAssignment.StaticPrefix.ValueString())
		prefixAssignmentPayload.SetStaticApplianceIp6(prefixAssignment.StaticApplianceIp6.ValueString())
		prefixAssignmentPayload.SetAutonomous(prefixAssignment.Autonomous.ValueBool())

		var origin NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin
		prefixAssignment.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})

		originPayload, err := origin.ToAPIPayload(ctx)
		if err != nil {
			return nil, err.Errors()
		}

		prefixAssignmentPayload.SetOrigin(*originPayload)

		payload.PrefixAssignments = append(payload.PrefixAssignments, prefixAssignmentPayload)
	}

	return payload, nil
}

// FromAPIResponse transforms an API response into the NetworksApplianceVLANsDataSourceModelIpv6 Terraform structure.
func (m *NetworksApplianceVLANsDataSourceModelIpv6) FromAPIResponse(ctx context.Context, apiResponse *openApiClient.GetNetworkApplianceVlans200ResponseInnerIpv6) diag.Diagnostics {
	if apiResponse == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API response for IPv6")}
	}

	m.Enabled = types.BoolValue(apiResponse.GetEnabled())

	var prefixAssignments []NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment
	for _, apiPA := range apiResponse.PrefixAssignments {
		var pa NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment

		diags := pa.FromAPIResponse(ctx, &apiPA)
		if diags.HasError() {
			tflog.Warn(ctx, "failed to extract FromAPIResponse to PrefixAssignments")
			return diags
		}

		prefixAssignments = append(prefixAssignments, pa)
	}

	p, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentAttrTypes()}, prefixAssignments)

	m.PrefixAssignments = p
	return nil
}

// NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentAttrTypes returns the attribute types for a prefix assignment.
func NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"autonomous":           types.BoolType,
		"static_prefix":        types.StringType,
		"static_appliance_ip6": types.StringType,
		"origin":               types.ObjectType{AttrTypes: NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrTypes()},
	}
}

// NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment represents a prefix assignment for an IPv6 configuration in the VLAN resource model.
type NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment struct {
	Autonomous         types.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       types.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 types.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             types.Object `tfsdk:"origin" json:"origin"`
}

func (pa *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner, diag.Diagnostics) {
	apiPA := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
		Autonomous:         pa.Autonomous.ValueBoolPointer(),
		StaticPrefix:       pa.StaticPrefix.ValueStringPointer(),
		StaticApplianceIp6: pa.StaticApplianceIp6.ValueStringPointer(),
	}

	var originObject NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin
	originObjectDiags := pa.Origin.As(ctx, originObject, basetypes.ObjectAsOptions{})
	if originObjectDiags.HasError() {
		return nil, originObjectDiags
	}

	// If 'Origin' is a nested structure, convert it too
	if !pa.Origin.IsNull() {
		originAPIModel, diags := originObject.ToAPIPayload(ctx)
		if diags.HasError() {
			return nil, diags
		}
		apiPA.Origin = originAPIModel
	}

	return apiPA, nil
}

// ToAPIPayload converts the Terraform prefix assignment into the API prefix assignment payload.
func (pa *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner, diag.Diagnostics) {

	paPayload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{}

	// Autonomous
	paPayload.Autonomous = pa.Autonomous.ValueBoolPointer()

	// StaticPrefix
	paPayload.StaticPrefix = pa.StaticPrefix.ValueStringPointer()

	// StaticApplianceIp6
	paPayload.StaticApplianceIp6 = pa.StaticApplianceIp6.ValueStringPointer()

	// Origin
	var origin NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin
	diags := pa.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	originPayload, diags := origin.ToAPIPayload(ctx)
	if diags.HasError() {
		return nil, diags
	}

	paPayload.Origin = originPayload

	return paPayload, nil
}

// ToTerraformObject converts the NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment instance to a map suitable for ObjectValueFrom.
func (pa *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment) ToTerraformObject(ctx context.Context) (map[string]attr.Value, diag.Diagnostics) {
	return map[string]attr.Value{
		"autonomous":           pa.Autonomous,
		"static_prefix":        pa.StaticPrefix,
		"static_appliance_ip6": pa.StaticApplianceIp6,
		"origin":               pa.Origin,
	}, nil
}

// FromAPIResponse fills the NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment with data from the API response.
// This method transforms the OpenAPI response into the format expected by the Terraform provider.
func (pa *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment) FromAPIResponse(ctx context.Context, apiPA *openApiClient.GetNetworkApplianceVlans200ResponseInnerIpv6PrefixAssignmentsInner) diag.Diagnostics {
	if apiPA == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("API Prefix Assignment Error", "Received nil API prefix assignment")}
	}

	pa.Autonomous = types.BoolValue(apiPA.GetAutonomous())

	// staticPrefix
	staticPrefix := apiPA.GetStaticPrefix()

	if staticPrefix == "" {
		// Handle the null scenario
		pa.StaticPrefix = types.StringNull()
	} else {

		pa.StaticPrefix = types.StringValue(apiPA.GetStaticPrefix())

		if pa.StaticPrefix.IsUnknown() {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic("Invalid Assignment", "The 'staticPrefix' field assignment resulted in an unknown value"),
			}
		}
	}

	// staticApplianceIp6
	staticApplianceIp6 := apiPA.GetStaticApplianceIp6()

	if staticApplianceIp6 == "" {
		// Handle the null scenario
		pa.StaticApplianceIp6 = types.StringNull()
	} else {
		pa.StaticApplianceIp6 = types.StringValue(staticApplianceIp6)

		if pa.StaticApplianceIp6.IsUnknown() {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic("Invalid Assignment", "The 'staticApplianceIp6' field assignment resulted in an unknown value"),
			}
		}
	}

	var origin NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin
	originDiags := origin.FromAPIResponse(ctx, apiPA.Origin)
	if originDiags.HasError() {
		return originDiags
	}

	// Use the predefined functions for attribute types and map
	originAttrTypes := NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrTypes()
	originAttrMap := NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrMap(&origin)

	originTf, diags := types.ObjectValue(originAttrTypes, originAttrMap)
	if diags.HasError() {
		tflog.Warn(ctx, "failed to create object from PrefixAssignment Origin")
		return diags
	}

	pa.Origin = originTf

	return nil
}

// NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin represents the origin data structure for a VLAN resource model.
type NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin struct {
	Type       types.String `tfsdk:"type" json:"type"`
	Interfaces types.Set    `tfsdk:"interfaces" json:"interfaces"`
}

// NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrTypes returns the attribute types for the origin.
// This function is useful to define the schema of the origin in a consistent manner.
func NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":       types.StringType,
		"interfaces": types.SetType{ElemType: types.StringType},
	}
}

// NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrMap returns the attribute map for a given origin.
// It converts a NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin instance to a map suitable for ObjectValueFrom.
func NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOriginAttrMap(origin *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin) map[string]attr.Value {
	return map[string]attr.Value{
		"type":       origin.Type,
		"interfaces": origin.Interfaces,
	}
}

// ToAPIPayload converts the Terraform origin into the API origin payload.
func (o *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin, diag.Diagnostics) {

	originPayload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin{}

	originPayload.SetType(o.Type.ValueString())

	var interfaces []string
	// The ElementsAs function usually takes a pointer to the slice
	if diags := o.Interfaces.ElementsAs(ctx, &interfaces, false); diags.HasError() {
		return nil, diags
	}

	// Process interfaces to remove extra quotes if necessary
	for i, iface := range interfaces {
		interfaces[i] = strings.Trim(iface, "\"")
	}

	originPayload.SetInterfaces(interfaces)

	return originPayload, nil
}

// FromAPIResponse fills the NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin with data from the API response.
// This method transforms the OpenAPI response into the format expected by the Terraform provider.
func (o *NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentOrigin) FromAPIResponse(ctx context.Context, apiOrigin *openApiClient.CreateNetworkAppliancePrefixesDelegatedStaticRequestOrigin) diag.Diagnostics {
	// Get the type from API response
	apiType := apiOrigin.GetType()

	// Validate the apiType
	// (Add any specific validation logic here. For example, checking if it's non-empty, or if it matches certain criteria)
	if apiType == "" {
		// Handle the invalid scenario
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("Invalid API Origin Type", "The 'type' field from the API origin is empty"),
		}
	}

	if apiOrigin == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("API Origin Error", "Received nil API origin")}
	}

	// Assuming the value is valid, set the Type
	o.Type = types.StringValue(apiType)

	// Additional check: Verify if the assignment is as expected
	// This part depends on the implementation details of types.StringValue and your specific validation needs
	if o.Type.IsUnknown() || o.Type.IsNull() {
		// Handle the scenario where the assignment didn't work as expected
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("Invalid Assignment", "The 'type' field assignment resulted in an unknown or null value"),
		}
	}

	var interfaces []types.String
	for _, iface := range apiOrigin.Interfaces {
		interfaces = append(interfaces, types.StringValue(iface))
	}

	var diags diag.Diagnostics
	o.Interfaces, diags = types.SetValueFrom(ctx, types.StringType, interfaces)
	if diags.HasError() {
		tflog.Warn(ctx, "failed to create list from Origin interfaces")
		return diags
	}

	return nil
}

type NetworksApplianceVLANsDataSourceModelMandatoryDhcp struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

func NetworksApplianceVLANsDatasourceReadHttpResponse(ctx context.Context, data *NetworksApplianceVLANsDataSourceModel, response *openApiClient.GetNetworkApplianceVlans200ResponseInner) diag.Diagnostics {

	resp := diag.Diagnostics{}

	// Id field only returns "", this is a bug in the HTTP client

	// InterfaceId
	if response.HasInterfaceId() {
		data.InterfaceId = types.StringValue(response.GetInterfaceId())
	} else {
		if data.InterfaceId.IsUnknown() {
			data.InterfaceId = types.StringNull()
		}
	}

	// Name
	if response.HasName() {
		data.Name = types.StringValue(response.GetName())
	} else {
		if data.Name.IsUnknown() {
			data.Name = types.StringNull()
		}
	}

	// Subnet
	if response.HasSubnet() {
		data.Subnet = types.StringValue(response.GetSubnet())
	} else {
		if data.Subnet.IsUnknown() {
			data.Subnet = types.StringNull()
		}
	}

	// ApplianceIp
	if response.HasApplianceIp() {
		data.ApplianceIp = types.StringValue(response.GetApplianceIp())
	} else {
		if data.ApplianceIp.IsUnknown() {
			data.ApplianceIp = types.StringNull()
		}
	}

	// GroupPolicyId
	if response.HasGroupPolicyId() {
		data.GroupPolicyId = types.StringValue(response.GetGroupPolicyId())
	} else {
		if data.GroupPolicyId.IsUnknown() {
			data.GroupPolicyId = types.StringNull()
		}
	}

	// TemplateVlanType
	if response.HasTemplateVlanType() {
		data.TemplateVlanType = types.StringValue(response.GetTemplateVlanType())
	} else {
		if data.TemplateVlanType.IsUnknown() {
			data.TemplateVlanType = types.StringNull()
		}
	}

	// Cidr
	if response.HasCidr() {
		data.Cidr = types.StringValue(response.GetCidr())
	} else {
		if data.Cidr.IsUnknown() {
			data.Cidr = types.StringNull()
		}
	}

	// Mask
	if response.HasMask() {
		data.Mask = types.Int64Value(int64(response.GetMask()))
	} else {
		if data.Mask.IsUnknown() {
			data.Mask = types.Int64Null()
		}
	}

	// DhcpRelayServerIps
	data.DhcpRelayServerIps = basetypes.NewListNull(types.StringType)
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

	// DhcpHandling
	if response.HasDhcpHandling() {
		data.DhcpHandling = types.StringValue(response.GetDhcpHandling())

	} else {
		if data.DhcpHandling.IsUnknown() {
			data.DhcpHandling = types.StringNull()
		}
	}

	// DhcpLeaseTime
	if response.HasDhcpLeaseTime() {
		data.DhcpLeaseTime = types.StringValue(response.GetDhcpLeaseTime())

	} else {
		if data.DhcpLeaseTime.IsUnknown() {
			data.DhcpLeaseTime = types.StringNull()
		}
	}

	// DhcpBootOptionsEnabled
	if response.HasDhcpBootOptionsEnabled() {
		data.DhcpBootOptionsEnabled = types.BoolValue(response.GetDhcpBootOptionsEnabled())
	} else {
		if data.DhcpBootOptionsEnabled.IsUnknown() {
			data.DhcpBootOptionsEnabled = types.BoolValue(false)
		}
	}

	// DhcpBootNextServer
	if response.HasDhcpBootNextServer() {
		data.DhcpBootNextServer = types.StringValue(response.GetDhcpBootNextServer())
	} else {
		if data.DhcpBootNextServer.IsUnknown() {
			data.DhcpBootNextServer = types.StringNull()
		}
	}

	// DhcpBootFilename
	if response.HasDhcpBootFilename() {
		data.DhcpBootFilename = types.StringValue(response.GetDhcpBootFilename())

	} else {
		if data.DhcpBootFilename.IsUnknown() {
			data.DhcpBootFilename = types.StringNull()
		}
	}

	// FixedIpAssignments
	if response.HasFixedIpAssignments() {
		fixedIpAssignmentsMap := map[string]attr.Value{}

		fixedIpAssignmentAttrTypes := map[string]attr.Type{
			"ip":   types.StringType,
			"name": types.StringType,
		}

		for macAddress, assignmentInterface := range response.GetFixedIpAssignments() {

			fixedIpAssignmentObject := map[string]attr.Value{
				"ip":   types.StringValue(assignmentInterface.GetIp()),
				"name": types.StringValue(assignmentInterface.GetName()),
			}

			fixedIpAssignmentValue, fixedIpAssignmentsDiags := types.ObjectValue(fixedIpAssignmentAttrTypes, fixedIpAssignmentObject)
			if fixedIpAssignmentsDiags.HasError() {
				resp.Append(fixedIpAssignmentsDiags...)
				continue
			}

			fixedIpAssignmentsMap[macAddress] = fixedIpAssignmentValue
		}

		fixedIpAssignmentsValue, fixedIpAssignmentsDiags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: fixedIpAssignmentAttrTypes}, fixedIpAssignmentsMap)
		if fixedIpAssignmentsDiags.HasError() {
			resp.Append(fixedIpAssignmentsDiags...)
		}

		data.FixedIpAssignments = fixedIpAssignmentsValue

	} else {
		if data.FixedIpAssignments.IsUnknown() {
			fixedIpAssignmentsAttrTypes := map[string]attr.Type{
				"ip":   types.StringType,
				"name": types.StringType,
			}

			data.FixedIpAssignments = types.MapNull(
				types.ObjectType{AttrTypes: fixedIpAssignmentsAttrTypes},
			)
		}
	}

	if response.HasReservedIpRanges() {

		reservedIpRangeAttrTypes := map[string]attr.Type{
			"comment": types.StringType,
			"end":     types.StringType,
			"start":   types.StringType,
		}

		// Define the ObjectType
		objectType := types.ObjectType{AttrTypes: reservedIpRangeAttrTypes}

		// Create a slice to hold the ObjectValues
		var objectValues []attr.Value

		for _, reservedIpRange := range response.GetReservedIpRanges() {

			// Construct the map for the current ObjectValue
			valuesMap, valuesMapErrs := basetypes.NewObjectValue(reservedIpRangeAttrTypes, map[string]attr.Value{
				"start":   basetypes.NewStringValue(reservedIpRange.GetStart()),
				"end":     basetypes.NewStringValue(reservedIpRange.GetEnd()),
				"comment": basetypes.NewStringValue(reservedIpRange.GetComment()),
			})

			if valuesMapErrs.HasError() {
				for _, valuesMapErr := range valuesMapErrs.Errors() {
					tflog.Error(ctx, valuesMapErr.Summary()+valuesMapErr.Detail())
				}
				resp.Append(valuesMapErrs...)

				continue
			}

			// Create the ObjectValue for the current DHCP option
			reservedIpRangeValue, reservedIpRangeDiags := types.ObjectValueFrom(ctx, reservedIpRangeAttrTypes, valuesMap)
			if reservedIpRangeDiags.HasError() {
				for _, dhcpOptionReservedIpRangeDiag := range reservedIpRangeDiags.Errors() {
					tflog.Error(ctx, dhcpOptionReservedIpRangeDiag.Summary()+dhcpOptionReservedIpRangeDiag.Detail())
				}
				resp.Append(reservedIpRangeDiags...)

				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, reservedIpRangeValue)

		}

		// Create a ListValue from the slice of ObjectValue
		reservedIpRangesValue, reservedIpRangesListDiags := types.ListValueFrom(ctx, objectType, objectValues)
		if reservedIpRangesListDiags.HasError() {
			resp.Append(reservedIpRangesListDiags...)
		}

		data.ReservedIpRanges = reservedIpRangesValue

	} else {
		if data.ReservedIpRanges.IsUnknown() {
			data.ReservedIpRanges = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"comment": types.StringType,
					"end":     types.StringType,
					"start":   types.StringType,
				},
			})
		}
	}

	// DnsNameservers
	if response.HasDnsNameservers() {
		data.DnsNameservers = types.StringValue(response.GetDnsNameservers())
	} else {
		if data.DnsNameservers.IsUnknown() {
			data.DnsNameservers = types.StringNull()
		}
	}

	// VpnNatSubnet
	if response.HasVpnNatSubnet() {
		data.VpnNatSubnet = types.StringValue(response.GetVpnNatSubnet())

	} else {
		if data.VpnNatSubnet.IsUnknown() {
			data.VpnNatSubnet = types.StringNull()
		}
	}

	// DhcpOptions
	if response.HasDhcpOptions() {
		// Define the structure of each object in the list
		dhcpOptionsAttrTypes := map[string]attr.Type{
			"code":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		}

		// Define the ObjectType
		objectType := types.ObjectType{AttrTypes: dhcpOptionsAttrTypes}

		// Create a slice to hold the ObjectValues
		var objectValues []attr.Value

		for _, dhcpOption := range response.GetDhcpOptions() {
			valuesMap := map[string]attr.Value{
				"code":  types.StringValue(dhcpOption.GetCode()),
				"type":  types.StringValue(dhcpOption.GetType()),
				"value": types.StringValue(dhcpOption.GetValue()),
			}

			// Create the ObjectValue for the current DHCP option
			dhcpOptionValue, dhcpOptionsDiags := types.ObjectValue(dhcpOptionsAttrTypes, valuesMap)
			if dhcpOptionsDiags.HasError() {
				resp.Append(dhcpOptionsDiags...)
				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, dhcpOptionValue)
		}

		// Create a ListValue from the slice of ObjectValues
		dhcpOptionsListValue, dhcpOptionsListDiags := types.ListValue(objectType, objectValues)
		if dhcpOptionsListDiags.HasError() {
			resp.Append(dhcpOptionsListDiags...)
		}

		data.DhcpOptions = dhcpOptionsListValue
	} else {
		if data.DhcpOptions.IsUnknown() {

			data.DhcpOptions = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"code":  types.StringType,
					"type":  types.StringType,
					"value": types.StringType,
				},
			})
		}
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANsDataSourceModelMandatoryDhcp{}

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
		if data.MandatoryDhcp.IsUnknown() {
			data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
				"enabled": types.BoolType,
			})
		}
	}

	// IPv6
	// Assuming response is a structure containing your data
	if response.HasIpv6() {
		ipv6Instance := NetworksApplianceVLANsDataSourceModelIpv6{}
		diags := ipv6Instance.FromAPIResponse(ctx, response.Ipv6)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Object, diags := types.ObjectValueFrom(ctx, NetworksApplianceVLANsDataSourceModelIpv6AttrTypes(), ipv6Instance)

		if diags.HasError() {
			resp.Append(diags...)
		}

		data.IPv6 = ipv6Object
	} else {
		if data.IPv6.IsUnknown() {
			ipv6Instance := NetworksApplianceVLANsDataSourceModelIpv6{}
			ipv6Prefixes := []NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignment{}

			ipv6PrefixesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NetworksApplianceVLANsDataSourceModelIpv6PrefixAssignmentAttrTypes()}, ipv6Prefixes)
			if diags.HasError() {
				resp.Append(diags...)
			}

			ipv6Instance.PrefixAssignments = ipv6PrefixesList

			ipv6Object, diags := types.ObjectValueFrom(ctx, NetworksApplianceVLANsDataSourceModelIpv6AttrTypes(), ipv6Instance)
			if diags.HasError() {
				resp.Append(diags...)
			}

			data.IPv6 = ipv6Object
		}
	}

	return resp
}

func (r *NetworksApplianceVLANsDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *NetworksApplianceVLANsDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ".",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"vlan_id": schema.Int64Attribute{
						Computed: true,
						Optional: true,
					},
					"network_id": schema.StringAttribute{
						MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
						Computed:            true,
						Optional:            true,
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
						Computed:            true,
					},
					"vpn_nat_subnet": schema.StringAttribute{
						MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
						Optional:            true,
						Computed:            true,
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
						Computed:    true,
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
						Computed:            true,
					},
					"dhcp_boot_filename": schema.StringAttribute{
						MarkdownDescription: "DHCP boot option for boot filename ",
						Optional:            true,
						Computed:            true,
					},
					"fixed_ip_assignments": schema.MapNestedAttribute{
						Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"ip": schema.StringAttribute{
									MarkdownDescription: "Enable IPv6 on VLAN.",
									Required:            true,
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
									Validators: []validator.String{
										stringvalidator.OneOf("text", "ip", "hex", "integer"),
									},
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
						Validators: []validator.String{
							stringvalidator.OneOf("same", "unique"),
						},
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
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Enable IPv6 on VLAN.",
								Optional:            true,
								Computed:            true,
							},
							"prefix_assignments": schema.ListNestedAttribute{
								Optional:    true,
								Computed:    true,
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
													Validators: []validator.String{
														stringvalidator.OneOf("independent", "internet"),
													},
												},
												"interfaces": schema.SetAttribute{
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
				}}},
		},
	}
}

func (r *NetworksApplianceVLANsDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NetworksApplianceVLANsDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksApplianceVLANsDatasourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlans(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
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
		return
	} else {
		resp.Diagnostics.Append()
	}

	// Assuming httpResp is your *http.Response object
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error: unable to read the response body
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read HTTP response body: %v", err))
		return
	}

	// Define a struct to specifically capture the ID from the JSON data
	type HttpRespID struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	var jsonResponse []HttpRespID
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		// Handle error: JSON parsing error
		resp.Diagnostics.AddError("JSON Parsing Error", fmt.Sprintf("Error parsing JSON data for ID field: %v", err))
	}

	for _, inRespData := range inlineResp {

		vlanData := NetworksApplianceVLANsDataSourceModel{}
		vlanData.NetworkId = types.StringValue(data.NetworkId.ValueString())

		// Workaround for Id bug in client.GetNetworkApplianceVlans200ResponseInner
		for _, jsonInRespData := range jsonResponse {
			if jsonInRespData.Name == inRespData.GetName() {

				/*
					// Convert string to int64
							vlanId, err := strconv.ParseInt(idStr, 10, 64)
							if err != nil {
								resp.AddError("VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s' to int64: %v", idStr, err))

				*/
				vlanData.VlanId = types.Int64Value(jsonInRespData.ID)
				data.Id = types.StringValue(fmt.Sprintf("%s,%v", data.NetworkId.ValueString(), strconv.FormatInt(jsonInRespData.ID, 10)))
			}
		}

		payloadRespDiags := NetworksApplianceVLANsDatasourceReadHttpResponse(ctx, &vlanData, &inRespData)
		if payloadRespDiags != nil {
			resp.Diagnostics.Append(payloadRespDiags...)
		}

		data.List = append(data.List, vlanData)

	}

	data.Id = types.StringValue(data.NetworkId.ValueString())

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
