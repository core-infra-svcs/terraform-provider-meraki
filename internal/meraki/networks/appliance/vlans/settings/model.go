package settings

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// resourceModel describes the resource data model.
type resourceModel struct {
	Id           jsontypes.String `tfsdk:"id"`
	NetworkId    jsontypes.String `tfsdk:"network_id" json:"network_id"`
	VlansEnabled jsontypes.Bool   `tfsdk:"vlans_enabled"`
}

// datasourceModel resourceModel describes the resource data model.
type datasourceModel struct {
	Id           jsontypes.String `tfsdk:"id"`
	NetworkId    jsontypes.String `tfsdk:"network_id" json:"network_id"`
	VlansEnabled jsontypes.Bool   `tfsdk:"vlans_enabled"  json:"vlansEnabled"`
}
