package provider

import (
	"context"
	"fmt"
	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithConfigure   = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
)

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

// OrganizationResource defines the resource implementation.
type OrganizationResource struct {
	client *openApiClient.APIClient
}

// OrganizationResourceModel describes the resource data model.
type OrganizationResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	ApiEnabled             types.Bool   `tfsdk:"api_enabled"`
	CloudRegionName        types.String `tfsdk:"cloud_region_name"`
	ManagementDetailsName  types.String `tfsdk:"management_details_name"`
	ManagementDetailsValue types.String `tfsdk:"management_details_value"`
	OrgId                  types.String `tfsdk:"organization_id"`
	LicensingModel         types.String `tfsdk:"licensing_model"`
	Name                   types.String `tfsdk:"name"`
	Url                    types.String `tfsdk:"url"`
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the organizations that the user has privileges on",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"api_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable API access",
				Optional:            true,
				Computed:            true,
			},
			"cloud_region_name": schema.StringAttribute{
				MarkdownDescription: "Name of region",
				Optional:            true,
				Computed:            true,
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
			"licensing_model": schema.StringAttribute{
				MarkdownDescription: "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"co-term", "per-device", "subscription"}...),
				},
			},
			"management_details_name": schema.StringAttribute{
				MarkdownDescription: "Name of management data",
				Optional:            true,
				Computed:            true,
			},
			"management_details_value": schema.StringAttribute{
				MarkdownDescription: "Value of management data",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Organization name",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Organization URL",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganization := *openApiClient.NewInlineObject165(data.Name.ValueString())

	// Set management details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()
	detail := openApiClient.OrganizationsManagementDetails{
		Name:  &name,
		Value: &value,
	}
	details := []openApiClient.OrganizationsManagementDetails{detail}
	organizationsManagement := openApiClient.OrganizationsManagement{Details: details}
	createOrganization.SetManagement(organizationsManagement)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganization(context.Background()).CreateOrganization(createOrganization).Execute()
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

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.CloudRegionName = types.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() == true {
			data.ManagementDetailsName = types.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = types.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() == true {
			data.ManagementDetailsValue = types.StringValue(managementDetailValue)
		} else {
			data.ManagementDetailsValue = types.StringNull()
		}

	} else {
		data.ManagementDetailsName = types.StringNull()
		data.ManagementDetailsValue = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganization(context.Background(), data.OrgId.ValueString()).Execute()
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
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.CloudRegionName = types.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() == true {
			data.ManagementDetailsName = types.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = types.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() == true {
			data.ManagementDetailsValue = types.StringValue(managementDetailValue)
		} else {
			data.ManagementDetailsValue = types.StringNull()
		}
	} else {
		data.ManagementDetailsName = types.StringNull()
		data.ManagementDetailsValue = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Create HTTP request body
	updateOrganization := openApiClient.NewInlineObject166()
	updateOrganization.SetName(data.Name.ValueString())

	// Set enabled attribute
	var enabled = data.ApiEnabled.ValueBool()
	Api := openApiClient.OrganizationsOrganizationIdApi{Enabled: &enabled}
	updateOrganization.SetApi(Api)

	// Set management details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()
	detail := openApiClient.OrganizationsManagementDetails{
		Name:  &name,
		Value: &value,
	}
	details := []openApiClient.OrganizationsManagementDetails{detail}
	organizationsManagement := openApiClient.OrganizationsManagement{Details: details}
	updateOrganization.SetManagement(organizationsManagement)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.UpdateOrganization(context.Background(),
		data.OrgId.ValueString()).UpdateOrganization(*updateOrganization).Execute()
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

	// save inlineResp data into Terraform state
	data.Id = types.StringValue("example-id")
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.CloudRegionName = types.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() == true {
			data.ManagementDetailsName = types.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = types.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() == true {
			data.ManagementDetailsValue = types.StringValue(managementDetailValue)
		} else {
			data.ManagementDetailsValue = types.StringNull()
		}
	} else {
		data.ManagementDetailsName = types.StringNull()
		data.ManagementDetailsValue = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	httpResp, err := r.client.OrganizationsApi.DeleteOrganization(context.Background(), data.OrgId.ValueString()).Execute()
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
	if httpResp.StatusCode != 204 {
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

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), req.ID)...)
}
