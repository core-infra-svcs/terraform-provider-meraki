package provider

import (
	"context"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	Id                      types.String `tfsdk:"id"`
	EnrollmentString        types.String `tfsdk:"enrollment_string"`
	NetworkId               types.String `tfsdk:"network_id"`
	IsBoundToConfigTemplate types.Bool   `tfsdk:"is_bound_to_config_template"`
	Name                    types.String `tfsdk:"name"`
	Notes                   types.String `tfsdk:"notes"`
	OrganizationId          types.String `tfsdk:"organization_id"`
	ProductTypes            types.Set    `tfsdk:"product_types"`
	Tags                    types.Set    `tfsdk:"tags"`
	Timezone                types.String `tfsdk:"timezone"`
	Url                     types.String `tfsdk:"url"`
	CopyFromNetworkId       types.String `tfsdk:"copy_from_network_id"`
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
				Type:                types.SetType{ElemType: types.StringType},
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
				Type:                types.SetType{ElemType: types.StringType},
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

	// validate inputs and update product types list
	//expectedProductTypes := []string{"appliance", "switch", "systemsManager", "camera", "cellularGateway", "sensor"}
	var productTypes []string

	if len(data.ProductTypes.Elements()) > 0 {

		// compare with list of expected product type strings
		for _, attribute := range data.ProductTypes.Elements() {
			switch string(attribute.String()) {
			case
				"appliance",
				"switch",
				"systemsManager",
				"camera",
				"cellularGateway",
				"sensor":
				productTypes = append(productTypes, attribute.String())
			default:
				resp.Diagnostics.AddError("Invalid Entry", fmt.Sprintf("The input: %s, is not a valid product type", attribute.String()))
			}
		}
	}

	// check for product types input
	if len(productTypes) < 1 {
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
	if data.Tags.IsNull() != true {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tags = append(tags, attribute.String())
		}
		createOrganizationNetwork.SetTags(tags)
	}

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
	data.OrganizationId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	if products := inlineResp.ProductTypes; products != nil {
		var productList []attr.Value
		for _, v := range products {
			product := types.StringValue(v)
			productList = append(productList, product)
		}

		data.ProductTypes, _ = types.SetValue(types.StringType, productList)
	} else {
		data.ProductTypes = types.SetNull(types.StringType)
	}

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	// tags attribute
	if tags := inlineResp.Tags; tags != nil {
		var tagList []attr.Value
		for _, v := range tags {
			tag := types.StringValue(v)
			tagList = append(tagList, tag)
		}

		data.Tags, _ = types.SetValue(types.StringType, tagList)
	} else {
		data.Tags = types.SetNull(types.StringType)
	}

	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

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
	if len(data.NetworkId.ValueString()) < 1 {
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
	data.OrganizationId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	if products := inlineResp.ProductTypes; products != nil {
		var productList []attr.Value
		for _, v := range products {
			product := types.StringValue(v)
			productList = append(productList, product)
		}

		data.ProductTypes, _ = types.SetValue(types.StringType, productList)
	} else {
		data.ProductTypes = types.SetNull(types.StringType)
	}

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	// tags attribute
	if tags := inlineResp.Tags; tags != nil {
		var tagList []attr.Value
		for _, v := range tags {
			tag := types.StringValue(v)
			tagList = append(tagList, tag)
		}

		data.Tags, _ = types.SetValue(types.StringType, tagList)
	} else {
		data.Tags = types.SetNull(types.StringType)
	}

	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworkResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// check for required parameters
	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing network_Id", fmt.Sprintf("network_id: %s", data.NetworkId.ValueString()))
		return
	}

	// Create HTTP request body
	updateNetwork := apiclient.NewInlineObject25()
	updateNetwork.SetName(data.Name.ValueString())
	updateNetwork.SetTimeZone(data.Timezone.ValueString())

	// Tags
	if data.Tags.IsNull() != true {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tags = append(tags, attribute.String())
		}
		updateNetwork.SetTags(tags)
	}

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
	data.OrganizationId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	if products := inlineResp.ProductTypes; products != nil {
		var productList []attr.Value
		for _, v := range products {
			product := types.StringValue(v)
			productList = append(productList, product)
		}

		data.ProductTypes, _ = types.SetValue(types.StringType, productList)
	} else {
		data.ProductTypes = types.SetNull(types.StringType)
	}

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	// tags attribute
	if tags := inlineResp.Tags; tags != nil {
		var tagList []attr.Value
		for _, v := range tags {
			tag := types.StringValue(v)
			tagList = append(tagList, tag)
		}

		data.Tags, _ = types.SetValue(types.StringType, tagList)
	} else {
		data.Tags = types.SetNull(types.StringType)
	}

	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

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
}
