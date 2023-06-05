package provider

import (
	"context"
	"fmt"

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

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksSettingsResource{}
var _ resource.ResourceWithImportState = &NetworksSettingsResource{}

func NewNetworksSettingsResource() resource.Resource {
	return &NetworksSettingsResource{}
}

// NetworksSettingsResource defines the resource implementation.
type NetworksSettingsResource struct {
	client *openApiClient.APIClient
}

// NetworksSettingsResourceModel describes the resource data model.
type NetworksSettingsResourceModel struct {
	Id                                    jsontypes.String `tfsdk:"id"`
	NetworkId                             jsontypes.String `tfsdk:"network_id" json:"network_id"`
	LocalStatusPageEnabled                jsontypes.Bool   `tfsdk:"local_status_page_enabled" json:"localStatusPageEnabled"`
	RemoteStatusPageEnabled               jsontypes.Bool   `tfsdk:"remote_status_page_enabled" json:"remoteStatusPageEnabled"`
	SecurePortEnabled                     SecurePort       `tfsdk:"secure_port_enabled" json:"securePort"`
	LocalStatusPage                       LocalStatusPage  `tfsdk:"local_status_page" json:"localStatusPage"`
	LocalStatusPageAuthenticationPassword jsontypes.String `tfsdk:"local_status_page_authentication_password" json:"local_status_page_authentication_password"`
	FipsEnabled                           jsontypes.Bool   `tfsdk:"fips_enabled"`
	NamedVlansEnabled                     jsontypes.Bool   `tfsdk:"named_vlans_enabled"`
	ClientPrivacyExpireDataOlderThan      jsontypes.Int64  `tfsdk:"client_privacy_expire_data_older_than"`
	ClientPrivacyExpireDataBefore         jsontypes.String `tfsdk:"client_privacy_expire_data_before"`
}

type SecurePort struct {
	Enabled bool `tfsdk:"enabled" json:"enabled"`
}

type LocalStatusPage struct {
	Authentication AuthenticationInfo `tfsdk:"authentication" json:"authentication"`
}

type AuthenticationInfo struct {
	Enabled  bool   `tfsdk:"enabled" json:"enabled"`
	Username string `tfsdk:"username" json:"username"`
}

type Fips struct {
	Enabled bool `tfsdk:"enabled"`
}

type NamedVlans struct {
	Enabled bool `tfsdk:"enabled"`
}
type ClientPrivacy struct {
	ExpireDataOlderThan jsontypes.Int64 `tfsdk:"expireDataOlderThan"`
	ExpireDataBefore    string          `tfsdk:"expireDataBefore"`
}

func (r *NetworksSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_settings"
}

func (r *NetworksSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSettings resource for updating network settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"local_status_page": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"authentication": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Enables / disables the authentication on Local Status page(s).",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
							"username": schema.StringAttribute{
								MarkdownDescription: "The username used for Local Status Page(s).",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
			"local_status_page_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"remote_status_page_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"secure_port_enabled": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
			"local_status_page_authentication_password": schema.StringAttribute{
				MarkdownDescription: "The password used for Local Status Page(s). Set this to null to clear the password.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"fips_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables FIPS on the network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"named_vlans_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables Named VLANs on the Network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"client_privacy_expire_data_older_than": schema.Int64Attribute{
				MarkdownDescription: "The number of days, weeks, or months in Epoch time to expire the data before",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"client_privacy_expire_data_before": schema.StringAttribute{
				MarkdownDescription: "The date to expire the data before",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
		},
	}
}

func (r *NetworksSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NetworksSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSettings := *openApiClient.NewInlineObject98()
	updateNetworkSettings.SetLocalStatusPageEnabled(data.LocalStatusPageEnabled.ValueBool())
	updateNetworkSettings.SetRemoteStatusPageEnabled(data.RemoteStatusPageEnabled.ValueBool())
	var v openApiClient.InlineResponse20041SecurePort

	v.SetEnabled(data.SecurePortEnabled.Enabled)
	updateNetworkSettings.SetSecurePort(v)
	var l openApiClient.NetworksNetworkIdSettingsLocalStatusPage
	var a openApiClient.NetworksNetworkIdSettingsLocalStatusPageAuthentication
	a.SetEnabled(data.LocalStatusPage.Authentication.Enabled)
	a.SetPassword(data.LocalStatusPage.Authentication.Username)
	l.SetAuthentication(a)
	updateNetworkSettings.SetLocalStatusPage(l)

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettings(updateNetworkSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.Id = jsontypes.StringValue("example-id")
	data.LocalStatusPage.Authentication.Enabled = *inlineResp.GetLocalStatusPage().Authentication.Enabled
	data.LocalStatusPage.Authentication.Username = *inlineResp.GetLocalStatusPage().Authentication.Username
	data.LocalStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetLocalStatusPageEnabled())
	data.RemoteStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetRemoteStatusPageEnabled())
	if inlineResp.Fips.GetEnabled() {
		data.FipsEnabled = jsontypes.BoolValue(inlineResp.Fips.GetEnabled())
	} else {
		data.FipsEnabled = jsontypes.BoolNull()
	}
	if inlineResp.NamedVlans.GetEnabled() {
		data.NamedVlansEnabled = jsontypes.BoolValue(inlineResp.NamedVlans.GetEnabled())
	} else {
		data.NamedVlansEnabled = jsontypes.BoolNull()
	}

	if len(inlineResp.ClientPrivacy.GetExpireDataBefore().String()) > 0 {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringValue(inlineResp.ClientPrivacy.GetExpireDataBefore().String())

	} else {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringNull()

	}

	if inlineResp.ClientPrivacy.GetExpireDataOlderThan() != 0 {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Value(int64(inlineResp.ClientPrivacy.GetExpireDataOlderThan()))
	} else {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkSettings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	data.Id = jsontypes.StringValue("example-id")
	data.LocalStatusPage.Authentication.Enabled = *inlineResp.GetLocalStatusPage().Authentication.Enabled
	data.LocalStatusPage.Authentication.Username = *inlineResp.GetLocalStatusPage().Authentication.Username
	data.LocalStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetLocalStatusPageEnabled())
	data.RemoteStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetRemoteStatusPageEnabled())
	if inlineResp.Fips.GetEnabled() {
		data.FipsEnabled = jsontypes.BoolValue(inlineResp.Fips.GetEnabled())
	} else {
		data.FipsEnabled = jsontypes.BoolNull()
	}
	if inlineResp.NamedVlans.GetEnabled() {
		data.NamedVlansEnabled = jsontypes.BoolValue(inlineResp.NamedVlans.GetEnabled())
	} else {
		data.NamedVlansEnabled = jsontypes.BoolNull()
	}

	if len(inlineResp.ClientPrivacy.GetExpireDataBefore().String()) > 0 {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringValue(inlineResp.ClientPrivacy.GetExpireDataBefore().String())

	} else {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringNull()

	}

	if inlineResp.ClientPrivacy.GetExpireDataOlderThan() != 0 {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Value(int64(inlineResp.ClientPrivacy.GetExpireDataOlderThan()))
	} else {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSettings := *openApiClient.NewInlineObject98()
	updateNetworkSettings.SetLocalStatusPageEnabled(data.LocalStatusPageEnabled.ValueBool())
	updateNetworkSettings.SetRemoteStatusPageEnabled(data.RemoteStatusPageEnabled.ValueBool())
	var v openApiClient.InlineResponse20041SecurePort
	v.SetEnabled(data.SecurePortEnabled.Enabled)
	updateNetworkSettings.SetSecurePort(v)
	var l openApiClient.NetworksNetworkIdSettingsLocalStatusPage
	var a openApiClient.NetworksNetworkIdSettingsLocalStatusPageAuthentication
	a.SetEnabled(data.LocalStatusPage.Authentication.Enabled)
	a.SetPassword(data.LocalStatusPage.Authentication.Username)
	l.SetAuthentication(a)
	updateNetworkSettings.SetLocalStatusPage(l)

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettings(updateNetworkSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.Id = jsontypes.StringValue("example-id")
	data.LocalStatusPage.Authentication.Enabled = *inlineResp.GetLocalStatusPage().Authentication.Enabled
	data.LocalStatusPage.Authentication.Username = *inlineResp.GetLocalStatusPage().Authentication.Username
	data.LocalStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetLocalStatusPageEnabled())
	data.RemoteStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetRemoteStatusPageEnabled())
	if inlineResp.Fips.GetEnabled() {
		data.FipsEnabled = jsontypes.BoolValue(inlineResp.Fips.GetEnabled())
	} else {
		data.FipsEnabled = jsontypes.BoolNull()
	}
	if inlineResp.NamedVlans.GetEnabled() {
		data.NamedVlansEnabled = jsontypes.BoolValue(inlineResp.NamedVlans.GetEnabled())
	} else {
		data.NamedVlansEnabled = jsontypes.BoolNull()
	}

	if len(inlineResp.ClientPrivacy.GetExpireDataBefore().String()) > 0 {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringValue(inlineResp.ClientPrivacy.GetExpireDataBefore().String())

	} else {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringNull()

	}

	if inlineResp.ClientPrivacy.GetExpireDataOlderThan() != 0 {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Value(int64(inlineResp.ClientPrivacy.GetExpireDataOlderThan()))
	} else {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSettings := *openApiClient.NewInlineObject98()
	updateNetworkSettings.SetLocalStatusPageEnabled(data.LocalStatusPageEnabled.ValueBool())
	updateNetworkSettings.SetRemoteStatusPageEnabled(data.RemoteStatusPageEnabled.ValueBool())
	var v openApiClient.InlineResponse20041SecurePort
	v.SetEnabled(data.SecurePortEnabled.Enabled)
	updateNetworkSettings.SetSecurePort(v)
	var l openApiClient.NetworksNetworkIdSettingsLocalStatusPage
	var a openApiClient.NetworksNetworkIdSettingsLocalStatusPageAuthentication
	a.SetEnabled(data.LocalStatusPage.Authentication.Enabled)
	a.SetPassword(data.LocalStatusPage.Authentication.Username)
	l.SetAuthentication(a)
	updateNetworkSettings.SetLocalStatusPage(l)

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettings(updateNetworkSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.Id = jsontypes.StringValue("example-id")
	data.LocalStatusPage.Authentication.Enabled = *inlineResp.GetLocalStatusPage().Authentication.Enabled
	data.LocalStatusPage.Authentication.Username = *inlineResp.GetLocalStatusPage().Authentication.Username
	data.LocalStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetLocalStatusPageEnabled())
	data.RemoteStatusPageEnabled = jsontypes.BoolValue(inlineResp.GetRemoteStatusPageEnabled())
	if inlineResp.Fips.GetEnabled() {
		data.FipsEnabled = jsontypes.BoolValue(inlineResp.Fips.GetEnabled())
	} else {
		data.FipsEnabled = jsontypes.BoolNull()
	}
	if inlineResp.NamedVlans.GetEnabled() {
		data.NamedVlansEnabled = jsontypes.BoolValue(inlineResp.NamedVlans.GetEnabled())
	} else {
		data.NamedVlansEnabled = jsontypes.BoolNull()
	}

	if len(inlineResp.ClientPrivacy.GetExpireDataBefore().String()) > 0 {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringValue(inlineResp.ClientPrivacy.GetExpireDataBefore().String())

	} else {
		data.ClientPrivacyExpireDataBefore = jsontypes.StringNull()

	}

	if inlineResp.ClientPrivacy.GetExpireDataOlderThan() != 0 {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Value(int64(inlineResp.ClientPrivacy.GetExpireDataOlderThan()))
	} else {
		data.ClientPrivacyExpireDataOlderThan = jsontypes.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
