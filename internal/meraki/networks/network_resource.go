package networks

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
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
	NetworkId               types.String `tfsdk:"network_id" json:"id"`
	OrganizationId          types.String `tfsdk:"organization_id" json:"organizationId"`
	Name                    types.String `tfsdk:"name"`
	ProductTypes            types.Set    `tfsdk:"product_types" json:"productTypes"`
	Timezone                types.String `tfsdk:"timezone" json:"timeZone"`
	Tags                    types.Set    `tfsdk:"tags"`
	EnrollmentString        types.String `tfsdk:"enrollment_string" json:"enrollmentString"`
	Url                     types.String `tfsdk:"url"`
	Notes                   types.String `tfsdk:"notes"`
	IsBoundToConfigTemplate types.Bool   `tfsdk:"is_bound_to_config_template" json:"IsBoundToConfigTemplate"`
	CopyFromNetworkId       types.String `tfsdk:"copy_from_network_id" json:"copyFromNetworkId"`
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
					stringvalidator.LengthBetween(1, 31),
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
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Network name",
				Optional:            true,
				Computed:            true,
			},
			"product_types": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf([]string{"appliance", "switch", "wireless", "systemsManager", "camera", "cellularGateway", "sensor", "cloudGateway"}...),
					),
				},
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "Timezone of the network",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
					stringvalidator.LengthBetween(1, 50),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to the network Dashboard UI",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the network",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"copy_from_network_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the network to copy configuration from. Other provided parameters will override the copied configuration, except type which must match this network's type exactly.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_bound_to_config_template": schema.BoolAttribute{
				MarkdownDescription: "If the network is bound to a config template",
				Computed:            true,
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

func updateNetworksNetworksResourceCreatePayload(plan *NetworkResourceModel) (openApiClient.CreateOrganizationNetworkRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// name
	name := plan.Name.ValueString()

	// ProductTypes
	var productTypes []string
	if !plan.ProductTypes.IsNull() && !plan.ProductTypes.IsUnknown() {
		for _, product := range plan.ProductTypes.Elements() {
			pt := fmt.Sprint(strings.Trim(product.String(), "\""))
			productTypes = append(productTypes, pt)
		}
	}

	// Create HTTP request body
	payload := openApiClient.NewCreateOrganizationNetworkRequest(name, productTypes)

	// Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		for _, tag := range plan.Tags.Elements() {
			t := fmt.Sprint(strings.Trim(tag.String(), "\""))
			tags = append(tags, t)
		}
		payload.SetTags(tags)
	}

	//    TimeZone
	if !plan.Timezone.IsNull() && !plan.Timezone.IsUnknown() {
		payload.SetTimeZone(plan.Timezone.ValueString())
	}

	// CopyFromNetworkId
	if !plan.CopyFromNetworkId.IsNull() && !plan.CopyFromNetworkId.IsUnknown() {
		payload.SetCopyFromNetworkId(plan.CopyFromNetworkId.ValueString())
	}

	// Notes
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		payload.SetNotes(plan.Notes.ValueString())
	}

	return *payload, diags

}

func updateNetworksNetworksResourceUpdatePayload(plan *NetworkResourceModel) (openApiClient.UpdateNetworkRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload := openApiClient.NewUpdateNetworkRequest()

	//   Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())
	}

	//    TimeZone
	if !plan.Timezone.IsNull() && !plan.Timezone.IsUnknown() {
		payload.SetTimeZone(plan.Timezone.ValueString())
	}

	//    Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		for _, tag := range plan.Tags.Elements() {
			tags = append(tags, tag.String())
		}
		payload.SetTags(tags)
	}

	//    EnrollmentString
	if !plan.EnrollmentString.IsNull() && !plan.EnrollmentString.IsUnknown() {
		payload.SetEnrollmentString(plan.EnrollmentString.ValueString())
	}

	//    Notes
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		payload.SetNotes(plan.Notes.ValueString())
	}

	return *payload, diags

}

