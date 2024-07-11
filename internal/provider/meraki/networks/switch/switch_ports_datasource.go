package _switch

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"time"
)

var _ datasource.DataSource = &DevicesSwitchPortsStatusesDataSource{}

func NewDevicesSwitchPortsStatusesDataSource() datasource.DataSource {
	return &DevicesSwitchPortsStatusesDataSource{}
}

// DevicesSwitchPortsStatusesDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
type DevicesSwitchPortsStatusesDataSource struct {
	client *openApiClient.APIClient
}

// The DevicesSwitchPortsStatusesDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type DevicesSwitchPortsStatusesDataSourceModel struct {
	Id     jsontypes.String                                `tfsdk:"id"`
	Serial jsontypes.String                                `tfsdk:"serial"`
	List   []DevicesSwitchPortsStatusesDataSourceModelList `tfsdk:"list"`
}

type DevicesSwitchPortsStatusesDataSourceModelList struct {
	PortId                  jsontypes.String                                 `tfsdk:"port_id"`
	Name                    jsontypes.String                                 `tfsdk:"name"`
	Tags                    []jsontypes.String                               `tfsdk:"tags"`
	Enabled                 jsontypes.Bool                                   `tfsdk:"enabled"`
	PoeEnabled              jsontypes.Bool                                   `tfsdk:"poe_enabled"`
	Type                    jsontypes.String                                 `tfsdk:"type"`
	Vlan                    jsontypes.Int64                                  `tfsdk:"vlan"`
	VoiceVlan               jsontypes.Int64                                  `tfsdk:"voice_vlan"`
	AllowedVlans            jsontypes.String                                 `tfsdk:"allowed_vlans"`
	IsolationEnabled        jsontypes.Bool                                   `tfsdk:"isolation_enabled"`
	RstpEnabled             jsontypes.Bool                                   `tfsdk:"rstp_enabled"`
	StpGuard                jsontypes.String                                 `tfsdk:"stp_guard"`
	AccessPolicyNumber      jsontypes.Int64                                  `tfsdk:"access_policy_number"`
	AccessPolicyType        jsontypes.String                                 `tfsdk:"access_policy_type"`
	LinkNegotiation         jsontypes.String                                 `tfsdk:"link_negotiation"`
	PortScheduleId          jsontypes.String                                 `tfsdk:"port_schedule_id"`
	Udld                    jsontypes.String                                 `tfsdk:"udld"`
	StickyMacWhitelistLimit jsontypes.Int64                                  `tfsdk:"sticky_mac_white_list_limit"`
	StormControlEnabled     jsontypes.Bool                                   `tfsdk:"storm_control_enabled"`
	MacWhitelist            []jsontypes.String                               `tfsdk:"mac_white_list"`
	StickyMacWhitelist      []jsontypes.String                               `tfsdk:"sticky_mac_white_list"`
	AdaptivePolicyGroupId   jsontypes.String                                 `tfsdk:"adaptive_policy_group_id"`
	PeerSgtCapable          jsontypes.Bool                                   `tfsdk:"peer_sgt_capable"`
	FlexibleStackingEnabled jsontypes.Bool                                   `tfsdk:"flexible_stacking_enabled"`
	DaiTrusted              jsontypes.Bool                                   `tfsdk:"dai_trusted"`
	Profile                 DevicesSwitchPortsStatusesDataSourceModelProfile `tfsdk:"profile"`
}

type DevicesSwitchPortsStatusesDataSourceModelProfile struct {
	Enabled jsontypes.Bool   `tfsdk:"enabled"`
	Id      jsontypes.String `tfsdk:"id"`
	Iname   jsontypes.String `tfsdk:"iname"`
}

func (d *DevicesSwitchPortsStatusesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_devices_switch_ports"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *DevicesSwitchPortsStatusesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		MarkdownDescription: "List the switch ports for a switch",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			// Every data source must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "A list of serial numbers. The returned devices will be filtered to only include these serials.",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "List of switch ports",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port_id": schema.StringAttribute{
							MarkdownDescription: "The identifier of the switch port.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the switch port.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "The list of tags of the switch port.",
							ElementType:         jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "The status of the switch port.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"poe_enabled": schema.BoolAttribute{
							MarkdownDescription: "The PoE status of the switch port.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the switch port ('trunk' or 'access').",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"vlan": schema.Int64Attribute{
							MarkdownDescription: "The VLAN of the switch port. A null value will clear the value set for trunk ports.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"voice_vlan": schema.Int64Attribute{
							MarkdownDescription: "The voice VLAN of the switch port. Only applicable to access ports.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"allowed_vlans": schema.StringAttribute{
							MarkdownDescription: "The VLANs allowed on the switch port. Only applicable to trunk ports.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"isolation_enabled": schema.BoolAttribute{
							MarkdownDescription: "The isolation status of the switch port.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"rstp_enabled": schema.BoolAttribute{
							MarkdownDescription: "The rapid spanning tree protocol status.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"stp_guard": schema.StringAttribute{
							MarkdownDescription: "The state of the STP guard ('disabled', 'root guard', 'bpdu guard' or 'loop guard').",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"access_policy_type": schema.StringAttribute{
							MarkdownDescription: "The type of the access policy of the switch port. Only applicable to access ports. Can be one of 'Open', 'Custom access policy', 'MAC allow list' or 'Sticky MAC allow list'.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"access_policy_number": schema.Int64Attribute{
							MarkdownDescription: "The number of a custom access policy to configure on the switch port. Only applicable when 'accessPolicyType' is 'Custom access policy'.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"link_negotiation": schema.StringAttribute{
							MarkdownDescription: "The link speed for the switch port.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"port_schedule_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the port schedule. A value of null will clear the port schedule.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"udld": schema.StringAttribute{
							MarkdownDescription: "The action to take when Unidirectional Link is detected (Alert only, Enforce). Default configuration is Alert only.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"sticky_mac_white_list_limit": schema.Int64Attribute{
							MarkdownDescription: "The maximum number of MAC addresses for sticky MAC allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"storm_control_enabled": schema.BoolAttribute{
							MarkdownDescription: "The storm control status of the switch port.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"mac_white_list": schema.SetAttribute{
							MarkdownDescription: "Only devices with MAC addresses specified in this list will have access to this port. Up to 20 MAC addresses can be defined. Only applicable when 'accessPolicyType' is 'MAC allow list'.",
							ElementType:         jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"sticky_mac_white_list": schema.SetAttribute{
							MarkdownDescription: "The initial list of MAC addresses for sticky Mac allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
							ElementType:         jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"adaptive_policy_group_id": schema.StringAttribute{
							MarkdownDescription: "The adaptive policy group ID that will be used to tag traffic through this switch port. This ID must pre-exist during the configuration, else needs to be created using adaptivePolicy/groups API. Cannot be applied to a port on a switch bound to profile.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"peer_sgt_capable": schema.BoolAttribute{
							MarkdownDescription: "If true, Peer SGT is enabled for traffic through this switch port. Applicable to trunk port only, not access port. Cannot be applied to a port on a switch bound to profile.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"flexible_stacking_enabled": schema.BoolAttribute{
							MarkdownDescription: "For supported switches (e.g. MS420/MS425), whether or not the port has flexible stacking enabled.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"dai_trusted": schema.BoolAttribute{
							MarkdownDescription: "If true, ARP packets for this port will be considered trusted, and Dynamic ARP Inspection will allow the traffic.",
							CustomType:          jsontypes.BoolType,
							Optional:            true,
							Computed:            true,
						},
						"profile": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									MarkdownDescription: "When enabled, override this port's configuration with a port profile.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
								},
								"id": schema.StringAttribute{
									MarkdownDescription: "When enabled, the ID of the port profile used to override the port's configuration.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"iname": schema.StringAttribute{
									MarkdownDescription: "When enabled, the IName of the profile.",
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
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *DevicesSwitchPortsStatusesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// This allows the data source to use the configured provider for any API calls it needs to make.
	d.client = client
}

// Read method is responsible for reading an existing data source's state.
func (d *DevicesSwitchPortsStatusesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DevicesSwitchPortsStatusesDataSourceModel
	var diags diag.Diagnostics
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := d.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(d.client.GetConfig().Retry4xxErrorWaitTime)

	// usage of CustomHttpRequestRetry with a slice of strongly typed structs
	apiCallSlice := func() ([]openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := d.client.SwitchApi.GetDeviceSwitchPorts(ctx, data.Serial.ValueString()).Execute()
		return inline, httpResp, err, diags
	}

	resultSlice, httpRespSlice, errSlice, tfDiags := utils.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCallSlice)
	if errSlice != nil {

		if tfDiags.HasError() {
			fmt.Printf(": %s\n", tfDiags.Errors())
		}

		fmt.Printf("Error creating group policy: %s\n", errSlice)
		if httpRespSlice != nil {
			var responseBody string
			if httpRespSlice.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpRespSlice.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			fmt.Printf("Failed to create resource. HTTP Status Code: %d, Response Body: %s\n", httpRespSlice.StatusCode, responseBody)
		}
		return
	}

	// Type assert apiResp to the expected []openApiClient.GetDeviceSwitchPorts200ResponseInner type
	inlineResp, ok := any(resultSlice).([]openApiClient.GetDeviceSwitchPorts200ResponseInner)
	if !ok {
		fmt.Println("Failed to assert API response type to []openApiClient.GetDeviceSwitchPorts200ResponseInner. Please ensure the API response structure matches the expected type.")
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpRespSlice.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpRespSlice.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	for _, switchData := range inlineResp {
		var devicesSwitchPortData DevicesSwitchPortsStatusesDataSourceModelList
		devicesSwitchPortData.Name = jsontypes.StringValue(switchData.GetName())
		devicesSwitchPortData.PortId = jsontypes.StringValue(switchData.GetPortId())
		devicesSwitchPortData.Enabled = jsontypes.BoolValue(switchData.GetEnabled())
		devicesSwitchPortData.PoeEnabled = jsontypes.BoolValue(switchData.GetPoeEnabled())
		devicesSwitchPortData.Type = jsontypes.StringValue(switchData.GetType())
		devicesSwitchPortData.Vlan = jsontypes.Int64Value(int64(switchData.GetVlan()))
		devicesSwitchPortData.VoiceVlan = jsontypes.Int64Value(int64(switchData.GetVoiceVlan()))
		devicesSwitchPortData.AllowedVlans = jsontypes.StringValue(switchData.GetAllowedVlans())
		devicesSwitchPortData.IsolationEnabled = jsontypes.BoolValue(switchData.GetIsolationEnabled())
		devicesSwitchPortData.RstpEnabled = jsontypes.BoolValue(switchData.GetRstpEnabled())
		devicesSwitchPortData.StpGuard = jsontypes.StringValue(switchData.GetStpGuard())
		devicesSwitchPortData.AccessPolicyNumber = jsontypes.Int64Value(int64(switchData.GetAccessPolicyNumber()))
		devicesSwitchPortData.AccessPolicyType = jsontypes.StringValue(switchData.GetAccessPolicyType())
		devicesSwitchPortData.LinkNegotiation = jsontypes.StringValue(switchData.GetLinkNegotiation())
		devicesSwitchPortData.PortScheduleId = jsontypes.StringValue(switchData.GetPortScheduleId())
		devicesSwitchPortData.Udld = jsontypes.StringValue(switchData.GetUdld())
		devicesSwitchPortData.StickyMacWhitelistLimit = jsontypes.Int64Value(int64(switchData.GetStickyMacAllowListLimit()))
		devicesSwitchPortData.StormControlEnabled = jsontypes.BoolValue(switchData.GetStormControlEnabled())
		devicesSwitchPortData.AdaptivePolicyGroupId = jsontypes.StringValue(switchData.GetAdaptivePolicyGroupId())
		devicesSwitchPortData.PeerSgtCapable = jsontypes.BoolValue(switchData.GetPeerSgtCapable())
		devicesSwitchPortData.FlexibleStackingEnabled = jsontypes.BoolValue(switchData.GetFlexibleStackingEnabled())
		devicesSwitchPortData.DaiTrusted = jsontypes.BoolValue(switchData.GetDaiTrusted())
		devicesSwitchPortData.Profile.Enabled = jsontypes.BoolValue(switchData.Profile.GetEnabled())
		devicesSwitchPortData.Profile.Id = jsontypes.StringValue(switchData.Profile.GetId())
		devicesSwitchPortData.Profile.Iname = jsontypes.StringValue(switchData.Profile.GetIname())
		for _, attribute := range switchData.GetStickyMacAllowList() {
			devicesSwitchPortData.StickyMacWhitelist = append(devicesSwitchPortData.StickyMacWhitelist, jsontypes.StringValue(attribute))
		}
		for _, attribute := range switchData.GetTags() {
			devicesSwitchPortData.Tags = append(devicesSwitchPortData.Tags, jsontypes.StringValue(attribute))
		}
		for _, attribute := range switchData.GetMacAllowList() {
			devicesSwitchPortData.MacWhitelist = append(devicesSwitchPortData.MacWhitelist, jsontypes.StringValue(attribute))
		}
		data.List = append(data.List, devicesSwitchPortData)
	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
