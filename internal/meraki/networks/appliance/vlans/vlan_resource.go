package vlans

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksApplianceVLANResource{}
	_ resource.ResourceWithConfigure   = &NetworksApplianceVLANResource{}
	_ resource.ResourceWithImportState = &NetworksApplianceVLANResource{}
)

func NewNetworksApplianceVLANResource() resource.Resource {
	return &NetworksApplianceVLANResource{}
}

// NetworksApplianceVLANResource defines the resource implementation.
type NetworksApplianceVLANResource struct {
	client *openApiClient.APIClient
}

func (r *NetworksApplianceVLANResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlan"
}

func (r *NetworksApplianceVLANResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage the VLANs for an MX network",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"vlan_id": schema.Int64Attribute{
				Required: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"interface_id": schema.StringAttribute{
				MarkdownDescription: "The Interface ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new VLAN",
				Optional:            true,
				Computed:            true,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The subnet of the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"appliance_ip": schema.StringAttribute{
				MarkdownDescription: "The local IP of the appliance on the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"group_policy_id": schema.StringAttribute{
				MarkdownDescription: " desired group policy to apply to the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"vpn_nat_subnet": schema.StringAttribute{
				MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_handling": schema.StringAttribute{
				MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_relay_server_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
				Optional:    true,
				Computed:    true,
			},
			"dhcp_lease_time": schema.StringAttribute{
				MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_options_enabled": schema.BoolAttribute{
				MarkdownDescription: "Use DHCP boot options specified in other properties",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_next_server": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option to direct boot clients to the server to load the boot file from",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_filename": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option for boot filename ",
				Optional:            true,
				Computed:            true,
			},
			"fixed_ip_assignments": schema.MapNestedAttribute{
				Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Optional:            true,
						},
					},
				},
			},
			"reserved_ip_ranges": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DHCP reserved IP ranges on the VLAN",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							MarkdownDescription: "The first IP in the reserved range",
							Optional:            true,
							Computed:            true,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "The last IP in the reserved range",
							Optional:            true,
							Computed:            true,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "A text comment for the reserved range",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"dns_nameservers": schema.StringAttribute{
				MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_options": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The list of DHCP options that will be included in DHCP responses. Each object in the list should have \"code\", \"type\", and \"value\" properties.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							MarkdownDescription: "The code for the DHCP option. This should be an integer between 2 and 254.",
							Optional:            true,
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type for the DHCP option. One of: 'text', 'ip', 'hex' or 'integer'",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("text", "ip", "hex", "integer"),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value for the DHCP option",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"template_vlan_type": schema.StringAttribute{
				MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("same", "unique"),
				},
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
				Optional:            true,
				Computed:            true,
			},
			"mask": schema.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
			},
			"ipv6": schema.SingleNestedAttribute{
				Description: "IPv6 configuration on the VLAN",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
					},
					"prefix_assignments": schema.ListNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Prefix assignments on the VLAN",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"autonomous": schema.BoolAttribute{
									MarkdownDescription: "Auto assign a /64 prefix from the origin to the VLAN",
									Optional:            true,
									Computed:            true,
								},
								"static_prefix": schema.StringAttribute{
									MarkdownDescription: "Manual configuration of a /64 prefix on the VLAN",
									Optional:            true,
									Computed:            true,
								},
								"static_appliance_ip6": schema.StringAttribute{
									MarkdownDescription: "Manual configuration of the IPv6 Appliance IP",
									Optional:            true,
									Computed:            true,
								},
								"origin": schema.SingleNestedAttribute{
									MarkdownDescription: "The origin of the prefix",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "Type of the origin",
											Optional:            true,
											Computed:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("independent", "internet"),
											},
										},
										"interfaces": schema.SetAttribute{
											ElementType: types.StringType,
											Description: "Interfaces associated with the prefix",
											Optional:    true,
											Computed:    true,
										},
									},
								},
							}},
					},
				},
			},
			"mandatory_dhcp": schema.SingleNestedAttribute{
				Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *NetworksApplianceVLANResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceVLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] CREATE Function Call")
	tflog.Trace(ctx, "Create Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initial create API call
	payload, payloadReqDiags := CreateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlanRequest(payload).Execute()

	// Meraki API seems to return http status code 201 as an error.
	if err != nil && httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"HTTP Client Create Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	payloadRespDiags := CreateHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Update to capture config items not accessible in HTTP POST
	updatePayload, updatePayloadReqDiags := UpdateHttpReqPayload(ctx, data)
	if updatePayloadReqDiags != nil {
		resp.Diagnostics.Append(updatePayloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns this as string, openAPI spec has set as Integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	updateInlineResp, updateHttpResp, updateErr := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).UpdateNetworkApplianceVlanRequest(*updatePayload).Execute()
	if updateErr != nil && updateHttpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
			utils.HttpDiagnostics(updateHttpResp),
		)
		return
	}

	// Check for API success response code
	if updateHttpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", updateHttpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	updatePayloadRespDiags := ReadHttpResponse(ctx, data, updateInlineResp)
	if updatePayloadRespDiags != nil {
		resp.Diagnostics.Append(updatePayloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] CREATE Function Call")
	tflog.Trace(ctx, "Create function completed", map[string]interface{}{
		"data": data,
	})
}

func (r *NetworksApplianceVLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] READ Function Call")
	tflog.Trace(ctx, "Read Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns string, OpenAPI defines integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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

	payloadRespDiags := ReadHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] READ Function Call")
	tflog.Trace(ctx, "Read Function", map[string]interface{}{
		"data": data,
	})
}

func (r *NetworksApplianceVLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] UPDATE Function Call")
	tflog.Trace(ctx, "Update Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadReqDiags := UpdateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns this as string, openAPI spec has set as Integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).UpdateNetworkApplianceVlanRequest(*payload).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
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

	payloadRespDiags := ReadHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] UPDATE Function Call")
	tflog.Trace(ctx, "Update Function", map[string]interface{}{
		"data": data,
	})
}

func (r *NetworksApplianceVLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] DELETE Function Call")
	tflog.Trace(ctx, "Delete Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns a string, OpenAPI spec defines an integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	httpResp, err := r.client.ApplianceApi.DeleteNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).Execute()
	if err != nil && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"HTTP Client Delete Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
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

	// Log the response data
	tflog.Info(ctx, "[finish] DELETE Function Call")
	tflog.Trace(ctx, "Delete Function", map[string]interface{}{
		"data": data,
	})

}

func (r *NetworksApplianceVLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, vlan_id. Got: %q", req.ID),
		)
		return
	}

	// ensure vlanId is formatted properly
	str := idParts[1]

	// Convert the string to int64
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert vlanId to integer",
			fmt.Sprintf("Expected import identifier with format: network_id, vlan_id. Got: %q", req.ID),
		)
	}

	// Convert the int64 to types.Int64Value
	vlanId := types.Int64Value(i)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vlan_id"), vlanId)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
