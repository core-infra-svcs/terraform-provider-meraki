package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
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

type FixedIpAssignmentTerraform struct {
	IP   types.String `tfsdk:"ip"`
	Name types.String `tfsdk:"name"`
}

type FixedIpAssignment struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

func (n *NetworksApplianceVLANsResourceModelReservedIpRange) FromTerraformValue(ctx context.Context, val tftypes.Value) error {
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

type NetworksApplianceVLANsResourceModelDhcpOption struct {
	Code  types.String `tfsdk:"code" json:"code"`
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

func (n *NetworksApplianceVLANsResourceModelDhcpOption) FromTerraformValue(ctx context.Context, val tftypes.Value) error {
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

type NetworksApplianceVLANsResourceModelFixedIpAssignment struct {
	Ip   types.String `tfsdk:"ip"`
	Name types.String `tfsdk:"name"`
}

// NetworksApplianceVLANsResourceModelIpv6 represents the IPv6 configuration for a VLAN resource model.
type NetworksApplianceVLANsResourceModelIpv6 struct {
	Enabled           types.Bool `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments types.List `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

func NetworksApplianceVLANsResourceModelIpv6AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":            types.BoolType,
		"prefix_assignments": types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentAttrTypes()}},
	}
}

// ToAPIPayload converts the Terraform resource data model into the API payload.
func (m *NetworksApplianceVLANsResourceModelIpv6) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6, diag.Diagnostics) {
	tflog.Info(ctx, "[start] NetworksApplianceVLANsResourceModelIpv6 ToAPIPayload")

	payload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6{}

	// Convert 'Enabled' field
	payload.Enabled = m.Enabled.ValueBoolPointer()

	var prefixAssignments []NetworksApplianceVLANsResourceModelIpv6PrefixAssignment

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

		var origin NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
		prefixAssignment.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})

		originPayload, err := origin.ToAPIPayload(ctx)
		if err != nil {
			return nil, err.Errors()
		}

		prefixAssignmentPayload.SetOrigin(*originPayload)

		payload.PrefixAssignments = append(payload.PrefixAssignments, prefixAssignmentPayload)
	}

	tflog.Info(ctx, "[finish] NetworksApplianceVLANsResourceModelIpv6 ToAPIPayload")
	return payload, nil
}

// FromAPIResponse transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *NetworksApplianceVLANsResourceModelIpv6) FromAPIResponse(ctx context.Context, apiResponse *openApiClient.GetNetworkApplianceVlans200ResponseInnerIpv6) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksApplianceVLANsResourceModelIpv6 FromAPIResponse")
	tflog.Trace(ctx, "NetworksApplianceVLANsResourceModelIpv6 FromAPIResponse", map[string]interface{}{
		"apiResponse": apiResponse,
	})
	if apiResponse == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API response for IPv6")}
	}

	m.Enabled = types.BoolValue(apiResponse.GetEnabled())

	var prefixAssignments []NetworksApplianceVLANsResourceModelIpv6PrefixAssignment
	for _, apiPA := range apiResponse.PrefixAssignments {
		var pa NetworksApplianceVLANsResourceModelIpv6PrefixAssignment

		diags := pa.FromAPIResponse(ctx, &apiPA)
		if diags.HasError() {
			tflog.Warn(ctx, "failed to extract FromAPIResponse to PrefixAssignments")
			return diags
		}

		prefixAssignments = append(prefixAssignments, pa)
	}

	p, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentAttrTypes()}, prefixAssignments)

	m.PrefixAssignments = p
	tflog.Info(ctx, "[finish] NetworksApplianceVLANsResourceModelIpv6 FromAPIResponse")
	return nil
}

// NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentAttrTypes returns the attribute types for a prefix assignment.
func NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"autonomous":           types.BoolType,
		"static_prefix":        types.StringType,
		"static_appliance_ip6": types.StringType,
		"origin":               types.ObjectType{AttrTypes: NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrTypes()},
	}
}

