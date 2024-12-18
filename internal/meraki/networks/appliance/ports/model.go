package ports

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

type dataSourceModel struct {
	Id        jsontypes.String      `tfsdk:"id"`
	NetworkId jsontypes.String      `tfsdk:"network_id"`
	List      []dataSourceListModel `tfsdk:"list"`
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                  jsontypes.String `tfsdk:"id"`
	NetworkId           jsontypes.String `tfsdk:"network_id"`
	PortId              jsontypes.String `tfsdk:"port_id"`
	Accesspolicy        jsontypes.String `tfsdk:"access_policy" json:"access_policy"`
	Allowedvlans        jsontypes.String `tfsdk:"allowed_vlans" json:"allowed_vlans"`
	Dropuntaggedtraffic jsontypes.Bool   `tfsdk:"drop_untagged_traffic" json:"drop_untagged_traffic"`
	Enabled             jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Number              jsontypes.Int64  `tfsdk:"number" json:"number"`
	Type                jsontypes.String `tfsdk:"type" json:"type"`
	Vlan                jsontypes.Int64  `tfsdk:"vlan" json:"vlan"`
}

type dataSourceListModel struct {
	Accesspolicy        jsontypes.String `tfsdk:"access_policy" json:"access_policy"`
	Allowedvlans        jsontypes.String `tfsdk:"allowed_vlans" json:"allowed_vlans"`
	Dropuntaggedtraffic jsontypes.Bool   `tfsdk:"drop_untagged_traffic" json:"drop_untagged_traffic"`
	Enabled             jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Number              jsontypes.Int64  `tfsdk:"number" json:"number"`
	Type                jsontypes.String `tfsdk:"type" json:"type"`
	Vlan                jsontypes.Int64  `tfsdk:"vlan" json:"vlan"`
}
