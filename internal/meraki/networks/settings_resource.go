package networks

import (
	"context"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	Id                      jsontypes2.String `tfsdk:"id"`
	NetworkId               jsontypes2.String `tfsdk:"network_id" json:"network_id"`
	LocalStatusPageEnabled  jsontypes2.Bool   `tfsdk:"local_status_page_enabled" json:"localStatusPageEnabled"`
	RemoteStatusPageEnabled jsontypes2.Bool   `tfsdk:"remote_status_page_enabled" json:"remoteStatusPageEnabled"`
	LocalStatusPage         types.Object      `tfsdk:"local_status_page" json:"localStatusPage"`
	SecurePortEnabled       jsontypes2.Bool   `tfsdk:"secure_port_enabled" json:"securePort"`
	FipsEnabled             jsontypes2.Bool   `tfsdk:"fips_enabled" json:"fipsEnabled"`
	NamedVlansEnabled       jsontypes2.Bool   `tfsdk:"named_vlans_enabled" json:"namedVlansEnabled"`
	//ClientPrivacyExpireDataOlderThan      jsontypes.Int64                              `tfsdk:"client_privacy_expire_data_older_than"`
	//ClientPrivacyExpireDataBefore         jsontypes.String                             `tfsdk:"client_privacy_expire_data_before"`
}

// FromGetNetworkSettings200Response transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *NetworksSettingsResourceModel) FromGetNetworkSettings200Response(ctx context.Context, data *NetworksSettingsResourceModel, networkSettings200Response *openApiClient.GetNetworkSettings200Response) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksSettingsResourceModel FromGetNetworkSettings200Response")
	if networkSettings200Response == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API networkSettings200Response for IPv6")}
	}

	// Id
	m.Id = data.NetworkId

	// NetworkId
	m.NetworkId = data.NetworkId

	// LocalStatusPageEnabled
	if networkSettings200Response.HasLocalStatusPageEnabled() {
		m.LocalStatusPageEnabled = jsontypes2.BoolValue(networkSettings200Response.GetLocalStatusPageEnabled())
	} else {
		m.LocalStatusPageEnabled = jsontypes2.BoolNull()
	}

	// RemoteStatusPageEnabled
	if networkSettings200Response.HasRemoteStatusPageEnabled() {
		m.RemoteStatusPageEnabled = jsontypes2.BoolValue(networkSettings200Response.GetRemoteStatusPageEnabled())
	} else {
		m.RemoteStatusPageEnabled = jsontypes2.BoolNull()
	}

	// SecurePortEnabled
	if networkSettings200Response.SecurePort.HasEnabled() {
		m.SecurePortEnabled = jsontypes2.BoolValue(networkSettings200Response.SecurePort.GetEnabled())
	} else {
		m.SecurePortEnabled = jsontypes2.BoolNull()
	}

	// FipsEnabled
	if networkSettings200Response.Fips.GetEnabled() {
		m.FipsEnabled = jsontypes2.BoolValue(networkSettings200Response.Fips.GetEnabled())
	} else {
		m.FipsEnabled = jsontypes2.BoolValue(false)
	}

	// NamedVlans
	if networkSettings200Response.NamedVlans.GetEnabled() {
		m.NamedVlansEnabled = jsontypes2.BoolValue(networkSettings200Response.NamedVlans.GetEnabled())
	} else {
		m.NamedVlansEnabled = jsontypes2.BoolNull()
	}

	// LocalStatusPage
	var localStatusPage NetworksSettingsResourceModelLocalStatusPage
	localStatusPageDiags := localStatusPage.FromGetNetworkSettings200Response(ctx, data, networkSettings200Response)
	if localStatusPageDiags.HasError() {
		tflog.Error(ctx, "[fail] NetworksSettingsResourceModel localStatusPage Variable")
		return localStatusPageDiags
	}

	/*
		localStatusPage := NetworksSettingsResourceModelLocalStatusPage{
					Authentication: authenticationObject,
				}

				localStatusPageObject, localStatusPageObjectDiags := types.ObjectValueFrom(ctx, NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAttrTypes(), localStatusPage)
				if localStatusPageObjectDiags.HasError() {
					tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPage localStatusPage object")
					return localStatusPageObjectDiags
				}
	*/

	localStatusPageObject, localStatusPageObjectDiags := types.ObjectValueFrom(ctx, NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAttrTypes(), localStatusPage)
	if localStatusPageObjectDiags.HasError() {
		tflog.Error(ctx, "[fail] NetworksSettingsResourceModel localStatusPageObject")
		return localStatusPageObjectDiags
	}

	m.LocalStatusPage = localStatusPageObject

	tflog.Info(ctx, "[end] NetworksSettingsResourceModel FromGetNetworkSettings200Response")

	return nil
}

