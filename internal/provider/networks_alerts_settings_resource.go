package provider

import (
	"context"
	"encoding/json"
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
	Id                  jsontypes.String    `tfsdk:"id" json:"id"`
	NetworkId           jsontypes.String    `tfsdk:"network_id" json:"networkId"`
	DefaultDestinations DefaultDestinations `tfsdk:"default_destinations" json:"defaultDestinations"`
	Alerts              []Alert             `tfsdk:"alerts" json:"alerts"`
}

type DefaultDestinations struct {
	Emails        jsontypes.Set[jsontypes.String] `tfsdk:"emails" json:"emails"`
	Snmp          jsontypes.Bool                  `tfsdk:"snmp" json:"snmp"`
	AllAdmins     jsontypes.Bool                  `tfsdk:"all_admins" json:"allAdmins"`
	HttpServerIds jsontypes.Set[jsontypes.String] `tfsdk:"http_server_ids" json:"httpServerIds"`
}

type AlertDestinations struct {
	Emails        jsontypes.Set[jsontypes.String] `tfsdk:"emails" json:"emails"`
	Snmp          jsontypes.Bool                  `tfsdk:"snmp" json:"snmp"`
	AllAdmins     jsontypes.Bool                  `tfsdk:"all_admins" json:"allAdmins"`
	HttpServerIds jsontypes.Set[jsontypes.String] `tfsdk:"http_server_ids" json:"httpServerIds"`
}

type Filter struct {
	Timeout   jsontypes.Int64                 `tfsdk:"timeout" json:"timeout"`
	Selector  jsontypes.String                `tfsdk:"selector" json:"selector"`
	Threshold jsontypes.Int64                 `tfsdk:"threshold" json:"threshold"`
	Period    jsontypes.Int64                 `tfsdk:"period" json:"period"`
	Clients   jsontypes.Set[jsontypes.String] `tfsdk:"clients" json:"clients"`
}

type Alert struct {
	Type              jsontypes.String  `tfsdk:"type" json:"type"`
	Enabled           jsontypes.Bool    `tfsdk:"enabled" json:"enabled"`
	AlertDestinations AlertDestinations `tfsdk:"alert_destinations" json:"alertDestinations"`
	Filters           Filter            `tfsdk:"filters" json:"filters"`
}

func (r *NetworksAlertsSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_alerts_settings"
}

