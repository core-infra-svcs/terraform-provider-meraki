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
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksSwitchAccessPoliciesResource{}
	_ resource.ResourceWithConfigure   = &NetworksSwitchAccessPoliciesResource{}
	_ resource.ResourceWithImportState = &NetworksSwitchAccessPoliciesResource{}
)

func NewNetworksSwitchAccessPoliciesResource() resource.Resource {
	return &NetworksSwitchAccessPoliciesResource{}
}

// NetworksSwitchAccessPoliciesResource defines the resource implementation.
type NetworksSwitchAccessPoliciesResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchAccessPoliciesResourceModel describes the resource data model.
type NetworksSwitchAccessPoliciesResourceModel struct {
	Id                             jsontypes.String                                           `tfsdk:"id"`
	AccessPolicyNumber             jsontypes.String                                           `tfsdk:"access_policy_number"`
	NetworkId                      jsontypes.String                                           `tfsdk:"network_id"`
	AccessPolicyType               jsontypes.String                                           `tfsdk:"access_policy_type"`
	Dot1xControlDirection          jsontypes.String                                           `tfsdk:"dot_1x_control_direction"`
	GuestVLANId                    jsontypes.Int64                                            `tfsdk:"guest_vlan_id"`
	HostMode                       jsontypes.String                                           `tfsdk:"host_mode"`
	IncreaseAccessSpeed            jsontypes.Bool                                             `tfsdk:"increase_access_speed"`
	Name                           jsontypes.String                                           `tfsdk:"name"`
	RadiusAccountingEnabled        jsontypes.Bool                                             `tfsdk:"radius_accounting_enabled"`
	RadiusAccountingServers        []NetworkSwitchAccessPolicyRadiusServersResourceModelRules `tfsdk:"radius_accounting_servers"`
	RadiusCoaSupportEnabled        jsontypes.Bool                                             `tfsdk:"radius_coa_support_enabled"`
	RadiusGroupAttribute           jsontypes.String                                           `tfsdk:"radius_group_attribute"`
	RadiusServers                  []NetworkSwitchAccessPolicyRadiusServersResourceModelRules `tfsdk:"radius_servers"`
	RadiusTestingEnabled           jsontypes.Bool                                             `tfsdk:"radius_testing_enabled"`
	RadiusCriticalAuth             RadiusCriticalAuth                                         `tfsdk:"radius_critical_auth"`
	RadiusFailedAuthVlanId         jsontypes.Int64                                            `tfsdk:"radius_failed_auth_vlan_id"`
	RadiusAuthenticationInterval   jsontypes.Int64                                            `tfsdk:"radius_authentication_interval"`
	UrlRedirectWalledGardenEnabled jsontypes.Bool                                             `tfsdk:"url_redirect_walled_garden_enabled"`
	UrlRedirectWalledGardenRanges  jsontypes.Set[jsontypes.String]                            `tfsdk:"url_redirect_walled_garden_ranges"`
	VoiceVlanClients               jsontypes.Bool                                             `tfsdk:"voice_vlan_clients"`
}

type NetworkSwitchAccessPolicyRadiusServersResourceModelRules struct {
	Host   jsontypes.String `tfsdk:"host"`
	Port   jsontypes.Int64  `tfsdk:"port"`
	Secret jsontypes.String `tfsdk:"secret"`
}

type RadiusCriticalAuth struct {
	DataVlanId        jsontypes.Int64 `tfsdk:"data_vlan_id"`
	VoiceVlanId       jsontypes.Int64 `tfsdk:"voice_vlan_id"`
	SuspendPortBounce jsontypes.Bool  `tfsdk:"suspend_port_bounce"`
}

func (r *NetworksSwitchAccessPoliciesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_access_policies"
}

func (r *NetworksSwitchAccessPoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchAccessPolicies",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"access_policy_number": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
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

			"access_policy_type": schema.StringAttribute{
				MarkdownDescription: "Access Type of the policy. Automatically 'Hybrid authentication' when hostMode is 'Multi-Domain'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},

			"dot_1x_control_direction": schema.StringAttribute{
				MarkdownDescription: "Supports either 'both' or 'inbound'. Set to 'inbound' to allow unauthorized egress on the switchport. Set to 'both' to control both traffic directions with authorization. Defaults to 'both'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},

			"guest_vlan_id": schema.Int64Attribute{
				MarkdownDescription: "ID for the guest VLAN allow unauthorized devices access to limited network resources",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},

			"host_mode": schema.StringAttribute{
				MarkdownDescription: "Choose the Host Mode for the access policy.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},

			"increase_access_speed": schema.BoolAttribute{
				MarkdownDescription: "Enabling this option will make switches execute 802.1X and MAC-bypass authentication simultaneously so that clients authenticate faster. Only required when accessPolicyType is 'Hybrid Authentication.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the access policy",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},

			"radius_accounting_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable to send start, interim-update and stop messages to a configured RADIUS accounting server for tracking connected clients",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},

			"radius_accounting_servers": schema.SetNestedAttribute{
				MarkdownDescription: "List of RADIUS accounting servers to require connecting devices to authenticate against before granting network access",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},

			"radius_coa_support_enabled": schema.BoolAttribute{
				MarkdownDescription: "Change of authentication for RADIUS re-authentication and disconnection",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},

			"radius_group_attribute": schema.StringAttribute{
				MarkdownDescription: `Acceptable values are "" for None, or "11" for Group Policies ACL`,
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},

			"radius_servers": schema.SetNestedAttribute{
				MarkdownDescription: "List of RADIUS servers to require connecting devices to authenticate against before granting network access",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},

			"radius_testing_enabled": schema.BoolAttribute{
				MarkdownDescription: "If enabled, Meraki devices will periodically send access-request messages to these RADIUS servers",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},

			"radius_critical_auth": schema.SingleNestedAttribute{
				MarkdownDescription: "Critical auth settings for when authentication is rejected by the RADIUS server",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"data_vlan_id": schema.Int64Attribute{
						MarkdownDescription: "",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"voice_vlan_id": schema.Int64Attribute{
						MarkdownDescription: "",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"suspend_port_bounce": schema.BoolAttribute{
						MarkdownDescription: "",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},

			"radius_failed_auth_vlan_id": schema.Int64Attribute{
				MarkdownDescription: "VLAN that clients will be placed on when RADIUS authentication fails. Will be null if hostMode is Multi-Auth",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},

			"radius_authentication_interval": schema.Int64Attribute{
				MarkdownDescription: "Re-authentication period in seconds. Will be null if hostMode is Multi-Auth",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},

			"url_redirect_walled_garden_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable to restrict access for clients to a response_objectific set of IP addresses or hostnames prior to authentication",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},

			"url_redirect_walled_garden_ranges": schema.SetAttribute{
				MarkdownDescription: "IP address ranges, in CIDR notation, to restrict access for clients to a specific set of IP addresses or hostnames prior to authentication",
				//ElementType: types.StringType,
				CustomType: jsontypes.SetType[jsontypes.String](),
				Optional:   true,
				Computed:   true,
			},

			"voice_vlan_clients": schema.BoolAttribute{
				MarkdownDescription: "CDP/LLDP capable voice clients will be able to use this VLAN. Automatically true when hostMode is 'Multi-Domain'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
		},
	}
}

func (r *NetworksSwitchAccessPoliciesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchAccessPoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchAccessPoliciesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	radiusServers := []openApiClient.NetworksNetworkIdSwitchAccessPoliciesRadiusServers1{}
	for _, radiusServer := range data.RadiusServers {
		port := int32(radiusServer.Port.Int64Value.ValueInt64())
		newServer := openApiClient.NewNetworksNetworkIdSwitchAccessPoliciesRadiusServers1(radiusServer.Host.String(), port, radiusServer.Secret.ValueString())
		radiusServers = append(radiusServers, *newServer)
	}

	radiusAccountingServers := []openApiClient.NetworksNetworkIdSwitchAccessPoliciesRadiusAccountingServers1{}
	for _, radiusServer := range data.RadiusAccountingServers {
		port := int32(radiusServer.Port.Int64Value.ValueInt64())
		newServer := openApiClient.NewNetworksNetworkIdSwitchAccessPoliciesRadiusAccountingServers1(radiusServer.Host.String(), port, radiusServer.Secret.ValueString())
		radiusAccountingServers = append(radiusAccountingServers, *newServer)
	}

	createNetworkSwitchAccessPolicy := openApiClient.NewInlineObject110(data.Name.String(), radiusServers, data.RadiusTestingEnabled.ValueBool(), data.RadiusCoaSupportEnabled.ValueBool(), data.RadiusAccountingEnabled.ValueBool(), data.HostMode.ValueString(), data.UrlRedirectWalledGardenEnabled.ValueBool())
	createNetworkSwitchAccessPolicy.RadiusAccountingServers = radiusAccountingServers

	i := int32(data.RadiusCriticalAuth.DataVlanId.ValueInt64())
	voiceVlanID := int32(data.RadiusCriticalAuth.VoiceVlanId.ValueInt64())
	criticalAuth := openApiClient.NewNetworksNetworkIdSwitchAccessPoliciesRadiusCriticalAuth()
	criticalAuth.SetDataVlanId(i)
	criticalAuth.SetVoiceVlanId(voiceVlanID)
	criticalAuth.SetSuspendPortBounce(data.RadiusCriticalAuth.SuspendPortBounce.ValueBool())

	radius := openApiClient.NewNetworksNetworkIdSwitchAccessPoliciesRadius()
	radius.SetCriticalAuth(*criticalAuth)
	createNetworkSwitchAccessPolicy.SetRadius(*radius)
	var walledGardenRanges []string
	for _, urlRange := range data.UrlRedirectWalledGardenRanges.Elements() {
		walledGardenRanges = append(walledGardenRanges, urlRange.String())
	}
	createNetworkSwitchAccessPolicy.SetUrlRedirectWalledGardenRanges(walledGardenRanges)

	_, httpResp, err := r.client.SwitchApi.CreateNetworkSwitchAccessPolicy(context.Background(), data.NetworkId.ValueString()).CreateNetworkSwitchAccessPolicy(*createNetworkSwitchAccessPolicy).Execute()
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

func (r *NetworksSwitchAccessPoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchAccessPoliciesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SwitchApi.GetNetworkSwitchAccessPolicy(ctx, data.NetworkId.String(), data.AccessPolicyNumber.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
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
}

func (r *NetworksSwitchAccessPoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksSwitchAccessPoliciesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	radiusServers := []openApiClient.NetworksNetworkIdSwitchAccessPoliciesRadiusServers1{}
	for _, radiusServer := range data.RadiusServers {
		port := int32(radiusServer.Port.Int64Value.ValueInt64())
		newServer := openApiClient.NewNetworksNetworkIdSwitchAccessPoliciesRadiusServers1(radiusServer.Host.String(), port, radiusServer.Secret.ValueString())
		radiusServers = append(radiusServers, *newServer)
	}

	updateNetworkSwitchAccessPolicy := openApiClient.NewInlineObject111()
	updateNetworkSwitchAccessPolicy.RadiusServers = radiusServers
	radiusTestingEnabled := data.RadiusTestingEnabled.ValueBool()
	updateNetworkSwitchAccessPolicy.RadiusTestingEnabled = &radiusTestingEnabled
	radiusCoaSupportEnabled := data.RadiusCoaSupportEnabled.ValueBool()
	updateNetworkSwitchAccessPolicy.RadiusCoaSupportEnabled = &radiusCoaSupportEnabled
	radiusAccountingEnbled := data.RadiusAccountingEnabled.ValueBool()
	updateNetworkSwitchAccessPolicy.RadiusAccountingEnabled = &radiusAccountingEnbled
	hostMode := data.HostMode.ValueString()
	updateNetworkSwitchAccessPolicy.HostMode = &hostMode
	urlRedirectWalledGardenEnabled := data.UrlRedirectWalledGardenEnabled.ValueBool()
	radiusReauthInterval := int32(data.RadiusAuthenticationInterval.Int64Value.ValueInt64())
	updateNetworkSwitchAccessPolicy.Radius.ReAuthenticationInterval = &radiusReauthInterval
	updateNetworkSwitchAccessPolicy.UrlRedirectWalledGardenEnabled = &urlRedirectWalledGardenEnabled
	accessPolicyType := data.AccessPolicyType.String()
	updateNetworkSwitchAccessPolicy.AccessPolicyType = &accessPolicyType
	increaseAccessSpeed := data.IncreaseAccessSpeed.ValueBool()
	updateNetworkSwitchAccessPolicy.IncreaseAccessSpeed = &increaseAccessSpeed
	name := data.Name.ValueString()
	updateNetworkSwitchAccessPolicy.Name = &name
	radiusFailedAuthVlanId := int32(data.RadiusFailedAuthVlanId.Int64Value.ValueInt64())
	updateNetworkSwitchAccessPolicy.Radius.FailedAuthVlanId = &radiusFailedAuthVlanId
	controlDirection := data.Dot1xControlDirection.String()
	updateNetworkSwitchAccessPolicy.Dot1x.ControlDirection = &controlDirection
	guestVLANID := int32(data.GuestVLANId.Int64Value.ValueInt64())
	updateNetworkSwitchAccessPolicy.GuestVlanId = &guestVLANID
	radiusGroupAttribute := data.RadiusGroupAttribute.ValueString()
	updateNetworkSwitchAccessPolicy.RadiusGroupAttribute = &radiusGroupAttribute
	voiceVlanClients := data.VoiceVlanClients.ValueBool()
	updateNetworkSwitchAccessPolicy.VoiceVlanClients = &voiceVlanClients
	var walledGardenRanges []string
	for _, urlRange := range data.UrlRedirectWalledGardenRanges.Elements() {
		walledGardenRanges = append(walledGardenRanges, urlRange.String())
	}
	updateNetworkSwitchAccessPolicy.SetUrlRedirectWalledGardenRanges(walledGardenRanges)

	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchAccessPolicy(context.Background(), data.NetworkId.ValueString(), data.AccessPolicyNumber.String()).UpdateNetworkSwitchAccessPolicy(*updateNetworkSwitchAccessPolicy).Execute()
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

func (r *NetworksSwitchAccessPoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksSwitchAccessPoliciesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.SwitchApi.DeleteNetworkSwitchAccessPolicy(context.Background(), data.NetworkId.ValueString(), data.AccessPolicyType.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *NetworksSwitchAccessPoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, acl_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("access_policy_number"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
