package cellular

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*

// Sample API Response v1.52.0

{
  "sims": [
    {
      "slot": "sim1",
      "isPrimary": false,
      "apns": [
        {
          "name": "internet",
          "allowedIpTypes": [
            "ipv4",
            "ipv6"
          ],
          "authentication": {
            "type": "pap",
            "username": "milesmeraki",
            "password": "secret"
          }
        }
      ]
    }
  ],
  "simOrdering": [
    "sim1",
    "sim2",
    "sim3"
  ],
  "simFailover": {
    "enabled": true,
    "timeout": 300
  }
}


*/

// resourceModel represents the top-level resource structure.
type resourceModel struct {
	Id          types.String `tfsdk:"id" json:"-"`
	Serial      types.String `tfsdk:"serial" json:"serial"`
	Sims        types.Set    `tfsdk:"sims" json:"sims"` // Updated to match schema
	SimFailOver types.Object `tfsdk:"sim_failover" json:"simFailover"`
	SimOrdering types.Set    `tfsdk:"sim_ordering" json:"simOrdering"`
}

// resourceModelAttrTypes defines the attribute types for the top-level resourceModel.
func resourceModelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"serial":       types.StringType,
		"sims":         types.ListType{ElemType: types.ObjectType{AttrTypes: ResourceModelSimAttrTypes()}},
		"sim_failover": types.ObjectType{AttrTypes: ResourceModelSimFailOverAttrTypes()},
		"sim_ordering": types.SetType{ElemType: types.StringType},
	}
}

// ResourceModelSim represents an individual SIM configuration.
type ResourceModelSim struct {
	Slot      types.String `tfsdk:"slot" json:"slot"`
	IsPrimary types.Bool   `tfsdk:"is_primary" json:"isPrimary"`
	Apns      types.Set    `tfsdk:"apns" json:"apns"` // Updated to match schema
}

// ResourceModelSimAttrTypes defines the attribute types for ResourceModelSim.
func ResourceModelSimAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"slot":       types.StringType,
		"is_primary": types.BoolType,
		"apns":       types.SetType{ElemType: types.ObjectType{AttrTypes: ResourceModelApnsAttrTypes()}},
	}
}

// ResourceModelApns represents an APN configuration for a SIM.
type ResourceModelApns struct {
	Name           types.String `tfsdk:"name" json:"name"`
	AllowedIpTypes types.Set    `tfsdk:"allowed_ip_types" json:"allowedIpTypes"`
	Authentication types.Object `tfsdk:"authentication" json:"authentication"`
}

// ResourceModelApnsAttrTypes defines the attribute types for ResourceModelApns.
func ResourceModelApnsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":             types.StringType,
		"allowed_ip_types": types.SetType{ElemType: types.StringType},
		"authentication":   types.ObjectType{AttrTypes: ResourceModelAuthenticationAttrTypes()},
	}
}

// ResourceModelAuthentication represents APN authentication details.
type ResourceModelAuthentication struct {
	Password types.String `tfsdk:"password" json:"password"`
	Username types.String `tfsdk:"username" json:"username"`
	Type     types.String `tfsdk:"type" json:"type"`
}

// ResourceModelAuthenticationAttrTypes defines the attribute types for ResourceModelAuthentication.
func ResourceModelAuthenticationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"password": types.StringType,
		"username": types.StringType,
		"type":     types.StringType,
	}
}

// ResourceModelSimFailOver represents SIM failover configuration.
type ResourceModelSimFailOver struct {
	Enabled types.Bool  `tfsdk:"enabled" json:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout" json:"timeout"`
}

// ResourceModelSimFailOverAttrTypes defines the attribute types for ResourceModelSimFailOver.
func ResourceModelSimFailOverAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"timeout": types.Int64Type,
	}
}
