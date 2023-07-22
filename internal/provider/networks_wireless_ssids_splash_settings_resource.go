package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
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
	Id                              jsontypes.String   `tfsdk:"id"`
	NetworkId                       jsontypes.String   `tfsdk:"network_id"`
	Number                          jsontypes.String   `tfsdk:"number"`
	SplashUrl                       jsontypes.String   `tfsdk:"splash_url"`
	UseSplashUrl                    jsontypes.Bool     `tfsdk:"use_splash_url"`
	SplashTimeout                   jsontypes.Int64    `tfsdk:"splash_timeout"`
	WelcomeMessage                  jsontypes.String   `tfsdk:"welcome_message"`
	RedirectUrl                     jsontypes.String   `tfsdk:"redirect_url"`
	UseRedirectUrl                  jsontypes.Bool     `tfsdk:"use_redirect_url"`
	BlockAllTrafficBeforeSignOn     jsontypes.Bool     `tfsdk:"blockall_trafficbefore_signon"`
	ControllerDisconnectionBehavior jsontypes.String   `tfsdk:"controller_disconnection_behavior"`
	AllowSimultaneousLogins         jsontypes.Bool     `tfsdk:"allow_simultaneous_logins"`
	Billing                         Billing            `tfsdk:"billing"`
	GuestSponsorship                GuestSponsorship   `tfsdk:"guest_sponsorship"`
	SentryEnrollment                SentryEnrollment   `tfsdk:"sentry_enrollment"`
	SplashImage                     SplashImage        `tfsdk:"splash_image"`
	SplashLogo                      SplashLogo         `tfsdk:"splash_logo"`
	SplashPrepaidFront              SplashPrepaidFront `tfsdk:"splash_prepaid_front"`
}

type Billing struct {
	ReplyToEmailAddress           jsontypes.String `tfsdk:"reply_to_email_address"`
	PrepaidAccessFastLoginEnabled jsontypes.Bool   `tfsdk:"prepaid_access_fast_login_enabled"`
	FreeAccess                    FreeAccess       `tfsdk:"free_access"`
}

type FreeAccess struct {
	DurationInMinutes jsontypes.Int64 `tfsdk:"duration_in_minutes"`
	Enabled           jsontypes.Bool  `tfsdk:"enabled"`
}

type GuestSponsorship struct {
	DurationInMinutes        jsontypes.Int64 `tfsdk:"duration_in_minutes"`
	GuestCanRequestTimeframe jsontypes.Bool  `tfsdk:"guest_can_request_time_frame"`
}

type SplashImage struct {
	Extension jsontypes.String `tfsdk:"extension"`
	Md5       jsontypes.String `tfsdk:"md5"`
	Image     Image            `tfsdk:"image"`
}

type SplashLogo struct {
	Extension jsontypes.String `tfsdk:"extension"`
	Md5       jsontypes.String `tfsdk:"md5"`
	Image     Image            `tfsdk:"image"`
}

type SplashPrepaidFront struct {
	Extension jsontypes.String `tfsdk:"extension"`
	Md5       jsontypes.String `tfsdk:"md5"`
	Image     Image            `tfsdk:"image"`
}

type Image struct {
	Contents jsontypes.String `tfsdk:"contents"`
	Format   jsontypes.String `tfsdk:"format"`
}

type SentryEnrollment struct {
	Strength              jsontypes.String      `tfsdk:"strength"`
	EnforcedSystems       []string              `tfsdk:"enforced_systems"`
	SystemsManagerNetwork SystemsManagerNetwork `tfsdk:"systems_manager_network"`
}

