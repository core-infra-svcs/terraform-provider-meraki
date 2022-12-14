package provider

import (
	"context"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworkResource{}
var _ resource.ResourceWithImportState = &NetworkResource{}

func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}

// NetworkResource defines the resource implementation.
type NetworkResource struct {
	client *apiclient.APIClient
}

// NetworkResourceModel describes the resource data model.
type NetworkResourceModel struct {
	Id                      types.String   `tfsdk:"id"`
	NetworkId               types.String   `tfsdk:"network_id"`
	OrganizationId          types.String   `tfsdk:"organization_id"`
	Name                    types.String   `tfsdk:"name"`
	ProductTypes            []types.String `tfsdk:"product_types"`
	Timezone                types.String   `tfsdk:"timezone"`
	Tags                    []types.String `tfsdk:"tags"`
	EnrollmentString        types.String   `tfsdk:"enrollment_string"`
	Url                     types.String   `tfsdk:"url"`
	Notes                   types.String   `tfsdk:"notes"`
	IsBoundToConfigTemplate types.Bool     `tfsdk:"is_bound_to_config_template"`
	CopyFromNetworkId       types.String   `tfsdk:"copy_from_network_id"`
}

func (r *NetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *NetworkResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Network resource - Meraki network resource.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"network_id": {
				Description:         "Network ID",
				MarkdownDescription: "Network ID",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"organization_id": {
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"name": {
				Description:         "Network name",
				MarkdownDescription: "Network name",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"product_types": {
				Description:         "List of the product types that the network supports",
				MarkdownDescription: "List of the product types that the network supports",
				Type:                types.ListType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"timezone": {
				Description:         "Timezone of the network",
				MarkdownDescription: "Timezone of the network",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"tags": {
				Description:         "Network tags",
				MarkdownDescription: "Network tags",
				Type:                types.ListType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"enrollment_string": {
				Description:         "Enrollment string for the network",
				MarkdownDescription: "Enrollment string for the network",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"url": {
				Description:         "URL to the network Dashboard UI",
				MarkdownDescription: "URL to the network Dashboard UI",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"notes": {
				Description:         "Notes for the network",
				MarkdownDescription: "Notes for the network",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"is_bound_to_config_template": {
				Description:         "If the network is bound to a config template",
				MarkdownDescription: "If the network is bound to a config template",
				Type:                types.BoolType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"copy_from_network_id": {
				Description:         "URL to the network Dashboard UI",
				MarkdownDescription: "URL to the network Dashboard UI",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
		},
	}, nil
}

func (r *NetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworkResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("organization_id: %s", data.OrganizationId.ValueString()))
		return
	}

	var productTypes []string
	var validProductTypes = []string{"appliance", "switch", "wireless", "systemsManager", "camera", "cellularGateway", "sensor"}

	// ProductTypes
	for _, product := range data.ProductTypes {
		for _, productType := range validProductTypes {
			if product.ValueString() == productType {
				productTypes = append(productTypes, productType)
			}
		}
	}

	// check for product types input
	if len(validProductTypes) < 1 {
		resp.Diagnostics.AddError("Missing required input product_types", fmt.Sprintf("Must be one of 'wireless', "+
			"'appliance','switch', 'systemsManager', 'camera', 'cellularGateway' or 'sensor'"))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganizationNetwork := apiclient.NewInlineObject207(data.Name.ValueString(), productTypes) // InlineObject207 |

	createOrganizationNetwork.SetTimeZone(data.Timezone.ValueString())

	// Tags
	var tags []string
	for _, attribute := range data.Tags {
		tags = append(tags, attribute.String())
	}
	createOrganizationNetwork.SetTags(tags)

	createOrganizationNetwork.SetNotes(data.Notes.ValueString())

	if len(data.CopyFromNetworkId.ValueString()) > 0 {
		createOrganizationNetwork.SetCopyFromNetworkId(data.CopyFromNetworkId.ValueString())
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganizationNetwork(ctx, data.OrganizationId.ValueString()).CreateOrganizationNetwork(*createOrganizationNetwork).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append()

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.NetworkId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	data.ProductTypes = nil
	for _, product := range inlineResp.ProductTypes {
		data.ProductTypes = append(data.ProductTypes, types.StringValue(product))
	}

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	// tags attribute
	data.Tags = nil
	for _, tag := range inlineResp.Tags {
		trimmedTag := strings.Trim(tag, "\"")                        // BUG: string wrapped in double quotes...
		data.Tags = append(data.Tags, types.StringValue(trimmedTag)) //
	}
	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())

	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

	if len(data.CopyFromNetworkId.ValueString()) < 1 {
		data.CopyFromNetworkId = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworkResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// check for required parameters
	if len(data.NetworkId.ValueString()) == 0 {
		resp.Diagnostics.AddError("Missing network_Id", fmt.Sprintf("network_id: %s", data.NetworkId.ValueString()))
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.NetworksApi.GetNetwork(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.NetworkId = types.StringValue(inlineResp.GetId())
	data.OrganizationId = types.StringValue(inlineResp.GetOrganizationId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	data.ProductTypes = nil
	for _, product := range inlineResp.ProductTypes {
		data.ProductTypes = append(data.ProductTypes, types.StringValue(product))
	}

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	// tags attribute
	data.Tags = nil
	for _, tag := range inlineResp.Tags {
		trimmedTag := strings.Trim(tag, "\"")                        // BUG: string wrapped in double quotes...
		data.Tags = append(data.Tags, types.StringValue(trimmedTag)) //
	}

	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())
	data.CopyFromNetworkId = types.StringNull()

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworkResourceModel
	var stateData *NetworkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	// check for required parameters
	if len(data.NetworkId.ValueString()) == 0 {
		data.NetworkId = stateData.NetworkId
	}

	if len(data.NetworkId.ValueString()) == 0 {
		resp.Diagnostics.AddError("Missing network_Id", fmt.Sprintf("network_id: %s", data.NetworkId.ValueString()))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	updateNetwork := apiclient.NewInlineObject25()
	updateNetwork.SetName(data.Name.ValueString())
	updateNetwork.SetTimeZone(data.Timezone.ValueString())

	// Tags
	var tags []string
	for _, attribute := range data.Tags {
		tags = append(tags, attribute.String())
	}
	updateNetwork.SetTags(tags)

	// check for enrollment state
	updateNetwork.SetEnrollmentString(data.EnrollmentString.ValueString())

	updateNetwork.SetNotes(data.Notes.ValueString())

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetwork(context.Background(),
		data.NetworkId.ValueString()).UpdateNetwork(*updateNetwork).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.NetworkId = types.StringValue(inlineResp.GetId())
	data.OrganizationId = types.StringValue(inlineResp.GetOrganizationId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	data.ProductTypes = stateData.ProductTypes

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	// tags attribute
	data.Tags = nil
	for _, tag := range inlineResp.Tags {
		trimmedTag := strings.Trim(tag, "\"")                        // BUG: string wrapped in double quotes...
		data.Tags = append(data.Tags, types.StringValue(trimmedTag)) //
	}

	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())
	data.CopyFromNetworkId = types.StringNull()

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworkResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// check for required parameters
	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing network_Id", fmt.Sprintf("network_id: %s", data.NetworkId.ValueString()))
		return
	}

	// Initialize provider client and make API call
	httpResp, err := r.client.NetworksApi.DeleteNetwork(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, network_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