// NetworksApplianceVLANsResourceModelIpv6PrefixAssignment represents a prefix assignment for an IPv6 configuration in the VLAN resource model.
type NetworksApplianceVLANsResourceModelIpv6PrefixAssignment struct {
	Autonomous         types.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       types.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 types.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             types.Object `tfsdk:"origin" json:"origin"`
}

func (pa *NetworksApplianceVLANsResourceModelIpv6PrefixAssignment) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner, diag.Diagnostics) {
	apiPA := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
		Autonomous:         pa.Autonomous.ValueBoolPointer(),
		StaticPrefix:       pa.StaticPrefix.ValueStringPointer(),
		StaticApplianceIp6: pa.StaticApplianceIp6.ValueStringPointer(),
	}

	var originObject NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
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
func (pa *NetworksApplianceVLANsResourceModelIpv6PrefixAssignment) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner, diag.Diagnostics) {
	tflog.Info(ctx, "[start] NetworksApplianceVLANsResourceModelIpv6PrefixAssignment ToAPIPayload")

	paPayload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{}

	// Autonomous
	paPayload.Autonomous = pa.Autonomous.ValueBoolPointer()

	// StaticPrefix
	paPayload.StaticPrefix = pa.StaticPrefix.ValueStringPointer()

	// StaticApplianceIp6
	paPayload.StaticApplianceIp6 = pa.StaticApplianceIp6.ValueStringPointer()

	// Origin
	var origin NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
	diags := pa.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	originPayload, diags := origin.ToAPIPayload(ctx)
	if diags.HasError() {
		return nil, diags
	}

	paPayload.Origin = originPayload

	tflog.Info(ctx, "[finish] NetworksApplianceVLANsResourceModelIpv6PrefixAssignment ToAPIPayload")
	return paPayload, nil
}

// ToTerraformObject converts the NetworksApplianceVLANsResourceModelIpv6PrefixAssignment instance to a map suitable for ObjectValueFrom.
func (pa *NetworksApplianceVLANsResourceModelIpv6PrefixAssignment) ToTerraformObject(ctx context.Context) (map[string]attr.Value, diag.Diagnostics) {
	return map[string]attr.Value{
		"autonomous":           pa.Autonomous,
		"static_prefix":        pa.StaticPrefix,
		"static_appliance_ip6": pa.StaticApplianceIp6,
		"origin":               pa.Origin,
	}, nil
}

// FromAPIResponse fills the NetworksApplianceVLANsResourceModelIpv6PrefixAssignment with data from the API response.
// This method transforms the OpenAPI response into the format expected by the Terraform provider.
func (pa *NetworksApplianceVLANsResourceModelIpv6PrefixAssignment) FromAPIResponse(ctx context.Context, apiPA *openApiClient.GetNetworkApplianceVlans200ResponseInnerIpv6PrefixAssignmentsInner) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksApplianceVLANsResourceModelIpv6PrefixAssignment FromAPIResponse")
	tflog.Trace(ctx, "NetworksApplianceVLANsResourceModelIpv6 FromAPIResponse", map[string]interface{}{
		"apiPA": apiPA,
	})
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

	var origin NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
	originDiags := origin.FromAPIResponse(ctx, apiPA.Origin)
	if originDiags.HasError() {
		return originDiags
	}

	// Use the predefined functions for attribute types and map
	originAttrTypes := NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrTypes()
	originAttrMap := NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrMap(&origin)

	originTf, diags := types.ObjectValue(originAttrTypes, originAttrMap)
	if diags.HasError() {
		tflog.Warn(ctx, "failed to create object from PrefixAssignment Origin")
		return diags
	}

	pa.Origin = originTf

	tflog.Info(ctx, "[finish] NetworksApplianceVLANsResourceModelIpv6PrefixAssignment FromAPIResponse")
	return nil
}

// NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin represents the origin data structure for a VLAN resource model.
type NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin struct {
	Type       types.String `tfsdk:"type" json:"type"`
	Interfaces types.Set    `tfsdk:"interfaces" json:"interfaces"`
}

// NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrTypes returns the attribute types for the origin.
// This function is useful to define the schema of the origin in a consistent manner.
func NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":       types.StringType,
		"interfaces": types.SetType{ElemType: types.StringType},
	}
}

// NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrMap returns the attribute map for a given origin.
// It converts a NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin instance to a map suitable for ObjectValueFrom.
func NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOriginAttrMap(origin *NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin) map[string]attr.Value {
	return map[string]attr.Value{
		"type":       origin.Type,
		"interfaces": origin.Interfaces,
	}
}

// ToAPIPayload converts the Terraform origin into the API origin payload.
func (o *NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin, diag.Diagnostics) {
	tflog.Info(ctx, "[start] NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin ToAPIPayload")

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

	tflog.Info(ctx, "[finish] NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin ToAPIPayload")
	return originPayload, nil
}

// FromAPIResponse fills the NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin with data from the API response.
// This method transforms the OpenAPI response into the format expected by the Terraform provider.
func (o *NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin) FromAPIResponse(ctx context.Context, apiOrigin *openApiClient.CreateNetworkAppliancePrefixesDelegatedStaticRequestOrigin) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin FromAPIResponse")
	tflog.Trace(ctx, "NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin FromAPIResponse", map[string]interface{}{
		"apiOrigin": apiOrigin,
	})

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

	tflog.Info(ctx, "[finish] NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin FromAPIResponse")
	return nil
}

type NetworksApplianceVLANsResourceModelMandatoryDhcp struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

func CreateHttpReqPayload(ctx context.Context, data *NetworksApplianceVLANsResourceModel) (openApiClient.CreateNetworkApplianceVlanRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	// Log the received request
	tflog.Info(ctx, "[start] Create HTTP Request Payload Call")
	tflog.Trace(ctx, "Create Request Payload", map[string]interface{}{
		"data": data,
	})

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

		var ipv6 NetworksApplianceVLANsResourceModelIpv6
		diags := data.IPv6.As(ctx, &ipv6, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return *payload, diags
		}

		ipv6Payload, ipv6PayloadErr := ipv6.ToAPIPayload(ctx)
		if ipv6PayloadErr.HasError() {
			return *payload, ipv6PayloadErr
		}

		// Enabled
		ipv6Payload.SetEnabled(ipv6.Enabled.ValueBool())

		// Handle Prefix Assignments
		var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
		for _, prefixAssignmentData := range ipv6.PrefixAssignments.Elements() {

			// Convert the prefixAssignmentData (which is of type attr.Value) to your struct
			var pa NetworksApplianceVLANsResourceModelIpv6PrefixAssignment

			prefixAssignmentDataDiags := tfsdk.ValueAs(ctx, prefixAssignmentData, &pa)
			if prefixAssignmentDataDiags.HasError() {
				return *payload, prefixAssignmentDataDiags
			}

			// Now create your API client struct
			var prefixAssignment openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
			prefixAssignment.SetAutonomous(pa.Autonomous.ValueBool())
			prefixAssignment.SetStaticPrefix(pa.StaticPrefix.ValueString())
			prefixAssignment.SetStaticApplianceIp6(pa.StaticApplianceIp6.ValueString())

			// Assuming 'Origin' is another struct that you need to convert similarly
			var originData NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
			originDiags := pa.Origin.As(ctx, &originData, basetypes.ObjectAsOptions{})
			if originDiags.HasError() {
				return *payload, originDiags
			}

			// Populate originData into the prefixAssignment's Origin field
			origin, originDiags := originData.ToAPIPayload(ctx)
			if originDiags.HasError() {
				return *payload, originDiags
			}

			prefixAssignment.SetOrigin(*origin)

			prefixAssignments = append(prefixAssignments, prefixAssignment)
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

	tflog.Info(ctx, "[finish] Create HTTP Request Payload")
	tflog.Trace(ctx, "Create Request Payload", map[string]interface{}{
		"payload": payload,
	})

	return *payload, nil
}

func CreateHttpResponse(ctx context.Context, data *NetworksApplianceVLANsResourceModel, response *openApiClient.CreateNetworkApplianceVlan201Response) diag.Diagnostics {

	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] CreatePayloadResponse Call")
	tflog.Trace(ctx, "Create Payload Response", map[string]interface{}{
		"response": response,
	})

	// Set to Ids needed for importing resource
	if data.Id.IsUnknown() {
		data.Id = types.StringValue(fmt.Sprintf("%s,%s", data.NetworkId.ValueString(), data.VlanId.String()))
	}

	// VlanId
	if response.HasId() {
		// API returns string, openAPI spec defines int
		idStr := response.GetId()

		// Check if the string is not empty
		if idStr != "" {

			// Convert string to int
			vlanId, err := strconv.Atoi(idStr)
			if err != nil {
				// Handle the error if conversion fails
				resp.AddError("CreateHttpResponse VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s': %v", idStr, err))
			}

			data.VlanId = types.Int64Value(int64(vlanId))
		}
	}

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
		if data.MandatoryDhcp.IsUnknown() {
			data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
				"enabled": types.BoolType,
			})
		}
	}

	if response.HasIpv6() {
		ipv6Instance := NetworksApplianceVLANsResourceModelIpv6{}
		diags := ipv6Instance.FromAPIResponse(ctx, response.Ipv6)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Object, diags := types.ObjectValueFrom(ctx, NetworksApplianceVLANsResourceModelIpv6AttrTypes(), ipv6Instance)
		if diags.HasError() {
			resp.Append(diags...)
		}
		tflog.Warn(ctx, fmt.Sprintf("CreateHttpResponse: %v", ipv6Object.String()))

		data.IPv6 = ipv6Object

		tflog.Warn(ctx, fmt.Sprintf("CREATE PAYLOAD: %v", data.IPv6.String()))

	} else {
		if data.IPv6.IsUnknown() {
			tflog.Info(ctx, fmt.Sprintf("Empty IPv6 response: %v", data.IPv6))

			ipv6Instance := NetworksApplianceVLANsResourceModelIpv6{}
			ipv6Prefixes := []NetworksApplianceVLANsResourceModelIpv6PrefixAssignment{}

			ipv6PrefixesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentAttrTypes()}, ipv6Prefixes)
			if diags.HasError() {
				resp.Append(diags...)
			}

			ipv6Instance.PrefixAssignments = ipv6PrefixesList

			ipv6Object, diags := types.ObjectValueFrom(ctx, NetworksApplianceVLANsResourceModelIpv6AttrTypes(), ipv6Instance)
			if diags.HasError() {
				resp.Append(diags...)
			}

			data.IPv6 = ipv6Object
		}
	}

	tflog.Info(ctx, "[finish] CreateResponsePayloadResponse Call")

	return resp
}

