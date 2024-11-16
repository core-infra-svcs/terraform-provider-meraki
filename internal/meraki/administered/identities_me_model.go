package administered

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*

// Sample API Response v1.52.0
{
  "name": "Miles Meraki",
  "email": "miles@meraki.com",
  "lastUsedDashboardAt": "2018-02-11T00:00:00.090210Z",
  "authentication": {
    "mode": "email",
    "api": {
      "key": {
        "created": true
      }
    },
    "twoFactor": {
      "enabled": false
    },
    "saml": {
      "enabled": false
    }
  }
}

*/

// Data Models

type keyAttrModel struct {
	Created types.Bool `tfsdk:"created" json:"created"`
}

var keyAttrType = map[string]attr.Type{
	"created": types.BoolType,
}

type apiAttrModel struct {
	Key types.Object `tfsdk:"key" json:"key"`
}

var apiAttrType = map[string]attr.Type{
	"key": types.ObjectType{AttrTypes: keyAttrType},
}

type samlAttrModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

var samlAttrType = map[string]attr.Type{
	"enabled": types.BoolType,
}

type twofactorAttrModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

var twoFactorAttrType = map[string]attr.Type{
	"enabled": types.BoolType,
}

type authenticationAttrModel struct {
	ApiKeyCreated types.Bool   `tfsdk:"api_key_created" json:"api_key_created"`
	Mode          types.String `tfsdk:"mode" json:"mode"`
	Api           types.Object `tfsdk:"api" json:"api"`
	Saml          types.Object `tfsdk:"saml_enabled" json:"saml"`
	Twofactor     types.Object `tfsdk:"two_factor" json:"twoFactor"`
}

var authenticationAttrType = map[string]attr.Type{
	"api_key_created": types.BoolType,
	"mode":            types.StringType,
	"api":             types.ObjectType{AttrTypes: apiAttrType},
	"saml_enabled":    types.ObjectType{AttrTypes: samlAttrType},
	"two_factor":      types.ObjectType{AttrTypes: twoFactorAttrType},
}

type identitiesMeAttrModel struct {
	Name                types.String `tfsdk:"name" json:"name"`
	Email               types.String `tfsdk:"email" json:"email"`
	LastUsedDashboardAt types.String `tfsdk:"last_used_dashboard_at" json:"lastUsedDashboardAt"`
	Authentication      types.Object `tfsdk:"authentication" json:"authentication"`
}

var identitiesMeAttrType = map[string]attr.Type{
	"name":                   types.StringType,
	"email":                  types.StringType,
	"last_used_dashboard_at": types.StringType,
	"authentication":         types.ObjectType{AttrTypes: authenticationAttrType},
}
