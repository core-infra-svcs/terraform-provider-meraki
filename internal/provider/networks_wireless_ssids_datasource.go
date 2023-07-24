package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// NetworksWirelessSsidsDataSource struct. If not, implement them.
var _ datasource.DataSource = &NetworksWirelessSsidsDataSource{}

// The NewNetworksWirelessSsidsDataSource function is a constructor for the data source. This function needs
// to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) DataSources.
// If it's not added, the provider won't be aware of this data source's existence.
func NewNetworksWirelessSsidsDataSource() datasource.DataSource {
	return &NetworksWirelessSsidsDataSource{}
}

// NetworksWirelessSsidsDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
// If additional fields are required (e.g., for caching or for tracking internal state), add them here.
type NetworksWirelessSsidsDataSource struct {
	client *openApiClient.APIClient
}

// The NetworksWirelessSsidsDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type NetworksWirelessSsidsDataSourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`

	List []NetworksWirelessSSIDData `tfsdk:"list" json:"-"`
}

type DNSRewrite struct {
	Enabled              jsontypes.Bool                  `tfsdk:"enabled" json:"enabled"`
	DNSCustomNameservers jsontypes.Set[jsontypes.String] `tfsdk:"dns_custom_nameservers" json:"dnsCustomNameservers"`
}

type SpeedBurst struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

type NetworksWirelessSSIDData struct {
	Number                       jsontypes.Int64                 `tfsdk:"number" json:"number"`
	Name                         jsontypes.String                `tfsdk:"name" json:"name"`
	Enabled                      jsontypes.Bool                  `tfsdk:"enabled" json:"enabled"`
	SplashPage                   jsontypes.String                `tfsdk:"splash_page" json:"splashPage"`
	SsidAdminAccessible          jsontypes.Bool                  `tfsdk:"ssid_admin_accessible" json:"ssidAdminAccessible"`
	AuthMode                     jsontypes.String                `tfsdk:"auth_mode" json:"authMode"`
	IpAssignmentMode             jsontypes.String                `tfsdk:"ip_assignment_mode" json:"ipAssignmentMode"`
	AdultContentFilteringEnabled jsontypes.Bool                  `tfsdk:"adult_content_filtering_enabled" json:"adultContentFilteringEnabled"`
	DNSRewrite                   DNSRewrite                      `tfsdk:"dns_rewrite" json:"dnsRewrite"`
	MinBitrate                   jsontypes.Int64                 `tfsdk:"min_bitrate" json:"minBitrate"`
	BandSelection                jsontypes.String                `tfsdk:"band_selection" json:"bandSelection"`
	PerClientBandwidthLimitUp    jsontypes.Int64                 `tfsdk:"per_client_bandwidth_limit_up" json:"perClientBandwidthLimitUp"`
	PerClientBandwidthLimitDown  jsontypes.Int64                 `tfsdk:"per_client_bandwidth_limit_down" json:"perClientBandwidthLimitDown"`
	PerSsidBandwidthLimitUp      jsontypes.Int64                 `tfsdk:"per_ssid_bandwidth_limit_up" json:"perSsidBandwidthLimitUp"`
	PerSsidBandwidthLimitDown    jsontypes.Int64                 `tfsdk:"per_ssid_bandwidth_limit_down" json:"perSsidBandwidthLimitDown"`
	MandatoryDhcpEnabled         jsontypes.Bool                  `tfsdk:"mandatory_dhcp_enabled" json:"mandatoryDhcpEnabled"`
	Visible                      jsontypes.Bool                  `tfsdk:"visible" json:"visible"`
	AvailableOnAllAps            jsontypes.Bool                  `tfsdk:"available_on_all_aps" json:"availableOnAllAps"`
	AvailabilityTags             jsontypes.Set[jsontypes.String] `tfsdk:"availability_tags" json:"availabilityTags"`
	SpeedBurst                   SpeedBurst                      `tfsdk:"speed_burst" json:"speedBurst"`
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
				Computed:    false,
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
							Required:            true,
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
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"ssid_admin_accessible": schema.BoolAttribute{
							MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"auth_mode": schema.StringAttribute{
							MarkdownDescription: "The association control method for the SSID ('open', 'open-enhanced', 'psk', 'open-with-radius', 'open-with-nac', '8021x-meraki', '8021x-nac', '8021x-radius', '8021x-google', '8021x-localradius', 'ipsk-with-radius' or 'ipsk-without-radius')",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"ip_assignment_mode": schema.StringAttribute{
							MarkdownDescription: "The client IP assignment mode ('NAT mode', 'Bridge mode', 'Layer 3 roaming', 'Ethernet over GRE', 'Layer 3 roaming with a concentrator' or 'VPN')",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"adult_content_filtering_enabled": schema.BoolAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"dns_rewrite": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether or not the SSID is enabled.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
								},
								"dns_custom_nameservers": schema.SetAttribute{
									CustomType:  jsontypes.SetType[jsontypes.String](),
									ElementType: jsontypes.StringType,
									Description: "The IPs of the DHCP servers that DHCP requests should be relayed to",
									Computed:    true,
									Optional:    true,
								},
							},
						},
						"min_bitrate": schema.Int64Attribute{
							MarkdownDescription: "The minimum bitrate in Mbps of this SSID in the default indoor RF profile. ('1', '2', '5.5', '6', '9', '11', '12', '18', '24', '36', '48' or '54').",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"band_selection": schema.StringAttribute{
							MarkdownDescription: "This policy determines how authentication requests should be handled in the event that all of the configured RADIUS servers are unreachable ('Deny access' or 'Allow access').",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"per_client_bandwidth_limit_up": schema.Int64Attribute{
							MarkdownDescription: "The upload bandwidth limit in Kbps. (0 represents no limit.)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"per_client_bandwidth_limit_down": schema.Int64Attribute{
							MarkdownDescription: "The download bandwidth limit in Kbps.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
							MarkdownDescription: "The total upload bandwidth limit in Kbps. (0 represents no limit.)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
							MarkdownDescription: "The total download bandwidth limit in Kbps. (0 represents no limit.).",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"mandatory_dhcp_enabled": schema.BoolAttribute{
							MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"visible": schema.BoolAttribute{
							MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"available_on_all_aps": schema.BoolAttribute{
							MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"availability_tags": schema.SetAttribute{
							CustomType:  jsontypes.SetType[jsontypes.String](),
							ElementType: jsontypes.StringType,
							Description: "The IPs of the DHCP servers that DHCP requests should be relayed to",
							Computed:    true,
							Optional:    true,
						},
						"speed_burst": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									MarkdownDescription: "Whether or not the SSID is enabled.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
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

	// Remember to handle any potential errors.
	_, httpResp, err := d.client.WirelessApi.GetNetworkWirelessSsids(ctx, data.NetworkId.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read data source",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
