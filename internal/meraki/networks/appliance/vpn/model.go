package vpn

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resourceModel describes the resource data model.
type resourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id" json:"network_id"`
	Mode      jsontypes.String `tfsdk:"mode" json:"mode"`
	Hubs      types.List       `tfsdk:"hubs" json:"hubs"`
	Subnets   types.List       `tfsdk:"subnets" json:"subnets"`
}

type resourceModelHubs struct {
	HubId           jsontypes.String `tfsdk:"hub_id" json:"hubId"`
	UseDefaultRoute jsontypes.Bool   `tfsdk:"use_default_route" json:"useDefaultRoute"`
}

type resourceModelSubnets struct {
	LocalSubnet jsontypes.String `tfsdk:"local_subnet" json:"localSubnet"`
	UseVpn      jsontypes.Bool   `tfsdk:"use_vpn" json:"useVpn"`
}

// datasourceModel resourceModel describes the resource data model.
type datasourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id" json:"network_id"`
	Mode      jsontypes.String `tfsdk:"mode" json:"mode"`
	Hubs      types.List       `tfsdk:"hubs" json:"hubs"`
	Subnets   types.List       `tfsdk:"subnets" json:"subnets"`
}

/*
type datasourceModelHubs struct {
	HubId           jsontypes.String `tfsdk:"hub_id" json:"hubId"`
	UseDefaultRoute jsontypes.Bool   `tfsdk:"use_default_route" json:"useDefaultRoute"`
}

type datasourceModelSubnets struct {
	LocalSubnet jsontypes.String `tfsdk:"local_subnet" json:"localSubnet"`
	UseVpn      jsontypes.Bool   `tfsdk:"use_vpn" json:"useVpn"`
}

*/
