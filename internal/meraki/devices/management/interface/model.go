package _interface

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*

// Sample API Response v1.52.0

{
  "ddnsHostnames": {
    "activeDdnsHostname": "mx1-sample.dynamic-m.com",
    "ddnsHostnameWan1": "mx1-sample-1.dynamic-m.com",
    "ddnsHostnameWan2": "mx1-sample-2.dynamic-m.com"
  },
  "wan1": {
    "wanEnabled": "not configured",
    "usingStaticIp": true,
    "staticIp": "1.2.3.4",
    "staticSubnetMask": "255.255.255.0",
    "staticGatewayIp": "1.2.3.1",
    "staticDns": [
      "1.2.3.2",
      "1.2.3.3"
    ],
    "vlan": 7
  },
  "wan2": {
    "wanEnabled": "enabled",
    "usingStaticIp": false,
    "staticIp": "1.2.3.4",
    "staticSubnetMask": "255.255.255.0",
    "staticGatewayIp": "1.2.3.1",
    "staticDns": [
      "1.2.3.2",
      "1.2.3.3"
    ],
    "vlan": 2
  }
}

*/

// DdnsHostnamesModel represents the attributes for DDNSHostnames configuration.
type DdnsHostnamesModel struct {
	ActiveDdnsHostname types.String `tfsdk:"active_ddns_hostname" json:"activeDdnsHostname"`
	DdnsHostnameWan1   types.String `tfsdk:"ddns_hostname_wan1" json:"ddnsHostnameWan1"`
	DdnsHostnameWan2   types.String `tfsdk:"ddns_hostname_wan2" json:"ddnsHostnameWan2"`
}

// DdnsHostnamesType defines the attribute types for DDNSHostnamesModel.
var DdnsHostnamesType = map[string]attr.Type{
	"active_ddns_hostname": types.StringType,
	"ddns_hostname_wan1":   types.StringType,
	"ddns_hostname_wan2":   types.StringType,
}

// WANModel represents the attributes for WAN configuration.
type WANModel struct {
	WanEnabled       types.String `tfsdk:"wan_enabled" json:"wanEnabled"`
	UsingStaticIp    types.Bool   `tfsdk:"using_static_ip" json:"usingStaticIp"`
	StaticIp         types.String `tfsdk:"static_ip" json:"staticIp"`
	StaticSubnetMask types.String `tfsdk:"static_subnet_mask" json:"staticSubnetMask"`
	StaticGatewayIp  types.String `tfsdk:"static_gateway_ip" json:"staticGatewayIp"`
	StaticDns        types.List   `tfsdk:"static_dns" json:"staticDns"` // List of strings
	Vlan             types.Int64  `tfsdk:"vlan" json:"vlan"`
}

// WANType defines the attribute types for WANModel.
var WANType = map[string]attr.Type{
	"wan_enabled":        types.StringType,
	"using_static_ip":    types.BoolType,
	"static_ip":          types.StringType,
	"static_subnet_mask": types.StringType,
	"static_gateway_ip":  types.StringType,
	"static_dns":         types.ListType{ElemType: types.StringType},
	"vlan":               types.Int64Type,
}

// ResourceModel represents the resource's main data model.
type ResourceModel struct {
	Id            types.String `tfsdk:"id" json:"-"`
	Serial        types.String `tfsdk:"serial" json:"serial"`
	DDNSHostnames types.Object `tfsdk:"ddns_hostnames" json:"ddnsHostnames"`
	Wan1          types.Object `tfsdk:"wan1" json:"wan1"`
	Wan2          types.Object `tfsdk:"wan2" json:"wan2"`
}

// DataSourceModel represents the data source's main data model.
type DataSourceModel ResourceModel