type SystemsManagerNetwork struct {
	Id jsontypes.String `tfsdk:"id"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *NetworksWirelessSsidsSplashSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source and it's also used in the acceptance tests.
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
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
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
				CustomType:          jsontypes.StringType,
			},
			"splash_url": schema.StringAttribute{
				MarkdownDescription: "The custom splash URL of the click-through splash page. Note that the URL can be configured without necessarily being used. In order to enable the custom URL, see 'useSplashUrl'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"use_splash_url": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"splash_timeout": schema.Int64Attribute{
				MarkdownDescription: "Splash timeout in minutes. This will determine how often users will see the splash page.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"welcome_message": schema.StringAttribute{
				MarkdownDescription: "The welcome message for the users on the splash page.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"redirect_url": schema.StringAttribute{
				MarkdownDescription: "The custom redirect URL where the users will go after the splash page.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"use_redirect_url": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"blockall_trafficbefore_signon": schema.BoolAttribute{
				MarkdownDescription: "How restricted allowing traffic should be. If true, all traffic types are blocked until the splash page is acknowledged. If false, all non-HTTP traffic is allowed before the splash page is acknowledged.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"controller_disconnection_behavior": schema.StringAttribute{
				MarkdownDescription: "How login attempts should be handled when the controller is unreachable. Can be either 'open', 'restricted', or 'default'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"allow_simultaneous_logins": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to allow simultaneous logins from different devices.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
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
						CustomType:          jsontypes.StringType,
					},
					"prepaid_access_fast_login_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not billing uses the fast login prepaid access option.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"free_access": schema.SingleNestedAttribute{
						MarkdownDescription: "Details associated with a free access plan with limits.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"duration_in_minutes": schema.Int64Attribute{
								MarkdownDescription: "How long a device can use a network for free..",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether or not free access is enabled.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
				},
			},
			"guest_sponsorship": schema.SingleNestedAttribute{
				MarkdownDescription: "Details associated with guest sponsored splash.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"duration_in_minutes": schema.Int64Attribute{
						MarkdownDescription: "Duration in minutes of sponsored guest authorization. Must be between 1 and 60480 (6 weeks).",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"guest_can_request_time_frame": schema.BoolAttribute{
						MarkdownDescription: "Whether or not guests can specify how much time they are requesting.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
			"sentry_enrollment": schema.SingleNestedAttribute{
				MarkdownDescription: "Systems Manager sentry enrollment splash settings.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"strength": schema.StringAttribute{
						MarkdownDescription: "The strength of the enforcement of selected system types. Must be one of: 'focused', 'click-through', and 'strict'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"enforced_systems": schema.SetAttribute{
						MarkdownDescription: "The system types that the Sentry enforces. Must be included in: 'iOS, 'Android', 'macOS', and 'Windows'.",
						Optional:            true,
						Computed:            true,
						ElementType:         jsontypes.StringType,
					},
					"systems_manager_network": schema.SingleNestedAttribute{
						MarkdownDescription: "Systems Manager network targeted for sentry enrollment..",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "The network ID of the Systems Manager network.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
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
						CustomType:          jsontypes.StringType,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The system types that the Sentry enforces. Must be included in: 'iOS, 'Android', 'macOS', and 'Windows'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
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
								CustomType:          jsontypes.StringType,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
			"splash_logo": schema.SingleNestedAttribute{
				MarkdownDescription: "The logo used in the splash page.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the logo file.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the logo file. Setting this to null will remove the logo from the splash page.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new logo.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
			"splash_prepaid_front": schema.SingleNestedAttribute{
				MarkdownDescription: "The prepaid front image used in the splash page.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the prepaid front image file.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the prepaid front image file. Setting this to null will remove the prepaid front from the splash page.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new prepaid.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
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

// Create method is responsible for creating a new resource.
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsSplashSettingsResourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidSplashSettings := *openApiClient.NewInlineObject164()

	if !data.SplashUrl.IsUnknown() {
		if !data.SplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashUrl(data.SplashUrl.ValueString())
		}
	}
	if !data.UseSplashUrl.IsUnknown() {
		if !data.UseSplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetUseSplashUrl(data.UseSplashUrl.ValueBool())
		}
	}
	if !data.SplashTimeout.IsUnknown() {
		if !data.SplashTimeout.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashTimeout(int32(data.SplashTimeout.ValueInt64()))
		}
	}
	if !data.WelcomeMessage.IsUnknown() {
		if !data.WelcomeMessage.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetWelcomeMessage(data.WelcomeMessage.ValueString())
		}
	}
	if !data.RedirectUrl.IsUnknown() {
		if !data.RedirectUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetRedirectUrl(data.RedirectUrl.ValueString())
		}
	}
	if !data.UseRedirectUrl.IsUnknown() {
		if !data.UseRedirectUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetUseRedirectUrl(data.UseRedirectUrl.ValueBool())
		}
	}
	if !data.AllowSimultaneousLogins.IsUnknown() {
		if !data.AllowSimultaneousLogins.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetAllowSimultaneousLogins(data.AllowSimultaneousLogins.ValueBool())
		}
	}
	if !data.ControllerDisconnectionBehavior.IsUnknown() {
		if !data.ControllerDisconnectionBehavior.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetControllerDisconnectionBehavior(data.ControllerDisconnectionBehavior.ValueString())
		}
	}
	if !data.BlockAllTrafficBeforeSignOn.IsUnknown() {
		if !data.BlockAllTrafficBeforeSignOn.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetBlockAllTrafficBeforeSignOn(data.BlockAllTrafficBeforeSignOn.ValueBool())
		}
	}
	var billing openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsBilling
	if !data.Billing.ReplyToEmailAddress.IsUnknown() {
		billing.SetReplyToEmailAddress(data.Billing.ReplyToEmailAddress.ValueString())
	}
	if !data.Billing.PrepaidAccessFastLoginEnabled.IsUnknown() {
		billing.SetPrepaidAccessFastLoginEnabled(data.Billing.PrepaidAccessFastLoginEnabled.ValueBool())
	}
	var freeaccess openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsBillingFreeAccess
	if !data.Billing.FreeAccess.DurationInMinutes.IsUnknown() {
		freeaccess.SetDurationInMinutes(int32(data.Billing.FreeAccess.DurationInMinutes.ValueInt64()))
	}
	if !data.Billing.FreeAccess.Enabled.IsUnknown() {
		freeaccess.SetEnabled(data.Billing.FreeAccess.Enabled.ValueBool())
	}
	updateNetworkWirelessSsidSplashSettings.SetBilling(billing)
	updateNetworkWirelessSsidSplashSettings.Billing.SetFreeAccess(freeaccess)

	var guestSponsorship openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsGuestSponsorship
	if !data.GuestSponsorship.DurationInMinutes.IsUnknown() {
		guestSponsorship.SetDurationInMinutes(int32(data.GuestSponsorship.DurationInMinutes.ValueInt64()))
	}
	if !data.GuestSponsorship.GuestCanRequestTimeframe.IsUnknown() {
		guestSponsorship.SetGuestCanRequestTimeframe(data.GuestSponsorship.GuestCanRequestTimeframe.ValueBool())
	}
	updateNetworkWirelessSsidSplashSettings.SetGuestSponsorship(guestSponsorship)

	if !data.SentryEnrollment.SystemsManagerNetwork.Id.IsUnknown() {
		var systemsManagerNetwork openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSentryEnrollmentSystemsManagerNetwork
		systemsManagerNetwork.SetId(data.SentryEnrollment.SystemsManagerNetwork.Id.ValueString())
		var sentryEnrollment openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSentryEnrollment
		if len(data.SentryEnrollment.EnforcedSystems) > 0 {
			sentryEnrollment.SetEnforcedSystems(data.SentryEnrollment.EnforcedSystems)
		}
		if !data.SentryEnrollment.Strength.IsUnknown() {
			sentryEnrollment.SetStrength(data.SentryEnrollment.Strength.ValueString())
		}
		sentryEnrollment.SetSystemsManagerNetwork(systemsManagerNetwork)
		updateNetworkWirelessSsidSplashSettings.SetSentryEnrollment(sentryEnrollment)
	}
	var splashImage openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashImage
	if !data.SplashImage.Extension.IsUnknown() {
		splashImage.SetExtension(data.SplashImage.Extension.ValueString())
	}
	if !data.SplashImage.Md5.IsUnknown() {
		splashImage.SetMd5(data.SplashImage.Md5.ValueString())
	}
	var image openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashImageImage
	if !data.SplashImage.Image.Contents.IsUnknown() {
		image.SetContents(data.SplashImage.Image.Contents.ValueString())
	}
	if !data.SplashImage.Image.Format.IsUnknown() {
		image.SetContents(data.SplashImage.Image.Format.ValueString())
	}
	splashImage.SetImage(image)
	updateNetworkWirelessSsidSplashSettings.SetSplashImage(splashImage)

	var splashLogo openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashLogo
	if !data.SplashLogo.Extension.IsUnknown() {
		splashLogo.SetExtension(data.SplashLogo.Extension.ValueString())
	}
	if !data.SplashLogo.Md5.IsUnknown() {
		splashLogo.SetMd5(data.SplashLogo.Md5.ValueString())
	}
	var imageLogo openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashLogoImage
	if !data.SplashLogo.Image.Contents.IsUnknown() {
		imageLogo.SetContents(data.SplashLogo.Image.Contents.ValueString())
	}
	if !data.SplashLogo.Image.Format.IsUnknown() {
		imageLogo.SetContents(data.SplashLogo.Image.Format.ValueString())
	}
	splashLogo.SetImage(imageLogo)
	updateNetworkWirelessSsidSplashSettings.SetSplashLogo(splashLogo)

	var splashPrepaidFront openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashPrepaidFront
	if !data.SplashPrepaidFront.Extension.IsUnknown() {
		splashPrepaidFront.SetExtension(data.SplashPrepaidFront.Extension.ValueString())
	}
	if !data.SplashPrepaidFront.Md5.IsUnknown() {
		splashPrepaidFront.SetMd5(data.SplashPrepaidFront.Md5.ValueString())
	}
	var imagePrepaidFront openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashPrepaidFrontImage
	if !data.SplashPrepaidFront.Image.Contents.IsUnknown() {
		imagePrepaidFront.SetContents(data.SplashPrepaidFront.Image.Contents.ValueString())
	}
	if !data.SplashPrepaidFront.Image.Format.IsUnknown() {
		imagePrepaidFront.SetContents(data.SplashPrepaidFront.Image.Format.ValueString())
	}
	splashPrepaidFront.SetImage(imagePrepaidFront)
	updateNetworkWirelessSsidSplashSettings.SetSplashPrepaidFront(splashPrepaidFront)

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettings(updateNetworkWirelessSsidSplashSettings).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	data.SplashUrl = jsontypes.StringValue(inlineResp.GetSplashUrl())
	data.UseSplashUrl = jsontypes.BoolValue(inlineResp.GetUseSplashUrl())
	data.SplashTimeout = jsontypes.Int64Value(int64(inlineResp.GetSplashTimeout()))
	data.WelcomeMessage = jsontypes.StringValue(inlineResp.GetWelcomeMessage())
	data.RedirectUrl = jsontypes.StringValue(inlineResp.GetRedirectUrl())
	data.AllowSimultaneousLogins = jsontypes.BoolValue(inlineResp.GetAllowSimultaneousLogins())
	data.BlockAllTrafficBeforeSignOn = jsontypes.BoolValue(inlineResp.GetBlockAllTrafficBeforeSignOn())
	data.ControllerDisconnectionBehavior = jsontypes.StringValue(inlineResp.GetControllerDisconnectionBehavior())
	data.UseRedirectUrl = jsontypes.BoolValue(inlineResp.GetUseRedirectUrl())

	data.Billing.FreeAccess.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.Billing.FreeAccess.GetDurationInMinutes()))
	data.Billing.FreeAccess.Enabled = jsontypes.BoolValue(inlineResp.Billing.FreeAccess.GetEnabled())
	data.Billing.PrepaidAccessFastLoginEnabled = jsontypes.BoolValue(inlineResp.Billing.GetPrepaidAccessFastLoginEnabled())
	data.Billing.ReplyToEmailAddress = jsontypes.StringValue(inlineResp.Billing.GetReplyToEmailAddress())

	data.GuestSponsorship.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.GuestSponsorship.GetDurationInMinutes()))
	data.GuestSponsorship.GuestCanRequestTimeframe = jsontypes.BoolValue(inlineResp.GuestSponsorship.GetGuestCanRequestTimeframe())

	data.SentryEnrollment.EnforcedSystems = inlineResp.SentryEnrollment.GetEnforcedSystems()
	if len(data.SentryEnrollment.EnforcedSystems) == 0 {
		data.SentryEnrollment.EnforcedSystems = make([]string, 0)
	}
	data.SentryEnrollment.Strength = jsontypes.StringValue(inlineResp.SentryEnrollment.GetStrength())
	data.SentryEnrollment.SystemsManagerNetwork.Id = jsontypes.StringValue(inlineResp.SentryEnrollment.SystemsManagerNetwork.GetId())

	data.SplashImage.Extension = jsontypes.StringValue(inlineResp.SplashImage.GetExtension())
	data.SplashImage.Md5 = jsontypes.StringValue(inlineResp.SplashImage.GetMd5())
	if data.SplashImage.Image.Contents.IsUnknown() {
		data.SplashImage.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashImage.Image.Format.IsUnknown() {
		data.SplashImage.Image.Format = jsontypes.StringNull()
	}
	data.SplashLogo.Extension = jsontypes.StringValue(inlineResp.SplashLogo.GetExtension())
	data.SplashLogo.Md5 = jsontypes.StringValue(inlineResp.SplashLogo.GetMd5())
	if data.SplashLogo.Image.Contents.IsUnknown() {
		data.SplashLogo.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashLogo.Image.Format.IsUnknown() {
		data.SplashLogo.Image.Format = jsontypes.StringNull()
	}

	data.SplashPrepaidFront.Extension = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetExtension())
	data.SplashPrepaidFront.Md5 = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetMd5())
	if data.SplashPrepaidFront.Image.Contents.IsUnknown() {
		data.SplashPrepaidFront.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashPrepaidFront.Image.Format.IsUnknown() {
		data.SplashPrepaidFront.Image.Format = jsontypes.StringNull()
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	// TODO: Now you need to make an API call to get the current state of the resource.
	// Remember to handle any potential errors.

	inlineResp, httpResp, err := r.client.SettingsApi.GetNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	data.SplashUrl = jsontypes.StringValue(inlineResp.GetSplashUrl())
	data.UseSplashUrl = jsontypes.BoolValue(inlineResp.GetUseSplashUrl())
	data.SplashTimeout = jsontypes.Int64Value(int64(inlineResp.GetSplashTimeout()))
	data.WelcomeMessage = jsontypes.StringValue(inlineResp.GetWelcomeMessage())
	data.RedirectUrl = jsontypes.StringValue(inlineResp.GetRedirectUrl())
	data.AllowSimultaneousLogins = jsontypes.BoolValue(inlineResp.GetAllowSimultaneousLogins())
	data.BlockAllTrafficBeforeSignOn = jsontypes.BoolValue(inlineResp.GetBlockAllTrafficBeforeSignOn())
	data.ControllerDisconnectionBehavior = jsontypes.StringValue(inlineResp.GetControllerDisconnectionBehavior())
	data.UseRedirectUrl = jsontypes.BoolValue(inlineResp.GetUseRedirectUrl())
	data.Billing.FreeAccess.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.Billing.FreeAccess.GetDurationInMinutes()))
	data.Billing.FreeAccess.Enabled = jsontypes.BoolValue(inlineResp.Billing.FreeAccess.GetEnabled())
	data.Billing.PrepaidAccessFastLoginEnabled = jsontypes.BoolValue(inlineResp.Billing.GetPrepaidAccessFastLoginEnabled())
	data.Billing.ReplyToEmailAddress = jsontypes.StringValue(inlineResp.Billing.GetReplyToEmailAddress())

	data.GuestSponsorship.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.GuestSponsorship.GetDurationInMinutes()))
	data.GuestSponsorship.GuestCanRequestTimeframe = jsontypes.BoolValue(inlineResp.GuestSponsorship.GetGuestCanRequestTimeframe())

	data.SentryEnrollment.EnforcedSystems = inlineResp.SentryEnrollment.GetEnforcedSystems()
	if len(data.SentryEnrollment.EnforcedSystems) == 0 {
		data.SentryEnrollment.EnforcedSystems = make([]string, 0)
	}
	data.SentryEnrollment.Strength = jsontypes.StringValue(inlineResp.SentryEnrollment.GetStrength())
	data.SentryEnrollment.SystemsManagerNetwork.Id = jsontypes.StringValue(inlineResp.SentryEnrollment.SystemsManagerNetwork.GetId())
	data.SplashImage.Extension = jsontypes.StringValue(inlineResp.SplashImage.GetExtension())
	data.SplashImage.Md5 = jsontypes.StringValue(inlineResp.SplashImage.GetMd5())
	if data.SplashImage.Image.Contents.IsUnknown() {
		data.SplashImage.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashImage.Image.Format.IsUnknown() {
		data.SplashImage.Image.Format = jsontypes.StringNull()
	}
	data.SplashLogo.Extension = jsontypes.StringValue(inlineResp.SplashLogo.GetExtension())
	data.SplashLogo.Md5 = jsontypes.StringValue(inlineResp.SplashLogo.GetMd5())
	if data.SplashLogo.Image.Contents.IsUnknown() {
		data.SplashLogo.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashLogo.Image.Format.IsUnknown() {
		data.SplashLogo.Image.Format = jsontypes.StringNull()
	}

	data.SplashPrepaidFront.Extension = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetExtension())
	data.SplashPrepaidFront.Md5 = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetMd5())
	if data.SplashPrepaidFront.Image.Contents.IsUnknown() {
		data.SplashPrepaidFront.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashPrepaidFront.Image.Format.IsUnknown() {
		data.SplashPrepaidFront.Image.Format = jsontypes.StringNull()
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsSplashSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksWirelessSsidsSplashSettingsResourceModel

	// TODO: Make sure the plan data matches the structure of the data model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidSplashSettings := *openApiClient.NewInlineObject164()

	if !data.SplashUrl.IsUnknown() {
		if !data.SplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashUrl(data.SplashUrl.ValueString())
		}
	}
	if !data.UseSplashUrl.IsUnknown() {
		if !data.UseSplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetUseSplashUrl(data.UseSplashUrl.ValueBool())
		}
	}
	if !data.SplashTimeout.IsUnknown() {
		if !data.SplashTimeout.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashTimeout(int32(data.SplashTimeout.ValueInt64()))
		}
	}
	if !data.WelcomeMessage.IsUnknown() {
		if !data.WelcomeMessage.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetWelcomeMessage(data.WelcomeMessage.ValueString())
		}
	}
	if !data.RedirectUrl.IsUnknown() {
		if !data.RedirectUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetRedirectUrl(data.RedirectUrl.ValueString())
		}
	}
	if !data.UseRedirectUrl.IsUnknown() {
		if !data.UseRedirectUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetUseRedirectUrl(data.UseRedirectUrl.ValueBool())
		}
	}
	if !data.AllowSimultaneousLogins.IsUnknown() {
		if !data.AllowSimultaneousLogins.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetAllowSimultaneousLogins(data.AllowSimultaneousLogins.ValueBool())
		}
	}
	if !data.ControllerDisconnectionBehavior.IsUnknown() {
		if !data.ControllerDisconnectionBehavior.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetControllerDisconnectionBehavior(data.ControllerDisconnectionBehavior.ValueString())
		}
	}
	if !data.BlockAllTrafficBeforeSignOn.IsUnknown() {
		if !data.BlockAllTrafficBeforeSignOn.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetBlockAllTrafficBeforeSignOn(data.BlockAllTrafficBeforeSignOn.ValueBool())
		}
	}
	var billing openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsBilling
	if !data.Billing.ReplyToEmailAddress.IsUnknown() {
		billing.SetReplyToEmailAddress(data.Billing.ReplyToEmailAddress.ValueString())
	}
	if !data.Billing.PrepaidAccessFastLoginEnabled.IsUnknown() {
		billing.SetPrepaidAccessFastLoginEnabled(data.Billing.PrepaidAccessFastLoginEnabled.ValueBool())
	}
	var freeaccess openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsBillingFreeAccess
	if !data.Billing.FreeAccess.DurationInMinutes.IsUnknown() {
		freeaccess.SetDurationInMinutes(int32(data.Billing.FreeAccess.DurationInMinutes.ValueInt64()))
	}
	if !data.Billing.FreeAccess.Enabled.IsUnknown() {
		freeaccess.SetEnabled(data.Billing.FreeAccess.Enabled.ValueBool())
	}
	updateNetworkWirelessSsidSplashSettings.SetBilling(billing)
	updateNetworkWirelessSsidSplashSettings.Billing.SetFreeAccess(freeaccess)

	var guestSponsorship openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsGuestSponsorship
	if !data.GuestSponsorship.DurationInMinutes.IsUnknown() {
		guestSponsorship.SetDurationInMinutes(int32(data.GuestSponsorship.DurationInMinutes.ValueInt64()))
	}
	if !data.GuestSponsorship.GuestCanRequestTimeframe.IsUnknown() {
		guestSponsorship.SetGuestCanRequestTimeframe(data.GuestSponsorship.GuestCanRequestTimeframe.ValueBool())
	}
	updateNetworkWirelessSsidSplashSettings.SetGuestSponsorship(guestSponsorship)

	if !data.SentryEnrollment.SystemsManagerNetwork.Id.IsUnknown() {
		var systemsManagerNetwork openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSentryEnrollmentSystemsManagerNetwork
		systemsManagerNetwork.SetId(data.SentryEnrollment.SystemsManagerNetwork.Id.ValueString())
		var sentryEnrollment openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSentryEnrollment
		if len(data.SentryEnrollment.EnforcedSystems) > 0 {
			sentryEnrollment.SetEnforcedSystems(data.SentryEnrollment.EnforcedSystems)
		}
		if !data.SentryEnrollment.Strength.IsUnknown() {
			sentryEnrollment.SetStrength(data.SentryEnrollment.Strength.ValueString())
		}
		sentryEnrollment.SetSystemsManagerNetwork(systemsManagerNetwork)
		updateNetworkWirelessSsidSplashSettings.SetSentryEnrollment(sentryEnrollment)
	}

	var splashImage openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashImage
	if !data.SplashImage.Extension.IsUnknown() {
		splashImage.SetExtension(data.SplashImage.Extension.ValueString())
	}
	if !data.SplashImage.Md5.IsUnknown() {
		splashImage.SetMd5(data.SplashImage.Md5.ValueString())
	}
	var image openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashImageImage
	if !data.SplashImage.Image.Contents.IsUnknown() {
		image.SetContents(data.SplashImage.Image.Contents.ValueString())
	}
	if !data.SplashImage.Image.Format.IsUnknown() {
		image.SetContents(data.SplashImage.Image.Format.ValueString())
	}
	splashImage.SetImage(image)
	updateNetworkWirelessSsidSplashSettings.SetSplashImage(splashImage)

	var splashLogo openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashLogo
	if !data.SplashLogo.Extension.IsUnknown() {
		splashLogo.SetExtension(data.SplashLogo.Extension.ValueString())
	}
	if !data.SplashLogo.Md5.IsUnknown() {
		splashLogo.SetMd5(data.SplashLogo.Md5.ValueString())
	}
	var imageLogo openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashLogoImage
	if !data.SplashLogo.Image.Contents.IsUnknown() {
		imageLogo.SetContents(data.SplashLogo.Image.Contents.ValueString())
	}
	if !data.SplashLogo.Image.Format.IsUnknown() {
		imageLogo.SetContents(data.SplashLogo.Image.Format.ValueString())
	}
	splashLogo.SetImage(imageLogo)
	updateNetworkWirelessSsidSplashSettings.SetSplashLogo(splashLogo)

	var splashPrepaidFront openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashPrepaidFront
	if !data.SplashPrepaidFront.Extension.IsUnknown() {
		splashPrepaidFront.SetExtension(data.SplashPrepaidFront.Extension.ValueString())
	}
	if !data.SplashPrepaidFront.Md5.IsUnknown() {
		splashPrepaidFront.SetMd5(data.SplashPrepaidFront.Md5.ValueString())
	}
	var imagePrepaidFront openApiClient.NetworksNetworkIdWirelessSsidsNumberSplashSettingsSplashPrepaidFrontImage
	if !data.SplashPrepaidFront.Image.Contents.IsUnknown() {
		imagePrepaidFront.SetContents(data.SplashPrepaidFront.Image.Contents.ValueString())
	}
	if !data.SplashPrepaidFront.Image.Format.IsUnknown() {
		imagePrepaidFront.SetContents(data.SplashPrepaidFront.Image.Format.ValueString())
	}
	splashPrepaidFront.SetImage(imagePrepaidFront)
	updateNetworkWirelessSsidSplashSettings.SetSplashPrepaidFront(splashPrepaidFront)

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettings(updateNetworkWirelessSsidSplashSettings).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	data.SplashUrl = jsontypes.StringValue(inlineResp.GetSplashUrl())
	data.UseSplashUrl = jsontypes.BoolValue(inlineResp.GetUseSplashUrl())
	data.SplashTimeout = jsontypes.Int64Value(int64(inlineResp.GetSplashTimeout()))
	data.WelcomeMessage = jsontypes.StringValue(inlineResp.GetWelcomeMessage())
	data.RedirectUrl = jsontypes.StringValue(inlineResp.GetRedirectUrl())
	data.AllowSimultaneousLogins = jsontypes.BoolValue(inlineResp.GetAllowSimultaneousLogins())
	data.BlockAllTrafficBeforeSignOn = jsontypes.BoolValue(inlineResp.GetBlockAllTrafficBeforeSignOn())
	data.ControllerDisconnectionBehavior = jsontypes.StringValue(inlineResp.GetControllerDisconnectionBehavior())
	data.UseRedirectUrl = jsontypes.BoolValue(inlineResp.GetUseRedirectUrl())
	data.Billing.FreeAccess.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.Billing.FreeAccess.GetDurationInMinutes()))
	data.Billing.FreeAccess.Enabled = jsontypes.BoolValue(inlineResp.Billing.FreeAccess.GetEnabled())
	data.Billing.PrepaidAccessFastLoginEnabled = jsontypes.BoolValue(inlineResp.Billing.GetPrepaidAccessFastLoginEnabled())
	data.Billing.ReplyToEmailAddress = jsontypes.StringValue(inlineResp.Billing.GetReplyToEmailAddress())

	data.GuestSponsorship.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.GuestSponsorship.GetDurationInMinutes()))
	data.GuestSponsorship.GuestCanRequestTimeframe = jsontypes.BoolValue(inlineResp.GuestSponsorship.GetGuestCanRequestTimeframe())

	data.SentryEnrollment.EnforcedSystems = inlineResp.SentryEnrollment.GetEnforcedSystems()
	if len(data.SentryEnrollment.EnforcedSystems) == 0 {
		data.SentryEnrollment.EnforcedSystems = make([]string, 0)
	}
	data.SentryEnrollment.Strength = jsontypes.StringValue(inlineResp.SentryEnrollment.GetStrength())
	data.SentryEnrollment.SystemsManagerNetwork.Id = jsontypes.StringValue(inlineResp.SentryEnrollment.SystemsManagerNetwork.GetId())
	data.SplashImage.Extension = jsontypes.StringValue(inlineResp.SplashImage.GetExtension())
	data.SplashImage.Md5 = jsontypes.StringValue(inlineResp.SplashImage.GetMd5())
	if data.SplashImage.Image.Contents.IsUnknown() {
		data.SplashImage.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashImage.Image.Format.IsUnknown() {
		data.SplashImage.Image.Format = jsontypes.StringNull()
	}
	data.SplashLogo.Extension = jsontypes.StringValue(inlineResp.SplashLogo.GetExtension())
	data.SplashLogo.Md5 = jsontypes.StringValue(inlineResp.SplashLogo.GetMd5())
	if data.SplashLogo.Image.Contents.IsUnknown() {
		data.SplashLogo.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashLogo.Image.Format.IsUnknown() {
		data.SplashLogo.Image.Format = jsontypes.StringNull()
	}

	data.SplashPrepaidFront.Extension = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetExtension())
	data.SplashPrepaidFront.Md5 = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetMd5())
	if data.SplashPrepaidFront.Image.Contents.IsUnknown() {
		data.SplashPrepaidFront.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashPrepaidFront.Image.Format.IsUnknown() {
		data.SplashPrepaidFront.Image.Format = jsontypes.StringNull()
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

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

	updateNetworkWirelessSsidSplashSettings := *openApiClient.NewInlineObject164()

	if !data.SplashUrl.IsUnknown() {
		if !data.SplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashUrl(data.SplashUrl.ValueString())
		}
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettings(updateNetworkWirelessSsidSplashSettings).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	data.SplashUrl = jsontypes.StringValue(inlineResp.GetSplashUrl())
	data.UseSplashUrl = jsontypes.BoolValue(inlineResp.GetUseSplashUrl())
	data.SplashTimeout = jsontypes.Int64Value(int64(inlineResp.GetSplashTimeout()))
	data.WelcomeMessage = jsontypes.StringValue(inlineResp.GetWelcomeMessage())
	data.RedirectUrl = jsontypes.StringValue(inlineResp.GetRedirectUrl())
	data.AllowSimultaneousLogins = jsontypes.BoolValue(inlineResp.GetAllowSimultaneousLogins())
	data.BlockAllTrafficBeforeSignOn = jsontypes.BoolValue(inlineResp.GetBlockAllTrafficBeforeSignOn())
	data.ControllerDisconnectionBehavior = jsontypes.StringValue(inlineResp.GetControllerDisconnectionBehavior())
	data.UseRedirectUrl = jsontypes.BoolValue(inlineResp.GetUseRedirectUrl())

	data.GuestSponsorship.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.GuestSponsorship.GetDurationInMinutes()))
	data.GuestSponsorship.GuestCanRequestTimeframe = jsontypes.BoolValue(inlineResp.GuestSponsorship.GetGuestCanRequestTimeframe())

	data.Billing.FreeAccess.DurationInMinutes = jsontypes.Int64Value(int64(inlineResp.Billing.FreeAccess.GetDurationInMinutes()))
	data.Billing.FreeAccess.Enabled = jsontypes.BoolValue(inlineResp.Billing.FreeAccess.GetEnabled())
	data.Billing.PrepaidAccessFastLoginEnabled = jsontypes.BoolValue(inlineResp.Billing.GetPrepaidAccessFastLoginEnabled())
	data.Billing.ReplyToEmailAddress = jsontypes.StringValue(inlineResp.Billing.GetReplyToEmailAddress())
	data.SentryEnrollment.EnforcedSystems = inlineResp.SentryEnrollment.GetEnforcedSystems()
	if len(data.SentryEnrollment.EnforcedSystems) == 0 {
		data.SentryEnrollment.EnforcedSystems = make([]string, 0)
	}
	data.SentryEnrollment.Strength = jsontypes.StringValue(inlineResp.SentryEnrollment.GetStrength())
	data.SentryEnrollment.SystemsManagerNetwork.Id = jsontypes.StringValue(inlineResp.SentryEnrollment.SystemsManagerNetwork.GetId())
	data.SplashImage.Extension = jsontypes.StringValue(inlineResp.SplashImage.GetExtension())
	data.SplashImage.Md5 = jsontypes.StringValue(inlineResp.SplashImage.GetMd5())
	if data.SplashImage.Image.Contents.IsUnknown() {
		data.SplashImage.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashImage.Image.Format.IsUnknown() {
		data.SplashImage.Image.Format = jsontypes.StringNull()
	}
	data.SplashLogo.Extension = jsontypes.StringValue(inlineResp.SplashLogo.GetExtension())
	data.SplashLogo.Md5 = jsontypes.StringValue(inlineResp.SplashLogo.GetMd5())
	if data.SplashLogo.Image.Contents.IsUnknown() {
		data.SplashLogo.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashLogo.Image.Format.IsUnknown() {
		data.SplashLogo.Image.Format = jsontypes.StringNull()
	}

	data.SplashPrepaidFront.Extension = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetExtension())
	data.SplashPrepaidFront.Md5 = jsontypes.StringValue(inlineResp.SplashPrepaidFront.GetMd5())
	if data.SplashPrepaidFront.Image.Contents.IsUnknown() {
		data.SplashPrepaidFront.Image.Contents = jsontypes.StringNull()
	}
	if data.SplashPrepaidFront.Image.Format.IsUnknown() {
		data.SplashPrepaidFront.Image.Format = jsontypes.StringNull()
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// TODO: The resource has been deleted, so remove it from the state.
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), idParts[1])...)

	// If there were any errors setting the attributes, return early.
	if resp.Diagnostics.HasError() {
		return
	}

}
