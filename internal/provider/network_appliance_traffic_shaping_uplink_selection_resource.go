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
	PreferredUplink   jsontypes.String     `tfsdk:"preferred_uplink"`
	FailOverCriterion jsontypes.String     `tfsdk:"failover_criterion"`
	PerformanceClass  PerformanceClass     `tfsdk:"performance_class"`
	TrafficFilters    []VlanTrafficFilters `tfsdk:"traffic_filters"`
}

type PerformanceClass struct {
	BuiltinPerformanceClassName jsontypes.String `tfsdk:"builtin_performance_class_name"`
	CustomPerformanceClassId    jsontypes.String `tfsdk:"custom_performance_class_id"`
	Type                        jsontypes.String `tfsdk:"type"`
}

type FailoverAndFailback struct {
	Immediate Immediate `tfsdk:"immediate"`
}

type Immediate struct {
	Enabled jsontypes.Bool `tfsdk:"enabled"`
}

type WanTrafficUplinkPreferences struct {
	TrafficFilters  []TrafficFilters `tfsdk:"traffic_filters"`
	PreferredUplink jsontypes.String `tfsdk:"preferred_uplink"`
}

type TrafficFilters struct {
	Type  jsontypes.String `tfsdk:"type"`
	Value Value            `tfsdk:"value"`
}

type VlanTrafficFilters struct {
	Type      jsontypes.String `tfsdk:"type"`
	VlanValue VlanValue        `tfsdk:"value"`
}

type VlanValue struct {
	id          jsontypes.String `tfsdk:"id"`
	Protocol    jsontypes.String `tfsdk:"protocol"`
	Source      Source           `tfsdk:"source"`
	Destination Destination      `tfsdk:"destination"`
}

type Value struct {
	Protocol    jsontypes.String `tfsdk:"protocol"`
	Source      Source           `tfsdk:"source"`
	Destination Destination      `tfsdk:"destination"`
}

type Source struct {
	Port jsontypes.String `tfsdk:"port"`
	Cidr jsontypes.String `tfsdk:"cidr"`
	Vlan jsontypes.Int64  `tfsdk:"vlan"`
	Host jsontypes.Int64  `tfsdk:"host"`
}