func createNetworksNetworksResourceState(ctx context.Context, state *NetworkResourceModel, inlineResp *openApiClient.GetNetwork200Response, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	//  Id (NetworkId)
	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		state.NetworkId = types.StringValue(inlineResp.GetId())
	}

	orgId := fmt.Sprint(strings.Trim(inlineResp.GetOrganizationId(), "\""))
	state.OrganizationId = types.StringValue(orgId)

	//  Id (Terraform Resource)
	if !state.NetworkId.IsNull() || !state.NetworkId.IsUnknown() && !state.OrganizationId.IsNull() || !state.OrganizationId.IsUnknown() {
		importId := state.OrganizationId.String() + "," + inlineResp.GetId()
		state.Id = types.StringValue(importId)
	} else {
		state.Id = types.StringNull()
	}

	//    Name
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name = types.StringValue(inlineResp.GetName())
	}

	//    ProductTypes
	if state.ProductTypes.IsNull() || state.ProductTypes.IsUnknown() {

		var productTypesList []string

		productTypesList = append(productTypesList, inlineResp.ProductTypes...)

		productTypesListObj, err := types.SetValueFrom(ctx, types.StringType, productTypesList)
		if err.HasError() {
			diags.Append(err...)
		}

		state.ProductTypes = productTypesListObj

	}

	//  	Timezone
	if state.Timezone.IsNull() || state.Timezone.IsUnknown() {
		state.Timezone = types.StringValue(inlineResp.GetTimeZone())
	}

	//    Tags
	if state.Tags.IsNull() || state.Tags.IsUnknown() {

		// Tags
		var tagsList []string
		for _, tag := range inlineResp.Tags {
			// Strip any extra quotes from the tags
			tagsList = append(tagsList, strings.Trim(tag, `"`))
		}
		tagsListObj, err := types.SetValueFrom(ctx, types.StringType, tagsList)
		if err.HasError() {
			diags.Append(err...)
		}
		state.Tags = tagsListObj

	}

	//    EnrollmentString
	if state.EnrollmentString.IsNull() || state.EnrollmentString.IsUnknown() {
		if inlineResp.GetEnrollmentString() == "" {
			state.EnrollmentString = types.StringNull()
		} else {
			state.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
		}

	}

	//    Url
	if state.Url.IsNull() || state.Url.IsUnknown() {
		state.Url = types.StringValue(inlineResp.GetUrl())
	}

	//    Notes
	if state.Notes.IsNull() || state.Notes.IsUnknown() {
		state.Notes = types.StringValue(inlineResp.GetNotes())
	}

	//    IsBoundToConfigTemplate
	if state.IsBoundToConfigTemplate.IsNull() || state.IsBoundToConfigTemplate.IsUnknown() {
		state.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())
	}

	// CopyFromNetworkId
	if state.CopyFromNetworkId.IsNull() || state.CopyFromNetworkId.IsUnknown() {
		state.CopyFromNetworkId = types.StringNull()
	}

	return diags
}

func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	createPayload, createPayloadDiags := updateNetworksNetworksResourceCreatePayload(&plan)
	if createPayloadDiags.HasError() {

		resp.Diagnostics.AddError(
			"Error creating payload",
			fmt.Sprintf("Unexpected error: %s", createPayloadDiags.Errors()),
		)
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganizationNetwork(ctx, plan.OrganizationId.ValueString()).CreateOrganizationNetworkRequest(createPayload).Execute()
	if err != nil {

		resp.Diagnostics.AddError(
			"Create Network HTTP Client Failure",
			err.Error(),
		)

		return
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
	}

	diags := createNetworksNetworksResourceState(ctx, &plan, inlineResp, httpResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.NetworksApi.GetNetwork(context.Background(), state.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Network HTTP Client Failure",
			err.Error(),
		)

		return
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
	}

	// unmarshal Payload into state
	diags = createNetworksNetworksResourceState(ctx, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	updatePayload, updatePayloadDiags := updateNetworksNetworksResourceUpdatePayload(&plan)
	if updatePayloadDiags.HasError() {
		tflog.Error(ctx, "Failed to update resource payload", map[string]interface{}{
			"error": updatePayloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error updating network payload",
			fmt.Sprintf("Unexpected error: %s", updatePayloadDiags),
		)
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetwork(ctx, plan.NetworkId.ValueString()).UpdateNetworkRequest(updatePayload).Execute()
	if err != nil {

		resp.Diagnostics.AddError(
			"Update Network HTTP Client Failure",
			err.Error(),
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
	}

	// Save updated data into Terraform state
	diags = createNetworksNetworksResourceState(ctx, &plan, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkResourceModel

	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		time.Sleep(retryDelay * time.Millisecond)

		httpResp, err := r.client.NetworksApi.DeleteNetwork(context.Background(), state.NetworkId.ValueString()).Execute()
		return nil, httpResp, err
	}

	// HTTP DELETE METHOD does not leverage the retry-after header and throws 400 errors.
	_, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to delete resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error deleting network",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
			resp.Diagnostics.AddError(
				"Error deleting network",
				err.Error(),
			)
		}
		return
	}

	// check for HTTP errors
	if httpResp.StatusCode != 204 {
		if err != nil {
			resp.Diagnostics.AddError(
				"Delete Network HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
		}
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
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