// ReadHttpResponse - used by READ, UPDATE & DELETE funcs
func ReadHttpResponse(ctx context.Context, data *NetworksApplianceVLANsResourceModel, response *openApiClient.GetNetworkApplianceVlans200ResponseInner) diag.Diagnostics {

	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] ReadHttpResponse Call")
	tflog.Trace(ctx, "Read Response Payload ", map[string]interface{}{
		"response": response,
	})

	// Set to Ids needed for importing resource
	if data.Id.IsUnknown() {
		data.Id = types.StringValue(fmt.Sprintf("%s,%s", data.NetworkId.ValueString(), data.VlanId.String()))
	}

	// VlanId
	if response.HasId() {
		// API returns string, openAPI spec defines int
		idStr := response.GetId()

		// Check if the string is not empty
		if idStr != "" {

			// Convert string to int
			vlanId, err := strconv.Atoi(idStr)
			if err != nil {
				// Handle the error if conversion fails
				resp.AddError("CreateHttpResponse VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s': %v", idStr, err))
			}

			data.VlanId = types.Int64Value(int64(vlanId))
		}
	}

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
	} else {
		if data.DhcpRelayServerIps.IsUnknown() {
			data.DhcpRelayServerIps = basetypes.NewListNull(types.StringType)
		}
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
			data.DhcpBootOptionsEnabled = types.BoolNull()
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

			fixedIpAssignmentValue, fixedIpAssignmentsDiags := types.ObjectValueFrom(ctx, fixedIpAssignmentAttrTypes, fixedIpAssignmentObject)
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
			// Construct the map for the current ObjectValue
			valuesMap, valuesMapErrs := basetypes.NewObjectValue(dhcpOptionsAttrTypes, map[string]attr.Value{
				"code":  basetypes.NewStringValue(dhcpOption.GetCode()),
				"type":  basetypes.NewStringValue(dhcpOption.GetType()),
				"value": basetypes.NewStringValue(dhcpOption.GetValue()),
			})

			if valuesMapErrs.HasError() {
				for _, valuesMapErr := range valuesMapErrs.Errors() {
					tflog.Error(ctx, valuesMapErr.Summary()+valuesMapErr.Detail())
				}
				resp.Append(valuesMapErrs...)

				continue
			}

			// Create the ObjectValue for the current DHCP option
			dhcpOptionValue, dhcpOptionsDiags := types.ObjectValueFrom(ctx, dhcpOptionsAttrTypes, valuesMap)
			if dhcpOptionsDiags.HasError() {
				for _, dhcpOptionsDiag := range dhcpOptionsDiags.Errors() {
					tflog.Error(ctx, dhcpOptionsDiag.Summary()+dhcpOptionsDiag.Detail())
				}
				resp.Append(dhcpOptionsDiags...)

				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, dhcpOptionValue)
		}

		// Create a ListValue from the slice of ObjectValue
		dhcpOptionsListValue, dhcpOptionsListDiags := types.ListValueFrom(ctx, objectType, objectValues)
		if dhcpOptionsListDiags.HasError() {
			resp.Append(dhcpOptionsListDiags...)
		}

		data.DhcpOptions = dhcpOptionsListValue
	} else {
		if data.DhcpOptions.IsUnknown() {
			data.DhcpOptions = types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"code":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			}})
		}
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
		if data.MandatoryDhcp.IsUnknown() {
			data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
				"enabled": types.BoolType,
			})
		}
	}

	// IPv6
	// Assuming response is a structure containing your data
	if response.HasIpv6() {
		ipv6Instance := NetworksApplianceVLANsResourceModelIpv6{}
		diags := ipv6Instance.FromAPIResponse(ctx, response.Ipv6)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Object, diags := types.ObjectValueFrom(ctx, NetworksApplianceVLANsResourceModelIpv6AttrTypes(), ipv6Instance)

		if diags.HasError() {
			resp.Append(diags...)
		}

		data.IPv6 = ipv6Object
	} else {
		if data.IPv6.IsUnknown() {
			tflog.Info(ctx, fmt.Sprintf("Empty IPv6 response: %v", data.IPv6))

			ipv6Instance := NetworksApplianceVLANsResourceModelIpv6{}
			ipv6Prefixes := []NetworksApplianceVLANsResourceModelIpv6PrefixAssignment{}

			ipv6PrefixesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentAttrTypes()}, ipv6Prefixes)
			if diags.HasError() {
				resp.Append(diags...)
			}

			ipv6Instance.PrefixAssignments = ipv6PrefixesList

			ipv6Object, diags := types.ObjectValueFrom(ctx, NetworksApplianceVLANsResourceModelIpv6AttrTypes(), ipv6Instance)
			if diags.HasError() {
				resp.Append(diags...)
			}

			data.IPv6 = ipv6Object
		}
	}

	tflog.Info(ctx, "[finish] ReadHttpResponse Call")

	return resp
}

