package admins

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceModel struct {
	Id                   jsontypes.String       `tfsdk:"id"`
	OrgId                jsontypes.String       `tfsdk:"organization_id" json:"organizationId"`
	AdminId              jsontypes.String       `tfsdk:"admin_id" json:"id"`
	Name                 jsontypes.String       `tfsdk:"name"`
	Email                jsontypes.String       `tfsdk:"email"`
	OrgAccess            jsontypes.String       `tfsdk:"org_access" json:"orgAccess"`
	AccountStatus        jsontypes.String       `tfsdk:"account_status" json:"accountStatus"`
	TwoFactorAuthEnabled jsontypes.Bool         `tfsdk:"two_factor_auth_enabled" json:"twoFactorAuthEnabled"`
	HasApiKey            jsontypes.Bool         `tfsdk:"has_api_key" json:"hasApiKey"`
	LastActive           jsontypes.String       `tfsdk:"last_active" json:"lastActive"`
	Tags                 []resourceModelTag     `tfsdk:"tags" json:"tags"`
	Networks             []resourceModelNetwork `tfsdk:"networks" json:"networks"`
	AuthenticationMethod jsontypes.String       `tfsdk:"authentication_method" json:"authenticationMethod"`
}

type resourceModelTag struct {
	Tag    jsontypes.String `tfsdk:"tag" json:"tag"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

type resourceModelNetwork struct {
	Id     jsontypes.String `tfsdk:"id" json:"id"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

type dataSourceModel struct {
	Id             types.String          `tfsdk:"id" json:"-"`
	OrganizationId types.String          `tfsdk:"organization_id" json:"organizationId"`
	List           []dataSourceModelList `tfsdk:"list"`
}

type dataSourceModelList struct {
	Id                   jsontypes.String          `tfsdk:"id" json:"id"`
	Name                 jsontypes.String          `tfsdk:"name"`
	Email                jsontypes.String          `tfsdk:"email"`
	OrgAccess            jsontypes.String          `tfsdk:"org_access" json:"orgAccess"`
	AccountStatus        jsontypes.String          `tfsdk:"account_status" json:"accountStatus"`
	TwoFactorAuthEnabled jsontypes.Bool            `tfsdk:"two_factor_auth_enabled" json:"twoFactorAuthEnabled"`
	HasApiKey            jsontypes.Bool            `tfsdk:"has_api_key" json:"hasApiKey"`
	LastActive           jsontypes.String          `tfsdk:"last_active" json:"lastActive"`
	Tags                 []dataSourceModelTags     `tfsdk:"tags"`
	Networks             []dataSourceModelNetworks `tfsdk:"networks"`
	AuthenticationMethod jsontypes.String          `tfsdk:"authentication_method" json:"authenticationMethod"`
}

type dataSourceModelNetworks struct {
	Id     jsontypes.String `tfsdk:"id"`
	Access jsontypes.String `tfsdk:"access"`
}

type dataSourceModelTags struct {
	Tag    jsontypes.String `tfsdk:"tag"`
	Access jsontypes.String `tfsdk:"access"`
}
