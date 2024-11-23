package vlans

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksApplianceVLANsDatasource{}

func NewNetworksApplianceVLANsDatasource() datasource.DataSource {
	return &NetworksApplianceVLANsDatasource{}
}

// NetworksApplianceVLANsDatasource defines the resource implementation.
type NetworksApplianceVLANsDatasource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceVLANsDatasourceModel NetworksApplianceVLANModel describes the resource data model.
type NetworksApplianceVLANsDatasourceModel struct {
	Id        types.String                 `tfsdk:"id" json:"-"`
	NetworkId types.String                 `tfsdk:"network_id" json:"network_id"`
	List      []NetworksApplianceVLANModel `tfsdk:"list"`
}

func (r *NetworksApplianceVLANsDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *NetworksApplianceVLANsDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ".",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"vlan_id": schema.Int64Attribute{
						Computed: true,
						Optional: true,
					},
					"network_id": schema.StringAttribute{
						MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
						Computed:            true,
						Optional:            true,
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
				}}},
		},
	}
}

func (r *NetworksApplianceVLANsDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NetworksApplianceVLANsDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksApplianceVLANsDatasourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlans(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
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
		return
	} else {
		resp.Diagnostics.Append()
	}

	// Assuming httpResp is your *http.Response object
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error: unable to read the response body
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read HTTP response body: %v", err))
		return
	}

	// Define a struct to specifically capture the ID from the JSON data
	type HttpRespID struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	var jsonResponse []HttpRespID
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		// Handle error: JSON parsing error
		resp.Diagnostics.AddError("JSON Parsing Error", fmt.Sprintf("Error parsing JSON data for ID field: %v", err))
	}

	for _, inRespData := range inlineResp {

		vlanData := NetworksApplianceVLANModel{}
		vlanData.NetworkId = types.StringValue(data.NetworkId.ValueString())

		// Workaround for Id bug in client.GetNetworkApplianceVlans200ResponseInner
		for _, jsonInRespData := range jsonResponse {
			if jsonInRespData.Name == inRespData.GetName() {

				/*
					// Convert string to int64
							vlanId, err := strconv.ParseInt(idStr, 10, 64)
							if err != nil {
								resp.AddError("VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s' to int64: %v", idStr, err))

				*/
				vlanData.VlanId = types.Int64Value(jsonInRespData.ID)
				data.Id = types.StringValue(fmt.Sprintf("%s,%v", data.NetworkId.ValueString(), strconv.FormatInt(jsonInRespData.ID, 10)))
			}
		}

		payloadRespDiags := DatasourceReadHttpResponse(ctx, &vlanData, &inRespData)
		if payloadRespDiags != nil {
			resp.Diagnostics.Append(payloadRespDiags...)
		}

		data.List = append(data.List, vlanData)

	}

	data.Id = types.StringValue(data.NetworkId.ValueString())

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