func (r *NetworksAlertsSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksAlertsSettings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"default_destinations": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"emails": schema.SetAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            true,
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
			"alerts": schema.ListNestedAttribute{
				Description: "Exceptions on a per switch basis to &quot;useCombinedPower&quot;",
				Required:    true,
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
									Computed:            true,
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
								"selector": schema.StringAttribute{
									MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"threshold": schema.Int64Attribute{
									MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"period": schema.Int64Attribute{
									MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"clients": schema.SetAttribute{
									MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.SetType[jsontypes.String](),
								},
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

	object27 := openApiClient.NewInlineObject27()
	destinations := openApiClient.NewNetworksNetworkIdAlertsSettingsDefaultDestinations()
	destinations.SetAllAdmins(data.DefaultDestinations.AllAdmins.ValueBool())
	adminEmails := []string{}
	for _, email := range data.DefaultDestinations.Emails.Elements() {
		adminEmails = append(adminEmails, email.String())
	}
	destinations.SetEmails(adminEmails)
	destinations.SetSnmp(data.DefaultDestinations.Snmp.ValueBool())
	serverIDs := []string{}
	for _, serverID := range data.DefaultDestinations.HttpServerIds.Elements() {
		serverIDs = append(serverIDs, serverID.String())
	}
	destinations.SetHttpServerIds(serverIDs)
	object27.SetDefaultDestinations(*destinations)
	alerts := []openApiClient.NetworksNetworkIdAlertsSettingsAlerts{}
	for _, alert := range data.Alerts {
		settingsAlerts := openApiClient.NewNetworksNetworkIdAlertsSettingsAlerts(alert.Type.ValueString())
		settingsAlerts.SetEnabled(alert.Enabled.ValueBool())
		alertDestinations := openApiClient.NewNetworksNetworkIdAlertsSettingsAlertDestinations()
		serverIDs := []string{}
		for _, serverID := range alert.AlertDestinations.HttpServerIds.Elements() {
			serverIDs = append(serverIDs, serverID.String())
		}
		adminEmails := []string{}
		for _, email := range alert.AlertDestinations.Emails.Elements() {
			adminEmails = append(adminEmails, email.String())
		}
		alertDestinations.SetEmails(adminEmails)
		alertDestinations.SetHttpServerIds(serverIDs)
		alertDestinations.SetSnmp(alert.AlertDestinations.Snmp.ValueBool())
		alertDestinations.SetAllAdmins(alert.AlertDestinations.AllAdmins.ValueBool())
		settingsAlerts.SetAlertDestinations(*alertDestinations)
		clients := []string{}
		for _, client := range alert.Filters.Clients.Elements() {
			clients = append(clients, client.String())
		}
		alert.Filters.Clients.Elements()
		filters := map[string]interface{}{}
		filters["selector"] = alert.Filters.Selector.ValueString()
		filters["period"] = alert.Filters.Period.ValueInt64()
		filters["threshold"] = alert.Filters.Threshold.ValueInt64()
		filters["timeout"] = alert.Filters.Timeout.ValueInt64()
		filters["clients"] = clients
		settingsAlerts.SetFilters(filters)

		alerts = append(alerts, *settingsAlerts)
	}
	object27.SetAlerts(alerts)
	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).UpdateNetworkAlertsSettings(*object27).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
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

	// save into the Terraform state.
	data.Id = jsontypes.StringValue("example-id")
	marshal, err := json.Marshal(inlineResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	err = json.Unmarshal(marshal, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
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

	inlineResp, httpResp, err := r.client.SettingsApi.GetNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).Execute()
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

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	marshal, err := json.Marshal(inlineResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	err = json.Unmarshal(marshal, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
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

	object27 := openApiClient.NewInlineObject27()
	destinations := openApiClient.NewNetworksNetworkIdAlertsSettingsDefaultDestinations()
	destinations.SetAllAdmins(data.DefaultDestinations.AllAdmins.ValueBool())
	adminEmails := []string{}
	for _, email := range data.DefaultDestinations.Emails.Elements() {
		adminEmails = append(adminEmails, email.String())
	}
	destinations.SetEmails(adminEmails)
	destinations.SetSnmp(data.DefaultDestinations.Snmp.ValueBool())
	serverIDs := []string{}
	for _, serverID := range data.DefaultDestinations.HttpServerIds.Elements() {
		serverIDs = append(serverIDs, serverID.String())
	}
	destinations.SetHttpServerIds(serverIDs)
	object27.SetDefaultDestinations(*destinations)
	alerts := []openApiClient.NetworksNetworkIdAlertsSettingsAlerts{}
	for _, alert := range data.Alerts {
		settingsAlerts := openApiClient.NewNetworksNetworkIdAlertsSettingsAlerts(alert.Type.ValueString())
		settingsAlerts.SetEnabled(alert.Enabled.ValueBool())
		alertDestinations := openApiClient.NewNetworksNetworkIdAlertsSettingsAlertDestinations()
		serverIDs := []string{}
		for _, serverID := range alert.AlertDestinations.HttpServerIds.Elements() {
			serverIDs = append(serverIDs, serverID.String())
		}
		adminEmails := []string{}
		for _, email := range alert.AlertDestinations.Emails.Elements() {
			adminEmails = append(adminEmails, email.String())
		}
		alertDestinations.SetEmails(adminEmails)
		alertDestinations.SetHttpServerIds(serverIDs)
		alertDestinations.SetSnmp(alert.AlertDestinations.Snmp.ValueBool())
		alertDestinations.SetAllAdmins(alert.AlertDestinations.AllAdmins.ValueBool())
		settingsAlerts.SetAlertDestinations(*alertDestinations)
		clients := []string{}
		for _, client := range alert.Filters.Clients.Elements() {
			clients = append(clients, client.String())
		}
		alert.Filters.Clients.Elements()
		filters := map[string]interface{}{}
		filters["selector"] = alert.Filters.Selector.ValueString()
		filters["period"] = alert.Filters.Period.ValueInt64()
		filters["threshold"] = alert.Filters.Threshold.ValueInt64()
		filters["timeout"] = alert.Filters.Timeout.ValueInt64()
		filters["clients"] = clients
		settingsAlerts.SetFilters(filters)

		alerts = append(alerts, *settingsAlerts)
	}
	object27.SetAlerts(alerts)
	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).UpdateNetworkAlertsSettings(*object27).Execute()
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
	if httpResp.StatusCode != 201 {
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

	// save into the Terraform state.
	data.Id = jsontypes.StringValue("example-id")
	marshal, err := json.Marshal(inlineResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	err = json.Unmarshal(marshal, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
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

	object27 := openApiClient.NewInlineObject27()
	destinations := openApiClient.NewNetworksNetworkIdAlertsSettingsDefaultDestinations()
	destinations.SetAllAdmins(data.DefaultDestinations.AllAdmins.ValueBool())
	adminEmails := []string{}
	for _, email := range data.DefaultDestinations.Emails.Elements() {
		adminEmails = append(adminEmails, email.String())
	}
	destinations.SetEmails(adminEmails)
	destinations.SetSnmp(data.DefaultDestinations.Snmp.ValueBool())
	serverIDs := []string{}
	for _, serverID := range data.DefaultDestinations.HttpServerIds.Elements() {
		serverIDs = append(serverIDs, serverID.String())
	}
	destinations.SetHttpServerIds(serverIDs)
	object27.SetDefaultDestinations(*destinations)
	alerts := []openApiClient.NetworksNetworkIdAlertsSettingsAlerts{}
	for _, alert := range data.Alerts {
		settingsAlerts := openApiClient.NewNetworksNetworkIdAlertsSettingsAlerts(alert.Type.ValueString())
		settingsAlerts.SetEnabled(alert.Enabled.ValueBool())
		alertDestinations := openApiClient.NewNetworksNetworkIdAlertsSettingsAlertDestinations()
		serverIDs := []string{}
		for _, serverID := range alert.AlertDestinations.HttpServerIds.Elements() {
			serverIDs = append(serverIDs, serverID.String())
		}
		adminEmails := []string{}
		for _, email := range alert.AlertDestinations.Emails.Elements() {
			adminEmails = append(adminEmails, email.String())
		}
		alertDestinations.SetSnmp(alert.AlertDestinations.Snmp.ValueBool())
		alertDestinations.SetAllAdmins(alert.AlertDestinations.AllAdmins.ValueBool())
		alertDestinations.SetEmails(adminEmails)
		alertDestinations.SetHttpServerIds(serverIDs)
		settingsAlerts.SetAlertDestinations(*alertDestinations)
		clients := []string{}
		for _, client := range alert.Filters.Clients.Elements() {
			clients = append(clients, client.String())
		}
		filters := map[string]interface{}{}
		filters["selector"] = alert.Filters.Selector.ValueString()
		filters["period"] = alert.Filters.Period.ValueInt64()
		filters["threshold"] = alert.Filters.Threshold.ValueInt64()
		filters["timeout"] = alert.Filters.Timeout.ValueInt64()
		filters["clients"] = clients
		settingsAlerts.SetFilters(filters)

		alerts = append(alerts, *settingsAlerts)
	}
	object27.SetAlerts(alerts)
	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkAlertsSettings(ctx, data.NetworkId.ValueString()).UpdateNetworkAlertsSettings(*object27).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to remove resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
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

	// save into the Terraform state.
	data.Id = jsontypes.StringValue("example-id")
	marshal, err := json.Marshal(inlineResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to remove resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	err = json.Unmarshal(marshal, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to remove resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksAlertsSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)
}
