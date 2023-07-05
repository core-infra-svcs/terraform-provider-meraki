package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &NetworksSwitchDscpToCosMappingsResource{}
	_ resource.ResourceWithConfigure   = &NetworksSwitchDscpToCosMappingsResource{}
	_ resource.ResourceWithImportState = &NetworksSwitchDscpToCosMappingsResource{}
)

func NewNetworksSwitchDscpToCosMappingsResource() resource.Resource {
	return &NetworksSwitchDscpToCosMappingsResource{}
}

// NetworksSwitchDscpToCosMappingsResource defines the resource implementation.
type NetworksSwitchDscpToCosMappingsResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchDscpToCosMappingsResourceModel describes the resource data model.
type NetworksSwitchDscpToCosMappingsResourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`
	Mappings  []Mapping        `tfsdk:"mappings" json:"mappings"`
}

type Mapping struct {
	Dscp jsontypes.Int64 `tfsdk:"dscp" json:"dscp"`
	Cos  jsontypes.Int64 `tfsdk:"cos" json:"cos"`
}

func (r *NetworksSwitchDscpToCosMappingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_dscp_to_cos_mappings"
}

func (r *NetworksSwitchDscpToCosMappingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Networks Switch Dscp to cos mappings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				Required:   true,
				CustomType: jsontypes.StringType,
			},
			"mappings": schema.ListNestedAttribute{
				Description: "An array of DSCP to CoS mappings. An empty array will reset the mappings to default.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dscp": schema.Int64Attribute{
							MarkdownDescription: "The Differentiated Services Code Point (DSCP) tag in the IP header that will be mapped to a particular Class-of-Service (CoS) queue. Value can be in the range of 0 to 63 inclusive.",
							Required:            true,
							CustomType:          jsontypes.Int64Type,
							Validators: []validator.Int64{
								int64validator.Between(0, 63),
							},
						},
						"cos": schema.Int64Attribute{
							MarkdownDescription: "The actual layer-2 CoS queue the DSCP value is mapped to. These are not bits set on outgoing frames. Value can be in the range of 0 to 5 inclusive.",
							Required:            true,
							CustomType:          jsontypes.Int64Type,
							Validators: []validator.Int64{
								int64validator.Between(0, 5),
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSwitchDscpToCosMappingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchDscpToCosMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchDscpToCosMappingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mappings := []openApiClient.UpdateNetworkSwitchDscpToCosMappingsRequestMappingsInner{}
	for _, mapping := range data.Mappings {
		mappingsMappings := openApiClient.UpdateNetworkSwitchDscpToCosMappingsRequestMappingsInner{
			Dscp: int32(mapping.Dscp.ValueInt64()),
			Cos:  int32(mapping.Cos.ValueInt64()),
		}
		mappings = append(mappings, mappingsMappings)
	}

	// Create Payload
	networkMappings := *openApiClient.NewUpdateNetworkSwitchDscpToCosMappingsRequest(mappings)

	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkSwitchDscpToCosMappings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchDscpToCosMappingsRequest(networkMappings).Execute()
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

func (r *NetworksSwitchDscpToCosMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchDscpToCosMappingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ConfigureApi.GetNetworkSwitchDscpToCosMappings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get resource",
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

func (r *NetworksSwitchDscpToCosMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksSwitchDscpToCosMappingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mappings := []openApiClient.UpdateNetworkSwitchDscpToCosMappingsRequestMappingsInner{}
	for _, mapping := range data.Mappings {
		mappings = append(mappings, openApiClient.UpdateNetworkSwitchDscpToCosMappingsRequestMappingsInner{
			Dscp: int32(mapping.Dscp.ValueInt64()),
			Cos:  int32(mapping.Cos.ValueInt64()),
		})
	}

	// Create Payload
	networkMappings := *openApiClient.NewUpdateNetworkSwitchDscpToCosMappingsRequest(mappings)

	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkSwitchDscpToCosMappings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchDscpToCosMappingsRequest(networkMappings).Execute()
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

func (r *NetworksSwitchDscpToCosMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksSwitchDscpToCosMappingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mappings := []openApiClient.UpdateNetworkSwitchDscpToCosMappingsRequestMappingsInner{}
	for _, mapping := range data.Mappings {
		mappings = append(mappings, openApiClient.UpdateNetworkSwitchDscpToCosMappingsRequestMappingsInner{
			Dscp: int32(mapping.Dscp.ValueInt64()),
			Cos:  int32(mapping.Cos.ValueInt64()),
		})
	}

	// Create Payload
	networkMappings := *openApiClient.NewUpdateNetworkSwitchDscpToCosMappingsRequest(mappings)

	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkSwitchDscpToCosMappings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchDscpToCosMappingsRequest(networkMappings).Execute()
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
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksSwitchDscpToCosMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
