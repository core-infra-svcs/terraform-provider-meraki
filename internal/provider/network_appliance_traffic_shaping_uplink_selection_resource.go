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
var _ resource.Resource = &NetworksApplianceTrafficShapingUplinkSelectionResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceTrafficShapingUplinkSelectionResource{}

func NewNetworksApplianceTrafficShapingUplinkSelectionResource() resource.Resource {
	return &NetworksApplianceTrafficShapingUplinkSelectionResource{}
}

// NetworksApplianceTrafficShapingUplinkSelectionResource defines the resource implementation.
type NetworksApplianceTrafficShapingUplinkSelectionResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceTrafficShapingUplinkSelectionResourceModel describes the resource data model.
type NetworksApplianceTrafficShapingUplinkSelectionResourceModel struct {
	Id                          jsontypes.String              `tfsdk:"id"`
	NetworkId                   jsontypes.String              `tfsdk:"network_id" json:"network_id"`
	ActiveActiveAutoVpnEnabled  jsontypes.Bool                `tfsdk:"active_active_auto_vpn_enabled"`
	DefaultUplink               jsontypes.String              `tfsdk:"default_uplink"`
	LoadBalancingEnabled        jsontypes.Bool                `tfsdk:"load_balancing_enabled"`
	VpnTrafficUplinkPreferences []VpnTrafficUplinkPreferences `tfsdk:"vpn_traffic_uplink_preferences"`
	FailoverAndFailback         FailoverAndFailback           `tfsdk:"failover_and_failback"`
	WanTrafficUplinkPreferences []WanTrafficUplinkPreferences `tfsdk:"wan_traffic_uplink_preferences"`
}

type VpnTrafficUplinkPreferences struct {
	PreferredUplink                             jsontypes.String     `tfsdk:"preferred_uplink"`
	FailOverCriterion                           jsontypes.String     `tfsdk:"failover_criterion"`
	PerformanceClassBuiltinPerformanceClassName jsontypes.String     `tfsdk:"performance_class_builtin_performance_class_name"`
	PerformanceClassCustomPerformanceClassId    jsontypes.String     `tfsdk:"performance_class_custom_performance_class_id"`
	PerformanceClassType                        jsontypes.String     `tfsdk:"performance_class_type"`
	TrafficFilters                              []VlanTrafficFilters `tfsdk:"traffic_filters"`
}

type FailoverAndFailback struct {
	Immediate Immediate `tfsdk:"immediate"`
}

type Immediate struct {
	Enabled jsontypes.Bool `tfsdk:"enabled"`
}

type WanTrafficUplinkPreferences struct {
	TrafficFilters  []WanTrafficFilters `tfsdk:"traffic_filters"`
	PreferredUplink jsontypes.String    `tfsdk:"preferred_uplink"`
}

type WanTrafficFilters struct {
	Type                 jsontypes.String `tfsdk:"type"`
	ValueProtocol        jsontypes.String `tfsdk:"value_protocol"`
	ValueSourcePort      jsontypes.String `tfsdk:"value_source_port"`
	ValueSourceCidr      jsontypes.String `tfsdk:"value_source_cidr"`
	ValueDestinationPort jsontypes.String `tfsdk:"value_destination_port"`
	ValueDestinationCidr jsontypes.String `tfsdk:"value_destination_cidr"`
}

