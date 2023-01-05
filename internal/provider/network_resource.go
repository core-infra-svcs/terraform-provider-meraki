package provider

import (
	"context"
	"fmt"
	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworkResource{}
	_ resource.ResourceWithConfigure   = &NetworkResource{}
	_ resource.ResourceWithImportState = &NetworkResource{}
)

func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}

// NetworkResource defines the resource implementation.
type NetworkResource struct {
	client *openApiClient.APIClient
}

// NetworkResourceModel describes the resource data model.
type NetworkResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	NetworkId               types.String `tfsdk:"network_id"`
	OrganizationId          types.String `tfsdk:"organization_id"`
	Name                    types.String `tfsdk:"name"`
	ProductTypes            types.Set    `tfsdk:"product_types"`
	Timezone                types.String `tfsdk:"timezone"`
	Tags                    types.Set    `tfsdk:"tags"`
	EnrollmentString        types.String `tfsdk:"enrollment_string"`
	Url                     types.String `tfsdk:"url"`
	Notes                   types.String `tfsdk:"notes"`
	IsBoundToConfigTemplate types.Bool   `tfsdk:"is_bound_to_config_template"`
	CopyFromNetworkId       types.String `tfsdk:"copy_from_network_id"`
}

func (r *NetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *NetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the networks that the user has privileges on in an organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Network name",
				Optional:            true,
				Computed:            true,
			},
			"product_types": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf([]string{"appliance", "switch", "wireless", "systemsManager", "camera", "cellularGateway", "sensor"}...),
						stringvalidator.LengthAtLeast(5),
					),
				},
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "Timezone of the network",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				Description: "Network tags",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enrollment_string": schema.StringAttribute{
				MarkdownDescription: "A unique identifier which can be used for device enrollment or easy access through the Meraki SM Registration page or the Self Service Portal. Once enabled, a network enrollment strings can be changed but they cannot be deleted.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.NoneOf([]string{";", ":", "@", "=", "&", "$", "!", "‘", "“", ",", "?", ".", "(", ")", "{", "}", "[", "]", "\\", "*", "+", "/", "#", "<", ">", "|", "^", "%"}...),
					stringvalidator.LengthBetween(4, 50),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to the network Dashboard UI",
				Optional:            true,
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the network",
				Optional:            true,
				Computed:            true,
			},
			"is_bound_to_config_template": schema.BoolAttribute{
				MarkdownDescription: "If the network is bound to a config template",
				Optional:            true,
				Computed:            true,
			},
			"copy_from_network_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the network to copy configuration from. Other provided parameters will override the copied configuration, except type which must match this network's type exactly.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
		},
	}
}

func (r *NetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworkResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganizationNetwork := openApiClient.NewInlineObject207(data.Name.ValueString(), nil)
	createOrganizationNetwork.SetTimeZone(data.Timezone.ValueString())

	// ProductTypes
	var productTypes []string
	for _, product := range data.ProductTypes.Elements() {
		pt := fmt.Sprint(strings.Trim(product.String(), "\""))
		productTypes = append(productTypes, pt)
	}
	createOrganizationNetwork.SetProductTypes(productTypes)

	if resp.Diagnostics.HasError() {
		return
	}

	createOrganizationNetwork.SetProductTypes(productTypes)

	// Tags
	if !data.Tags.IsUnknown() {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tags = append(tags, attribute.String())
		}
		createOrganizationNetwork.SetTags(tags)
	}

	// Notes
	if !data.Notes.IsUnknown() {
		createOrganizationNetwork.SetNotes(data.Notes.ValueString())
	}

	// CopyFromNetworkId
	if !data.CopyFromNetworkId.IsUnknown() {
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
		return
	} else {
		resp.Diagnostics.Append()
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.NetworkId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.Timezone = types.StringValue(inlineResp.GetTimeZone())
	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

	// product types attribute
	if inlineResp.ProductTypes != nil {
		var pt []attr.Value
		for _, productTypeResp := range inlineResp.ProductTypes {
			pt = append(pt, types.StringValue(productTypeResp))
		}
		data.ProductTypes, _ = types.SetValue(types.StringType, pt)
	}

	if data.CopyFromNetworkId.IsUnknown() {
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

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.NetworksApi.GetNetwork(context.Background(), data.NetworkId.ValueString()).Execute()
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

	// Check for API success inlineResp code
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

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.NetworkId = types.StringValue(inlineResp.GetId())
	data.OrganizationId = types.StringValue(inlineResp.GetOrganizationId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.Timezone = types.StringValue(inlineResp.GetTimeZone())
	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

	// product types attribute
	if inlineResp.ProductTypes != nil {
		var pt []attr.Value
		for _, productTypeResp := range inlineResp.ProductTypes {
			pt = append(pt, types.StringValue(productTypeResp))
		}
		data.ProductTypes, _ = types.SetValue(types.StringType, pt)
	}

	if data.CopyFromNetworkId.IsUnknown() {
		data.CopyFromNetworkId = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	updateNetwork := openApiClient.NewInlineObject25()
	updateNetwork.SetName(data.Name.ValueString())
	updateNetwork.SetTimeZone(data.Timezone.ValueString())

	// Tags
	if !data.Tags.IsUnknown() {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tags = append(tags, attribute.String())
		}
		updateNetwork.SetTags(tags)
	}

	// Enrollment String
	if !data.EnrollmentString.IsUnknown() {
		updateNetwork.SetEnrollmentString(data.EnrollmentString.ValueString())
	}

	// Notes
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

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.NetworkId = types.StringValue(inlineResp.GetId())
	data.OrganizationId = types.StringValue(inlineResp.GetOrganizationId())
	data.Name = types.StringValue(inlineResp.GetName())

	// product types attribute
	if inlineResp.ProductTypes != nil {
		var pt []attr.Value
		for _, productTypeResp := range inlineResp.ProductTypes {
			pt = append(pt, types.StringValue(productTypeResp))
		}
		data.ProductTypes, _ = types.SetValue(types.StringType, pt)
	}

	data.Timezone = types.StringValue(inlineResp.GetTimeZone())

	data.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.Notes = types.StringValue(inlineResp.GetNotes())
	data.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())

	if data.CopyFromNetworkId.IsUnknown() {
		data.CopyFromNetworkId = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworkResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// HTTP DELETE METHOD does not leverage the retry-after header and throws 400 errors.
	retries := 100 // This is the only way to ensure a scaled result.
	wait := 1
	var deletedFromMerakiPortal bool
	deletedFromMerakiPortal = false

	for retries > 0 {

		// Initialize provider client and make API call
		httpResp, err := r.client.NetworksApi.DeleteNetwork(context.Background(), data.NetworkId.ValueString()).Execute()

		if httpResp.StatusCode == 204 {

			// check for HTTP errors
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

			// Check for errors after diagnostics collected
			if resp.Diagnostics.HasError() {
				return
			} else {
				resp.Diagnostics.Append()
			}

			deletedFromMerakiPortal = true

			// escape loop
			break

		} else {

			// decrement retry counter
			retries -= 1

			// exponential wait
			time.Sleep(time.Duration(wait) * time.Second)
			wait += 1
		}
	}

	// Check if deleted from Meraki Portal was successful
	if deletedFromMerakiPortal == true {
		resp.State.RemoveResource(ctx)

		// Write logs using the tflog package
		tflog.Trace(ctx, "removed resource")
	} else {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			"Failed to delete resource",
		)
		return
	}

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
