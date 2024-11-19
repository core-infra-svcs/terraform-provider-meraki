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

// keyModel represents the API key attributes.
type keyModel struct {
	Created types.Bool `tfsdk:"created" json:"created"`
}

// keyType defines the attribute types for KeyModel.
var keyType = map[string]attr.Type{
	"created": types.BoolType,
}

// apiModel represents API-related attributes.
type apiModel struct {
	Key types.Object `tfsdk:"key" json:"key"`
}

// apiType defines the attribute types for APIModel.
var apiType = map[string]attr.Type{
	"key": types.ObjectType{AttrTypes: keyType},
}

// SAML Attribute Model
type samlModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// SAML Attribute Type
var samlType = map[string]attr.Type{
	"enabled": types.BoolType,
}

// Two-Factor Authentication Attribute Model
type twoFactorModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// Two-Factor Authentication Attribute Type
var twoFactorType = map[string]attr.Type{
	"enabled": types.BoolType,
}

// authenticationModel represents authentication attributes for a user.
type authenticationModel struct {
	Mode      types.String `tfsdk:"mode" json:"mode"`
	API       types.Object `tfsdk:"api" json:"api"`
	SAML      types.Object `tfsdk:"saml" json:"saml"`
	TwoFactor types.Object `tfsdk:"two_factor" json:"twoFactor"`
}

// authenticationType defines the attribute types for AuthenticationModel.
var authenticationType = map[string]attr.Type{
	"mode":       types.StringType,
	"api":        types.ObjectType{AttrTypes: apiType},
	"saml":       types.ObjectType{AttrTypes: samlType},
	"two_factor": types.ObjectType{AttrTypes: twoFactorType},
}

// dataSourceModel represents the main data model for identities.
type dataSourceModel struct {
	Id                  types.String `tfsdk:"id" json:"-"`
	Name                types.String `tfsdk:"name" json:"name"`
	Email               types.String `tfsdk:"email" json:"email"`
	LastUsedDashboardAt types.String `tfsdk:"last_used_dashboard_at" json:"lastUsedDashboardAt"`
	Authentication      types.Object `tfsdk:"authentication" json:"authentication"`
}

// dataSourceType defines the attribute types for IdentitiesMeModel.
var dataSourceType = map[string]attr.Type{
	"id":                     types.StringType,
	"name":                   types.StringType,
	"email":                  types.StringType,
	"last_used_dashboard_at": types.StringType,
	"authentication":         types.ObjectType{AttrTypes: authenticationType},
}
