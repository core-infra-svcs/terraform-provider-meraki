package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
var (
	_ resource.Resource                = &NetworksAlertsSettingsResource{}
	_ resource.ResourceWithConfigure   = &NetworksAlertsSettingsResource{}
	_ resource.ResourceWithImportState = &NetworksAlertsSettingsResource{}
)

func NewNetworksAlertsSettingsResource() resource.Resource {
	return &NetworksAlertsSettingsResource{}
}

// NetworksAlertsSettingsResource defines the resource implementation.
type NetworksAlertsSettingsResource struct {
	client *openApiClient.APIClient
}

// NetworksAlertsSettingsResourceModel describes the resource data model.
type NetworksAlertsSettingsResourceModel struct {
	Id                  jsontypes.String `tfsdk:"id"`
	NetworkId           jsontypes.String `tfsdk:"network_id" json:"networkId"`
	DefaultDestinations types.Object     `tfsdk:"default_destinations" json:"defaultDestinations,omitempty"`
	Alerts              types.Set        `tfsdk:"alerts" json:"alerts,omitempty"`
	Muting              types.Object     `tfsdk:"muting" json:"muting,omitempty"`
}

type NetworksAlertsSettingsResourceModelMuting struct {
	ByPortSchedules types.Object `tfsdk:"by_port_schedules" json:"byPortSchedules,omitempty"`
}

type NetworksAlertsSettingsResourceModelMutingByPortSchedules struct {
	Enabled jsontypes.Bool `tfsdk:" enabled" json:"enabled"`
}

type NetworksAlertsSettingsResourceModelDefaultDestinations struct {
	Emails        jsontypes.Set[jsontypes.String] `tfsdk:"emails" json:"emails"`
	Snmp          jsontypes.Bool                  `tfsdk:"snmp" json:"snmp"`
	AllAdmins     jsontypes.Bool                  `tfsdk:"all_admins" json:"allAdmins"`
	HttpServerIds jsontypes.Set[jsontypes.String] `tfsdk:"http_server_ids" json:"httpServerIds"`
}

type NetworksAlertsSettingsResourceModelAlert struct {
	Type              jsontypes.String `tfsdk:"type" json:"type"`
	Enabled           jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	AlertDestinations types.Object     `tfsdk:"alert_destinations" json:"alertDestinations"`
	Filters           types.Object     `tfsdk:"filters" json:"filters"`
}

type NetworksAlertsSettingsResourceModelAlertDestinations struct {
	Emails        jsontypes.Set[jsontypes.String] `tfsdk:"emails" json:"emails"`
	Snmp          jsontypes.Bool                  `tfsdk:"snmp" json:"snmp"`
	AllAdmins     jsontypes.Bool                  `tfsdk:"all_admins" json:"allAdmins"`
	HttpServerIds jsontypes.Set[jsontypes.String] `tfsdk:"http_server_ids" json:"httpServerIds"`
}

type NetworksAlertsSettingsResourceModelFilter struct {
	Timeout jsontypes.Int64 `tfsdk:"timeout" json:"timeout"`
}

func (r *NetworksAlertsSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_alerts_settings"
}