type NetworksSettingsResourceModelLocalStatusPage struct {
	Authentication types.Object `tfsdk:"authentication" json:"authentication"`
}

func NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"authentication": types.ObjectType{AttrTypes: NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes()},
	}

}

// FromGetNetworkSettings200Response transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *NetworksSettingsResourceModelLocalStatusPage) FromGetNetworkSettings200Response(ctx context.Context, data *NetworksSettingsResourceModel, networkSettings200Response *openApiClient.GetNetworkSettings200Response) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksSettingsResourceModelLocalStatusPage FromGetNetworkSettings200Response")

	if networkSettings200Response == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API networkSettings200Response for IPv6")}
	}

	if networkSettings200Response.HasLocalStatusPage() {

		var authentication NetworksSettingsResourceModelLocalStatusPageAuthentication
		authenticationDiags := authentication.FromGetNetworkSettings200Response(ctx, data, networkSettings200Response)
		if authenticationDiags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPage authentication variable")
			return authenticationDiags
		}

		authenticationObject, authenticationObjectDiags := types.ObjectValueFrom(ctx, NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes(), authentication)
		if authenticationObjectDiags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthentication authentication object")
			return authenticationObjectDiags
		}

		m.Authentication = authenticationObject

	} else {

		authenticationObject := types.ObjectNull(NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes())

		m.Authentication = authenticationObject

	}

	tflog.Info(ctx, "[end] NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthentication FromGetNetworkSettings200Response")

	return nil
}

type NetworksSettingsResourceModelLocalStatusPageAuthentication struct {
	Enabled  jsontypes2.Bool   `tfsdk:"enabled" json:"enabled"`
	Username jsontypes2.String `tfsdk:"username" json:"username"`
	Password jsontypes2.String `tfsdk:"password" json:"password"`
}

func NetworksSettingsResourceModelNetworksSettingsResourceModelLocalStatusPageAuthenticationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":  jsontypes2.BoolType,
		"username": jsontypes2.StringType,
		"password": jsontypes2.StringType,
	}
}

// FromGetNetworkSettings200Response transforms an API response into the NetworksApplianceVLANsResourceModelIpv6 Terraform structure.
func (m *NetworksSettingsResourceModelLocalStatusPageAuthentication) FromGetNetworkSettings200Response(ctx context.Context, data *NetworksSettingsResourceModel, networkSettings200Response *openApiClient.GetNetworkSettings200Response) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksSettingsResourceModelLocalStatusPageAuthentication FromGetNetworkSettings200Response")

	if networkSettings200Response == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API networkSettings200Response for IPv6")}
	}

	// Authentication
	if networkSettings200Response.LocalStatusPage.HasAuthentication() {

		if networkSettings200Response.LocalStatusPage.Authentication.HasEnabled() {

			// Enabled
			if networkSettings200Response.LocalStatusPage.Authentication.HasEnabled() {
				m.Enabled = jsontypes2.BoolValue(networkSettings200Response.LocalStatusPage.Authentication.GetEnabled())
			}

			// Username
			if networkSettings200Response.LocalStatusPage.Authentication.HasUsername() {
				m.Username = jsontypes2.StringValue(networkSettings200Response.LocalStatusPage.Authentication.GetUsername())
			}
		}

	} else {
		m.Enabled = jsontypes2.BoolNull()
		m.Username = jsontypes2.StringNull()
		m.Password = jsontypes2.StringNull()
	}

	if !data.LocalStatusPage.IsNull() {
		// Check Terraform Plan for Password Value
		var LocalStatusPagePlanData NetworksSettingsResourceModelLocalStatusPage

		// Extract LocalStatusPage into Plan Data Variable
		diags := data.LocalStatusPage.As(ctx, &LocalStatusPagePlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPageAuthentication localStatusPage variable")
			return diags
		}

		// Further Extract Object Plan Data into LocalStatusPageAuthentication Variable
		var LocalStatusPageAuthenticationPlanData NetworksSettingsResourceModelLocalStatusPageAuthentication

		diags = LocalStatusPagePlanData.Authentication.As(ctx, &LocalStatusPageAuthenticationPlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			tflog.Error(ctx, "[fail] NetworksSettingsResourceModelLocalStatusPageAuthentication localStatusPageAuthentication variable")
			return diags
		}

		// Extract Password if in Plan
		if m.Password.IsNull() {
			if LocalStatusPageAuthenticationPlanData.Password.ValueString() != "" {
				m.Password = jsontypes2.StringValue(LocalStatusPageAuthenticationPlanData.Password.ValueString())
			}
		}

		// Extract Username if in Plan
		if m.Username.IsNull() {

			if !LocalStatusPageAuthenticationPlanData.Username.IsUnknown() {
				m.Username = LocalStatusPageAuthenticationPlanData.Username
			}
		}
	}

	tflog.Info(ctx, "[end] NetworksSettingsResourceModelLocalStatusPageAuthentication FromGetNetworkSettings200Response")

	return nil
}

