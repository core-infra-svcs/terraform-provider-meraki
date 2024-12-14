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

// Attribute key constants
const (
	DatasourceIdentitiesMeKey               = "identities_me"
	DatasourceIdentitiesMeAuthenticationKey = "authentication"
	DatasourceIdentitiesMeAPIKey            = "api"
	DatasourceIdentitiesMeSAMLKey           = "saml"
	DatasourceIdentitiesMeTwoFactorKey      = "two_factor"
)

// KeyModel represents the API key attributes.
type KeyModel struct {
	Created types.Bool `tfsdk:"created" json:"created"`
}

// KeyType defines the attribute types for KeyModel.
var KeyType = map[string]attr.Type{
	"created": types.BoolType,
}

// APIModel represents API-related attributes.
type APIModel struct {
	Key types.Object `tfsdk:"key" json:"key"`
}

// APIType defines the attribute types for APIModel.
var APIType = map[string]attr.Type{
	"key": types.ObjectType{AttrTypes: KeyType},
}

// SAMLModel represents SAML authentication attributes.
type SAMLModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// SAMLType defines the attribute types for SAMLModel.
var SAMLType = map[string]attr.Type{
	"enabled": types.BoolType,
}

// TwoFactorModel represents two-factor authentication attributes.
type TwoFactorModel struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// TwoFactorType defines the attribute types for TwoFactorModel.
var TwoFactorType = map[string]attr.Type{
	"enabled": types.BoolType,
}

// AuthenticationModel represents authentication attributes for a user.
type AuthenticationModel struct {
	Mode      types.String `tfsdk:"mode" json:"mode"`
	API       types.Object `tfsdk:"api" json:"api"`
	SAML      types.Object `tfsdk:"saml" json:"saml"`
	TwoFactor types.Object `tfsdk:"two_factor" json:"twoFactor"`
}

// AuthenticationType defines the attribute types for AuthenticationModel.
var AuthenticationType = map[string]attr.Type{
	"mode":       types.StringType,
	"api":        types.ObjectType{AttrTypes: APIType},
	"saml":       types.ObjectType{AttrTypes: SAMLType},
	"two_factor": types.ObjectType{AttrTypes: TwoFactorType},
}

// DataSourceModel represents the main data model for identities.
type DataSourceModel struct {
	Id                  types.String `tfsdk:"id" json:"-"`                                       // Unique identifier for the datasource.
	Name                types.String `tfsdk:"name" json:"name"`                                  // The user's name.
	Email               types.String `tfsdk:"email" json:"email"`                                // The user's email address.
	LastUsedDashboardAt types.String `tfsdk:"last_used_dashboard_at" json:"lastUsedDashboardAt"` // Last dashboard access time.
	Authentication      types.Object `tfsdk:"authentication" json:"authentication"`              // Authentication details.
}

// DataSourceType defines the attribute types for DataSourceModel.
var DataSourceType = map[string]attr.Type{
	"id":                     types.StringType,
	"name":                   types.StringType,
	"email":                  types.StringType,
	"last_used_dashboard_at": types.StringType,
	"authentication":         types.ObjectType{AttrTypes: AuthenticationType},
}
