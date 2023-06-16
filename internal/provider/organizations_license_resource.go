package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &OrganizationsLicenseResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &OrganizationsLicenseResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &OrganizationsLicenseResource{} // Interface for resources with import state functionality
)

func NewOrganizationsLicenseResource() resource.Resource {
	return &OrganizationsLicenseResource{}
}

type OrganizationsLicenseResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The OrganizationsLicenseResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type OrganizationsLicenseResourceModel struct {
	Id                        jsontypes.String                                                             `tfsdk:"id"`
	OrganizationId            jsontypes.String                                                             `tfsdk:"organization_id"`
	LicenseId                 jsontypes.String                                                             `tfsdk:"license_id"`
	LicenseType               jsontypes.String                                                             `tfsdk:"license_type"`
	LicenseKey                jsontypes.String                                                             `tfsdk:"license_key"`
	OrderNumber               jsontypes.String                                                             `tfsdk:"order_number"`
	DeviceSerial              jsontypes.String                                                             `tfsdk:"device_serial"`
	NetworkId                 jsontypes.String                                                             `tfsdk:"network_id"`
	State                     jsontypes.String                                                             `tfsdk:"state"`
	ClaimDate                 jsontypes.String                                                             `tfsdk:"claim_date"`
	ActivationDate            jsontypes.String                                                             `tfsdk:"activation_date"`
	ExpirationDate            jsontypes.String                                                             `tfsdk:"expiration_date"`
	HeadLicenseId             jsontypes.String                                                             `tfsdk:"head_license_id"`
	SeatCount                 jsontypes.Int64                                                              `tfsdk:"seat_count"`
	TotalDurationInDays       jsontypes.Int64                                                              `tfsdk:"total_duration_in_days"`
	DurationInDays            jsontypes.Int64                                                              `tfsdk:"duration_in_days"`
	PermanentlyQueuedLicenses []openApiClient.OrganizationsOrganizationIdLicensesPermanentlyQueuedLicenses `tfsdk:"permanently_queued_licenses"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *OrganizationsLicenseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_organizations_license"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *OrganizationsLicenseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Organizations License Updates a license",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"license_id": schema.StringAttribute{
				MarkdownDescription: "License ID",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"device_serial": schema.StringAttribute{
				MarkdownDescription: "The serial number of the device to assign this license to. Set this to null to unassign the license. If a different license is already active on the device, this parameter will control queueing/dequeuing this license.",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"license_type": schema.StringAttribute{
				MarkdownDescription: "License Type.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"license_key": schema.StringAttribute{
				MarkdownDescription: "License Key.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"order_number": schema.StringAttribute{
				MarkdownDescription: "Order Number.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "ID of the network the license is assigned to.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the license. All queued licenses have a status of `recentlyQueued`.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"claim_date": schema.StringAttribute{
				MarkdownDescription: "The date the license was claimed into the organization.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"activation_date": schema.StringAttribute{
				MarkdownDescription: "The date the license started burning.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"expiration_date": schema.StringAttribute{
				MarkdownDescription: "The date the license will expire.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"head_license_id": schema.StringAttribute{
				MarkdownDescription: "The id of the head license this license is queued behind. If there is no head license, it returns nil.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"seat_count": schema.Int64Attribute{
				MarkdownDescription: "The number of seats of the license. Only applicable to SM licenses.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"total_duration_in_days": schema.Int64Attribute{
				MarkdownDescription: "The duration of the license plus all permanently queued licenses associated with it.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"duration_in_days": schema.Int64Attribute{
				MarkdownDescription: "The duration of the individual license.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"permanently_queued_licenses": schema.SingleNestedAttribute{
				MarkdownDescription: "DEPRECATED List of permanently queued licenses attached to the license. Instead, use /organizations/{organizationId}/licenses?deviceSerial= to retrieved queued licenses for a given device.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Permanently queued license ID.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"license_type": schema.StringAttribute{
						MarkdownDescription: "License type.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"license_key": schema.StringAttribute{
						MarkdownDescription: "License key.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"order_number": schema.StringAttribute{
						MarkdownDescription: "Order number.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"duration_in_days": schema.Int64Attribute{
						MarkdownDescription: "The duration of the individual license.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *OrganizationsLicenseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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
func (r *OrganizationsLicenseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *OrganizationsLicenseResourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	updateOrganizationLicense := *openApiClient.NewInlineObject208()
	updateOrganizationLicense.SetDeviceSerial(data.DeviceSerial.ValueString())

	inlineResp, httpResp, err := r.client.LicensesApi.UpdateOrganizationLicense(context.Background(), data.OrganizationId.ValueString(), data.LicenseId.ValueString()).UpdateOrganizationLicense(updateOrganizationLicense).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.LicenseType = jsontypes.StringValue(inlineResp.GetLicenseType())
	data.LicenseKey = jsontypes.StringValue(inlineResp.GetLicenseKey())
	data.OrderNumber = jsontypes.StringValue(inlineResp.GetOrderNumber())
	data.DeviceSerial = jsontypes.StringValue(inlineResp.GetDeviceSerial())
	data.NetworkId = jsontypes.StringValue(inlineResp.GetNetworkId())
	data.State = jsontypes.StringValue(inlineResp.GetState())
	data.ClaimDate = jsontypes.StringValue(inlineResp.GetClaimDate())
	data.ActivationDate = jsontypes.StringValue(inlineResp.GetActivationDate())
	data.ExpirationDate = jsontypes.StringValue(inlineResp.GetExpirationDate())
	data.HeadLicenseId = jsontypes.StringValue(inlineResp.GetHeadLicenseId())
	data.SeatCount = jsontypes.Int64Value(int64(inlineResp.GetSeatCount()))
	data.TotalDurationInDays = jsontypes.Int64Value(int64(inlineResp.GetTotalDurationInDays()))
	data.DurationInDays = jsontypes.Int64Value(int64(inlineResp.GetDurationInDays()))
	data.PermanentlyQueuedLicenses = inlineResp.GetPermanentlyQueuedLicenses()

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *OrganizationsLicenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsLicenseResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.LicensesApi.GetOrganizationLicense(ctx, data.OrganizationId.ValueString(), data.LicenseId.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
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
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *OrganizationsLicenseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *OrganizationsLicenseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateOrganizationLicense := *openApiClient.NewInlineObject208()
	updateOrganizationLicense.SetDeviceSerial(data.DeviceSerial.ValueString())

	inlineResp, httpResp, err := r.client.LicensesApi.UpdateOrganizationLicense(context.Background(), data.OrganizationId.ValueString(), data.LicenseId.ValueString()).UpdateOrganizationLicense(updateOrganizationLicense).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.LicenseType = jsontypes.StringValue(inlineResp.GetLicenseType())
	data.LicenseKey = jsontypes.StringValue(inlineResp.GetLicenseKey())
	data.OrderNumber = jsontypes.StringValue(inlineResp.GetOrderNumber())
	data.DeviceSerial = jsontypes.StringValue(inlineResp.GetDeviceSerial())
	data.NetworkId = jsontypes.StringValue(inlineResp.GetNetworkId())
	data.State = jsontypes.StringValue(inlineResp.GetState())
	data.ClaimDate = jsontypes.StringValue(inlineResp.GetClaimDate())
	data.ActivationDate = jsontypes.StringValue(inlineResp.GetActivationDate())
	data.ExpirationDate = jsontypes.StringValue(inlineResp.GetExpirationDate())
	data.HeadLicenseId = jsontypes.StringValue(inlineResp.GetHeadLicenseId())
	data.SeatCount = jsontypes.Int64Value(int64(inlineResp.GetSeatCount()))
	data.TotalDurationInDays = jsontypes.Int64Value(int64(inlineResp.GetTotalDurationInDays()))
	data.DurationInDays = jsontypes.Int64Value(int64(inlineResp.GetDurationInDays()))
	data.PermanentlyQueuedLicenses = inlineResp.GetPermanentlyQueuedLicenses()

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *OrganizationsLicenseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *OrganizationsLicenseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateOrganizationLicense := *openApiClient.NewInlineObject208()
	updateOrganizationLicense.SetDeviceSerial(data.DeviceSerial.ValueString())

	inlineResp, httpResp, err := r.client.LicensesApi.UpdateOrganizationLicense(context.Background(), data.OrganizationId.ValueString(), data.LicenseId.ValueString()).UpdateOrganizationLicense(updateOrganizationLicense).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.LicenseType = jsontypes.StringValue(inlineResp.GetLicenseType())
	data.LicenseKey = jsontypes.StringValue(inlineResp.GetLicenseKey())
	data.OrderNumber = jsontypes.StringValue(inlineResp.GetOrderNumber())
	data.DeviceSerial = jsontypes.StringValue(inlineResp.GetDeviceSerial())
	data.NetworkId = jsontypes.StringValue(inlineResp.GetNetworkId())
	data.State = jsontypes.StringValue(inlineResp.GetState())
	data.ClaimDate = jsontypes.StringValue(inlineResp.GetClaimDate())
	data.ActivationDate = jsontypes.StringValue(inlineResp.GetActivationDate())
	data.ExpirationDate = jsontypes.StringValue(inlineResp.GetExpirationDate())
	data.HeadLicenseId = jsontypes.StringValue(inlineResp.GetHeadLicenseId())
	data.SeatCount = jsontypes.Int64Value(int64(inlineResp.GetSeatCount()))
	data.TotalDurationInDays = jsontypes.Int64Value(int64(inlineResp.GetTotalDurationInDays()))
	data.DurationInDays = jsontypes.Int64Value(int64(inlineResp.GetDurationInDays()))
	data.PermanentlyQueuedLicenses = inlineResp.GetPermanentlyQueuedLicenses()

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *OrganizationsLicenseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, license_id. Got: %q", req.ID),
		)
		return
	}

	// Set the attributes required for making a Read API call in the state.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("license_id"), idParts[1])...)

	// If there were any errors setting the attributes, return early.
	if resp.Diagnostics.HasError() {
		return
	}

}