type Destination struct {
	Port jsontypes.String `tfsdk:"port"`
	Cidr jsontypes.String `tfsdk:"cidr"`
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
							Optional:            true,
							Computed:            true,
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
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
									},
									"value": schema.SingleNestedAttribute{
										MarkdownDescription: "Value object of this traffic filter.",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												MarkdownDescription: "ID of 'applicationCategory' or 'application' type traffic filter",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
											"protocol": schema.StringAttribute{
												MarkdownDescription: "Protocol of this custom type traffic filter. Must be one of: 'tcp', 'udp', 'icmp6' or 'any'",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
											"source": schema.SingleNestedAttribute{
												MarkdownDescription: "Source of this custom type traffic filter",
												Required:            true,
												Attributes: map[string]schema.Attribute{
													"cidr": schema.StringAttribute{
														MarkdownDescription: "CIDR format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: any, 0 (also means any), 8080, 1-1024",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"host": schema.Int64Attribute{
														MarkdownDescription: "Host ID in the VLAN, should be used along with 'vlan', and not exceed the vlan subnet capacity. Currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
													"vlan": schema.Int64Attribute{
														MarkdownDescription: "VLAN ID of the configured VLAN in the Meraki network. Currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
												},
											},
											"destination": schema.SingleNestedAttribute{
												MarkdownDescription: "Destination of this custom type traffic filter",
												Required:            true,
												Attributes: map[string]schema.Attribute{
													"cidr": schema.StringAttribute{
														MarkdownDescription: "CIDR format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: any, 0 (also means any), 8080, 1-1024",
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
						"performance_class": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"builtin_performance_class_name": schema.StringAttribute{
									MarkdownDescription: "Name of builtin performance class, must be present when performanceClass type is 'builtin', and value must be one of: 'VoIP'",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"custom_performance_class_id": schema.StringAttribute{
									MarkdownDescription: "ID of created custom performance class, must be present when performanceClass type is 'custom'",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Type of this performance class. Must be one of: 'builtin' or 'custom'",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
							},
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
				Required:            true,
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
									"value": schema.SingleNestedAttribute{
										MarkdownDescription: "Value object of this traffic filter.",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"protocol": schema.StringAttribute{
												MarkdownDescription: "Protocol of this custom type traffic filter. Must be one of: 'tcp', 'udp', 'icmp6' or 'any'",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
											"source": schema.SingleNestedAttribute{
												MarkdownDescription: "Source of this custom type traffic filter",
												Required:            true,
												Attributes: map[string]schema.Attribute{
													"cidr": schema.StringAttribute{
														MarkdownDescription: "CIDR format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: any, 0 (also means any), 8080, 1-1024",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"host": schema.Int64Attribute{
														MarkdownDescription: "Host ID in the VLAN, should be used along with 'vlan', and not exceed the vlan subnet capacity. Currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
													"vlan": schema.Int64Attribute{
														MarkdownDescription: "VLAN ID of the configured VLAN in the Meraki network. Currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
												},
											},
											"destination": schema.SingleNestedAttribute{
												MarkdownDescription: "Destination of this custom type traffic filter",
												Required:            true,
												Attributes: map[string]schema.Attribute{
													"cidr": schema.StringAttribute{
														MarkdownDescription: "CIDR format address, or any. E.g.: 192.168.10.0/24, 192.168.10.1 (same as 192.168.10.1/32), 0.0.0.0/0 (same as any)",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: any, 0 (also means any), 8080, 1-1024",
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
			if !attribute.PerformanceClass.BuiltinPerformanceClassName.IsUnknown() {
				performanceClass.SetBuiltinPerformanceClassName(attribute.PerformanceClass.BuiltinPerformanceClassName.ValueString())
			}
			if !attribute.PerformanceClass.CustomPerformanceClassId.IsUnknown() {
				performanceClass.SetCustomPerformanceClassId(attribute.PerformanceClass.CustomPerformanceClassId.ValueString())
			}
			if !attribute.PerformanceClass.Type.IsUnknown() {
				performanceClass.SetType(attribute.PerformanceClass.Type.ValueString())
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
					if !trafficFilter.VlanValue.id.IsUnknown() {
						value.SetId(trafficFilter.VlanValue.id.ValueString())
					}
					if !trafficFilter.VlanValue.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.VlanValue.Protocol.ValueString())
					}
					if !(trafficFilter.VlanValue.Destination.Port.IsUnknown() || trafficFilter.VlanValue.Destination.Cidr.IsUnknown()) {
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Destination
						destination.SetPort(trafficFilter.VlanValue.Destination.Port.ValueString())
						destination.SetCidr(trafficFilter.VlanValue.Destination.Cidr.ValueString())
						value.SetDestination(destination)
					}
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Source
					if !trafficFilter.VlanValue.Source.Cidr.IsUnknown() || !trafficFilter.VlanValue.Source.Port.IsUnknown() {
						if !trafficFilter.VlanValue.Source.Host.Int64Value.IsUnknown() || !trafficFilter.VlanValue.Source.Vlan.Int64Value.IsUnknown() {
							source.SetHost(int32(trafficFilter.VlanValue.Source.Host.ValueInt64()))
							source.SetVlan(int32(trafficFilter.VlanValue.Source.Vlan.ValueInt64()))
							source.SetCidr(trafficFilter.VlanValue.Source.Cidr.ValueString())
							source.SetPort(trafficFilter.VlanValue.Source.Port.ValueString())
							value.SetSource(source)
						}
					}
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
					if !trafficFilter.Value.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.Value.Protocol.ValueString())
					}
					if !trafficFilter.Value.Destination.Port.IsUnknown() || !trafficFilter.Value.Destination.Cidr.IsUnknown() {
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueDestination
						destination.SetPort(trafficFilter.Value.Destination.Port.ValueString())
						destination.SetCidr(trafficFilter.Value.Destination.Cidr.ValueString())
						value.SetDestination(destination)
					}
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueSource
					if !trafficFilter.Value.Source.Cidr.IsUnknown() || !trafficFilter.Value.Source.Port.IsUnknown() {
						if !trafficFilter.Value.Source.Host.Int64Value.IsUnknown() || !trafficFilter.Value.Source.Vlan.Int64Value.IsUnknown() {
							source.SetCidr(trafficFilter.Value.Source.Cidr.ValueString())
							source.SetPort(trafficFilter.Value.Source.Port.ValueString())
							source.SetHost(int32(trafficFilter.Value.Source.Host.ValueInt64()))
							source.SetVlan(int32(trafficFilter.Value.Source.Vlan.ValueInt64()))
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
		vpnTrafficUplinkPreference.PerformanceClass.BuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
		vpnTrafficUplinkPreference.PerformanceClass.CustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
		vpnTrafficUplinkPreference.PerformanceClass.Type = jsontypes.StringValue(attribute.PerformanceClass.GetType())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficData VlanTrafficFilters
			trafficData.Type = jsontypes.StringValue(traffic.GetType())
			trafficData.VlanValue.id = jsontypes.StringValue(traffic.Value.GetId())
			trafficData.VlanValue.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficData.VlanValue.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficData.VlanValue.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficData.VlanValue.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficData.VlanValue.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficData.VlanValue.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficData.VlanValue.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
		}
		data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
	}

	for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
		var wanTrafficUplinkPreference WanTrafficUplinkPreferences
		wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficFilter TrafficFilters
			trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
			trafficFilter.Value.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficFilter.Value.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficFilter.Value.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficFilter.Value.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficFilter.Value.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficFilter.Value.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficFilter.Value.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
		}
		data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
	}

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

	for _, attribute := range inlineResp.GetVpnTrafficUplinkPreferences() {
		var vpnTrafficUplinkPreference VpnTrafficUplinkPreferences
		vpnTrafficUplinkPreference.FailOverCriterion = jsontypes.StringValue(attribute.GetFailOverCriterion())
		vpnTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		vpnTrafficUplinkPreference.PerformanceClass.BuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
		vpnTrafficUplinkPreference.PerformanceClass.CustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
		vpnTrafficUplinkPreference.PerformanceClass.Type = jsontypes.StringValue(attribute.PerformanceClass.GetType())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficData VlanTrafficFilters
			trafficData.Type = jsontypes.StringValue(traffic.GetType())
			trafficData.VlanValue.id = jsontypes.StringValue(traffic.Value.GetId())
			trafficData.VlanValue.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficData.VlanValue.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficData.VlanValue.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficData.VlanValue.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficData.VlanValue.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficData.VlanValue.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficData.VlanValue.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
		}
		data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
	}

	for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
		var wanTrafficUplinkPreference WanTrafficUplinkPreferences
		wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficFilter TrafficFilters
			trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
			trafficFilter.Value.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficFilter.Value.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficFilter.Value.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficFilter.Value.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficFilter.Value.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficFilter.Value.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficFilter.Value.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
		}
		data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
	}
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
			if !attribute.PerformanceClass.BuiltinPerformanceClassName.IsUnknown() {
				performanceClass.SetBuiltinPerformanceClassName(attribute.PerformanceClass.BuiltinPerformanceClassName.ValueString())
			}
			if !attribute.PerformanceClass.CustomPerformanceClassId.IsUnknown() {
				performanceClass.SetCustomPerformanceClassId(attribute.PerformanceClass.CustomPerformanceClassId.ValueString())
			}
			if !attribute.PerformanceClass.Type.IsUnknown() {
				performanceClass.SetType(attribute.PerformanceClass.Type.ValueString())
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
					if !trafficFilter.VlanValue.id.IsUnknown() {
						value.SetId(trafficFilter.VlanValue.id.ValueString())
					}
					if !trafficFilter.VlanValue.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.VlanValue.Protocol.ValueString())
					}
					if !(trafficFilter.VlanValue.Destination.Port.IsUnknown() || trafficFilter.VlanValue.Destination.Cidr.IsUnknown()) {
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Destination
						destination.SetPort(trafficFilter.VlanValue.Destination.Port.ValueString())
						destination.SetCidr(trafficFilter.VlanValue.Destination.Cidr.ValueString())
						value.SetDestination(destination)
					}
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Source
					if !trafficFilter.VlanValue.Source.Cidr.IsUnknown() || !trafficFilter.VlanValue.Source.Port.IsUnknown() {
						if !trafficFilter.VlanValue.Source.Host.Int64Value.IsUnknown() || !trafficFilter.VlanValue.Source.Vlan.Int64Value.IsUnknown() {
							source.SetHost(int32(trafficFilter.VlanValue.Source.Host.ValueInt64()))
							source.SetVlan(int32(trafficFilter.VlanValue.Source.Vlan.ValueInt64()))
							source.SetCidr(trafficFilter.VlanValue.Source.Cidr.ValueString())
							source.SetPort(trafficFilter.VlanValue.Source.Port.ValueString())
							value.SetSource(source)
						}
					}
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
					if !trafficFilter.Value.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.Value.Protocol.ValueString())
					}
					if !trafficFilter.Value.Destination.Port.IsUnknown() || !trafficFilter.Value.Destination.Cidr.IsUnknown() {
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueDestination
						destination.SetPort(trafficFilter.Value.Destination.Port.ValueString())
						destination.SetCidr(trafficFilter.Value.Destination.Cidr.ValueString())
						value.SetDestination(destination)
					}
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueSource
					if !trafficFilter.Value.Source.Cidr.IsUnknown() || !trafficFilter.Value.Source.Port.IsUnknown() {
						if !trafficFilter.Value.Source.Host.Int64Value.IsUnknown() || !trafficFilter.Value.Source.Vlan.Int64Value.IsUnknown() {
							source.SetCidr(trafficFilter.Value.Source.Cidr.ValueString())
							source.SetPort(trafficFilter.Value.Source.Port.ValueString())
							source.SetHost(int32(trafficFilter.Value.Source.Host.ValueInt64()))
							source.SetVlan(int32(trafficFilter.Value.Source.Vlan.ValueInt64()))
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
		vpnTrafficUplinkPreference.PerformanceClass.BuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
		vpnTrafficUplinkPreference.PerformanceClass.CustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
		vpnTrafficUplinkPreference.PerformanceClass.Type = jsontypes.StringValue(attribute.PerformanceClass.GetType())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficData VlanTrafficFilters
			trafficData.Type = jsontypes.StringValue(traffic.GetType())
			trafficData.VlanValue.id = jsontypes.StringValue(traffic.Value.GetId())
			trafficData.VlanValue.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficData.VlanValue.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficData.VlanValue.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficData.VlanValue.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficData.VlanValue.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficData.VlanValue.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficData.VlanValue.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
		}
		data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
	}
	for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
		var wanTrafficUplinkPreference WanTrafficUplinkPreferences
		wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficFilter TrafficFilters
			trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
			trafficFilter.Value.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficFilter.Value.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficFilter.Value.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficFilter.Value.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficFilter.Value.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficFilter.Value.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficFilter.Value.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
		}
		data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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
			var p openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionPerformanceClass
			if !attribute.PerformanceClass.BuiltinPerformanceClassName.IsUnknown() {
				p.SetBuiltinPerformanceClassName(attribute.PerformanceClass.BuiltinPerformanceClassName.ValueString())
			}
			if !attribute.PerformanceClass.CustomPerformanceClassId.IsUnknown() {
				p.SetCustomPerformanceClassId(attribute.PerformanceClass.CustomPerformanceClassId.ValueString())
			}
			if !attribute.PerformanceClass.Type.IsUnknown() {
				p.SetType(attribute.PerformanceClass.Type.ValueString())
			}
			vpnTrafficUplinkPreference.SetPerformanceClass(p)

			if len(attribute.TrafficFilters) > 0 {
				var trafficFilters []openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
				for _, trafficFilter := range attribute.TrafficFilters {
					var trafficFilterData openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionTrafficFilters1
					if !trafficFilter.Type.IsUnknown() {
						trafficFilterData.SetType(trafficFilter.Type.ValueString())
					}
					var value openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1
					if !trafficFilter.VlanValue.id.IsUnknown() {
						value.SetId(trafficFilter.VlanValue.id.ValueString())
					}
					if !trafficFilter.VlanValue.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.VlanValue.Protocol.ValueString())
					}
					if !(trafficFilter.VlanValue.Destination.Port.IsUnknown() || trafficFilter.VlanValue.Destination.Cidr.IsUnknown()) {
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Destination
						destination.SetPort(trafficFilter.VlanValue.Destination.Port.ValueString())
						destination.SetCidr(trafficFilter.VlanValue.Destination.Cidr.ValueString())
						value.SetDestination(destination)
					}
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValue1Source
					if !trafficFilter.VlanValue.Source.Cidr.IsUnknown() {
						source.SetCidr(trafficFilter.VlanValue.Source.Cidr.ValueString())
					}
					if !trafficFilter.VlanValue.Source.Port.IsUnknown() {
						source.SetPort(trafficFilter.VlanValue.Source.Port.ValueString())
					}
					if !trafficFilter.VlanValue.Source.Host.Int64Value.IsUnknown() {
						source.SetHost(int32(trafficFilter.VlanValue.Source.Host.ValueInt64()))
					}
					if !trafficFilter.VlanValue.Source.Vlan.Int64Value.IsUnknown() {
						source.SetVlan(int32(trafficFilter.VlanValue.Source.Vlan.ValueInt64()))
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
					if !trafficFilter.Value.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.Value.Protocol.ValueString())
					}
					if !trafficFilter.Value.Destination.Port.IsUnknown() || !trafficFilter.Value.Destination.Cidr.IsUnknown() {
						var destination openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueDestination
						destination.SetPort(trafficFilter.Value.Destination.Port.ValueString())
						destination.SetCidr(trafficFilter.Value.Destination.Cidr.ValueString())
						value.SetDestination(destination)
					}
					var source openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkSelectionValueSource
					if !trafficFilter.Value.Source.Cidr.IsUnknown() {
						source.SetCidr(trafficFilter.Value.Source.Cidr.ValueString())
					}
					if !trafficFilter.Value.Source.Port.IsUnknown() {
						source.SetPort(trafficFilter.Value.Source.Port.ValueString())
					}
					if !trafficFilter.Value.Source.Host.Int64Value.IsUnknown() {
						source.SetHost(int32(trafficFilter.Value.Source.Host.ValueInt64()))
					}
					if !trafficFilter.Value.Source.Vlan.Int64Value.IsUnknown() {
						source.SetVlan(int32(trafficFilter.Value.Source.Vlan.ValueInt64()))
					}
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
		vpnTrafficUplinkPreference.PerformanceClass.BuiltinPerformanceClassName = jsontypes.StringValue(attribute.PerformanceClass.GetBuiltinPerformanceClassName())
		vpnTrafficUplinkPreference.PerformanceClass.CustomPerformanceClassId = jsontypes.StringValue(attribute.PerformanceClass.GetCustomPerformanceClassId())
		vpnTrafficUplinkPreference.PerformanceClass.Type = jsontypes.StringValue(attribute.PerformanceClass.GetType())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficData VlanTrafficFilters
			trafficData.Type = jsontypes.StringValue(traffic.GetType())
			trafficData.VlanValue.id = jsontypes.StringValue(traffic.Value.GetId())
			trafficData.VlanValue.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficData.VlanValue.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficData.VlanValue.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficData.VlanValue.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficData.VlanValue.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficData.VlanValue.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficData.VlanValue.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			vpnTrafficUplinkPreference.TrafficFilters = append(vpnTrafficUplinkPreference.TrafficFilters, trafficData)
		}
		data.VpnTrafficUplinkPreferences = append(data.VpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
	}
	for _, attribute := range inlineResp.GetWanTrafficUplinkPreferences() {
		var wanTrafficUplinkPreference WanTrafficUplinkPreferences
		wanTrafficUplinkPreference.PreferredUplink = jsontypes.StringValue(attribute.GetPreferredUplink())
		for _, traffic := range attribute.GetTrafficFilters() {
			var trafficFilter TrafficFilters
			trafficFilter.Type = jsontypes.StringValue(traffic.GetType())
			trafficFilter.Value.Protocol = jsontypes.StringValue(traffic.Value.GetProtocol())
			trafficFilter.Value.Destination.Cidr = jsontypes.StringValue(traffic.Value.Destination.GetCidr())
			trafficFilter.Value.Destination.Port = jsontypes.StringValue(traffic.Value.Destination.GetPort())
			trafficFilter.Value.Source.Vlan = jsontypes.Int64Value(int64(traffic.Value.Source.GetVlan()))
			trafficFilter.Value.Source.Host = jsontypes.Int64Value(int64(traffic.Value.Source.GetHost()))
			trafficFilter.Value.Source.Cidr = jsontypes.StringValue(traffic.Value.Source.GetCidr())
			trafficFilter.Value.Source.Port = jsontypes.StringValue(traffic.Value.Source.GetPort())
			wanTrafficUplinkPreference.TrafficFilters = append(wanTrafficUplinkPreference.TrafficFilters, trafficFilter)
		}
		data.WanTrafficUplinkPreferences = append(data.WanTrafficUplinkPreferences, wanTrafficUplinkPreference)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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
