package organizations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// The below var block ensures that the provider defined types fully satisfy the required
// interfaces for a Terraform resource. This includes resource.Resource, resource.ResourceWithConfigure,
// and resource.ResourceWithImportState.

// OrganizationsClaimResource struct.
var (
	_ resource.Resource              = &OrganizationsClaimResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure = &OrganizationsClaimResource{} // Interface for resources with configuration methods
)

// The NewOrganizationsClaimResource function is a constructor for the resource. This function needs
// to be added to the list of Resources in provider.go: func (p *ScaffoldingProvider) Resources.
// If it's not added, the provider won't be aware of this resource's existence.
func NewOrganizationsClaimResource() resource.Resource {
	return &OrganizationsClaimResource{}
}

// OrganizationsClaimResource struct defines the structure for this resource.
// It includes an APIClient field for making requests to the Meraki API.
// If additional fields are required (e.g., for caching or for tracking internal state), add them here.
type OrganizationsClaimResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The OrganizationsClaimResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type OrganizationsClaimResourceModel struct {

	// The Id field is mandatory for all resources. It's used for resource identification and is required
	// for the acceptance tests to run.
	Id             jsontypes.String                         `tfsdk:"id"`
	OrganizationId jsontypes.String                         `tfsdk:"organization_id"`
	Orders         []jsontypes.String                       `tfsdk:"orders"`
	Serials        []jsontypes.String                       `tfsdk:"serials"`
	Licences       []OrganizationsClaimResourceModelLicence `tfsdk:"licences"`
}

type OrganizationsClaimResourceModelLicence struct {
	Key  jsontypes.String `tfsdk:"key"`
	Mode jsontypes.String `tfsdk:"mode"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *OrganizationsClaimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source, and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_organizations_claim"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *OrganizationsClaimResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		// It should provide a clear and concise description of the resource.
		MarkdownDescription: "Claim a list of devices, licenses, and/or orders into an organization. When claiming by order, all devices and licenses in the order will be claimed; licenses will be added to the organization and devices will be placed in the organization's inventory.",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				CustomType:          jsontypes.StringType,
				MarkdownDescription: "Organization ID",
			},
			"orders": schema.SetAttribute{
				MarkdownDescription: "The numbers of the orders that should be claimed",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"serials": schema.SetAttribute{
				MarkdownDescription: "The serials of the devices that should be claimed",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"licences": schema.ListNestedAttribute{
				MarkdownDescription: "The licenses that should be claimed",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "The key of the license",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"mode": schema.StringAttribute{
							MarkdownDescription: "Either 'renew' or 'addDevices'. 'addDevices' will increase the license limit, " +
								"while 'renew' will extend the amount of time until expiration. Defaults to 'addDevices'. " +
								"All licenses must be claimed with the same mode, and at most one renewal can be claimed at a time. " +
								"This parameter is legacy and does not apply to organizations with per-device licensing enabled.",
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf("addDevices", "renew"),
							},
						},
					},
				},
				Computed: true,
				Optional: true,
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *OrganizationsClaimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// This allows the resource to use the configured provider for any API calls it needs to make.
	r.client = client
}

// Create method is responsible for creating a new resource.
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *OrganizationsClaimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsClaimResourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	claimIntoOrganizationInventoryRequest := *openApiClient.NewClaimIntoOrganizationInventoryRequest()

	// orders
	var orders []string
	for _, order := range data.Orders {
		orders = append(orders, order.ValueString())
	}

	// serials
	var serials []string
	for _, serial := range data.Serials {
		serials = append(serials, serial.ValueString())
	}

	// licenses
	var licenses []openApiClient.ClaimIntoOrganizationInventoryRequestLicensesInner
	for _, license := range data.Licences {
		var licenseMap openApiClient.ClaimIntoOrganizationInventoryRequestLicensesInner
		licenseMap.SetKey(license.Key.ValueString())
		licenseMap.SetMode(license.Mode.ValueString())
		licenses = append(licenses, licenseMap)
	}

	claimIntoOrganizationInventoryRequest.SetSerials(orders)
	claimIntoOrganizationInventoryRequest.SetLicenses(licenses)
	claimIntoOrganizationInventoryRequest.SetSerials(serials)

	// Remember to handle any potential errors.
	_, httpResp, err := r.client.ConfigureApi.ClaimIntoOrganizationInventory(ctx,
		data.OrganizationId.ValueString()).ClaimIntoOrganizationInventoryRequest(claimIntoOrganizationInventoryRequest).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {

		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue(data.OrganizationId.ValueString())

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *OrganizationsClaimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	inlineRespSerials, httpResp, err := r.client.OrganizationsApi.GetOrganizationDevices(ctx, data.OrganizationId.ValueString()).Execute()
	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Create a list to store the extracted strings
	var extractedSerials []jsontypes.String

	// Iterate over each serial in data.Serials
	for _, serial := range data.Serials {

		// Iterate over each inner response in inlineRespSerials
		for _, innerResp := range inlineRespSerials {

			// Check if the serials match
			if innerResp.Serial != nil && *innerResp.Serial == serial.ValueString() {

				// Extract the desired serial and add it to the list of strings
				if innerResp.Serial != nil {
					extractedSerials = append(extractedSerials, jsontypes.StringValue(*innerResp.Serial))
				}
			}
		}
	}

	// Add to list if found
	if len(extractedSerials) != 0 {
		data.Serials = extractedSerials
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *OrganizationsClaimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsClaimResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *OrganizationsClaimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *OrganizationsClaimResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// serials
	var serials []string
	for _, serial := range data.Serials {
		serials = append(serials, serial.ValueString())
	}

	releaseFromOrganizationInventoryRequest := *openApiClient.NewReleaseFromOrganizationInventoryRequest() // ReleaseFromOrganizationInventoryRequest |  (optional)
	releaseFromOrganizationInventoryRequest.SetSerials(serials)

	_, httpResp, err := r.client.ConfigureApi.ReleaseFromOrganizationInventory(ctx, data.OrganizationId.ValueString()).ReleaseFromOrganizationInventoryRequest(releaseFromOrganizationInventoryRequest).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}
