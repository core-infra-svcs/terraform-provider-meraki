package _switch

import (
	"context"
	"encoding/json"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource              = &DevicesSwitchPortsCycleResource{}
	_ resource.ResourceWithConfigure = &DevicesSwitchPortsCycleResource{}
)

func NewDevicesSwitchPortsCycleResource() resource.Resource {
	return &DevicesSwitchPortsCycleResource{}
}

// DevicesSwitchPortsCycleResource defines the resource implementation.
type DevicesSwitchPortsCycleResource struct {
	client *openApiClient.APIClient
}

// DevicesSwitchPortsCycleResourceModel describes the resource data model.
type DevicesSwitchPortsCycleResourceModel struct {
	Id     jsontypes2.String   `tfsdk:"id"`
	Serial jsontypes2.String   `tfsdk:"serial"`
	Ports  []jsontypes2.String `tfsdk:"ports" json:"ports"`
}

func (r *DevicesSwitchPortsCycleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_switch_ports_cycle"
}

func (r *DevicesSwitchPortsCycleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cycle a set of switch ports",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes2.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.StringType,
			},
			"ports": schema.SetAttribute{
				MarkdownDescription: "URL that is consuming SAML Identity Provider (IdP)",
				ElementType:         jsontypes2.StringType,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DevicesSwitchPortsCycleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesSwitchPortsCycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesSwitchPortsCycleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ports []string

	for _, port := range data.Ports {
		ports = append(ports, port.ValueString())
	}

	payload := openApiClient.NewCycleDeviceSwitchPortsRequest(ports)

	_, httpResp, err := r.client.PortsApi.CycleDeviceSwitchPorts(ctx, data.Serial.ValueString()).CycleDeviceSwitchPortsRequest(*payload).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	data.Id = jsontypes2.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesSwitchPortsCycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesSwitchPortsCycleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *DevicesSwitchPortsCycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesSwitchPortsCycleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesSwitchPortsCycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesSwitchPortsCycleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}
