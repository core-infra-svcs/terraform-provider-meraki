package ssids

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// NetworksWirelessSsidsDataSource struct. If not, implement them.
var _ datasource.DataSource = &NetworksWirelessSsidsDataSource{}

type NetworksWirelessSsidsDataSource struct {
	client *openApiClient.APIClient
}

func NewNetworksWirelessSsidsDataSource() datasource.DataSource {
	return &NetworksWirelessSsidsDataSource{}
}

// The NetworksWirelessSsidsDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type NetworksWirelessSsidsDataSourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`

	List []NetworksWirelessSsidsDataSourceModelList `tfsdk:"list"`
}

type NetworksWirelessSsidsDataSourceModelRadiusServer struct {
	Host                     jsontypes.String `tfsdk:"host" json:"host"`
	Port                     jsontypes.Int64  `tfsdk:"port" json:"port"`
	OpenRoamingCertificateId jsontypes.Int64  `tfsdk:"open_roaming_certificate_id" json:"openRoamingCertificateId"`
	CaCertificate            jsontypes.String `tfsdk:"ca_certificate" json:"caCertificate"`
}

type NetworksWirelessSsidsDataSourceModelList struct {
	Number                          jsontypes.Int64                                    `tfsdk:"number" json:"number"`
	Name                            jsontypes.String                                   `tfsdk:"name" json:"name"`
	Enabled                         jsontypes.Bool                                     `tfsdk:"enabled" json:"enabled"`
	SplashPage                      jsontypes.String                                   `tfsdk:"splash_page" json:"splashPage"`
	SSIDAdminAccessible             jsontypes.Bool                                     `tfsdk:"ssid_admin_accessible" json:"ssidAdminAccessible"`
	LocalAuth                       jsontypes.Bool                                     `json:"localAuth" tfsdk:"local_auth"`
	AuthMode                        jsontypes.String                                   `tfsdk:"auth_mode" json:"authMode"`
	EncryptionMode                  jsontypes.String                                   `tfsdk:"encryption_mode" json:"encryptionMode"`
	WPAEncryptionMode               jsontypes.String                                   `tfsdk:"wpa_encryption_mode" json:"wpaEncryptionMode"`
	RadiusServers                   []NetworksWirelessSsidsDataSourceModelRadiusServer `tfsdk:"radius_servers" json:"radiusServers"`
	RadiusAccountingServers         []NetworksWirelessSsidsDataSourceModelRadiusServer `tfsdk:"radius_accounting_servers" json:"radiusAccountingServers"`
	RadiusAccountingEnabled         jsontypes.Bool                                     `tfsdk:"radius_accounting_enabled" json:"radiusAccountingEnabled"`
	RadiusEnabled                   jsontypes.Bool                                     `tfsdk:"radius_enabled" json:"radiusEnabled"`
	RadiusAttributeForGroupPolicies jsontypes.String                                   `tfsdk:"radius_attribute_for_group_policies" json:"radiusAttributeForGroupPolicies"`
	RadiusFailoverPolicy            jsontypes.String                                   `tfsdk:"radius_failover_policy" json:"radiusFailoverPolicy"`
	RadiusLoadBalancingPolicy       jsontypes.String                                   `tfsdk:"radius_load_balancing_policy" json:"radiusLoadBalancingPolicy"`
	IPAssignmentMode                jsontypes.String                                   `tfsdk:"ip_assignment_mode" json:"ipAssignmentMode"`
	AdminSplashURL                  jsontypes.String                                   `tfsdk:"admin_splash_url" json:"adminSplashUrl"`
	SplashTimeout                   jsontypes.String                                   `tfsdk:"splash_timeout" json:"splashTimeout"`
	WalledGardenEnabled             jsontypes.Bool                                     `tfsdk:"walled_garden_enabled" json:"walledGardenEnabled"`
	WalledGardenRanges              []jsontypes.String                                 `tfsdk:"walled_garden_ranges" json:"walledGardenRanges"`
	MinBitrate                      jsontypes.Int64                                    `tfsdk:"min_bitrate" json:"minBitrate"`
	BandSelection                   jsontypes.String                                   `tfsdk:"band_selection" json:"bandSelection"`
	PerClientBandwidthLimitUp       jsontypes.Int64                                    `tfsdk:"per_client_bandwidth_limit_up" json:"perClientBandwidthLimitUp"`
	PerClientBandwidthLimitDown     jsontypes.Int64                                    `tfsdk:"per_client_bandwidth_limit_down" json:"perClientBandwidthLimitDown"`
	Visible                         jsontypes.Bool                                     `tfsdk:"visible" json:"visible"`
	AvailableOnAllAPs               jsontypes.Bool                                     `tfsdk:"available_on_all_aps" json:"availableOnAllAps"`
	AvailabilityTags                []jsontypes.String                                 `tfsdk:"availability_tags" json:"availabilityTags"`
	PerSSIDBandwidthLimitUp         jsontypes.Int64                                    `tfsdk:"per_ssid_bandwidth_limit_up" json:"perSsidBandwidthLimitUp"`
	PerSSIDBandwidthLimitDown       jsontypes.Int64                                    `tfsdk:"per_ssid_bandwidth_limit_down" json:"perSsidBandwidthLimitDown"`
	MandatoryDHCPEnabled            jsontypes.Bool                                     `tfsdk:"mandatory_dhcp_enabled" json:"mandatoryDhcpEnabled"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *NetworksWirelessSsidsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *NetworksWirelessSsidsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		// It should provide a clear and concise description of the data source.
		MarkdownDescription: "NetworksWirelessSsids",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			// Every data source must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"number": schema.Int64Attribute{
							MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether or not the SSID is enabled.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"splash_page": schema.StringAttribute{
							MarkdownDescription: "The type of splash page for the SSID ('None', 'Click-through splash page', 'Billing', 'Password-protected with Meraki RADIUS', 'Password-protected with custom RADIUS', 'Password-protected with Active Directory', 'Password-protected with LDAP', 'SMS authentication', 'Systems Manager Sentry', 'Facebook Wi-Fi', 'Google OAuth', 'Sponsored guest', 'Cisco ISE' or 'Google Apps domain'). This attribute is not supported for template children.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"ssid_admin_accessible": schema.BoolAttribute{
							MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"local_auth": schema.BoolAttribute{
							MarkdownDescription: "Extended local auth flag for Enterprise NAC.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"auth_mode": schema.StringAttribute{
							MarkdownDescription: "The association control method for the SSID ('open', 'open-enhanced', 'psk', 'open-with-radius', 'open-with-nac', '8021x-meraki', '8021x-nac', '8021x-radius', '8021x-google', '8021x-localradius', 'ipsk-with-radius' or 'ipsk-without-radius')",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"encryption_mode": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"wpa_encryption_mode": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"radius_servers": schema.ListNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.StringType,
									},
									"port": schema.Int64Attribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.Int64Type,
									},
									"open_roaming_certificate_id": schema.Int64Attribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.Int64Type,
									},
									"ca_certificate": schema.StringAttribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.StringType,
									},
								},
							},
						},
						"radius_accounting_servers": schema.ListNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.StringType,
									},
									"port": schema.Int64Attribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.Int64Type,
									},
									"open_roaming_certificate_id": schema.Int64Attribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.Int64Type,
									},
									"ca_certificate": schema.StringAttribute{
										Optional:   true,
										Computed:   true,
										CustomType: jsontypes.StringType,
									},
								},
							},
						},
						"radius_accounting_enabled": schema.BoolAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.BoolType,
						},
						"radius_enabled": schema.BoolAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.BoolType,
						},
						"radius_attribute_for_group_policies": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"radius_failover_policy": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"radius_load_balancing_policy": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"ip_assignment_mode": schema.StringAttribute{
							MarkdownDescription: "The client IP assignment mode ('NAT mode', 'Bridge mode', 'Layer 3 roaming', 'Ethernet over GRE', 'Layer 3 roaming with a concentrator' or 'VPN')",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"admin_splash_url": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"splash_timeout": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"walled_garden_enabled": schema.BoolAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.BoolType,
						},
						"walled_garden_ranges": schema.SetAttribute{
							Optional:    true,
							Computed:    true,
							CustomType:  jsontypes.SetType[jsontypes.String](),
							ElementType: jsontypes.StringType,
						},
						"min_bitrate": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.Int64Type,
						},
						"band_selection": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"per_client_bandwidth_limit_up": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.Int64Type,
						},
						"per_client_bandwidth_limit_down": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.Int64Type,
						},
						"visible": schema.BoolAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.BoolType,
						},
						"available_on_all_aps": schema.BoolAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.BoolType,
						},
						"availability_tags": schema.SetAttribute{
							Optional:    true,
							Computed:    true,
							CustomType:  jsontypes.SetType[jsontypes.String](),
							ElementType: jsontypes.StringType,
						},
						"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.Int64Type,
						},
						"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.Int64Type,
						},
						"mandatory_dhcp_enabled": schema.BoolAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *NetworksWirelessSsidsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *NetworksWirelessSsidsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworksWirelessSsidsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.WirelessApi.GetNetworkWirelessSsids(ctx, data.NetworkId.ValueString()).Execute()

	// Check for errors API call
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
	}

	// Now iterate over the inlineResp slice
	for _, inlineRespData := range inlineResp {
		result := NetworksWirelessSsidsDataSourceModelList{
			Number:                          jsontypes.Int64Value(int64(inlineRespData.GetNumber())),
			Name:                            jsontypes.StringValue(inlineRespData.GetName()),
			Enabled:                         jsontypes.BoolValue(inlineRespData.GetEnabled()),
			SplashPage:                      jsontypes.StringValue(inlineRespData.GetSplashPage()),
			SSIDAdminAccessible:             jsontypes.BoolValue(inlineRespData.GetSsidAdminAccessible()),
			LocalAuth:                       jsontypes.BoolValue(inlineRespData.GetLocalAuth()),
			AuthMode:                        jsontypes.StringValue(inlineRespData.GetAuthMode()),
			EncryptionMode:                  jsontypes.StringValue(inlineRespData.GetEncryptionMode()),
			WPAEncryptionMode:               jsontypes.StringValue(inlineRespData.GetWpaEncryptionMode()),
			RadiusServers:                   make([]NetworksWirelessSsidsDataSourceModelRadiusServer, 0),
			RadiusAccountingServers:         make([]NetworksWirelessSsidsDataSourceModelRadiusServer, 0),
			RadiusAccountingEnabled:         jsontypes.BoolValue(inlineRespData.GetRadiusAccountingEnabled()),
			RadiusEnabled:                   jsontypes.BoolValue(inlineRespData.GetRadiusEnabled()),
			RadiusAttributeForGroupPolicies: jsontypes.StringValue(inlineRespData.GetRadiusAttributeForGroupPolicies()),
			RadiusFailoverPolicy:            jsontypes.StringValue(inlineRespData.GetRadiusFailoverPolicy()),
			RadiusLoadBalancingPolicy:       jsontypes.StringValue(inlineRespData.GetRadiusLoadBalancingPolicy()),
			IPAssignmentMode:                jsontypes.StringValue(inlineRespData.GetIpAssignmentMode()),
			AdminSplashURL:                  jsontypes.StringValue(inlineRespData.GetAdminSplashUrl()),
			SplashTimeout:                   jsontypes.StringValue(inlineRespData.GetSplashTimeout()),
			WalledGardenEnabled:             jsontypes.BoolValue(inlineRespData.GetWalledGardenEnabled()),
			WalledGardenRanges:              make([]jsontypes.String, 0),
			MinBitrate:                      jsontypes.Int64Value(int64(inlineRespData.GetMinBitrate())),
			BandSelection:                   jsontypes.StringValue(inlineRespData.GetBandSelection()),
			PerClientBandwidthLimitUp:       jsontypes.Int64Value(int64(inlineRespData.GetPerClientBandwidthLimitUp())),
			PerClientBandwidthLimitDown:     jsontypes.Int64Value(int64(inlineRespData.GetPerClientBandwidthLimitDown())),
			Visible:                         jsontypes.BoolValue(inlineRespData.GetVisible()),
			AvailableOnAllAPs:               jsontypes.BoolValue(inlineRespData.GetAvailableOnAllAps()),
			AvailabilityTags:                make([]jsontypes.String, 0),
			PerSSIDBandwidthLimitUp:         jsontypes.Int64Value(int64(inlineRespData.GetPerSsidBandwidthLimitUp())),
			PerSSIDBandwidthLimitDown:       jsontypes.Int64Value(int64(inlineRespData.GetPerSsidBandwidthLimitDown())),
			MandatoryDHCPEnabled:            jsontypes.BoolValue(inlineRespData.GetMandatoryDhcpEnabled()),
		}

		// Populate RadiusServers slice
		for _, radiusServer := range inlineRespData.RadiusServers {
			result.RadiusServers = append(result.RadiusServers, NetworksWirelessSsidsDataSourceModelRadiusServer{
				Host:                     jsontypes.StringValue(radiusServer.GetHost()),
				Port:                     jsontypes.Int64Value(int64(radiusServer.GetPort())),
				OpenRoamingCertificateId: jsontypes.Int64Value(int64(radiusServer.GetOpenRoamingCertificateId())),
				CaCertificate:            jsontypes.StringValue(radiusServer.GetCaCertificate()),
			})
		}

		// Populate RadiusAccountingServers slice
		for _, radiusAccountingServer := range inlineRespData.RadiusAccountingServers {
			result.RadiusAccountingServers = append(result.RadiusAccountingServers, NetworksWirelessSsidsDataSourceModelRadiusServer{
				Host:                     jsontypes.StringValue(radiusAccountingServer.GetHost()),
				Port:                     jsontypes.Int64Value(int64(radiusAccountingServer.GetPort())),
				OpenRoamingCertificateId: jsontypes.Int64Value(int64(radiusAccountingServer.GetOpenRoamingCertificateId())),
				CaCertificate:            jsontypes.StringValue(radiusAccountingServer.GetCaCertificate()),
			})
		}

		// Populate WalledGardenRanges slice
		for _, walledGardenRange := range inlineRespData.WalledGardenRanges {
			result.WalledGardenRanges = append(result.WalledGardenRanges, jsontypes.StringValue(walledGardenRange))
		}

		// Populate AvailabilityTags slice
		for _, availabilityTag := range inlineRespData.AvailabilityTags {
			result.AvailabilityTags = append(result.AvailabilityTags, jsontypes.StringValue(availabilityTag))
		}

		// Save data into Terraform state
		data.List = append(data.List, result)
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