func UpdateHttpReqPayload(ctx context.Context, data *NetworksApplianceVLANsResourceModel) (*openApiClient.UpdateNetworkApplianceVlanRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] UpdateHttpReqPayload Call")

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

		var fixedIpAssignments map[string]FixedIpAssignmentTerraform
		fixedIpAssignmentsDiags := data.FixedIpAssignments.ElementsAs(ctx, &fixedIpAssignments, true)
		if fixedIpAssignmentsDiags.HasError() {
			resp.AddError(
				"Create Payload Failure, FixedIpAssignments", fmt.Sprintf("%v", fixedIpAssignmentsDiags),
			)
		}

		fixedIpAssignmentsIf := map[string]interface{}{}
		for key, val := range fixedIpAssignments {
			apiAttrs := FixedIpAssignment{}
			if !val.IP.IsUnknown() {
				val := val.IP.ValueString()
				apiAttrs.IP = val
			}

			if !val.Name.IsUnknown() {
				apiAttrs.Name = val.Name.ValueString()
			}

			fixedIpAssignmentsIf[key] = apiAttrs
		}

		payload.SetFixedIpAssignments(fixedIpAssignmentsIf)
	}

	// ReservedIpRanges
	if !data.ReservedIpRanges.IsUnknown() && !data.ReservedIpRanges.IsNull() {

		var reservedIpRanges []openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner

		for _, reservedIpRange := range data.ReservedIpRanges.Elements() {

			var reservedIpRangePayload openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner
			var reservedIpRangeData NetworksApplianceVLANsResourceModelReservedIpRange

			// Convert dhcpOption (types.Object) to dhcpOptionData (struct)
			reservedIpRangeValue, reservedIpRangeDiags := reservedIpRange.ToTerraformValue(ctx)
			if reservedIpRangeDiags != nil {
				// Handle errors during conversion
				resp.AddError(
					"Error converting reservedIpRange",
					reservedIpRangeDiags.Error(),
				)
				continue
			}

			err := reservedIpRangeData.FromTerraformValue(ctx, reservedIpRangeValue)
			if err != nil {
				// Handle errors during conversion
				resp.AddError("Error converting reservedIpRange Value", err.Error())
				tflog.Warn(ctx, err.Error())
				continue
			}

			// Set Payload
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

		var dhcpOptionsPayload []openApiClient.GetNetworkApplianceVlans200ResponseInnerDhcpOptionsInner

		for _, dhcpOption := range data.DhcpOptions.Elements() {
			var dhcpOptionPayload openApiClient.GetNetworkApplianceVlans200ResponseInnerDhcpOptionsInner
			var dhcpOptionData NetworksApplianceVLANsResourceModelDhcpOption

			// Convert dhcpOption (types.Object) to dhcpOptionData (struct)
			dhcpOptionValue, dhcpOptionDiags := dhcpOption.ToTerraformValue(ctx)
			if dhcpOptionDiags != nil {
				// Handle errors during conversion
				resp.AddError(
					"Error converting DHCP option",
					dhcpOptionDiags.Error(),
				)
				continue
			}

			err := dhcpOptionData.FromTerraformValue(ctx, dhcpOptionValue)
			if err != nil {
				// Handle errors during conversion
				resp.AddError("Error converting DHCP option", err.Error())
				tflog.Warn(ctx, err.Error())
				continue
			}

			// Set Payload
			dhcpOptionPayload.SetCode(dhcpOptionData.Code.ValueString())
			dhcpOptionPayload.SetType(dhcpOptionData.Type.ValueString())
			dhcpOptionPayload.SetValue(dhcpOptionData.Value.ValueString())

			dhcpOptionsPayload = append(dhcpOptionsPayload, dhcpOptionPayload)

		}

		payload.SetDhcpOptions(dhcpOptionsPayload)
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

	// IPV6
	if !data.IPv6.IsUnknown() && !data.IPv6.IsNull() {

		var ipv6 NetworksApplianceVLANsResourceModelIpv6
		diags := data.IPv6.As(ctx, &ipv6, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		ipv6Payload, ipv6PayloadErr := ipv6.ToAPIPayload(ctx)
		if ipv6PayloadErr.HasError() {
			return nil, ipv6PayloadErr
		}

		// Enabled
		ipv6Payload.SetEnabled(ipv6.Enabled.ValueBool())

		// Handle Prefix Assignments
		var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
		for _, prefixAssignmentData := range ipv6.PrefixAssignments.Elements() {
			// Convert the prefixAssignmentData (which is of type attr.Value) to your struct
			var pa NetworksApplianceVLANsResourceModelIpv6PrefixAssignment

			prefixAssignmentDataDiags := tfsdk.ValueAs(ctx, prefixAssignmentData, &pa)
			if prefixAssignmentDataDiags.HasError() {
				return nil, prefixAssignmentDataDiags
			}

			// Now create your API client struct
			var prefixAssignment openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
			prefixAssignment.SetAutonomous(pa.Autonomous.ValueBool())
			prefixAssignment.SetStaticPrefix(pa.StaticPrefix.ValueString())
			prefixAssignment.SetStaticApplianceIp6(pa.StaticApplianceIp6.ValueString())

			// Assuming 'Origin' is another struct that you need to convert similarly
			var originData NetworksApplianceVLANsResourceModelIpv6PrefixAssignmentOrigin
			originDiags := pa.Origin.As(ctx, &originData, basetypes.ObjectAsOptions{})
			if originDiags.HasError() {
				return nil, originDiags
			}

			// Populate originData into the prefixAssignment's Origin field
			origin, originDiags := originData.ToAPIPayload(ctx)
			if originDiags.HasError() {
				return nil, originDiags
			}

			prefixAssignment.SetOrigin(*origin)

			prefixAssignments = append(prefixAssignments, prefixAssignment)
		}

		ipv6Payload.SetPrefixAssignments(prefixAssignments)

		payload.SetIpv6(*ipv6Payload)
	}

	tflog.Info(ctx, "[finish] UpdateHttpReqPayload Call")
	tflog.Trace(ctx, "Update Request Payload", map[string]interface{}{
		"payload": payload,
	})

	return payload, nil
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

	// Log the received request
	tflog.Info(ctx, "[start] CREATE Function Call")
	tflog.Trace(ctx, "Create Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initial create API call
	payload, payloadReqDiags := CreateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlanRequest(payload).Execute()

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

	payloadRespDiags := CreateHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Update to capture config items not accessible in HTTP POST
	updatePayload, updatePayloadReqDiags := UpdateHttpReqPayload(ctx, data)
	if updatePayloadReqDiags != nil {
		resp.Diagnostics.Append(updatePayloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns this as string, openAPI spec has set as Integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	updateInlineResp, updateHttpResp, updateErr := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).UpdateNetworkApplianceVlanRequest(*updatePayload).Execute()
	if updateErr != nil && updateHttpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
			tools.HttpDiagnostics(updateHttpResp),
		)
		return
	}

	// Check for API success response code
	if updateHttpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", updateHttpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	updatePayloadRespDiags := ReadHttpResponse(ctx, data, updateInlineResp)
	if updatePayloadRespDiags != nil {
		resp.Diagnostics.Append(updatePayloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] CREATE Function Call")
	tflog.Trace(ctx, "Create function completed", map[string]interface{}{
		"data": data,
	})
}

func (r *NetworksApplianceVLANsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Log the received request
	tflog.Info(ctx, "[start] READ Function Call")
	tflog.Trace(ctx, "Read Function Called", map[string]interface{}{
		"request": req,
	})

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

	payloadRespDiags := ReadHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] READ Function Call")
	tflog.Trace(ctx, "Read Function", map[string]interface{}{
		"data": data,
	})
}

func (r *NetworksApplianceVLANsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Log the received request
	tflog.Info(ctx, "[start] UPDATE Function Call")
	tflog.Trace(ctx, "Update Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadReqDiags := UpdateHttpReqPayload(ctx, data)
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

	payloadRespDiags := ReadHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] UPDATE Function Call")
	tflog.Trace(ctx, "Update Function", map[string]interface{}{
		"data": data,
	})
}

func (r *NetworksApplianceVLANsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Log the received request
	tflog.Info(ctx, "[start] DELETE Function Call")
	tflog.Trace(ctx, "Delete Function Called", map[string]interface{}{
		"request": req,
	})

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

	// Log the response data
	tflog.Info(ctx, "[finish] DELETE Function Call")
	tflog.Trace(ctx, "Delete Function", map[string]interface{}{
		"data": data,
	})

}

func (r *NetworksApplianceVLANsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, vlan_id. Got: %q", req.ID),
		)
		return
	}

	// ensure vlanId is formatted properly
	str := idParts[1]

	// Convert the string to int64
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert vlanId to integer",
			fmt.Sprintf("Expected import identifier with format: network_id, vlan_id. Got: %q", req.ID),
		)
	}

	// Convert the int64 to types.Int64Value
	vlanId := types.Int64Value(i)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vlan_id"), vlanId)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