type VlanTrafficFilters struct {
	Type                    jsontypes.String `tfsdk:"type"`
	ValueId                 jsontypes.String `tfsdk:"value_id"`
	ValueProtocol           jsontypes.String `tfsdk:"value_protocol"`
	ValueSourcePort         jsontypes.String `tfsdk:"value_source_port"`
	ValueSourceCidr         jsontypes.String `tfsdk:"value_source_cidr"`
	ValueSourceVlan         jsontypes.Int64  `tfsdk:"value_source_vlan"`
	ValueSourceHost         jsontypes.Int64  `tfsdk:"value_source_host"`
	ValueSourceNetwork      jsontypes.String `tfsdk:"value_source_network"`
	ValueDestinationPort    jsontypes.String `tfsdk:"value_destination_port"`
	ValueDestinationCidr    jsontypes.String `tfsdk:"value_destination_cidr"`
	ValueDestinationVlan    jsontypes.Int64  `tfsdk:"value_destination_vlan"`
	ValueDestinationFqdn    jsontypes.String `tfsdk:"value_destination_fqdn"`
	ValueDestinationNetwork jsontypes.String `tfsdk:"value_destination_network"`
	ValueDestinationHost    jsontypes.Int64  `tfsdk:"value_destination_host"`
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_traffic_shaping_uplink_selection"
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"active_active_auto_vpn_enabled": schema.BoolAttribute{
				MarkdownDescription: "Toggle for enabling or disabling active-active AutoVPN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"default_uplink": schema.StringAttribute{
				MarkdownDescription: "The default uplink. Must be one of: 'wan1' or 'wan2'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"load_balancing_enabled": schema.BoolAttribute{
				MarkdownDescription: "Toggle for enabling or disabling active-active AutoVPN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"vpn_traffic_uplink_preferences": schema.SetNestedAttribute{
				MarkdownDescription: "Array of uplink preference rules for VPN traffic",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"preferred_uplink": schema.StringAttribute{
							MarkdownDescription: "Preferred uplink for this uplink preference rule. Must be one of: 'wan1', 'wan2', 'bestForVoIP', 'loadBalancing' or 'defaultUplink'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"failover_criterion": schema.StringAttribute{
							MarkdownDescription: "Fail over criterion for this uplink preference rule. Must be one of: 'poorPerformance' or 'uplinkDown'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"traffic_filters": schema.SetNestedAttribute{
							MarkdownDescription: "Array of traffic filters for this uplink preference rule",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "type of this traffic filter. Must be one of: 'custom'",
										Required:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_id": schema.StringAttribute{
										MarkdownDescription: "ID value of 'applicationCategory' or 'application' type traffic filter",
										Required:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_protocol": schema.StringAttribute{
										MarkdownDescription: "Protocol value of this custom type traffic filter. Must be one of: 'tcp', 'udp', 'icmp6' or 'any'",
										Required:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_source_cidr": schema.StringAttribute{
										MarkdownDescription: "Source CIDR Value format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_source_port": schema.StringAttribute{
										MarkdownDescription: "Source Port Value E.g.: any, 0 (also means any), 8080, 1-1024",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_source_host": schema.Int64Attribute{
										MarkdownDescription: "Source Host ID Value in the VLAN, should be used along with 'vlan', and not exceed the vlan subnet capacity. Currently only available under a template network.",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.Int64Type,
									},
									"value_source_vlan": schema.Int64Attribute{
										MarkdownDescription: "Source VLAN ID Value  of the configured VLAN in the Meraki network. Currently only available under a template network.",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.Int64Type,
									},
									"value_source_network": schema.StringAttribute{
										MarkdownDescription: "Meraki network ID. Currently only available under a template network, and the value should be ID of either same template network, or another template network currently. E.g.: &quot;L_12345678&quot;.",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_fqdn": schema.StringAttribute{
										MarkdownDescription: "FQDN format address. Currently only availabe in 'destination' of 'vpnTrafficUplinkPreference' object. E.g.: 'www.google.com'",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_cidr": schema.StringAttribute{
										MarkdownDescription: "Destination CIDR Value  format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_host": schema.Int64Attribute{
										MarkdownDescription: "Host ID in the VLAN, should be used along with 'vlan', and not exceed the vlan subnet capacity. Currently only available under a template network.",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.Int64Type,
									},
									"value_destination_network": schema.StringAttribute{
										MarkdownDescription: "Meraki network ID. Currently only available under a template network, and the value should be ID of either same template network, or another template network currently. E.g.: &quot;L_12345678&quot;.",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_port": schema.StringAttribute{
										MarkdownDescription: "E.g.: &quot;any&quot;, &quot;0&quot; (also means &quot;any&quot;), &quot;8080&quot;, &quot;1-1024&quot;",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_vlan": schema.Int64Attribute{
										MarkdownDescription: "VLAN ID of the configured VLAN in the Meraki network. Currently only available under a template network.",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.Int64Type,
									},
								},
							},
						},
						"performance_class_builtin_performance_class_name": schema.StringAttribute{
							MarkdownDescription: "Name of builtin performance class, must be present when performanceClass type is 'builtin', and value must be one of: 'VoIP'",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"performance_class_custom_performance_class_id": schema.StringAttribute{
							MarkdownDescription: "ID of created custom performance class, must be present when performanceClass type is 'custom'",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"performance_class_type": schema.StringAttribute{
							MarkdownDescription: "Type of this performance class. Must be one of: 'builtin' or 'custom'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"failover_and_failback": schema.SingleNestedAttribute{
				MarkdownDescription: "WAN failover and failback behavior.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"immediate": schema.SingleNestedAttribute{
						MarkdownDescription: "Immediate WAN transition terminates all flows (new and existing) on current WAN when it is deemed unreliable.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Toggle for enabling or disabling immediate WAN failover and failback",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
				},
			},
			"wan_traffic_uplink_preferences": schema.SetNestedAttribute{
				MarkdownDescription: "Array of uplink preference rules for WAN traffic",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"preferred_uplink": schema.StringAttribute{
							MarkdownDescription: "Preferred uplink for this uplink preference rule. Must be one of: 'wan1' or 'wan2'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"traffic_filters": schema.SetNestedAttribute{
							MarkdownDescription: "Array of traffic filters for this uplink preference rule",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "type of this traffic filter. Must be one of: 'custom'",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_protocol": schema.StringAttribute{
										MarkdownDescription: "Protocol value of this custom type traffic filter. Must be one of: 'tcp', 'udp', 'icmp6' or 'any'",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_source_cidr": schema.StringAttribute{
										MarkdownDescription: "Source CIDR format address value, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_source_port": schema.StringAttribute{
										MarkdownDescription: "Value of Source Port E.g.: any, 0 (also means any), 8080, 1-1024",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_cidr": schema.StringAttribute{
										MarkdownDescription: "Value of Destination CIDR format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value_destination_port": schema.StringAttribute{
										MarkdownDescription: "Value of Destination Port E.g.: any, 0 (also means any), 8080, 1-1024",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceTrafficShapingUplinkSelectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceTrafficShapingUplinkSelection := *openApiClient.NewInlineObject55()

	if !data.ActiveActiveAutoVpnEnabled.IsUnknown() {
		updateNetworkApplianceTrafficShapingUplinkSelection.SetActiveActiveAutoVpnEnabled(data.ActiveActiveAutoVpnEnabled.ValueBool())
	}
	if !data.DefaultUplink.IsUnknown() {
		updateNetworkApplianceTrafficShapingUplinkSelection.SetDefaultUplink(data.DefaultUplink.ValueString())
	}
	if !data.LoadBalancingEnabled.IsUnknown() {
		updateNetworkApplianceTrafficShapingUplinkSelection.SetLoadBalancingEnabled(data.LoadBalancingEnabled.ValueBool())
	}
	if len(data.VpnTrafficUplinkPreferences) > 0 {

		var vpnTrafficUplinkPreferences []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionVpnTrafficUplinkPreferences

		for _, attribute := range data.VpnTrafficUplinkPreferences {

			var vpnTrafficUplinkPreference openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionVpnTrafficUplinkPreferences

			if !attribute.PreferredUplink.IsUnknown() {
				vpnTrafficUplinkPreference.SetPreferredUplink(attribute.PreferredUplink.ValueString())
			}

			if !attribute.FailOverCriterion.IsUnknown() {
				vpnTrafficUplinkPreference.SetFailOverCriterion(attribute.FailOverCriterion.ValueString())
			}

			var performanceClass openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionPerformanceClass

			if !attribute.PerformanceClassType.IsUnknown() {

				if attribute.PerformanceClassType == jsontypes.StringValue("builtin") {

					if !attribute.PerformanceClassBuiltinPerformanceClassName.IsUnknown() {
						performanceClass.SetBuiltinPerformanceClassName(attribute.PerformanceClassBuiltinPerformanceClassName.ValueString())
					} else {
						if resp.Diagnostics.HasError() {
							resp.Diagnostics.AddError("Error:", "missing performance_class_builtin_performance_class_name missing")
							return
						}
					}
				}

				performanceClass.SetType(attribute.PerformanceClassType.ValueString())
			}
			vpnTrafficUplinkPreference.SetPerformanceClass(performanceClass)

			if len(attribute.TrafficFilters) > 0 {
				var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
				for _, trafficFilter := range attribute.TrafficFilters {
					var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
					if !trafficFilter.Type.IsUnknown() {
						trafficFilterData.SetType(trafficFilter.Type.ValueString())
					}
					var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1
					if !trafficFilter.ValueId.IsUnknown() {
						value.SetId(trafficFilter.ValueId.ValueString())
					}
					if !trafficFilter.ValueProtocol.IsUnknown() {
						value.SetProtocol(trafficFilter.ValueProtocol.ValueString())
					}
					fmt.Println(trafficFilter)
					var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Destination
					if !trafficFilter.ValueDestinationPort.IsUnknown() {
						destination.SetPort(trafficFilter.ValueDestinationPort.ValueString())
					}
					if !trafficFilter.ValueDestinationCidr.IsUnknown() {
						destination.SetCidr(trafficFilter.ValueDestinationCidr.ValueString())
					}
					if !trafficFilter.ValueDestinationNetwork.IsUnknown() {
						destination.SetNetwork(trafficFilter.ValueDestinationNetwork.ValueString())
					}
					if !trafficFilter.ValueDestinationFqdn.IsUnknown() {
						destination.SetFqdn(trafficFilter.ValueDestinationFqdn.ValueString())
					}
					if !trafficFilter.ValueDestinationHost.IsUnknown() {
						destination.SetHost(int32(trafficFilter.ValueDestinationHost.ValueInt64()))
					}
					if !trafficFilter.ValueDestinationVlan.IsUnknown() {
						destination.SetVlan(int32(trafficFilter.ValueDestinationVlan.ValueInt64()))
					}
					fmt.Println(destination)
					value.SetDestination(destination)
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Source
					if !trafficFilter.ValueSourceCidr.IsNull() {
						source.SetCidr(trafficFilter.ValueSourceCidr.ValueString())
					}
					if !trafficFilter.ValueSourceHost.IsUnknown() {
						source.SetHost(int32(trafficFilter.ValueSourceHost.ValueInt64()))
					}
					if !trafficFilter.ValueSourcePort.IsNull() {
						source.SetPort(trafficFilter.ValueSourcePort.ValueString())
					}
					if !trafficFilter.ValueSourceVlan.IsUnknown() {
						source.SetVlan(int32(trafficFilter.ValueSourceVlan.ValueInt64()))
					}
					if !trafficFilter.ValueSourceNetwork.IsUnknown() {
						source.SetNetwork(trafficFilter.ValueSourceNetwork.ValueString())
					}
					value.SetSource(source)
					trafficFilterData.SetValue(value)
					trafficFilters = append(trafficFilters, trafficFilterData)
				}
				vpnTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
			}
			vpnTrafficUplinkPreferences = append(vpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
		}
		updateNetworkApplianceTrafficShapingUplinkSelection.SetVpnTrafficUplinkPreferences(vpnTrafficUplinkPreferences)
	}

	var failoverAndFailback openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionFailoverAndFailback
	var failoverAndFailbackImmediate openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionFailoverAndFailbackImmediate
	if !data.FailoverAndFailback.Immediate.Enabled.IsUnknown() {
		failoverAndFailbackImmediate.SetEnabled(data.FailoverAndFailback.Immediate.Enabled.ValueBool())
		failoverAndFailback.SetImmediate(failoverAndFailbackImmediate)
		updateNetworkApplianceTrafficShapingUplinkSelection.SetFailoverAndFailback(failoverAndFailback)
	}

	if len(data.WanTrafficUplinkPreferences) > 0 {
		var wanTrafficUplinkPreferences []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionWanTrafficUplinkPreferences
		for _, attribute := range data.WanTrafficUplinkPreferences {
			var wanTrafficUplinkPreference openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionWanTrafficUplinkPreferences
			if !attribute.PreferredUplink.IsUnknown() {
				wanTrafficUplinkPreference.SetPreferredUplink(attribute.PreferredUplink.ValueString())
			}
			if len(attribute.TrafficFilters) > 0 {
				var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters
				for _, trafficFilter := range attribute.TrafficFilters {
					var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters
					if !trafficFilter.Type.IsUnknown() {
						trafficFilterData.SetType(trafficFilter.Type.ValueString())
					}
					var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue
					if !trafficFilter.ValueProtocol.IsUnknown() {
						value.SetProtocol(trafficFilter.ValueProtocol.ValueString())
					}
					var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueDestination
					if !trafficFilter.ValueDestinationPort.IsUnknown() {
						destination.SetPort(trafficFilter.ValueDestinationPort.ValueString())
					}
					if !trafficFilter.ValueDestinationCidr.IsUnknown() {
						destination.SetCidr(trafficFilter.ValueDestinationCidr.ValueString())
					}
					value.SetDestination(destination)
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueSource
					if !trafficFilter.ValueSourceCidr.IsUnknown() {
						source.SetCidr(trafficFilter.ValueSourceCidr.ValueString())
					}
					if !trafficFilter.ValueSourcePort.IsUnknown() {
						source.SetPort(trafficFilter.ValueSourcePort.ValueString())
					}
					//if !trafficFilter.ValueSourceHost.Int64Value.IsUnknown() {
					//source.SetHost(int32(trafficFilter.ValueSourceHost.ValueInt64()))
					//}
					//if !trafficFilter.ValueSourceVlan.Int64Value.IsUnknown() || trafficFilter.ValueSourceVlan.Int64Value {
					//source.SetVlan(int32(trafficFilter.ValueSourceVlan.ValueInt64()))
					//}
					value.SetSource(source)
					trafficFilterData.SetValue(value)
					trafficFilters = append(trafficFilters, trafficFilterData)
				}
				wanTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
			}
			wanTrafficUplinkPreferences = append(wanTrafficUplinkPreferences, wanTrafficUplinkPreference)
		}
		updateNetworkApplianceTrafficShapingUplinkSelection.SetWanTrafficUplinkPreferences(wanTrafficUplinkPreferences)
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkSelection(updateNetworkApplianceTrafficShapingUplinkSelection).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	fmt.Println(httpResp.Body)

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
	data.ActiveActiveAutoVpnEnabled = jsontypes.BoolValue(inlineResp.GetActiveActiveAutoVpnEnabled())
	data.DefaultUplink = jsontypes.StringValue(inlineResp.GetDefaultUplink())
	data.LoadBalancingEnabled = jsontypes.BoolValue(inlineResp.GetLoadBalancingEnabled())
	data.FailoverAndFailback.Immediate.Enabled = jsontypes.BoolValue(inlineResp.FailoverAndFailback.Immediate.GetEnabled())
	data.VpnTrafficUplinkPreferences = nil
	for _, attribute := range inlineResp.GetVpnTrafficUplinkPreferences() {
		var vpnTrafficUplinkPreference VpnTrafficUplinkPreferences
		vpnTrafficUplinkPreference.FailOverCriterion = jsontypes.StringValue(attribute.GetFailOverCriterion())
		vpnTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		vpnTrafficUplinkPreference.PerformanceClassType = jsontypes.StringValue(attribute.PerformanceClass.GetType())
		vpnTrafficUplinkPreference.PerformanceClassBuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
		vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
		if vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId == jsontypes.StringValue("") {
			vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId = jsontypes.StringNull()
		}
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficData VlanTrafficFilters
			trafficData.Type = jsontypes.StringValue(traffic.GetType())
			trafficData.ValueId = jsontypes.StringValue(traffic.Value.GetId())
			trafficData.ValueProtocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficData.ValueDestinationCidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficData.ValueDestinationPort = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficData.ValueDestinationFqdn = jsontypes.StringValue(traffic.Value.Destination.GetFqdn())
			if trafficData.ValueDestinationFqdn == jsontypes.StringValue("") {
				trafficData.ValueDestinationFqdn = jsontypes.StringNull()
			}
			trafficData.ValueDestinationHost = jsontypes.Int64Value(int64(traffic.Value.Destination.GetHost()))
			if trafficData.ValueDestinationHost == jsontypes.Int64Value(0) {
				trafficData.ValueDestinationHost = jsontypes.Int64Null()
			}
			//trafficData.ValueDestinationNetwork = jsontypes.StringValue(traffic.Value.Destination.GetNetwork())
			trafficData.ValueDestinationVlan = jsontypes.Int64Value(int64(traffic.Value.Destination.GetVlan()))
			if trafficData.ValueDestinationVlan == jsontypes.Int64Value(0) {
				trafficData.ValueDestinationVlan = jsontypes.Int64Null()
			}
			trafficData.ValueSourceVlan = jsontypes.Int64Value(int64(traffic.Value.Destination.GetVlan()))
			if trafficData.ValueSourceVlan == jsontypes.Int64Value(0) {
				trafficData.ValueSourceVlan = jsontypes.Int64Null()
			}
			trafficData.ValueSourceHost = jsontypes.Int64Value(int64(traffic.Value.Destination.GetHost()))
			if trafficData.ValueSourceHost == jsontypes.Int64Value(0) {
				trafficData.ValueSourceHost = jsontypes.Int64Null()
			}
			trafficData.ValueSourceCidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficData.ValueSourcePort = jsontypes.StringValue(traffic.Value.Source.GetPort())
			//trafficData.ValueSourceNetwork = jsontypes.StringValue(traffic.Value.Source.GetNetwork())
			fmt.Println(trafficData)
			vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
		}
		data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
		fmt.Println(vpnTrafficUplinkPreference)
		fmt.Println(data.VpnTrafficUplinkPreferences)
	}

	for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
		var wanTrafficUplinkPreference WanTrafficUplinkPreferences
		wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficFilter WanTrafficFilters
			trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
			trafficFilter.ValueProtocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficFilter.ValueDestinationCidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficFilter.ValueDestinationPort = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficFilter.ValueSourceCidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficFilter.ValueSourcePort = jsontypes.StringValue(traffic.Value.Source.GetPort())
			wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
		}
		data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
	}

	fmt.Println(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceTrafficShapingUplinkSelectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).Execute()
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
	data.ActiveActiveAutoVpnEnabled = jsontypes.BoolValue(inlineResp.GetActiveActiveAutoVpnEnabled())
	data.DefaultUplink = jsontypes.StringValue(inlineResp.GetDefaultUplink())
	data.LoadBalancingEnabled = jsontypes.BoolValue(inlineResp.GetLoadBalancingEnabled())
	data.FailoverAndFailback.Immediate.Enabled = jsontypes.BoolValue(inlineResp.FailoverAndFailback.Immediate.GetEnabled())
	data.VpnTrafficUplinkPreferences = nil
	for _, attribute := range inlineResp.GetVpnTrafficUplinkPreferences() {
		var vpnTrafficUplinkPreference VpnTrafficUplinkPreferences
		vpnTrafficUplinkPreference.FailOverCriterion = jsontypes.StringValue(attribute.GetFailOverCriterion())
		vpnTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		vpnTrafficUplinkPreference.PerformanceClassType = jsontypes.StringValue(attribute.PerformanceClass.GetType())
		vpnTrafficUplinkPreference.PerformanceClassBuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
		vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
		if vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId.IsUnknown() {
			vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId = jsontypes.StringNull()
		}
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficData VlanTrafficFilters
			trafficData.Type = jsontypes.StringValue(traffic.GetType())
			trafficData.ValueId = jsontypes.StringValue(traffic.Value.GetId())
			trafficData.ValueProtocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficData.ValueDestinationCidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficData.ValueDestinationPort = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficData.ValueDestinationFqdn = jsontypes.StringValue(traffic.Value.Destination.GetFqdn())
			//trafficData.ValueDestinationHost = jsontypes.Int64Value(int64(traffic.Value.Destination.GetHost()))
			trafficData.ValueDestinationNetwork = jsontypes.StringValue(traffic.Value.Destination.GetNetwork())
			trafficData.ValueDestinationVlan = jsontypes.Int64Value(int64(traffic.Value.Destination.GetVlan()))
			trafficData.ValueSourceVlan = jsontypes.Int64Value(int64(traffic.Value.Destination.GetVlan()))
			trafficData.ValueSourceHost = jsontypes.Int64Value(int64(traffic.Value.Destination.GetHost()))
			trafficData.ValueSourceCidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficData.ValueSourcePort = jsontypes.StringValue(traffic.Value.Source.GetPort())
			trafficData.ValueSourceNetwork = jsontypes.StringValue(traffic.Value.Source.GetNetwork())
			vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
		}
		data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
	}

	for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
		var wanTrafficUplinkPreference WanTrafficUplinkPreferences
		wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficFilter WanTrafficFilters
			trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
			trafficFilter.ValueProtocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficFilter.ValueDestinationCidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficFilter.ValueDestinationPort = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficFilter.ValueSourceCidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficFilter.ValueSourcePort = jsontypes.StringValue(traffic.Value.Source.GetPort())
			wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
		}
		data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
	}

	fmt.Println(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceTrafficShapingUplinkSelectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	/*
		updateNetworkApplianceTrafficShapingUplinkSelection := *openApiClient.NewInlineObject55()

		if !data.ActiveActiveAutoVpnEnabled.IsUnknown() {
			updateNetworkApplianceTrafficShapingUplinkSelection.SetActiveActiveAutoVpnEnabled(data.ActiveActiveAutoVpnEnabled.ValueBool())
		}
		if !data.DefaultUplink.IsUnknown() {
			updateNetworkApplianceTrafficShapingUplinkSelection.SetDefaultUplink(data.DefaultUplink.ValueString())
		}
		if !data.LoadBalancingEnabled.IsUnknown() {
			updateNetworkApplianceTrafficShapingUplinkSelection.SetLoadBalancingEnabled(data.LoadBalancingEnabled.ValueBool())
		}
		if len(data.VpnTrafficUplinkPreferences) > 0 {
			var vpnTrafficUplinkPreferences []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionVpnTrafficUplinkPreferences
			for _, attribute := range data.VpnTrafficUplinkPreferences {
				var vpnTrafficUplinkPreference openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionVpnTrafficUplinkPreferences
				if !attribute.PreferredUplink.IsUnknown() {
					vpnTrafficUplinkPreference.SetPreferredUplink(attribute.PreferredUplink.ValueString())
				}
				if !attribute.FailOverCriterion.IsUnknown() {
					vpnTrafficUplinkPreference.SetFailOverCriterion(attribute.FailOverCriterion.ValueString())
				}
				var performanceClass openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionPerformanceClass
				if !attribute.PerformanceClassBuiltinPerformanceClassName.IsUnknown() {
					performanceClass.SetBuiltinPerformanceClassName(attribute.PerformanceClassBuiltinPerformanceClassName.ValueString())
				}
				if !attribute.PerformanceClassCustomPerformanceClassId.IsUnknown() {
					performanceClass.SetCustomPerformanceClassId(attribute.PerformanceClassCustomPerformanceClassId.ValueString())
				}
				if !attribute.PerformanceClassType.IsUnknown() {
					performanceClass.SetType(attribute.PerformanceClassType.ValueString())
				}
				vpnTrafficUplinkPreference.SetPerformanceClass(performanceClass)

				if len(attribute.TrafficFilters) > 0 {
					var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
					for _, trafficFilter := range attribute.TrafficFilters {
						var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
						if !trafficFilter.Type.IsUnknown() {
							trafficFilterData.SetType(trafficFilter.Type.ValueString())
						}
						var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1
						if !trafficFilter.ValueId.IsUnknown() {
							value.SetId(trafficFilter.ValueId.ValueString())
						}
						if !trafficFilter.ValueProtocol.IsUnknown() {
							value.SetProtocol(trafficFilter.ValueProtocol.ValueString())
						}
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Destination
						if !trafficFilter.ValueDestinationPort.IsUnknown() {
							destination.SetPort(trafficFilter.ValueDestinationPort.ValueString())
						}
						if !trafficFilter.ValueDestinationCidr.IsUnknown() {
							destination.SetCidr(trafficFilter.ValueDestinationCidr.ValueString())
						}
						value.SetDestination(destination)
						var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Source
						if !trafficFilter.ValueSourceCidr.IsUnknown() {
							source.SetCidr(trafficFilter.ValueSourceCidr.ValueString())
						}
						if !trafficFilter.ValueSourceHost.IsUnknown() {
							source.SetHost(int32(trafficFilter.ValueSourceHost.ValueInt64()))
						}
						if !trafficFilter.ValueSourcePort.IsUnknown() {
							source.SetPort(trafficFilter.ValueSourcePort.ValueString())
						}
						if !trafficFilter.ValueSourceVlan.IsUnknown() {
							source.SetVlan(int32(trafficFilter.ValueSourceVlan.ValueInt64()))
						}
						value.SetSource(source)
						trafficFilterData.SetValue(value)
						trafficFilters = append(trafficFilters, trafficFilterData)
					}
					vpnTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
				}
				vpnTrafficUplinkPreferences = append(vpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
			}
			updateNetworkApplianceTrafficShapingUplinkSelection.SetVpnTrafficUplinkPreferences(vpnTrafficUplinkPreferences)
		}

		var failoverAndFailback openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionFailoverAndFailback
		var failoverAndFailbackImmediate openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionFailoverAndFailbackImmediate
		if !data.FailoverAndFailback.Immediate.Enabled.IsUnknown() {
			failoverAndFailbackImmediate.SetEnabled(data.FailoverAndFailback.Immediate.Enabled.ValueBool())
			failoverAndFailback.SetImmediate(failoverAndFailbackImmediate)
			updateNetworkApplianceTrafficShapingUplinkSelection.SetFailoverAndFailback(failoverAndFailback)
		}

		if len(data.WanTrafficUplinkPreferences) > 0 {
			var wanTrafficUplinkPreferences []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionWanTrafficUplinkPreferences
			for _, attribute := range data.WanTrafficUplinkPreferences {
				var wanTrafficUplinkPreference openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionWanTrafficUplinkPreferences
				if !attribute.PreferredUplink.IsUnknown() {
					wanTrafficUplinkPreference.SetPreferredUplink(attribute.PreferredUplink.ValueString())
				}
				if len(attribute.TrafficFilters) > 0 {
					var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters
					for _, trafficFilter := range attribute.TrafficFilters {
						var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters
						if !trafficFilter.Type.IsUnknown() {
							trafficFilterData.SetType(trafficFilter.Type.ValueString())
						}
						var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue
						if !trafficFilter.ValueProtocol.IsUnknown() {
							value.SetProtocol(trafficFilter.ValueProtocol.ValueString())
						}
						if !trafficFilter.ValueDestinationPort.IsUnknown() || !trafficFilter.ValueDestinationCidr.IsUnknown() {
							var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueDestination
							destination.SetPort(trafficFilter.ValueDestinationPort.ValueString())
							destination.SetCidr(trafficFilter.ValueDestinationCidr.ValueString())
							value.SetDestination(destination)
						}
						var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueSource
						if !trafficFilter.ValueSourceCidr.IsUnknown() || !trafficFilter.ValueSourcePort.IsUnknown() {
							if !trafficFilter.ValueSourceHost.Int64Value.IsUnknown() || !trafficFilter.ValueSourceVlan.Int64Value.IsUnknown() {
								source.SetCidr(trafficFilter.ValueSourceCidr.ValueString())
								source.SetPort(trafficFilter.ValueSourcePort.ValueString())
								source.SetHost(int32(trafficFilter.ValueSourceHost.ValueInt64()))
								source.SetVlan(int32(trafficFilter.ValueSourceVlan.ValueInt64()))
								value.SetSource(source)
							}
						}
						trafficFilterData.SetValue(value)
						trafficFilters = append(trafficFilters, trafficFilterData)
					}
					wanTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
				}
				wanTrafficUplinkPreferences = append(wanTrafficUplinkPreferences, wanTrafficUplinkPreference)
			}
			updateNetworkApplianceTrafficShapingUplinkSelection.SetWanTrafficUplinkPreferences(wanTrafficUplinkPreferences)
		}

		inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkSelection(updateNetworkApplianceTrafficShapingUplinkSelection).Execute()
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
		data.ActiveActiveAutoVpnEnabled = jsontypes.BoolValue(inlineResp.GetActiveActiveAutoVpnEnabled())
		data.DefaultUplink = jsontypes.StringValue(inlineResp.GetDefaultUplink())
		data.LoadBalancingEnabled = jsontypes.BoolValue(inlineResp.GetLoadBalancingEnabled())
		data.FailoverAndFailback.Immediate.Enabled = jsontypes.BoolValue(inlineResp.FailoverAndFailback.Immediate.GetEnabled())
		for _, attribute := range inlineResp.GetVpnTrafficUplinkPreferences() {
			var vpnTrafficUplinkPreference VpnTrafficUplinkPreferences
			vpnTrafficUplinkPreference.FailOverCriterion = jsontypes.StringValue(attribute.GetFailOverCriterion())
			vpnTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
			vpnTrafficUplinkPreference.PerformanceClassBuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
			vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
			vpnTrafficUplinkPreference.PerformanceClassCustomPerformanceClassId = jsontypes.StringNull()
			vpnTrafficUplinkPreference.PerformanceClassType = jsontypes.StringValue(attribute.PerformanceClass.GetType())
			for _, traffic := range attribute.GetTrafficFilters() {
				var trafficData VlanTrafficFilters
				trafficData.Type = jsontypes.StringValue(traffic.GetType())
				trafficData.ValueId = jsontypes.StringValue(traffic.Value.GetId())
				trafficData.ValueProtocol = jsontypes.StringValue(traffic.Value.GetProtocol())
				trafficData.ValueDestinationCidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
				trafficData.ValueDestinationPort = jsontypes.StringValue(traffic.Value.Destination.GetPort())
				trafficData.ValueSourceVlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
				trafficData.ValueSourceHost = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
				trafficData.ValueSourceCidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
				trafficData.ValueSourcePort = jsontypes.StringValue(traffic.Value.Source.GetPort())
				vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
			}
			data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
		}
		for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
			var wanTrafficUplinkPreference WanTrafficUplinkPreferences
			wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
			for _, traffic := range attribute.GetTrafficFilters() {
				var trafficFilter WanTrafficFilters
				trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
				trafficFilter.ValueProtocol = jsontypes.StringValue(traffic.Value.GetProtocol())
				trafficFilter.ValueDestinationCidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
				trafficFilter.ValueDestinationPort = jsontypes.StringValue(traffic.Value.Destination.GetPort())
				trafficFilter.ValueSourceVlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
				trafficFilter.ValueSourceHost = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
				trafficFilter.ValueSourceCidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
				trafficFilter.ValueSourcePort = jsontypes.StringValue(traffic.Value.Source.GetPort())
				wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
			}
			data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	*/

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceTrafficShapingUplinkSelectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	/*

		updateNetworkApplianceTrafficShapingUplinkSelection := *openApiClient.NewInlineObject55()

		if len(data.VpnTrafficUplinkPreferences) > 0 {
			var vpnTrafficUplinkPreferences []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionVpnTrafficUplinkPreferences
			for _, attribute := range data.VpnTrafficUplinkPreferences {
				var vpnTrafficUplinkPreference openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionVpnTrafficUplinkPreferences
				var p openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionPerformanceClass
				vpnTrafficUplinkPreference.SetPerformanceClass(p)

				if len(attribute.TrafficFilters) > 0 {
					var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
					for _, trafficFilter := range attribute.TrafficFilters {
						var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
						if !trafficFilter.Type.IsUnknown() {
							trafficFilterData.SetType(trafficFilter.Type.ValueString())
						}
						var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1
						if !trafficFilter.ValueId.IsUnknown() {
							value.SetId(trafficFilter.ValueId.ValueString())
						}
						if !trafficFilter.ValueProtocol.IsUnknown() {
							value.SetProtocol(trafficFilter.ValueProtocol.ValueString())
						}
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Destination
						if !trafficFilter.ValueDestinationPort.IsUnknown() {
							destination.SetPort(trafficFilter.ValueDestinationPort.ValueString())
						}
						if !trafficFilter.ValueDestinationCidr.IsUnknown() {
							destination.SetCidr(trafficFilter.ValueDestinationCidr.ValueString())
						}
						value.SetDestination(destination)
						var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Source
						if !trafficFilter.ValueSourceCidr.IsUnknown() {
							source.SetCidr(trafficFilter.ValueSourceCidr.ValueString())
						}
						if !trafficFilter.ValueSourceHost.IsUnknown() {
							source.SetHost(int32(trafficFilter.ValueSourceHost.ValueInt64()))
						}
						if !trafficFilter.ValueSourcePort.IsUnknown() {
							source.SetPort(trafficFilter.ValueSourcePort.ValueString())
						}
						if !trafficFilter.ValueSourceVlan.IsUnknown() {
							source.SetVlan(int32(trafficFilter.ValueSourceVlan.ValueInt64()))
						}
						value.SetSource(source)
						trafficFilterData.SetValue(value)
						trafficFilters = append(trafficFilters, trafficFilterData)
					}
					vpnTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
				}
				vpnTrafficUplinkPreferences = append(vpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
			}
			updateNetworkApplianceTrafficShapingUplinkSelection.SetVpnTrafficUplinkPreferences(vpnTrafficUplinkPreferences)
		}

		var failoverAndFailback openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionFailoverAndFailback
		var failoverAndFailbackImmediate openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionFailoverAndFailbackImmediate
		if !data.FailoverAndFailback.Immediate.Enabled.IsUnknown() {
			failoverAndFailback.SetImmediate(failoverAndFailbackImmediate)
			updateNetworkApplianceTrafficShapingUplinkSelection.SetFailoverAndFailback(failoverAndFailback)
		}

		if len(data.WanTrafficUplinkPreferences) > 0 {
			var wanTrafficUplinkPreferences []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionWanTrafficUplinkPreferences
			for _, attribute := range data.WanTrafficUplinkPreferences {
				var wanTrafficUplinkPreference openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionWanTrafficUplinkPreferences
				if !attribute.PreferredUplink.IsUnknown() {
				}
				if len(attribute.TrafficFilters) > 0 {
					var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters
					for _, trafficFilter := range attribute.TrafficFilters {
						var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters
						if !trafficFilter.Type.IsUnknown() {
							trafficFilterData.SetType(trafficFilter.Type.ValueString())
						}
						var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue
						if !trafficFilter.ValueProtocol.IsUnknown() {
							value.SetProtocol(trafficFilter.ValueProtocol.ValueString())
						}
						if !trafficFilter.ValueDestinationPort.IsUnknown() || !trafficFilter.ValueDestinationCidr.IsUnknown() {
							var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueDestination
							destination.SetPort(trafficFilter.ValueDestinationPort.ValueString())
							destination.SetCidr(trafficFilter.ValueDestinationCidr.ValueString())
							value.SetDestination(destination)
						}
						var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueSource
						if !trafficFilter.ValueSourceCidr.IsUnknown() || !trafficFilter.ValueSourcePort.IsUnknown() {
							if !trafficFilter.ValueSourceHost.Int64Value.IsUnknown() || !trafficFilter.ValueSourceVlan.Int64Value.IsUnknown() {
								source.SetCidr(trafficFilter.ValueSourceCidr.ValueString())
								source.SetPort(trafficFilter.ValueSourcePort.ValueString())
								source.SetHost(int32(trafficFilter.ValueSourceHost.ValueInt64()))
								source.SetVlan(int32(trafficFilter.ValueSourceVlan.ValueInt64()))
								value.SetSource(source)
							}
						}
						trafficFilterData.SetValue(value)
						trafficFilters = append(trafficFilters, trafficFilterData)
					}
					wanTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
				}
				wanTrafficUplinkPreferences = append(wanTrafficUplinkPreferences, wanTrafficUplinkPreference)
			}
			updateNetworkApplianceTrafficShapingUplinkSelection.SetWanTrafficUplinkPreferences(wanTrafficUplinkPreferences)
		}
		_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkSelection(updateNetworkApplianceTrafficShapingUplinkSelection).Execute()
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
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	*/

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