func (r *NetworksAlertsSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage network alerts settings",
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
			"default_destinations": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"emails": schema.SetAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            false,
						CustomType:          jsontypes.SetType[jsontypes.String](),
						ElementType:         jsontypes.StringType,
					},
					"snmp": schema.BoolAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"all_admins": schema.BoolAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"http_server_ids": schema.SetAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.SetType[jsontypes.String](),
						ElementType:         jsontypes.StringType,
					},
				},
			},
			"alerts": schema.SetNestedAttribute{
				MarkdownDescription: "Exceptions on a per switch basis to &quot;useCombinedPower&quot;",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Serial number of the switch",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"alert_destinations": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"emails": schema.SetAttribute{
									MarkdownDescription: "Enables / disables the secure port.",
									Optional:            true,
									Computed:            false,
									CustomType:          jsontypes.SetType[jsontypes.String](),
									ElementType:         jsontypes.StringType,
								},
								"snmp": schema.BoolAttribute{
									MarkdownDescription: "Enables / disables the secure port.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
								},
								"all_admins": schema.BoolAttribute{
									MarkdownDescription: "Enables / disables the secure port.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
								},
								"http_server_ids": schema.SetAttribute{
									MarkdownDescription: "Enables / disables the secure port.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.SetType[jsontypes.String](),
									ElementType:         jsontypes.StringType,
								},
							},
						},
						"filters": schema.SingleNestedAttribute{
							Description: "Exceptions on a per switch basis to &quot;useCombinedPower&quot;",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"timeout": schema.Int64Attribute{
									MarkdownDescription: "Serial number of the switch",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.Int64Type,
								},
							},
						},
					},
				},
			},
			"muting": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"by_port_schedules": schema.SingleNestedAttribute{
						MarkdownDescription: "",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksAlertsSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksAlertsSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksAlertsSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadDiag := NetworksAlertSettingsResourcePayload(context.Background(), data)
	if payloadDiag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", payloadDiag.Errors()))
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", data))
		return
	}

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).UpdateNetworkAlertsSettingsRequest(payload).Execute()
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
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksAlertsSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksAlertsSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SettingsApi.GetNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).Execute()
	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Read HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksAlertsSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksAlertsSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadDiag := NetworksAlertSettingsResourcePayload(context.Background(), data)
	if payloadDiag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", payloadDiag))
		return
	}

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).UpdateNetworkAlertsSettingsRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksAlertsSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksAlertsSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAlertsSettingsRequest()

	var muting openApiClient.UpdateNetworkAlertsSettingsRequestMuting
	payload.SetMuting(muting)

	// Destinations
	var destinations openApiClient.UpdateNetworkAlertsSettingsRequestDefaultDestinations
	destinations.SetAllAdmins(false)
	destinations.SetSnmp(false)

	//emails
	var adminEmails []string
	destinations.SetEmails(adminEmails)

	// serverIds
	var serverIDs []string
	destinations.SetHttpServerIds(serverIDs)

	payload.SetDefaultDestinations(destinations)
	//Alerts
	var alerts []openApiClient.UpdateNetworkAlertsSettingsRequestAlertsInner
	payload.SetAlerts(alerts)

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).UpdateNetworkAlertsSettingsRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete HTTP Client Failure",
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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksAlertsSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)
}

func NetworksAlertSettingsResourcePayload(ctx context.Context, data *NetworksAlertsSettingsResourceModel) (openApiClient.UpdateNetworkAlertsSettingsRequest, diag.Diagnostics) {

	// Create HTTP request body
	payload := *openApiClient.NewUpdateNetworkAlertsSettingsRequest()

	// DefaultDestinations
	if !data.DefaultDestinations.IsUnknown() && !data.DefaultDestinations.IsNull() {
		var defaultDestinations openApiClient.UpdateNetworkAlertsSettingsRequestDefaultDestinations
		var defaultDestinationsData NetworksAlertsSettingsResourceModelDefaultDestinations
		data.DefaultDestinations.As(ctx, &defaultDestinationsData, basetypes.ObjectAsOptions{})
		payload.SetDefaultDestinations(defaultDestinations)
	}

	// Muting
	if !data.Muting.IsUnknown() && !data.Muting.IsNull() {
		var muting openApiClient.UpdateNetworkAlertsSettingsRequestMuting
		var mutingData NetworksAlertsSettingsResourceModelMuting

		data.Muting.As(ctx, &mutingData, basetypes.ObjectAsOptions{})

		var byPortSchedules openApiClient.UpdateNetworkAlertsSettingsRequestMutingByPortSchedules
		mutingData.ByPortSchedules.As(ctx, &mutingData.ByPortSchedules, basetypes.ObjectAsOptions{})

		muting.SetByPortSchedules(byPortSchedules)

		payload.SetMuting(muting)
	}

	// Alerts
	if !data.Alerts.IsUnknown() && !data.Alerts.IsNull() {
		var alerts []openApiClient.UpdateNetworkAlertsSettingsRequestAlertsInner
		var alertsData []NetworksAlertsSettingsResourceModelAlert
		diags := data.Alerts.ElementsAs(ctx, &alertsData, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetAlerts(alerts)
	}

	return payload, nil
}
