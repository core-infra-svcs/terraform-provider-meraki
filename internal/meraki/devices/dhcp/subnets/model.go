package subnets

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*

// Sample API Response v1.52.0

[
  {
    "subnet": "192.168.1.0/24",
    "vlanId": 100,
    "usedCount": 2,
    "freeCount": 251
  }
]
*/

// dataSourceModel represents the top-level data source structure.
type dataSourceModel struct {
	Id        types.String `tfsdk:"id" json:"id"`
	Serial    types.String `tfsdk:"serial" json:"serial"`
	Resources types.List   `tfsdk:"resources" json:"-"`
}

// dataSourceAttrTypes defines the attribute types for the data source model.
func dataSourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"serial":    types.StringType,
		"resources": types.ListType{ElemType: types.ObjectType{AttrTypes: resourceAttrTypes()}},
	}
}

// resourceModel represents an individual DHCP subnet configuration.
type resourceModel struct {
	Id        types.String `tfsdk:"id" json:"-"`
	Subnet    types.String `tfsdk:"subnet" json:"subnet"`
	VlanId    types.Int64  `tfsdk:"vlan_id" json:"vlanId"`
	UsedCount types.Int64  `tfsdk:"used_count" json:"usedCount"`
	FreeCount types.Int64  `tfsdk:"free_count" json:"freeCount"`
}

// resourceAttrTypes defines the attribute types for an individual resource (subnet).
func resourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"subnet":     types.StringType,
		"vlan_id":    types.Int64Type,
		"used_count": types.Int64Type,
		"free_count": types.Int64Type,
	}
}
