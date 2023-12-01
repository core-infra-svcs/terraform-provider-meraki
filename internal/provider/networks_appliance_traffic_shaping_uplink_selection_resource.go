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
var _ resource.Resource = &NetworksApplianceTrafficShapingUplinkSelectionResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceTrafficShapingUplinkSelectionResource{}

func NewNetworksApplianceTrafficShapingUplinkSelectionResource() resource.Resource {
	return &NetworksApplianceTrafficShapingUplinkSelectionResource{}
}

// NetworksApplianceTrafficShapingUplinkSelectionResource defines the resource implementation.
type NetworksApplianceTrafficShapingUplinkSelectionResource struct {
	client *openApiClient.APIClient
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                          jsontypes.String             `tfsdk:"id"`
	NetworkId                   jsontypes.String             `tfsdk:"network_id" json:"network_id"`
	ActiveActiveAutoVpnEnabled  jsontypes.Bool               `json:"activeActiveAutoVpnEnabled" tfsdk:"active_active_auto_vpn_enabled"`
	DefaultUplink               jsontypes.String             `json:"defaultUplink" tfsdk:"default_uplink"`
	LoadBalancingEnabled        jsontypes.Bool               `json:"loadBalancingEnabled" tfsdk:"load_balancing_enabled"`
	FailoverAndFailback         failoverAndFailback          `json:"failoverAndFailback" tfsdk:"failover_and_failback"`
	WanTrafficUplinkPreferences []wanTrafficUplinkPreference `json:"wanTrafficUplinkPreferences" tfsdk:"wan_traffic_uplink_preferences"`
	VpnTrafficUplinkPreferences []vpnTrafficUplinkPreference `json:"vpnTrafficUplinkPreferences" tfsdk:"vpn_traffic_uplink_preferences"`
}

type failoverAndFailback struct {
	Immediate struct {
		Enabled jsontypes.Bool `json:"enabled" tfsdk:"enabled"`
	} `json:"immediate" tfsdk:"immediate"`
}

type wanDestination struct {
	Port jsontypes.String `json:"port" tfsdk:"port"`
	CIDR jsontypes.String `json:"cidr" tfsdk:"cidr"`
}

type wanSource struct {
	Port jsontypes.String `json:"port" tfsdk:"port"`
	CIDR jsontypes.String `json:"cidr" tfsdk:"cidr"`
	VLAN jsontypes.Int64  `json:"vlan" tfsdk:"vlan"`
	Host jsontypes.Int64  `json:"host" tfsdk:"host"`
}

type wanValue struct {
	Protocol    jsontypes.String `json:"protocol" tfsdk:"protocol"`
	Source      wanSource        `json:"source" tfsdk:"source"`
	Destination wanDestination   `json:"destination" tfsdk:"destination"`
}

type wanTrafficFilter struct {
	Type  jsontypes.String `json:"type" tfsdk:"type"`
	Value wanValue         `json:"value" tfsdk:"value"`
}

type wanTrafficUplinkPreference struct {
	TrafficFilters  []wanTrafficFilter `json:"trafficFilters" tfsdk:"traffic_filters"`
	PreferredUplink jsontypes.String   `json:"preferredUplink" tfsdk:"preferred_uplink"`
}

type performanceClass struct {
	Type                   jsontypes.String `json:"type" tfsdk:"type"`
	BuiltinPerformanceName jsontypes.String `json:"builtinPerformanceClassName" tfsdk:"builtin_performance_class_name"`
	CustomPerformanceID    jsontypes.String `json:"customPerformanceClassId" tfsdk:"custom_performance_class_id"`
}

type vpnDestination struct {
	Port    jsontypes.String `json:"port" tfsdk:"port"`
	CIDR    jsontypes.String `json:"cidr" tfsdk:"cidr"`
	Network jsontypes.String `json:"network" tfsdk:"network"`
	VLAN    jsontypes.Int64  `json:"vlan" tfsdk:"vlan"`
	Host    jsontypes.Int64  `json:"host" tfsdk:"host"`
	FQDN    jsontypes.String `json:"fqdn" tfsdk:"fqdn"`
}

