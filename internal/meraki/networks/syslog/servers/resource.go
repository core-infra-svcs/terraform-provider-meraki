package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"strconv"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id        jsontypes.String      `tfsdk:"id"`
	NetworkId jsontypes.String      `tfsdk:"network_id"`
	Servers   []resourceModelServer `tfsdk:"servers"`
}

type resourceModelServer struct {
	Host  jsontypes.String   `tfsdk:"host"`
	Port  jsontypes.String   `tfsdk:"port"`
	Roles []jsontypes.String `tfsdk:"roles"`
}

/*
type resourceModelResponse struct {
	Servers []resourceModelServer `tfsdk:"servers"`
}
*/

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_syslog_servers"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSyslogServers resource for updating networks syslog servers.",
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
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"servers": schema.SetNestedAttribute{
				MarkdownDescription: "servers",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "The IP address of the syslog server",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"port": schema.StringAttribute{
							MarkdownDescription: "The port of the syslog server",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"roles": schema.SetAttribute{
							Description: "roles",
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

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var servers []openApiClient.UpdateNetworkSyslogServersRequestServersInner

	// Servers
	if len(data.Servers) > 0 {

		for _, attribute := range data.Servers {
			var server openApiClient.UpdateNetworkSyslogServersRequestServersInner
			server.SetHost(attribute.Host.ValueString())

			parsedStr, err := strconv.ParseInt(attribute.Port.ValueString(), 10, 32)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to convert port string to int32 payload",
					fmt.Sprintf("%v\n", err.Error()),
				)
			}
			server.SetPort(int32(parsedStr))

			var roles []string
			for _, role := range attribute.Roles {
				roles = append(roles, role.ValueString())

			}
			server.SetRoles(roles)
			servers = append(servers, server)
		}

	}

	updateSyslogServers := *openApiClient.NewUpdateNetworkSyslogServersRequest(servers)

	_, httpResp, err := r.client.SyslogServersApi.UpdateNetworkSyslogServers(ctx, data.NetworkId.ValueString()).UpdateNetworkSyslogServersRequest(updateSyslogServers).Execute()
	if err != nil && !strings.HasPrefix(err.Error(), "json:") {
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

	data.Id = jsontypes.StringValue("example-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data *resourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.NetworksApi.GetNetworkSyslogServers(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil && !strings.HasPrefix(err.Error(), "json:") {
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

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var servers []openApiClient.UpdateNetworkSyslogServersRequestServersInner

	// Servers
	if len(data.Servers) > 0 {

		for _, attribute := range data.Servers {
			var server openApiClient.UpdateNetworkSyslogServersRequestServersInner
			server.SetHost(attribute.Host.ValueString())

			parsedStr, err := strconv.ParseInt(attribute.Port.ValueString(), 10, 32)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to convert port string to int32 payload",
					fmt.Sprintf("%v\n", err.Error()),
				)
			}
			server.SetPort(int32(parsedStr))

			var roles []string
			for _, role := range attribute.Roles {
				roles = append(roles, role.ValueString())

			}
			server.SetRoles(roles)
			servers = append(servers, server)
		}

	}

	updateSyslogServers := *openApiClient.NewUpdateNetworkSyslogServersRequest(servers)

	_, httpResp, err := r.client.SyslogServersApi.UpdateNetworkSyslogServers(ctx, data.NetworkId.ValueString()).UpdateNetworkSyslogServersRequest(updateSyslogServers).Execute()
	if err != nil && !strings.HasPrefix(err.Error(), "json:") {
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

	data.Id = jsontypes.StringValue("example-id")

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateSyslogServers := *openApiClient.NewUpdateNetworkSyslogServersRequest(nil)

	_, httpResp, err := r.client.SyslogServersApi.UpdateNetworkSyslogServers(ctx, data.NetworkId.ValueString()).UpdateNetworkSyslogServersRequest(updateSyslogServers).Execute()
	if err != nil && !strings.HasPrefix(err.Error(), "json:") {
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

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
