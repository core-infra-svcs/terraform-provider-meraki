package mtu

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// dataSourceModel describes the resource data model.
type dataSourceModel struct {
	Id             jsontypes.String          `tfsdk:"id"`
	NetworkId      jsontypes.String          `tfsdk:"network_id" json:"network_id"`
	DefaultMtuSize jsontypes.Int64           `tfsdk:"default_mtu_size" json:"defaultMtuSize"`
	Overrides      []dataSourceModelOverride `tfsdk:"overrides" json:"overrides"`
}

type dataSourceModelOverride struct {
	Switches       []string        `tfsdk:"switches" json:"switches"`
	SwitchProfiles []string        `tfsdk:"switch_profiles" json:"switchProfiles"`
	MtuSize        jsontypes.Int64 `tfsdk:"mtu_size" json:"mtuSize"`
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id             jsontypes.String        `tfsdk:"id"`
	NetworkId      jsontypes.String        `tfsdk:"network_id" json:"network_id"`
	DefaultMtuSize jsontypes.Int64         `tfsdk:"default_mtu_size" json:"defaultMtuSize"`
	Overrides      []resourceModelOverride `tfsdk:"overrides" json:"overrides"`
}

type resourceModelOverride struct {
	Switches       []string        `tfsdk:"switches" json:"switches"`
	SwitchProfiles []string        `tfsdk:"switch_profiles" json:"switchProfiles"`
	MtuSize        jsontypes.Int64 `tfsdk:"mtu_size" json:"mtuSize"`
}
