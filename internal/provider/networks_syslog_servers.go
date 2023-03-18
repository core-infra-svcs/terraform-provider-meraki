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
var _ resource.Resource = &NetworksSyslogServersResource{}
var _ resource.ResourceWithImportState = &NetworksSyslogServersResource{}

func NewNetworksSyslogServersResource() resource.Resource {
	return &NetworksSyslogServersResource{}
}

// NetworksSyslogServersResource defines the resource implementation.
type NetworksSyslogServersResource struct {
	client *openApiClient.APIClient
}

// NetworksSyslogServersResourceModel describes the resource data model.
type NetworksSyslogServersResourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`
	Servers   []Server         `tfsdk:"servers"`
}

type Server struct {
	Host  string   `tfsdk:"host"`
	Port  int64    `tfsdk:"port"`
	Roles []string `tfsdk:"roles"`
}

func (r *NetworksSyslogServersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_syslog_servers"
}

func (r *NetworksSyslogServersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSyslogServers resource for updating networks syslog servers resource.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
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
			"servers": schema.SetNestedAttribute{
				MarkdownDescription: "servers",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
						"roles": schema.SetAttribute{
							Description: "Network tags",
							ElementType: jsontypes.StringType,
							CustomType:  jsontypes.SetType[jsontypes.String](),
							Computed:    true,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSyslogServersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSyslogServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSyslogServersResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var s []openApiClient.NetworksNetworkIdSyslogServersServers

	// Servers
	if len(data.Servers) > 0 {

		for _, attribute := range data.Servers {
			var t openApiClient.NetworksNetworkIdSyslogServersServers
			t.Host = attribute.Host
			t.Port = int32(attribute.Port)
			t.Roles = attribute.Roles
			s = append(s, t)
		}

	}

	updateNetworkNetflow := *openApiClient.NewInlineObject139(s)

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSyslogServers(ctx, data.NetworkId.ValueString()).UpdateNetworkSyslogServers(updateNetworkNetflow).Execute()
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

	for _, attribute := range inlineResp.GetServers() {
		var se Server
		se.Host = *attribute.Host
		se.Port = int64(*attribute.Port)
		se.Roles = attribute.Roles
		data.Servers = append(data.Servers, se)
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSyslogServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data *NetworksSyslogServersResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkSyslogServers(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get resource",
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
		return
	} else {
		resp.Diagnostics.Append()
	}

	for _, attribute := range inlineResp.GetServers() {
		var se Server
		se.Host = *attribute.Host
		se.Port = int64(*attribute.Port)
		se.Roles = attribute.Roles
		data.Servers = append(data.Servers, se)
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSyslogServersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSyslogServersResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var s []openApiClient.NetworksNetworkIdSyslogServersServers

	// Servers
	if len(data.Servers) > 0 {

		for _, attribute := range data.Servers {
			var t openApiClient.NetworksNetworkIdSyslogServersServers
			t.Host = attribute.Host
			t.Port = int32(attribute.Port)
			t.Roles = attribute.Roles
			s = append(s, t)
		}

	}

	updateNetworkNetflow := *openApiClient.NewInlineObject139(s)

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSyslogServers(ctx, data.NetworkId.ValueString()).UpdateNetworkSyslogServers(updateNetworkNetflow).Execute()
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

	for _, attribute := range inlineResp.GetServers() {
		var se Server
		se.Host = *attribute.Host
		se.Port = int64(*attribute.Port)
		se.Roles = attribute.Roles
		data.Servers = append(data.Servers, se)
	}

	data.Id = jsontypes.StringValue("example-id")

	// Write logs using the tflog package
	tflog.Trace(ctx, "update resource")
}

func (r *NetworksSyslogServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksSyslogServersResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var s []openApiClient.NetworksNetworkIdSyslogServersServers

	// Servers
	if len(data.Servers) > 0 {

		for _, attribute := range data.Servers {
			var t openApiClient.NetworksNetworkIdSyslogServersServers
			t.Host = attribute.Host
			t.Port = int32(attribute.Port)
			t.Roles = attribute.Roles
			s = append(s, t)
		}
	}

	updateNetworkNetflow := *openApiClient.NewInlineObject139(s)

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSyslogServers(ctx, data.NetworkId.ValueString()).UpdateNetworkSyslogServers(updateNetworkNetflow).Execute()
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

	for _, attribute := range inlineResp.GetServers() {
		var se Server
		se.Host = *attribute.Host
		se.Port = int64(*attribute.Port)
		se.Roles = attribute.Roles
		data.Servers = append(data.Servers, se)
	}
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksSyslogServersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