/*
type NetworksSettingsResourceModelClientPrivacy struct {
	ExpireDataOlderThan jsontypes.Int64 `tfsdk:"expireDataOlderThan"`
	ExpireDataBefore    string          `tfsdk:"expireDataBefore"`
}
*/

func createUpdateHttpReqPayload(ctx context.Context, data *NetworksSettingsResourceModel) (openApiClient.UpdateNetworkSettingsRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] createUpdateHttpReqPayload Function Call")

	payload := *openApiClient.NewUpdateNetworkSettingsRequest()

	tflog.Info(ctx, "[start] LocalStatusPageEnabled")
	if !data.LocalStatusPageEnabled.IsUnknown() && !data.LocalStatusPageEnabled.IsNull() {

		payload.SetLocalStatusPageEnabled(data.LocalStatusPageEnabled.ValueBool())

	}

	tflog.Info(ctx, "[start] RemoteStatusPageEnabled")
	if !data.RemoteStatusPageEnabled.IsUnknown() && !data.RemoteStatusPageEnabled.IsNull() {
		payload.SetRemoteStatusPageEnabled(data.RemoteStatusPageEnabled.ValueBool())
	}

	tflog.Info(ctx, "[start] LocalStatusPage")
	if !data.LocalStatusPage.IsUnknown() && !data.LocalStatusPage.IsNull() {

		var localStatusPage openApiClient.UpdateNetworkSettingsRequestLocalStatusPage

		var localStatusPagePlanData NetworksSettingsResourceModelLocalStatusPage
		diags := data.LocalStatusPage.As(ctx, &localStatusPagePlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"LocalStatusPagePlanDataReq",
				fmt.Sprintf("%s", diags.Errors()),
			)
			return payload, resp
		}

		var authenticationPlanData NetworksSettingsResourceModelLocalStatusPageAuthentication
		diags = localStatusPagePlanData.Authentication.As(ctx, &authenticationPlanData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"LocalStatusPageAuthenticationPlanDataReq",
				fmt.Sprintf("%s", diags.Errors()),
			)
			return payload, resp
		}

		if localStatusPage.HasAuthentication() {
			var authentication openApiClient.UpdateNetworkSettingsRequestLocalStatusPageAuthentication

			authentication.SetEnabled(authenticationPlanData.Enabled.ValueBool())
			authentication.SetPassword(authenticationPlanData.Password.ValueString())
			localStatusPage.SetAuthentication(authentication)
		}

		payload.SetLocalStatusPage(localStatusPage)

	}

	tflog.Info(ctx, "[end] LocalStatusPage")

	if !data.SecurePortEnabled.IsUnknown() && !data.SecurePortEnabled.IsNull() {
		var securePort openApiClient.GetNetworkSettings200ResponseSecurePort
		securePort.SetEnabled(data.SecurePortEnabled.ValueBool())
		payload.SetSecurePort(securePort)
	}

	//NamedVlans
	if !data.NamedVlansEnabled.IsUnknown() && !data.NamedVlansEnabled.IsNull() {
		var namedVlans openApiClient.UpdateNetworkSettingsRequestNamedVlans
		namedVlans.SetEnabled(data.NamedVlansEnabled.ValueBool())
		payload.SetNamedVlans(namedVlans)
	}

	tflog.Info(ctx, "[end] createUpdateHttpReqPayload Function Call")

	return payload, nil
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
				CustomType: jsontypes2.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes2.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"local_status_page_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
			"remote_status_page_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
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
								CustomType:          jsontypes2.BoolType,
							},
							"username": schema.StringAttribute{
								MarkdownDescription: "The username used for Local Status Page(s).",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes2.StringType,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password used for Local Status Page(s). Set this to null to clear the password.",
								Optional:            true,
								CustomType:          jsontypes2.StringType,
								Sensitive:           true,
							},
						}},
				},
			},
			"secure_port_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables the secure port.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
			"fips_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables FIPS on the network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
			"named_vlans_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables / disables Named VLANs on the Network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
			/*
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
			*/
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

	tflog.Info(ctx, "[start] Create Function Call")

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Initial create API call
	payload, payloadReqDiags := createUpdateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettingsRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	var NetworkSettings200Response NetworksSettingsResourceModel
	NetworkSettings200ResponseDiags := NetworkSettings200Response.FromGetNetworkSettings200Response(ctx, data, inlineResp)
	if NetworkSettings200ResponseDiags != nil {
		resp.Diagnostics.Append(NetworkSettings200ResponseDiags...)
		resp.Diagnostics.AddError("Create NetworkSettings200Response Error", fmt.Sprintf("\n%s", NetworkSettings200Response))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &NetworkSettings200Response)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
	tflog.Info(ctx, "[finish] Create Function Call")
}

