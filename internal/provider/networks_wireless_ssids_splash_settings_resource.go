package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &NetworksWirelessSsidsSplashSettingsResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsSplashSettingsResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsSplashSettingsResource{} // Interface for resources with import state functionality
)

func NewNetworksWirelessSsidsSplashSettingsResource() resource.Resource {
	return &NetworksWirelessSsidsSplashSettingsResource{}
}

type NetworksWirelessSsidsSplashSettingsResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The NetworksWirelessSsidsSplashSettingsResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type NetworksWirelessSsidsSplashSettingsResourceModel struct {
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

type NetworksWirelessSsidsSplashSettingsResourceModelBilling struct {
	ReplyToEmailAddress           types.String `tfsdk:"reply_to_email_address"`
	PrepaidAccessFastLoginEnabled types.Bool   `tfsdk:"prepaid_access_fast_login_enabled"`
	FreeAccess                    types.Object `tfsdk:"free_access"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelFreeAccess struct {
	DurationInMinutes types.Int64 `tfsdk:"duration_in_minutes"`
	Enabled           types.Bool  `tfsdk:"enabled"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelGuestSponsorship struct {
	DurationInMinutes        types.Int64 `tfsdk:"duration_in_minutes"`
	GuestCanRequestTimeframe types.Bool  `tfsdk:"guest_can_request_time_frame"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelSplashImage struct {
	Extension types.String `tfsdk:"extension"`
	Md5       types.String `tfsdk:"md5"`
	Image     types.Object `tfsdk:"image"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelSplashLogo struct {
	Extension types.String `tfsdk:"extension"`
	Md5       types.String `tfsdk:"md5"`
	Image     types.Object `tfsdk:"image"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelSplashPrepaidFront struct {
	Extension types.String `tfsdk:"extension"`
	Md5       types.String `tfsdk:"md5"`
	Image     types.Object `tfsdk:"image"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelImage struct {
	Contents types.String `tfsdk:"contents"`
	Format   types.String `tfsdk:"format"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelSentryEnrollment struct {
	Strength              types.String `tfsdk:"strength"`
	EnforcedSystems       types.List   `tfsdk:"enforced_systems"`
	SystemsManagerNetwork types.Object `tfsdk:"systems_manager_network"`
}

type NetworksWirelessSsidsSplashSettingsResourceModelSystemsManagerNetwork struct {
	Id types.String `tfsdk:"id"`
}

var NetworksWirelessSsidsSplashSettingsResourceModelAttrTypes = map[string]attr.Type{
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
	"billing":                           types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelBillingAttrTypes},
	"guest_sponsorship":                 types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelGuestSponsorshipAttrTypes},
	"sentry_enrollment":                 types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelSentryEnrollmentAttrTypes},
	"splash_image":                      types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelSplashImageAttrTypes},
	"splash_logo":                       types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelSplashLogoAttrTypes},
	"splash_prepaid_front":              types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelSplashPrepaidFrontAttrTypes},
}

var NetworksWirelessSsidsSplashSettingsResourceModelBillingAttrTypes = map[string]attr.Type{
	"reply_to_email_address":            types.StringType,
	"prepaid_access_fast_login_enabled": types.BoolType,
	"free_access":                       types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelFreeAccessAttrTypes},
}

var NetworksWirelessSsidsSplashSettingsResourceModelFreeAccessAttrTypes = map[string]attr.Type{
	"duration_in_minutes": types.Int64Type,
	"enabled":             types.BoolType,
}

var NetworksWirelessSsidsSplashSettingsResourceModelGuestSponsorshipAttrTypes = map[string]attr.Type{
	"duration_in_minutes":          types.Int64Type,
	"guest_can_request_time_frame": types.BoolType,
}

var NetworksWirelessSsidsSplashSettingsResourceModelSplashImageAttrTypes = map[string]attr.Type{
	"extension": types.StringType,
	"md5":       types.StringType,
	"image":     types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelImageAttrTypes},
}

var NetworksWirelessSsidsSplashSettingsResourceModelSplashLogoAttrTypes = map[string]attr.Type{
	"extension": types.StringType,
	"md5":       types.StringType,
	"image":     types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelImageAttrTypes},
}

var NetworksWirelessSsidsSplashSettingsResourceModelSplashPrepaidFrontAttrTypes = map[string]attr.Type{
	"extension": types.StringType,
	"md5":       types.StringType,
	"image":     types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelImageAttrTypes},
}

var NetworksWirelessSsidsSplashSettingsResourceModelImageAttrTypes = map[string]attr.Type{
	"contents": types.StringType,
	"format":   types.StringType,
}

var NetworksWirelessSsidsSplashSettingsResourceModelSentryEnrollmentAttrTypes = map[string]attr.Type{
	"strength":                types.StringType,
	"enforced_systems":        types.ListType{ElemType: types.StringType},
	"systems_manager_network": types.ObjectType{AttrTypes: NetworksWirelessSsidsSplashSettingsResourceModelSystemsManagerNetworkAttrTypes},
}

var NetworksWirelessSsidsSplashSettingsResourceModelSystemsManagerNetworkAttrTypes = map[string]attr.Type{
	"id": types.StringType,
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *NetworksWirelessSsidsSplashSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source, and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids_splash_settings"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *NetworksWirelessSsidsSplashSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksWirelessSsidsSplashSettings",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"number": schema.StringAttribute{
				MarkdownDescription: "SsIds Number",
				Required:            true,
			},
			"splash_url": schema.StringAttribute{
				MarkdownDescription: "The custom splash URL of the click-through splash page. Note that the URL can be configured without necessarily being used. In order to enable the custom URL, see 'useSplashUrl'",
				Optional:            true,
				Computed:            true,
			},
			"use_splash_url": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.",
				Optional:            true,
				Computed:            true,
			},
			"splash_timeout": schema.Int64Attribute{
				MarkdownDescription: "Splash timeout in minutes. This will determine how often users will see the splash page.",
				Optional:            true,
				Computed:            true,
			},
			"welcome_message": schema.StringAttribute{
				MarkdownDescription: "The welcome message for the users on the splash page.",
				Optional:            true,
				Computed:            true,
			},
			"redirect_url": schema.StringAttribute{
				MarkdownDescription: "The custom redirect URL where the users will go after the splash page.",
				Optional:            true,
				Computed:            true,
			},
			"use_redirect_url": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.",
				Optional:            true,
				Computed:            true,
			},
			"block_all_traffic_before_sign_on": schema.BoolAttribute{
				MarkdownDescription: "How restricted allowing traffic should be. If true, all traffic types are blocked until the splash page is acknowledged. If false, all non-HTTP traffic is allowed before the splash page is acknowledged.",
				Optional:            true,
				Computed:            true,
			},
			"controller_disconnection_behavior": schema.StringAttribute{
				MarkdownDescription: "How login attempts should be handled when the controller is unreachable. Can be either 'open', 'restricted', or 'default'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("open", "restricted", "default"),
				},
			},
			"allow_simultaneous_logins": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to allow simultaneous logins from different devices.",
				Optional:            true,
				Computed:            true,
			},
			"billing": schema.SingleNestedAttribute{
				MarkdownDescription: "Details associated with billing splash.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"reply_to_email_address": schema.StringAttribute{
						MarkdownDescription: "The email address that receives replies from clients.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"prepaid_access_fast_login_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not billing uses the fast login prepaid access option.",
						Optional:            true,
						Computed:            true,
					},
					"free_access": schema.SingleNestedAttribute{
						MarkdownDescription: "Details associated with a free access plan with limits.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"duration_in_minutes": schema.Int64Attribute{
								MarkdownDescription: "How long a device can use a network for free..",
								Optional:            true,
								Computed:            true,
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether or not free access is enabled.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
				},
			},
			"guest_sponsorship": schema.SingleNestedAttribute{
				MarkdownDescription: "Details associated with guest sponsored splash.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"duration_in_minutes": schema.Int64Attribute{
						MarkdownDescription: "Duration in minutes of sponsored guest authorization. Must be between 1 and 60480 (6 weeks).",
						Optional:            true,
						Computed:            true,
					},
					"guest_can_request_time_frame": schema.BoolAttribute{
						MarkdownDescription: "Whether or not guests can specify how much time they are requesting.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			"sentry_enrollment": schema.SingleNestedAttribute{
				MarkdownDescription: "Systems Manager sentry enrollment splash settings.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"strength": schema.StringAttribute{
						MarkdownDescription: "The strength of the enforcement of selected system types. Must be one of: 'focused', 'click-through', and 'strict'.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("click-through", "focused", "strict"),
						},
					},
					"enforced_systems": schema.ListAttribute{
						MarkdownDescription: "The system types that the Sentry enforces. Must be included in: 'iOS, 'Android', 'macOS', and 'Windows'.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"systems_manager_network": schema.SingleNestedAttribute{
						MarkdownDescription: "Systems Manager network targeted for sentry enrollment..",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "The network ID of the Systems Manager network.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"splash_image": schema.SingleNestedAttribute{
				MarkdownDescription: "The image used in the splash page.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new image.",
								Optional:            true,
								Computed:            true,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("gif", "jpg", "png"),
								},
							},
						},
					},
				},
			},
			"splash_logo": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new image.",
								Optional:            true,
								Computed:            true,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("gif", "jpg", "png"),
								},
							},
						},
					},
				},
			},
			"splash_prepaid_front": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new image.",
								Optional:            true,
								Computed:            true,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("gif", "jpg", "png"),
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *NetworksWirelessSsidsSplashSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// This allows the resource to use the configured provider for any API calls it needs to make.
	r.client = client
}

func updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx context.Context, state *NetworksWirelessSsidsSplashSettingsResourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	imageObj := types.ObjectNull(NetworksWirelessSsidsSplashSettingsResourceModelImageAttrTypes)

	// Handle SplashImageState
	var splashImageState NetworksWirelessSsidsSplashSettingsResourceModelSplashImage
	err := state.SplashImage.As(ctx, &splashImageState, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags.Append(err...)
	}

	if !splashImageState.Image.IsNull() && !splashImageState.Image.IsUnknown() {
		imageObj, err = types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelImageAttrTypes, splashImageState)
		if err.HasError() {
			diags.Append(err...)
		}
	}

	return imageObj, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceSplashImageState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *NetworksWirelessSsidsSplashSettingsResourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var splashImage NetworksWirelessSsidsSplashSettingsResourceModelSplashImage
	splashImageObjNull := types.ObjectNull(NetworksWirelessSsidsSplashSettingsResourceModelSplashImageAttrTypes)

	// Md5
	if inlineResp.SplashImage.HasMd5() {
		splashImage.Md5 = types.StringValue(inlineResp.SplashImage.GetMd5())
	} else {
		splashImage.Md5 = types.StringNull()
	}

	// extension
	if inlineResp.SplashImage.HasExtension() {
		splashImage.Extension = types.StringValue(inlineResp.SplashImage.GetExtension())
	} else {
		splashImage.Extension = types.StringNull()
	}

	// image
	imageObj, err := updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx, state)
	if err.HasError() {
		diags.Append(err...)
	}

	splashImage.Image = imageObj

	splashImageObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelSplashImageAttrTypes, splashImage)
	if err.HasError() {
		diags.Append(err...)
		return splashImageObj, diags
	}

	return splashImageObjNull, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceSplashLogoState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *NetworksWirelessSsidsSplashSettingsResourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var splashLogo NetworksWirelessSsidsSplashSettingsResourceModelSplashLogo

	// Md5
	if inlineResp.SplashLogo.HasMd5() {
		splashLogo.Md5 = types.StringValue(inlineResp.SplashLogo.GetMd5())
	} else {
		splashLogo.Md5 = types.StringNull()
	}

	// extension
	if inlineResp.SplashLogo.HasExtension() {
		splashLogo.Extension = types.StringValue(inlineResp.SplashLogo.GetExtension())
	} else {
		splashLogo.Extension = types.StringNull()
	}

	// image
	imageObj, err := updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx, state)
	if err.HasError() {
		diags.Append(err...)
	}

	splashLogo.Image = imageObj

	splashLogoObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelSplashLogoAttrTypes, splashLogo)
	if err.HasError() {
		diags.Append(err...)
	}

	return splashLogoObj, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceSplashPrepaidFrontState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *NetworksWirelessSsidsSplashSettingsResourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var splashPrepaidFront NetworksWirelessSsidsSplashSettingsResourceModelSplashPrepaidFront

	// Md5
	if inlineResp.SplashPrepaidFront.HasMd5() {
		splashPrepaidFront.Md5 = types.StringValue(inlineResp.SplashPrepaidFront.GetMd5())
	} else {
		splashPrepaidFront.Md5 = types.StringNull()
	}

	// extension
	if inlineResp.SplashPrepaidFront.HasExtension() {
		splashPrepaidFront.Extension = types.StringValue(inlineResp.SplashPrepaidFront.GetExtension())
	} else {
		splashPrepaidFront.Extension = types.StringNull()
	}

	// image
	imageObj, err := updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx, state)
	if err.HasError() {
		diags.Append(err...)
	}

	splashPrepaidFront.Image = imageObj

	splashPrepaidFrontObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelSplashPrepaidFrontAttrTypes, splashPrepaidFront)
	if err.HasError() {
		diags.Append(err...)
	}

	return splashPrepaidFrontObj, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *NetworksWirelessSsidsSplashSettingsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set ID for the new resource.
	if !state.NetworkId.IsNull() && !state.NetworkId.IsUnknown() && !state.Number.IsNull() && !state.Number.IsUnknown() {
		state.Id = types.StringValue(fmt.Sprintf("%s,%s", state.NetworkId.ValueString(), state.Number.ValueString()))
	} else {
		state.Id = types.StringNull()
	}

	// Simple Attributes
	state.SplashUrl = types.StringValue(inlineResp.GetSplashUrl())
	state.UseSplashUrl = types.BoolValue(inlineResp.GetUseSplashUrl())
	state.SplashTimeout = types.Int64Value(int64(inlineResp.GetSplashTimeout()))
	state.WelcomeMessage = types.StringValue(inlineResp.GetWelcomeMessage())
	state.RedirectUrl = types.StringValue(inlineResp.GetRedirectUrl())
	state.AllowSimultaneousLogins = types.BoolValue(inlineResp.GetAllowSimultaneousLogins())
	state.BlockAllTrafficBeforeSignOn = types.BoolValue(inlineResp.GetBlockAllTrafficBeforeSignOn())
	state.ControllerDisconnectionBehavior = types.StringValue(inlineResp.GetControllerDisconnectionBehavior())
	state.UseRedirectUrl = types.BoolValue(inlineResp.GetUseRedirectUrl())

	// Handle null values for Billing
	if inlineResp.Billing != nil {
		var freeAccess NetworksWirelessSsidsSplashSettingsResourceModelFreeAccess
		if inlineResp.Billing.FreeAccess != nil {
			freeAccess.DurationInMinutes = types.Int64Value(int64(inlineResp.Billing.FreeAccess.GetDurationInMinutes()))
			freeAccess.Enabled = types.BoolValue(inlineResp.Billing.FreeAccess.GetEnabled())
		}
		freeAccessObject, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelFreeAccessAttrTypes, freeAccess)
		if err.HasError() {
			diags.Append(err...)
		}

		billing := NetworksWirelessSsidsSplashSettingsResourceModelBilling{
			ReplyToEmailAddress:           types.StringValue(inlineResp.Billing.GetReplyToEmailAddress()),
			PrepaidAccessFastLoginEnabled: types.BoolValue(inlineResp.Billing.GetPrepaidAccessFastLoginEnabled()),
			FreeAccess:                    freeAccessObject,
		}

		billingObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelBillingAttrTypes, billing)
		if err.HasError() {
			diags.Append(err...)
		}

		state.Billing = billingObj
	}

	// Handle null values for GuestSponsorship
	if inlineResp.GuestSponsorship != nil {
		var guestSponsorship NetworksWirelessSsidsSplashSettingsResourceModelGuestSponsorship
		guestSponsorship.DurationInMinutes = types.Int64Value(int64(inlineResp.GuestSponsorship.GetDurationInMinutes()))
		guestSponsorship.GuestCanRequestTimeframe = types.BoolValue(inlineResp.GuestSponsorship.GetGuestCanRequestTimeframe())

		guestSponsorshipObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelGuestSponsorshipAttrTypes, guestSponsorship)
		if err.HasError() {
			diags.Append(err...)
		}
		state.GuestSponsorship = guestSponsorshipObj
	}

	// Handle null values for SentryEnrollment
	if inlineResp.SentryEnrollment != nil {
		var sentryEnrollment NetworksWirelessSsidsSplashSettingsResourceModelSentryEnrollment

		// enforcedSystems obj
		enforcedSystemsObj, err := types.ListValueFrom(ctx, types.StringType, inlineResp.SentryEnrollment.GetEnforcedSystems())
		if err.HasError() {
			diags.Append(err...)
		}
		sentryEnrollment.EnforcedSystems = enforcedSystemsObj

		sentryEnrollment.Strength = types.StringValue(inlineResp.SentryEnrollment.GetStrength())

		// systemsManagerNetwork
		if inlineResp.SentryEnrollment.SystemsManagerNetwork != nil {
			var systemsManagerNetwork NetworksWirelessSsidsSplashSettingsResourceModelSystemsManagerNetwork
			systemsManagerNetwork.Id = types.StringValue(inlineResp.SentryEnrollment.SystemsManagerNetwork.GetId())
			systemsManagerNetworkObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelSystemsManagerNetworkAttrTypes, systemsManagerNetwork)
			if err.HasError() {
				diags.Append(err...)
			}

			sentryEnrollment.SystemsManagerNetwork = systemsManagerNetworkObj
		}

		sentryEnrollmentObj, err := types.ObjectValueFrom(ctx, NetworksWirelessSsidsSplashSettingsResourceModelSentryEnrollmentAttrTypes, sentryEnrollment)
		if err.HasError() {
			diags.Append(err...)
		}
		state.SentryEnrollment = sentryEnrollmentObj
	}

	// Handle null values for SplashImage
	splashImage, err := updateNetworksWirelessSsidsSplashSettingsResourceSplashImageState(ctx, inlineResp, state)
	if err.HasError() {
		diags.Append(err...)
	}
	state.SplashImage = splashImage

	// Handle null values for SplashLogo
	splashLogo, err := updateNetworksWirelessSsidsSplashSettingsResourceSplashLogoState(ctx, inlineResp, state)
	if err.HasError() {
		diags.Append(err...)
	}
	state.SplashLogo = splashLogo

	// Handle null values for SplashPrepaidFront
	splashPrepaidFront, err := updateNetworksWirelessSsidsSplashSettingsResourceSplashPrepaidFrontState(ctx, inlineResp, state)
	if err.HasError() {
		diags.Append(err...)
	}
	state.SplashPrepaidFront = splashPrepaidFront

	return diags
}

func NetworksWirelessSsidsSplashSettingsResourcePayload(ctx context.Context, data *NetworksWirelessSsidsSplashSettingsResourceModel) (openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload := *openApiClient.NewUpdateNetworkWirelessSsidSplashSettingsRequest()
	payload.SetSplashUrl(data.SplashUrl.ValueString())
	payload.SetUseSplashUrl(data.UseSplashUrl.ValueBool())
	payload.SetSplashTimeout(int32(data.SplashTimeout.ValueInt64()))
	payload.SetWelcomeMessage(data.WelcomeMessage.ValueString())
	payload.SetRedirectUrl(data.RedirectUrl.ValueString())
	payload.SetUseRedirectUrl(data.UseRedirectUrl.ValueBool())
	payload.SetAllowSimultaneousLogins(data.AllowSimultaneousLogins.ValueBool())
	payload.SetControllerDisconnectionBehavior(data.ControllerDisconnectionBehavior.ValueString())
	payload.SetBlockAllTrafficBeforeSignOn(data.BlockAllTrafficBeforeSignOn.ValueBool())

	// Billing
	if !data.Billing.IsUnknown() && !data.Billing.IsNull() {
		var billing openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestBilling
		var billingData NetworksWirelessSsidsSplashSettingsResourceModelBilling

		err := data.Billing.As(ctx, &billingData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		billing.SetReplyToEmailAddress(billingData.ReplyToEmailAddress.ValueString())
		billing.SetPrepaidAccessFastLoginEnabled(billingData.PrepaidAccessFastLoginEnabled.ValueBool())

		if !billingData.FreeAccess.IsUnknown() && !billingData.FreeAccess.IsNull() {
			var freeAccess openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestBillingFreeAccess
			var freeAccessData NetworksWirelessSsidsSplashSettingsResourceModelFreeAccess

			err := billingData.FreeAccess.As(ctx, &freeAccessData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			freeAccess.SetDurationInMinutes(int32(freeAccessData.DurationInMinutes.ValueInt64()))
			freeAccess.SetEnabled(freeAccessData.Enabled.ValueBool())
			billing.SetFreeAccess(freeAccess)
		}
		payload.SetBilling(billing)
	}

	// Guest Sponsorship
	if !data.GuestSponsorship.IsUnknown() && !data.GuestSponsorship.IsNull() {
		var guestSponsorship openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestGuestSponsorship
		var guestSponsorshipData NetworksWirelessSsidsSplashSettingsResourceModelGuestSponsorship

		err := data.GuestSponsorship.As(ctx, &guestSponsorshipData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		guestSponsorship.SetDurationInMinutes(int32(guestSponsorshipData.DurationInMinutes.ValueInt64()))
		guestSponsorship.SetGuestCanRequestTimeframe(guestSponsorshipData.GuestCanRequestTimeframe.ValueBool())
		payload.SetGuestSponsorship(guestSponsorship)
	}

	// Sentry Enrollment
	if !data.SentryEnrollment.IsUnknown() && !data.SentryEnrollment.IsNull() {
		var sentryEnrollment openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSentryEnrollment
		var sentryEnrollmentData NetworksWirelessSsidsSplashSettingsResourceModelSentryEnrollment

		err := data.SentryEnrollment.As(ctx, &sentryEnrollmentData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		sentryEnrollment.SetStrength(sentryEnrollmentData.Strength.ValueString())

		var enforcedSystems []string
		for _, v := range sentryEnrollmentData.EnforcedSystems.Elements() {
			enforcedSystems = append(enforcedSystems, v.(types.String).ValueString())
		}
		sentryEnrollment.SetEnforcedSystems(enforcedSystems)

		if !sentryEnrollmentData.SystemsManagerNetwork.IsUnknown() && !sentryEnrollmentData.SystemsManagerNetwork.IsNull() {
			var systemsManagerNetwork openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSentryEnrollmentSystemsManagerNetwork
			var systemsManagerNetworkData NetworksWirelessSsidsSplashSettingsResourceModelSystemsManagerNetwork

			err := sentryEnrollmentData.SystemsManagerNetwork.As(ctx, &systemsManagerNetworkData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			systemsManagerNetwork.SetId(systemsManagerNetworkData.Id.ValueString())
			sentryEnrollment.SetSystemsManagerNetwork(systemsManagerNetwork)
		}
		payload.SetSentryEnrollment(sentryEnrollment)
	}

	// Splash Image
	if !data.SplashImage.IsUnknown() && !data.SplashImage.IsNull() {
		var splashImage openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashImage
		var splashImageData NetworksWirelessSsidsSplashSettingsResourceModelSplashImage

		err := data.SplashImage.As(ctx, &splashImageData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		splashImage.SetExtension(splashImageData.Extension.ValueString())
		splashImage.SetMd5(splashImageData.Md5.ValueString())

		if !splashImageData.Image.IsUnknown() && !splashImageData.Image.IsNull() {
			var image openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashImageImage
			var imageData NetworksWirelessSsidsSplashSettingsResourceModelImage

			err := splashImageData.Image.As(ctx, &imageData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			image.SetContents(imageData.Contents.ValueString())
			image.SetFormat(imageData.Format.ValueString())
			splashImage.SetImage(image)
		}
		payload.SetSplashImage(splashImage)
	}

	// Splash Logo
	if !data.SplashLogo.IsUnknown() && !data.SplashLogo.IsNull() {
		var splashLogo openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashLogo
		var splashLogoData NetworksWirelessSsidsSplashSettingsResourceModelSplashLogo

		err := data.SplashLogo.As(ctx, &splashLogoData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		splashLogo.SetExtension(splashLogoData.Extension.ValueString())
		splashLogo.SetMd5(splashLogoData.Md5.ValueString())

		if !splashLogoData.Image.IsUnknown() && !splashLogoData.Image.IsNull() {
			var imageLogo openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashLogoImage
			var imageLogoData NetworksWirelessSsidsSplashSettingsResourceModelImage

			err := splashLogoData.Image.As(ctx, &imageLogoData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			imageLogo.SetContents(imageLogoData.Contents.ValueString())
			imageLogo.SetFormat(imageLogoData.Format.ValueString())
			splashLogo.SetImage(imageLogo)
		}

		payload.SetSplashLogo(splashLogo)
	}

	// Splash Prepaid Front
	if !data.SplashPrepaidFront.IsUnknown() && !data.SplashPrepaidFront.IsNull() {
		var splashPrepaidFront openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashPrepaidFront
		var splashPrepaidFrontData NetworksWirelessSsidsSplashSettingsResourceModelSplashPrepaidFront

		err := data.SplashPrepaidFront.As(ctx, &splashPrepaidFrontData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		splashPrepaidFront.SetExtension(splashPrepaidFrontData.Extension.ValueString())
		splashPrepaidFront.SetMd5(splashPrepaidFrontData.Md5.ValueString())

		if !splashPrepaidFrontData.Image.IsUnknown() && !splashPrepaidFrontData.Image.IsNull() {
			var imagePrepaidFront openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashPrepaidFrontImage
			var imagePrepaidFrontData NetworksWirelessSsidsSplashSettingsResourceModelImage

			err := splashPrepaidFrontData.Image.As(ctx, &imagePrepaidFrontData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			imagePrepaidFront.SetContents(imagePrepaidFrontData.Contents.ValueString())
			imagePrepaidFront.SetFormat(imagePrepaidFrontData.Format.ValueString())
			splashPrepaidFront.SetImage(imagePrepaidFront)
		}
		payload.SetSplashPrepaidFront(splashPrepaidFront)
	}

	return payload, diags
}

// Create method is responsible for creating a new resource.
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *NetworksWirelessSsidsSplashSettingsResourceModel

	// Unmarshal the plan state into the internal state model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadErr := NetworksWirelessSsidsSplashSettingsResourcePayload(ctx, state)
	if payloadErr.HasError() {
		resp.Diagnostics.Append(payloadErr...)
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), state.NetworkId.ValueString(), state.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettingsRequest(payload).Execute()
	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// If there were any errors up to this point, log the plan state and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", state))
			return
		}
	}

	stateErr := updateNetworksWirelessSsidsSplashSettingsResourceState(ctx, *inlineResp, state)
	if stateErr.HasError() {
		resp.Diagnostics.Append(stateErr...)
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksWirelessSsidsSplashSettingsResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.SettingsApi.GetNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	stateErr := updateNetworksWirelessSsidsSplashSettingsResourceState(ctx, *inlineResp, data)
	if stateErr.HasError() {
		resp.Diagnostics.Append(stateErr...)
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksWirelessSsidsSplashSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadErr := NetworksWirelessSsidsSplashSettingsResourcePayload(ctx, data)
	if payloadErr.HasError() {
		resp.Diagnostics.Append(payloadErr...)
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettingsRequest(payload).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	stateErr := updateNetworksWirelessSsidsSplashSettingsResourceState(ctx, *inlineResp, data)
	if stateErr.HasError() {
		resp.Diagnostics.Append(stateErr...)
	}

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksWirelessSsidsSplashSettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidSplashSettings := *openApiClient.NewUpdateNetworkWirelessSsidSplashSettingsRequest()

	if !data.SplashUrl.IsUnknown() {
		if !data.SplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashUrl(data.SplashUrl.ValueString())
		}
	}

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettingsRequest(updateNetworkWirelessSsidSplashSettings).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, serial number. Got: %q", req.ID),
		)
		return
	}

	// Set the attributes required for making a Read API call in the state.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), idParts[1])...)

	// If there were any errors setting the attributes, return early.
	if resp.Diagnostics.HasError() {
		return
	}

}
