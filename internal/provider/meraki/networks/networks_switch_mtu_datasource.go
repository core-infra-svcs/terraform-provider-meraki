package networks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksSwitchMtuDataSource{}

func NewNetworksSwitchMtuDataSource() datasource.DataSource {
	return &NetworksSwitchMtuDataSource{}
}

// NetworksSwitchMtuDataSource defines the resource implementation.
type NetworksSwitchMtuDataSource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchMtuDataSourceModel describes the resource data model.
type NetworksSwitchMtuDataSourceModel struct {
	Id             jsontypes.String                           `tfsdk:"id"`
	NetworkId      jsontypes.String                           `tfsdk:"network_id" json:"network_id"`
	DefaultMtuSize jsontypes.Int64                            `tfsdk:"default_mtu_size" json:"defaultMtuSize"`
	Overrides      []NetworksSwitchMtuDataSourceModelOverride `tfsdk:"overrides" json:"overrides"`
}

type NetworksSwitchMtuDataSourceModelOverride struct {
	Switches       []string        `tfsdk:"switches" json:"switches"`
	SwitchProfiles []string        `tfsdk:"switch_profiles" json:"switchProfiles"`
	MtuSize        jsontypes.Int64 `tfsdk:"mtu_size" json:"mtuSize"`
}

func (r *NetworksSwitchMtuDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_mtu"
}

func (r *NetworksSwitchMtuDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Networks switch mtu resource for updating networks switch mtu.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"default_mtu_size": schema.Int64Attribute{
				MarkdownDescription: "MTU size for the entire network. Default value is 9578.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"overrides": schema.SetNestedAttribute{
				Description: "Override MTU size for individual switches or switch profiles. An empty array will clear overrides.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"switches": schema.SetAttribute{
							MarkdownDescription: "List of switch serials. Applicable only for switch network.",
							CustomType:          jsontypes.SetType[jsontypes.String](),
							Optional:            true,
							Computed:            true,
						},
						"switch_profiles": schema.SetAttribute{
							MarkdownDescription: "List of switch profile IDs. Applicable only for template network.",
							CustomType:          jsontypes.SetType[jsontypes.String](),
							Optional:            true,
							Computed:            true,
						},
						"mtu_size": schema.Int64Attribute{
							MarkdownDescription: "MTU size for the switches or switch profiles..",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSwitchMtuDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NetworksSwitchMtuDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksSwitchMtuDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.MtuApi.GetNetworkSwitchMtu(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
