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
var _ resource.Resource = &NetworksTrafficAnalysisResource{}
var _ resource.ResourceWithImportState = &NetworksTrafficAnalysisResource{}

func NewNetworksTrafficAnalysisResource() resource.Resource {
	return &NetworksTrafficAnalysisResource{}
}

// NetworksTrafficAnalysisResource defines the resource implementation.
type NetworksTrafficAnalysisResource struct {
	client *openApiClient.APIClient
}

// NetworksTrafficAnalysisResourceModel describes the resource data model.
type NetworksTrafficAnalysisResourceModel struct {
	Id                  jsontypes.String                                          `tfsdk:"id"`
	NetworkId           jsontypes.String                                          `tfsdk:"network_id" json:"network_id"`
	Mode                jsontypes.String                                          `tfsdk:"mode" json:"mode"`
	CustomPieChartItems []NetworksTrafficAnalysisResourceModelCustomPieChartItems `tfsdk:"custom_piechart_items" json:"customPieChartItems"`
}

type NetworksTrafficAnalysisResourceModelCustomPieChartItems struct {
	Name  jsontypes.String `tfsdk:"name" json:"name"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
	Value jsontypes.String `tfsdk:"value" json:"value"`
}

func (r *NetworksTrafficAnalysisResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_traffic_analysis"
}

func (r *NetworksTrafficAnalysisResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksTrafficAnalysis resource for updating networks traffic analysis.",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The traffic analysis mode for the network. Can be one of 'disabled' (do not collect traffic types) 'basic' (collect generic traffic categories), or 'detailed' (collect destination hostnames)",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"disabled", "basic", "detailed"}...),
					stringvalidator.LengthAtLeast(5),
				},
			},
			"custom_piechart_items": schema.ListNestedAttribute{
				Description: "The list of items that make up the custom pie chart for traffic reporting.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the custom pie chart item.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The signature type for the custom pie chart item. Can be one of 'host', 'port' or 'ipRange'.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"host", "port", "ipRange"}...),
								stringvalidator.LengthAtLeast(4),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value of the custom pie chart item. Valid syntax depends on the signature type of the chart item.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksTrafficAnalysisResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksTrafficAnalysisResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksTrafficAnalysisResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkTrafficAnalysis := *openApiClient.NewInlineObject140()
	updateNetworkTrafficAnalysis.SetMode(data.Mode.ValueString())

	if len(data.CustomPieChartItems) > 0 {
		var customPieChartItems []openApiClient.NetworksNetworkIdTrafficAnalysisCustomPieChartItems
		for _, attribute := range data.CustomPieChartItems {
			var customPieChartItem openApiClient.NetworksNetworkIdTrafficAnalysisCustomPieChartItems
			customPieChartItem.Name = attribute.Name.ValueString()
			customPieChartItem.Type = attribute.Type.ValueString()
			customPieChartItem.Value = attribute.Value.ValueString()
			customPieChartItems = append(customPieChartItems, customPieChartItem)
		}
		updateNetworkTrafficAnalysis.SetCustomPieChartItems(customPieChartItems)
	} else {
		data.CustomPieChartItems = nil
	}

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkTrafficAnalysis(ctx, data.NetworkId.ValueString()).UpdateNetworkTrafficAnalysis(updateNetworkTrafficAnalysis).Execute()
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
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksTrafficAnalysisResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksTrafficAnalysisResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.NetworksApi.GetNetworkTrafficAnalysis(ctx, data.NetworkId.ValueString()).Execute()
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

func (r *NetworksTrafficAnalysisResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksTrafficAnalysisResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkTrafficAnalysis := *openApiClient.NewInlineObject140()
	updateNetworkTrafficAnalysis.SetMode(data.Mode.ValueString())

	if len(data.CustomPieChartItems) > 0 {
		var customPieChartItems []openApiClient.NetworksNetworkIdTrafficAnalysisCustomPieChartItems
		for _, attribute := range data.CustomPieChartItems {
			var customPieChartItem openApiClient.NetworksNetworkIdTrafficAnalysisCustomPieChartItems
			customPieChartItem.Name = attribute.Name.ValueString()
			customPieChartItem.Type = attribute.Type.ValueString()
			customPieChartItem.Value = attribute.Value.ValueString()
			customPieChartItems = append(customPieChartItems, customPieChartItem)
		}
		updateNetworkTrafficAnalysis.SetCustomPieChartItems(customPieChartItems)
	} else {
		data.CustomPieChartItems = nil
	}

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkTrafficAnalysis(ctx, data.NetworkId.ValueString()).UpdateNetworkTrafficAnalysis(updateNetworkTrafficAnalysis).Execute()
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
	tflog.Trace(ctx, "update resource")
}

func (r *NetworksTrafficAnalysisResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksTrafficAnalysisResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkTrafficAnalysis := *openApiClient.NewInlineObject140()
	updateNetworkTrafficAnalysis.SetMode("disabled")
	updateNetworkTrafficAnalysis.SetCustomPieChartItems(nil)

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkTrafficAnalysis(ctx, data.NetworkId.ValueString()).UpdateNetworkTrafficAnalysis(updateNetworkTrafficAnalysis).Execute()
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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksTrafficAnalysisResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
