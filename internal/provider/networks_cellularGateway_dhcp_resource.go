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
	_ resource.Resource                = &NetworksCellularGatewayDhcpResource{}
	_ resource.ResourceWithConfigure   = &NetworksCellularGatewayDhcpResource{}
	_ resource.ResourceWithImportState = &NetworksCellularGatewayDhcpResource{}
)

func NewNetworksCellularGatewayDhcpResource() resource.Resource {
	return &NetworksCellularGatewayDhcpResource{}
}

// NetworksCellularGatewayDhcpResource defines the resource implementation.
type NetworksCellularGatewayDhcpResource struct {
	client *openApiClient.APIClient
}

// NetworksCellularGatewayDhcpResourceModel describes the resource data model.
type NetworksCellularGatewayDhcpResourceModel struct {
	Id                   jsontypes.String   `tfsdk:"id"`
	NetworkId            jsontypes.String   `tfsdk:"network_id"`
	DhcpLeaseTime        jsontypes.String   `tfsdk:"dhcp_lease_time"`
	DnsNameServers       jsontypes.String   `tfsdk:"dns_name_servers"`
	DnsCustomNameServers []jsontypes.String `tfsdk:"dns_custom_name_servers"`
}

func (r *NetworksCellularGatewayDhcpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_cellular_gateway_dhcp"
}

func (r *NetworksCellularGatewayDhcpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Networks Cellular Gateway DHCP",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"dhcp_lease_time": schema.StringAttribute{
				MarkdownDescription: " DHCP Lease time for all MG of the network. It can be '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dns_name_servers": schema.StringAttribute{
				MarkdownDescription: "'DNS name servers mode for all MG of the network. It can take 4 different values: ''upstream_dns'', '",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dns_custom_name_servers": schema.SetAttribute{
				MarkdownDescription: "list of fixed IP representing the the DNS Name servers when the mode is 'custom'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
			},
		},
	}
}

func (r *NetworksCellularGatewayDhcpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksCellularGatewayDhcpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksCellularGatewayDhcpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object68 := openApiClient.NewInlineObject68()
	object68.SetDnsNameservers(data.DnsNameServers.ValueString())
	object68.SetDhcpLeaseTime(data.DhcpLeaseTime.ValueString())
	var customNameServers []string
	{
	}
	for _, d := range data.DnsCustomNameServers {
		customNameServers = append(customNameServers, d.String())
	}
	object68.SetDnsCustomNameservers(customNameServers)
	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkCellularGatewayDhcp(ctx, data.NetworkId.ValueString()).UpdateNetworkCellularGatewayDhcp(*object68).Execute()

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

func (r *NetworksCellularGatewayDhcpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksCellularGatewayDhcpResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ConfigureApi.GetNetworkCellularGatewayDhcp(ctx, data.NetworkId.ValueString()).Execute()

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

func (r *NetworksCellularGatewayDhcpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksCellularGatewayDhcpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object68 := openApiClient.NewInlineObject68()
	object68.SetDnsNameservers(data.DnsNameServers.ValueString())
	object68.SetDhcpLeaseTime(data.DhcpLeaseTime.ValueString())
	var customNameServers []string
	{
	}
	for _, d := range data.DnsCustomNameServers {
		customNameServers = append(customNameServers, d.String())
	}
	object68.SetDnsCustomNameservers(customNameServers)
	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkCellularGatewayDhcp(ctx, data.NetworkId.ValueString()).UpdateNetworkCellularGatewayDhcp(*object68).Execute()

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
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksCellularGatewayDhcpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksCellularGatewayDhcpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object68 := openApiClient.NewInlineObject68()
	object68.SetDnsNameservers(data.DnsNameServers.ValueString())
	object68.SetDhcpLeaseTime(data.DhcpLeaseTime.ValueString())
	var customNameServers []string
	{
	}
	for _, d := range data.DnsCustomNameServers {
		customNameServers = append(customNameServers, d.String())
	}
	object68.SetDnsCustomNameservers(customNameServers)
	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkCellularGatewayDhcp(ctx, data.NetworkId.ValueString()).UpdateNetworkCellularGatewayDhcp(*object68).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
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
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksCellularGatewayDhcpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
