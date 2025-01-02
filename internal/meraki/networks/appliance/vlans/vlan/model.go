package vlan

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworksApplianceVLANModelReservedIpRange struct {
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

type NetworksApplianceVLANResourceModelDhcpOption struct {
	Code  types.String `tfsdk:"code" json:"code"`
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

// NetworksApplianceVLANModelIpv6 represents the IPv6 configuration for a VLAN resource model.
// For both the vlan_resource and vlan_datasource
type NetworksApplianceVLANModelIpv6 struct {
	Enabled           types.Bool `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments types.List `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

// Ipv6AttrTypes is Used in both the vlan_resource and vlan_datasource
func Ipv6AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":            types.BoolType,
		"prefix_assignments": types.ListType{ElemType: types.ObjectType{AttrTypes: Ipv6PrefixAssignmentAttrTypes()}},
	}
}

// Ipv6PrefixAssignmentAttrTypes returns the attribute types for a prefix assignment.
func Ipv6PrefixAssignmentAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"autonomous":           types.BoolType,
		"static_prefix":        types.StringType,
		"static_appliance_ip6": types.StringType,
		"origin":               types.ObjectType{AttrTypes: Ipv6PrefixAssignmentOriginAttrTypes()},
	}
}

// Ipv6PrefixAssignment represents a prefix assignment for an IPv6 configuration in the VLAN resource model.
type Ipv6PrefixAssignment struct {
	Autonomous         types.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       types.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 types.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             types.Object `tfsdk:"origin" json:"origin"`
}

// Ipv6PrefixAssignmentOrigin represents the origin data structure for a VLAN resource model.
type Ipv6PrefixAssignmentOrigin struct {
	Type       types.String `tfsdk:"type" json:"type"`
	Interfaces types.Set    `tfsdk:"interfaces" json:"interfaces"`
}

// Ipv6PrefixAssignmentOriginAttrTypes returns the attribute types for the origin.
// This function is useful to define the schema of the origin in a consistent manner.
func Ipv6PrefixAssignmentOriginAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":       types.StringType,
		"interfaces": types.SetType{ElemType: types.StringType},
	}
}

// Ipv6PrefixAssignmentOriginAttrMap returns the attribute map for a given origin.
// It converts a Ipv6PrefixAssignmentOrigin instance to a map suitable for ObjectValueFrom.
func Ipv6PrefixAssignmentOriginAttrMap(origin *Ipv6PrefixAssignmentOrigin) map[string]attr.Value {
	return map[string]attr.Value{
		"type":       origin.Type,
		"interfaces": origin.Interfaces,
	}
}

type NetworksApplianceVLANModelMandatoryDhcp struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// datasourceModel NetworksApplianceVLANModel describes the resource data model.
type datasourceModel struct {
	Id        types.String                 `tfsdk:"id" json:"-"`
	NetworkId types.String                 `tfsdk:"network_id" json:"network_id"`
	List      []NetworksApplianceVLANModel `tfsdk:"list"`
}

// NetworksApplianceVLANModel is Used in both the vlan_resource and vlan_datasource
type NetworksApplianceVLANModel struct {
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