func (r *NetworksSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSettingsResourceModel

	tflog.Info(ctx, "[start] Read Function Call")

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkSettings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	var NetworkSettings200Response NetworksSettingsResourceModel
	NetworkSettings200ResponseDiags := NetworkSettings200Response.FromGetNetworkSettings200Response(ctx, data, inlineResp)
	if NetworkSettings200ResponseDiags != nil {
		resp.Diagnostics.Append(NetworkSettings200ResponseDiags...)
		resp.Diagnostics.AddError("Read NetworkSettings200Response Error", fmt.Sprintf("\n%s", NetworkSettings200Response))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &NetworkSettings200Response)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
	tflog.Info(ctx, "[finish] Read Function Call")
}

func (r *NetworksSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Initial create API call
	payload, payloadReqDiags := createUpdateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettingsRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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

	var NetworkSettings200Response NetworksSettingsResourceModel
	NetworkSettings200ResponseDiags := NetworkSettings200Response.FromGetNetworkSettings200Response(ctx, data, inlineResp)
	if NetworkSettings200ResponseDiags != nil {
		resp.Diagnostics.Append(NetworkSettings200ResponseDiags...)
		resp.Diagnostics.AddError("Update NetworkSettings200Response Error", fmt.Sprintf("\n%s", NetworkSettings200Response))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &NetworkSettings200Response)...)

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

	updateNetworkSettings := *openApiClient.NewUpdateNetworkSettingsRequest()
	updateNetworkSettings.SetLocalStatusPageEnabled(true)
	updateNetworkSettings.SetRemoteStatusPageEnabled(false)
	var v openApiClient.GetNetworkSettings200ResponseSecurePort
	v.SetEnabled(false)
	updateNetworkSettings.SetSecurePort(v)
	var l openApiClient.UpdateNetworkSettingsRequestLocalStatusPage
	var a openApiClient.UpdateNetworkSettingsRequestLocalStatusPageAuthentication
	a.SetEnabled(false)
	a.SetPassword("")
	l.SetAuthentication(a)
	updateNetworkSettings.SetLocalStatusPage(l)

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettingsRequest(updateNetworkSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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
