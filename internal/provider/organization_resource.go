package provider

import (
	"context"
	"fmt"
	"net/http/httputil"

	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationResource{}
var _ resource.ResourceWithImportState = &OrganizationResource{}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

// OrganizationResource defines the resource implementation.
type OrganizationResource struct {
	client *apiclient.APIClient
}

// OrganizationResourceModel describes the resource data model.
type OrganizationResourceModel struct {
	ApiEnabled             types.Bool   `tfsdk:"api_enabled"`
	CloudRegion            types.String `tfsdk:"cloud_region"`
	ManagementDetailsName  types.String `tfsdk:"management_details_name"`
	ManagementDetailsValue types.String `tfsdk:"management_details_value"`
	Id                     types.String `tfsdk:"id"`
	OrgId                  types.String `tfsdk:"organization_id"`
	LicensingModel         types.String `tfsdk:"licensing_model"`
	Name                   types.String `tfsdk:"name"`
	Url                    types.String `tfsdk:"url"`
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Organization resource - Manage the organizations that the user has privileges on",
		Attributes: map[string]tfsdk.Attribute{
			"api_enabled": {
				Description:         "Enable API access",
				MarkdownDescription: "API-specific settings",
				Type:                types.BoolType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"cloud_region": {
				Description:         "Region info",
				MarkdownDescription: "",
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
			"id": {
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"organization_id": {
				Description:         "Organization Id",
				MarkdownDescription: "The Id of the organization",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"licensing_model": {
				Description:         "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
				MarkdownDescription: "",
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
			"management_details_name": {
				Description:         "The name of the organization's management system",
				MarkdownDescription: "The name of the organization's management system",
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
			"management_details_value": {
				Description:         "Information about the organization's management system",
				MarkdownDescription: "Information about the organization's management system",
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
				Description:         "Organization name",
				MarkdownDescription: "The name of the organization",
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
				Description:         "Organization URL",
				MarkdownDescription: "",
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

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganization := *apiclient.NewInlineObject165(data.Name.ValueString())

	// Set management details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()
	detail := apiclient.OrganizationsManagementDetails{
		Name:  &name,
		Value: &value,
	}
	details := []apiclient.OrganizationsManagementDetails{detail}
	organizationsManagement := apiclient.OrganizationsManagement{Details: details}
	createOrganization.SetManagement(organizationsManagement)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganization(context.Background()).CreateOrganization(createOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect HTTP request diagnostics
	reqDump, err := httputil.DumpRequestOut(httpResp.Request, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP request diagnostics", fmt.Sprintf("\n%s", err),
		)
	}

	// Collect HTTP response diagnostics
	respDump, err := httputil.DumpResponse(httpResp, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP inlineResp diagnostics", fmt.Sprintf("\n%s", err),
		)
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

		resp.Diagnostics.AddError(
			"Request Diagnostics:",
			fmt.Sprintf("\n%s", string(reqDump)),
		)

		resp.Diagnostics.AddError(
			"Response Diagnostics:",
			fmt.Sprintf("\n%s", string(respDump)),
		)
		return
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.CloudRegion = types.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() {
			data.ManagementDetailsName = types.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = types.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() {
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

	// check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("org_id: %s", data.OrgId.ValueString()))
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganization(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect HTTP request diagnostics
	reqDump, err := httputil.DumpRequestOut(httpResp.Request, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP request diagnostics", fmt.Sprintf("\n%s", err),
		)
	}

	// Collect HTTP response diagnostics
	respDump, err := httputil.DumpResponse(httpResp, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP inlineResp diagnostics", fmt.Sprintf("\n%s", err),
		)
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

		resp.Diagnostics.AddError(
			"Request Diagnostics:",
			fmt.Sprintf("\n%s", string(reqDump)),
		)

		resp.Diagnostics.AddError(
			"Response Diagnostics:",
			fmt.Sprintf("\n%s", string(respDump)),
		)
		return
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.CloudRegion = types.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() {
			data.ManagementDetailsName = types.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = types.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() {
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

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("org_id: %s", data.OrgId.ValueString()))
		return
	}

	// Create HTTP request body
	updateOrganization := apiclient.NewInlineObject166()
	updateOrganization.SetName(data.Name.ValueString())

	// Set enabled attribute
	var enabled = data.ApiEnabled.ValueBool()
	Api := apiclient.OrganizationsOrganizationIdApi{Enabled: &enabled}
	updateOrganization.SetApi(Api)

	// Set management details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()
	detail := apiclient.OrganizationsManagementDetails{
		Name:  &name,
		Value: &value,
	}
	details := []apiclient.OrganizationsManagementDetails{detail}
	organizationsManagement := apiclient.OrganizationsManagement{Details: details}
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

	// Collect HTTP request diagnostics
	reqDump, err := httputil.DumpRequestOut(httpResp.Request, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP request diagnostics", fmt.Sprintf("\n%s", err),
		)
	}

	// Collect HTTP response diagnostics
	respDump, err := httputil.DumpResponse(httpResp, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP inlineResp diagnostics", fmt.Sprintf("\n%s", err),
		)
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

		resp.Diagnostics.AddError(
			"Request Diagnostics:",
			fmt.Sprintf("\n%s", string(reqDump)),
		)

		resp.Diagnostics.AddError(
			"Response Diagnostics:",
			fmt.Sprintf("\n%s", string(respDump)),
		)
		return
	}

	// save inlineResp data into Terraform state
	data.Id = types.StringValue("example-id")
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.Name = types.StringValue(inlineResp.GetName())
	data.CloudRegion = types.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = types.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() {
			data.ManagementDetailsName = types.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = types.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() {
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

	// check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("org_id: %s", data.OrgId.ValueString()))
		return
	}

	// Initialize provider client and make API call
	httpResp, err := r.client.OrganizationsApi.DeleteOrganization(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect HTTP request diagnostics
	reqDump, err := httputil.DumpRequestOut(httpResp.Request, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP request diagnostics", fmt.Sprintf("\n%s", err),
		)
	}

	// Collect HTTP response diagnostics
	respDump, err := httputil.DumpResponse(httpResp, true)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Failed to gather HTTP inlineResp diagnostics", fmt.Sprintf("\n%s", err),
		)
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

		resp.Diagnostics.AddError(
			"Request Diagnostics:",
			fmt.Sprintf("\n%s", string(reqDump)),
		)

		resp.Diagnostics.AddError(
			"Response Diagnostics:",
			fmt.Sprintf("\n%s", string(respDump)),
		)
		return
	}

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