type vpnSource struct {
	Port    jsontypes.String `json:"port" tfsdk:"port"`
	CIDR    jsontypes.String `json:"cidr" tfsdk:"cidr"`
	Network jsontypes.String `json:"network" tfsdk:"network"`
	VLAN    jsontypes.Int64  `json:"vlan" tfsdk:"vlan"`
	Host    jsontypes.Int64  `json:"host" tfsdk:"host"`
}

type vpnValue struct {
	Protocol    jsontypes.String `json:"protocol" tfsdk:"protocol"`
	Source      vpnSource        `json:"source" tfsdk:"source"`
	Destination vpnDestination   `json:"destination" tfsdk:"destination"`
	Id          jsontypes.String `tfsdk:"id"`
}

type vpnTrafficFilter struct {
	Type  jsontypes.String `json:"type" tfsdk:"type"`
	Value vpnValue         `json:"value" tfsdk:"value"`
}

type vpnTrafficUplinkPreference struct {
	TrafficFilters    []vpnTrafficFilter `json:"trafficFilters" tfsdk:"traffic_filters"`
	PreferredUplink   jsontypes.String   `json:"preferredUplink" tfsdk:"preferred_uplink"`
	FailOverCriterion jsontypes.String   `json:"failOverCriterion" tfsdk:"fail_over_criterion"`
	PerformanceClass  performanceClass   `json:"performanceClass" tfsdk:"performance_class"`
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
				MarkdownDescription: "Whether active-active AutoVPN is enabled",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"default_uplink": schema.StringAttribute{
				MarkdownDescription: "The default uplink. Must be one of: 'wan1' or 'wan2'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("wan1", "wan2"),
				},
			},
			"load_balancing_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether load balancing is enabled",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"failover_and_failback": schema.SingleNestedAttribute{
				MarkdownDescription: "WAN failover and failback",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"immediate": schema.SingleNestedAttribute{
						MarkdownDescription: "Immediate WAN failover and failback",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether immediate WAN failover and failback is enabled",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
				},
			},
			"wan_traffic_uplink_preferences": schema.SetNestedAttribute{
				MarkdownDescription: "Uplink preference rules for WAN traffic",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"traffic_filters": schema.SetNestedAttribute{
							MarkdownDescription: "Array of traffic filters for this uplink preference rule",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "Traffic filter type. Must be \"custom\"",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
										Validators: []validator.String{
											stringvalidator.OneOf("custom"),
										},
									},
									"value": schema.SingleNestedAttribute{
										MarkdownDescription: "Value of traffic filter",
										Optional:            true,
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"protocol": schema.StringAttribute{
												MarkdownDescription: "Protocol value of this custom type traffic filter. Must be one of: 'tcp', 'udp', 'icmp6' or 'any'",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
												Validators: []validator.String{
													stringvalidator.OneOf("tcp", "udp", "icmp6", "any"),
												},
											},
											"source": schema.SingleNestedAttribute{
												MarkdownDescription: "Source of 'custom' type traffic filter",
												Optional:            true,
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: \"any\", \"0\" (also means \"any\"), \"8080\", \"1-1024\"",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"cidr": schema.StringAttribute{
														MarkdownDescription: "SCIDR format address (e.g.\"192.168.10.1\", which is the same as \"192.168.10.1/32\"), or \"any\". Cannot be used in combination with the \"vlan\" property",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"vlan": schema.Int64Attribute{
														MarkdownDescription: "VLAN ID of the configured VLAN in the Meraki network. Cannot be used in combination with the \"cidr\" property and is currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
													"host": schema.Int64Attribute{
														MarkdownDescription: "Host ID in the VLAN. Should not exceed the VLAN subnet capacity. Must be used along with the \"vlan\" property and is currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
												},
											},
											"destination": schema.SingleNestedAttribute{
												MarkdownDescription: "Destination of 'custom' type traffic filter",
												Optional:            true,
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: \"any\", \"0\" (also means \"any\"), \"8080\", \"1-1024\"",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"cidr": schema.StringAttribute{
														MarkdownDescription: "CIDR format address (e.g.\"192.168.10.1\", which is the same as \"192.168.10.1/32\"), or \"any\". Cannot be used in combination with the \"vlan\" property",
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
						"preferred_uplink": schema.StringAttribute{
							MarkdownDescription: "Preferred uplink for this uplink preference rule. Must be one of: 'wan1' or 'wan2'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf("wan1", "wan2"),
							},
						},
					},
				},
			},
			"vpn_traffic_uplink_preferences": schema.ListNestedAttribute{
				MarkdownDescription: "Array of uplink preference rules for VPN traffic",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"traffic_filters": schema.ListNestedAttribute{
							MarkdownDescription: "Array of traffic filters for this uplink preference rule",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "Traffic filter type. Must be one of: 'applicationCategory', 'application' or 'custom'",
										Optional:            true,
										Computed:            true,
										CustomType:          jsontypes.StringType,
										Validators: []validator.String{
											stringvalidator.OneOf("applicationCategory", "application", "custom"),
										},
									},
									"value": schema.SingleNestedAttribute{
										MarkdownDescription: "value of traffic filter",
										Optional:            true,
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"protocol": schema.StringAttribute{
												MarkdownDescription: "Protocol value of this custom type traffic filter. Must be one of: 'tcp', 'udp', 'icmp6' or 'any'",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
												Validators: []validator.String{
													stringvalidator.OneOf("tcp", "udp", "icmp6", "any"),
												},
											},
											"source": schema.SingleNestedAttribute{
												MarkdownDescription: "Source of 'custom' type traffic filter",
												Optional:            true,
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: \"any\", \"0\" (also means \"any\"), \"8080\", \"1-1024\"",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"cidr": schema.StringAttribute{
														MarkdownDescription: "SCIDR format address (e.g.\"192.168.10.1\", which is the same as \"192.168.10.1/32\"), or \"any\". Cannot be used in combination with the \"vlan\" property",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"network": schema.StringAttribute{
														MarkdownDescription: "Meraki network ID. Currently only available under a template network, and the value should be ID of either same template network, or another template network currently. E.g.: \"L_12345678\".",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"vlan": schema.Int64Attribute{
														MarkdownDescription: "VLAN ID of the configured VLAN in the Meraki network. Cannot be used in combination with the \"cidr\" property and is currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
													"host": schema.Int64Attribute{
														MarkdownDescription: "Host ID in the VLAN. Should not exceed the VLAN subnet capacity. Must be used along with the \"vlan\" property and is currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
												},
											},
											"destination": schema.SingleNestedAttribute{
												MarkdownDescription: "destination of 'custom' type traffic filter",
												Optional:            true,
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"port": schema.StringAttribute{
														MarkdownDescription: "E.g.: \"any\", \"0\" (also means \"any\"), \"8080\", \"1-1024\"",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"cidr": schema.StringAttribute{
														MarkdownDescription: "CIDR format address (e.g.\"192.168.10.1\", which is the same as \"192.168.10.1/32\"), or \"any\". Cannot be used in combination with the \"vlan\" property",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"network": schema.StringAttribute{
														MarkdownDescription: "Meraki network ID. Currently only available under a template network, and the value should be ID of either same template network, or another template network currently. E.g.: \"L_12345678\".",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
													"vlan": schema.Int64Attribute{
														MarkdownDescription: "VLAN ID of the configured VLAN in the Meraki network. Cannot be used in combination with the \"cidr\" property and is currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
													"host": schema.Int64Attribute{
														MarkdownDescription: "Host ID in the VLAN. Should not exceed the VLAN subnet capacity. Must be used along with the \"vlan\" property and is currently only available under a template network.",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.Int64Type,
													},
													"fqdn": schema.StringAttribute{
														MarkdownDescription: "FQDN format address. Cannot be used in combination with the \"cidr\" or \"fqdn\" property and is currently only available in the \"vpnDestination\" object of the \"vpnTrafficUplinkPreference\" object. E.g.: \"www.google.com\"",
														Optional:            true,
														Computed:            true,
														CustomType:          jsontypes.StringType,
													},
												},
											},
											"id": schema.StringAttribute{
												MarkdownDescription: "traffic filter id",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
										},
									},
								},
							},
						},
						"preferred_uplink": schema.StringAttribute{
							MarkdownDescription: "Preferred uplink for uplink preference rule. Must be one of: 'wan1', 'wan2', 'bestForVoIP', 'loadBalancing' or 'defaultUplink'",
							Optional:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf("bestForVoIP", "defaultUplink", "loadBalancing", "wan1", "wan2"),
							},
						},
						"fail_over_criterion": schema.StringAttribute{
							MarkdownDescription: "Fail over criterion for uplink preference rule. Must be one of: 'poorPerformance' or 'uplinkDown'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf("poorPerformance", "uplinkDown"),
							},
						},
						"performance_class": schema.SingleNestedAttribute{
							MarkdownDescription: "Type of this performance class. Must be one of: 'builtin' or 'custom'",
							Optional:            true,
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type of this performance class. Must be one of: 'builtin' or 'custom'",
									Optional:            true,
									CustomType:          jsontypes.StringType,
									Validators: []validator.String{
										stringvalidator.OneOf("builtin", "custom"),
									},
								},
								"builtin_performance_class_name": schema.StringAttribute{
									MarkdownDescription: "Name of builtin performance class. Must be present when performanceClass type is 'builtin' and value must be one of: 'VoIP'",
									Optional:            true,
									CustomType:          jsontypes.StringType,
									Validators: []validator.String{
										stringvalidator.OneOf("VoIP"),
									},
								},
								"custom_performance_class_id": schema.StringAttribute{
									MarkdownDescription: "ID of created custom performance class, must be present when performanceClass type is 'custom'",
									Optional:            true,
									CustomType:          jsontypes.StringType,
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

func updateNetworksApplianceTrafficShapingUplinkSelectionPayload(data *resourceModel) openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequest {

	updateNetworkApplianceTrafficShapingUplinkSelection := *openApiClient.NewUpdateNetworkApplianceTrafficShapingUplinkSelectionRequest()

	// ActiveActiveAutoVpnEnabled
	if !data.ActiveActiveAutoVpnEnabled.IsUnknown() {
		updateNetworkApplianceTrafficShapingUplinkSelection.SetActiveActiveAutoVpnEnabled(data.ActiveActiveAutoVpnEnabled.ValueBool())
	}

	// DefaultUplink
	if !data.DefaultUplink.IsUnknown() {
		updateNetworkApplianceTrafficShapingUplinkSelection.SetDefaultUplink(data.DefaultUplink.ValueString())
	}

	// LoadBalancingEnabled
	if !data.LoadBalancingEnabled.IsUnknown() {
		updateNetworkApplianceTrafficShapingUplinkSelection.SetLoadBalancingEnabled(data.LoadBalancingEnabled.ValueBool())
	}

	// failoverAndFailback
	if !data.FailoverAndFailback.Immediate.Enabled.IsUnknown() {
		var failoverAndFailback openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestFailoverAndFailback
		var failoverAndFailbackImmediate openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestFailoverAndFailbackImmediate

		failoverAndFailbackImmediate.SetEnabled(data.FailoverAndFailback.Immediate.Enabled.ValueBool())
		failoverAndFailback.SetImmediate(failoverAndFailbackImmediate)

		updateNetworkApplianceTrafficShapingUplinkSelection.SetFailoverAndFailback(failoverAndFailback)
	}

	// WanTrafficUplinkPreferences
	if len(data.WanTrafficUplinkPreferences) > 0 {
		var wanTrafficUplinkPreferences []openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInner
		for _, attribute := range data.WanTrafficUplinkPreferences {

			// PreferredUplink
			var wanTrafficUplinkPreference openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInner
			if !attribute.PreferredUplink.IsUnknown() {
				wanTrafficUplinkPreference.SetPreferredUplink(attribute.PreferredUplink.ValueString())
			}

			// TrafficFilters
			if len(attribute.TrafficFilters) > 0 {
				var trafficFilters []openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInnerTrafficFiltersInner
				for _, trafficFilter := range attribute.TrafficFilters {

					var trafficFilterData openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInnerTrafficFiltersInner

					// Type
					if !trafficFilter.Type.IsUnknown() {
						trafficFilterData.SetType(trafficFilter.Type.ValueString())
					}

					// value
					var value openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInnerTrafficFiltersInnerValue

					// Protocol
					if !trafficFilter.Value.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.Value.Protocol.ValueString())
					}

					// Source
					var source openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInnerTrafficFiltersInnerValueSource

					// Port
					if !trafficFilter.Value.Source.Port.IsUnknown() {
						source.SetPort(trafficFilter.Value.Source.Port.ValueString())
					}

					// CIDR
					if !trafficFilter.Value.Source.CIDR.IsUnknown() {
						source.SetCidr(trafficFilter.Value.Source.CIDR.ValueString())
					}

					// VLAN
					if !trafficFilter.Value.Source.VLAN.Int64Value.IsUnknown() {
						source.SetVlan(int32(trafficFilter.Value.Source.VLAN.ValueInt64()))
					}

					// Host
					if !trafficFilter.Value.Source.Host.Int64Value.IsUnknown() {
						source.SetHost(int32(trafficFilter.Value.Source.Host.ValueInt64()))
					}

					value.SetSource(source)
					trafficFilterData.SetValue(value)
					trafficFilters = append(trafficFilters, trafficFilterData)

					// destination
					var destination openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestWanTrafficUplinkPreferencesInnerTrafficFiltersInnerValueDestination

					// Port
					if !trafficFilter.Value.Destination.Port.IsUnknown() {
						destination.SetPort(trafficFilter.Value.Destination.Port.ValueString())
					}

					// CIDR
					if !trafficFilter.Value.Destination.CIDR.IsUnknown() {
						destination.SetCidr(trafficFilter.Value.Destination.CIDR.ValueString())
					}
					value.SetDestination(destination)

				}

				wanTrafficUplinkPreference.SetTrafficFilters(trafficFilters)

			}

			wanTrafficUplinkPreferences = append(wanTrafficUplinkPreferences, wanTrafficUplinkPreference)

		}

		updateNetworkApplianceTrafficShapingUplinkSelection.SetWanTrafficUplinkPreferences(wanTrafficUplinkPreferences)
	}

	// VpnTrafficUplinkPreferences
	if len(data.VpnTrafficUplinkPreferences) > 0 {
		var vpnTrafficUplinkPreferences []openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInner

		for _, attribute := range data.VpnTrafficUplinkPreferences {

			var vpnTrafficUplinkPreference openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInner

			// TrafficFilters
			if len(attribute.TrafficFilters) > 0 {

				// Traffic Filter
				var trafficFilters []openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInnerTrafficFiltersInner
				for _, trafficFilter := range attribute.TrafficFilters {
					var trafficFilterData openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInnerTrafficFiltersInner

					// Type
					if !trafficFilter.Type.IsUnknown() {
						trafficFilterData.SetType(trafficFilter.Type.ValueString())
					}

					// value
					var value openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInnerTrafficFiltersInnerValue

					// ID
					if !trafficFilter.Value.Id.IsUnknown() {
						value.SetId(trafficFilter.Value.Id.ValueString())
					}

					// Protool
					if !trafficFilter.Value.Protocol.IsUnknown() {
						value.SetProtocol(trafficFilter.Value.Protocol.ValueString())
					}

					// Source
					var source openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInnerTrafficFiltersInnerValueSource

					// Port
					if !trafficFilter.Value.Source.Port.IsNull() {
						source.SetPort(trafficFilter.Value.Source.Port.ValueString())
					}

					// Cidr
					if !trafficFilter.Value.Source.CIDR.IsNull() {
						source.SetCidr(trafficFilter.Value.Source.CIDR.ValueString())
					}

					// Network
					if !trafficFilter.Value.Source.Network.IsUnknown() {
						source.SetNetwork(trafficFilter.Value.Source.Network.ValueString())
					}

					// Vlan
					if !trafficFilter.Value.Source.VLAN.IsUnknown() {
						source.SetVlan(int32(trafficFilter.Value.Source.VLAN.ValueInt64()))
					}

					// Host
					if !trafficFilter.Value.Source.Host.IsUnknown() {
						source.SetHost(int32(trafficFilter.Value.Source.Host.ValueInt64()))
					}
					value.SetSource(source)

					// vpnDestination
					var destination openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInnerTrafficFiltersInnerValueDestination

					// Port
					if !trafficFilter.Value.Destination.Port.IsUnknown() {
						destination.SetPort(trafficFilter.Value.Destination.Port.ValueString())
					}

					// Cidr
					if !trafficFilter.Value.Destination.CIDR.IsUnknown() {
						destination.SetCidr(trafficFilter.Value.Destination.CIDR.ValueString())
					}

					// Network
					if !trafficFilter.Value.Destination.Network.IsUnknown() {
						destination.SetNetwork(trafficFilter.Value.Destination.Network.ValueString())
					}

					// Vlan
					if !trafficFilter.Value.Destination.VLAN.IsUnknown() {
						destination.SetVlan(int32(trafficFilter.Value.Destination.VLAN.ValueInt64()))
					}

					// Host
					if !trafficFilter.Value.Destination.Host.IsUnknown() {
						destination.SetHost(int32(trafficFilter.Value.Destination.Host.ValueInt64()))
					}

					// FQDN
					if !trafficFilter.Value.Destination.FQDN.IsUnknown() {
						destination.SetFqdn(trafficFilter.Value.Destination.FQDN.ValueString())
					}

					value.SetDestination(destination)

					trafficFilterData.SetValue(value)
					trafficFilters = append(trafficFilters, trafficFilterData)
				}

				vpnTrafficUplinkPreference.SetTrafficFilters(trafficFilters)
			}

			// PreferredUplink
			if !attribute.PreferredUplink.IsUnknown() {
				vpnTrafficUplinkPreference.SetPreferredUplink(attribute.PreferredUplink.ValueString())
			}

			// FailOverCriterion
			if !attribute.FailOverCriterion.IsUnknown() {
				vpnTrafficUplinkPreference.SetFailOverCriterion(attribute.FailOverCriterion.ValueString())
			}

			// performanceClass
			var performanceClass openApiClient.UpdateNetworkApplianceTrafficShapingUplinkSelectionRequestVpnTrafficUplinkPreferencesInnerPerformanceClass

			// type
			if !attribute.PerformanceClass.Type.IsUnknown() {
				performanceClass.SetType(attribute.PerformanceClass.Type.ValueString())
			}

			// BuiltinPerformanceName
			if !attribute.PerformanceClass.BuiltinPerformanceName.IsUnknown() {
				performanceClass.SetBuiltinPerformanceClassName(attribute.PerformanceClass.BuiltinPerformanceName.ValueString())
			}

			// CustomPerformanceID
			if !attribute.PerformanceClass.CustomPerformanceID.IsUnknown() {
				performanceClass.SetCustomPerformanceClassId(attribute.PerformanceClass.CustomPerformanceID.ValueString())
			}

			vpnTrafficUplinkPreference.SetPerformanceClass(performanceClass)

			vpnTrafficUplinkPreferences = append(vpnTrafficUplinkPreferences, vpnTrafficUplinkPreference)
		}

		updateNetworkApplianceTrafficShapingUplinkSelection.SetVpnTrafficUplinkPreferences(vpnTrafficUplinkPreferences)
	}

	return updateNetworkApplianceTrafficShapingUplinkSelection
}

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceTrafficShapingUplinkSelection := updateNetworksApplianceTrafficShapingUplinkSelectionPayload(data)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkSelectionRequest(updateNetworkApplianceTrafficShapingUplinkSelection).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).Execute()
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
		return
	} else {
		resp.Diagnostics.Append()
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

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceTrafficShapingUplinkSelection := updateNetworksApplianceTrafficShapingUplinkSelectionPayload(data)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkSelectionRequest(updateNetworkApplianceTrafficShapingUplinkSelection).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

func (r *NetworksApplianceTrafficShapingUplinkSelectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceTrafficShapingUplinkSelection := *openApiClient.NewUpdateNetworkApplianceTrafficShapingUplinkSelectionRequest()

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkSelection(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkSelectionRequest(updateNetworkApplianceTrafficShapingUplinkSelection).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

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
