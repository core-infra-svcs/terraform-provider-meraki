package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The resourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type resourceModel struct {
	Id                              types.String `tfsdk:"id"`
	NetworkId                       types.String `tfsdk:"network_id"`
	Number                          types.String `tfsdk:"number"`
	SplashUrl                       types.String `tfsdk:"splash_url"`
	UseSplashUrl                    types.Bool   `tfsdk:"use_splash_url"`
	SplashTimeout                   types.Int64  `tfsdk:"splash_timeout"`
	WelcomeMessage                  types.String `tfsdk:"welcome_message"`
	RedirectUrl                     types.String `tfsdk:"redirect_url"`
	UseRedirectUrl                  types.Bool   `tfsdk:"use_redirect_url"`
	BlockAllTrafficBeforeSignOn     types.Bool   `tfsdk:"block_all_traffic_before_sign_on"`
	ControllerDisconnectionBehavior types.String `tfsdk:"controller_disconnection_behavior"`
	AllowSimultaneousLogins         types.Bool   `tfsdk:"allow_simultaneous_logins"`
	Billing                         types.Object `tfsdk:"billing"`
	GuestSponsorship                types.Object `tfsdk:"guest_sponsorship"`
	SentryEnrollment                types.Object `tfsdk:"sentry_enrollment"`
	SplashImage                     types.Object `tfsdk:"splash_image"`
	SplashLogo                      types.Object `tfsdk:"splash_logo"`
	SplashPrepaidFront              types.Object `tfsdk:"splash_prepaid_front"`
}

type resourceModelBilling struct {
	ReplyToEmailAddress           types.String `tfsdk:"reply_to_email_address"`
	PrepaidAccessFastLoginEnabled types.Bool   `tfsdk:"prepaid_access_fast_login_enabled"`
	FreeAccess                    types.Object `tfsdk:"free_access"`
}

type resourceModelFreeAccess struct {
	DurationInMinutes types.Int64 `tfsdk:"duration_in_minutes"`
	Enabled           types.Bool  `tfsdk:"enabled"`
}

type resourceModelGuestSponsorship struct {
	DurationInMinutes        types.Int64 `tfsdk:"duration_in_minutes"`
	GuestCanRequestTimeframe types.Bool  `tfsdk:"guest_can_request_time_frame"`
}

type resourceModelSplashImage struct {
	Extension types.String `tfsdk:"extension"`
	Md5       types.String `tfsdk:"md5"`
	Image     types.Object `tfsdk:"image"`
}

type resourceModelSplashLogo struct {
	Extension types.String `tfsdk:"extension"`
	Md5       types.String `tfsdk:"md5"`
	Image     types.Object `tfsdk:"image"`
}

type resourceModelSplashPrepaidFront struct {
	Extension types.String `tfsdk:"extension"`
	Md5       types.String `tfsdk:"md5"`
	Image     types.Object `tfsdk:"image"`
}

type resourceModelImage struct {
	Contents types.String `tfsdk:"contents"`
	Format   types.String `tfsdk:"format"`
}

type resourceModelSentryEnrollment struct {
	Strength              types.String `tfsdk:"strength"`
	EnforcedSystems       types.List   `tfsdk:"enforced_systems"`
	SystemsManagerNetwork types.Object `tfsdk:"systems_manager_network"`
}

type resourceModelSystemsManagerNetwork struct {
	Id types.String `tfsdk:"id"`
}

var resourceModelAttrTypes = map[string]attr.Type{
	"id":                                types.StringType,
	"network_id":                        types.StringType,
	"number":                            types.StringType,
	"splash_url":                        types.StringType,
	"use_splash_url":                    types.BoolType,
	"splash_timeout":                    types.Int64Type,
	"welcome_message":                   types.StringType,
	"redirect_url":                      types.StringType,
	"use_redirect_url":                  types.BoolType,
	"block_all_traffic_before_sign_on":  types.BoolType,
	"controller_disconnection_behavior": types.StringType,
	"allow_simultaneous_logins":         types.BoolType,
	"billing":                           types.ObjectType{AttrTypes: resourceModelBillingAttrTypes},
	"guest_sponsorship":                 types.ObjectType{AttrTypes: resourceModelGuestSponsorshipAttrTypes},
	"sentry_enrollment":                 types.ObjectType{AttrTypes: resourceModelSentryEnrollmentAttrTypes},
	"splash_image":                      types.ObjectType{AttrTypes: resourceModelSplashImageAttrTypes},
	"splash_logo":                       types.ObjectType{AttrTypes: resourceModelSplashLogoAttrTypes},
	"splash_prepaid_front":              types.ObjectType{AttrTypes: resourceModelSplashPrepaidFrontAttrTypes},
}

var resourceModelBillingAttrTypes = map[string]attr.Type{
	"reply_to_email_address":            types.StringType,
	"prepaid_access_fast_login_enabled": types.BoolType,
	"free_access":                       types.ObjectType{AttrTypes: resourceModelFreeAccessAttrTypes},
}

var resourceModelFreeAccessAttrTypes = map[string]attr.Type{
	"duration_in_minutes": types.Int64Type,
	"enabled":             types.BoolType,
}

var resourceModelGuestSponsorshipAttrTypes = map[string]attr.Type{
	"duration_in_minutes":          types.Int64Type,
	"guest_can_request_time_frame": types.BoolType,
}

var resourceModelSplashImageAttrTypes = map[string]attr.Type{
	"extension": types.StringType,
	"md5":       types.StringType,
	"image":     types.ObjectType{AttrTypes: resourceModelImageAttrTypes},
}

var resourceModelSplashLogoAttrTypes = map[string]attr.Type{
	"extension": types.StringType,
	"md5":       types.StringType,
	"image":     types.ObjectType{AttrTypes: resourceModelImageAttrTypes},
}

var resourceModelSplashPrepaidFrontAttrTypes = map[string]attr.Type{
	"extension": types.StringType,
	"md5":       types.StringType,
	"image":     types.ObjectType{AttrTypes: resourceModelImageAttrTypes},
}

var resourceModelImageAttrTypes = map[string]attr.Type{
	"contents": types.StringType,
	"format":   types.StringType,
}

var resourceModelSentryEnrollmentAttrTypes = map[string]attr.Type{
	"strength":                types.StringType,
	"enforced_systems":        types.ListType{ElemType: types.StringType},
	"systems_manager_network": types.ObjectType{AttrTypes: resourceModelSystemsManagerNetworkAttrTypes},
}

var resourceModelSystemsManagerNetworkAttrTypes = map[string]attr.Type{
	"id": types.StringType,
}
