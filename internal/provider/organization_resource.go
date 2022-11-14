package provider

import (
	"context"
	"fmt"
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

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createOrganization := *apiclient.NewInlineObject165(data.Name.ValueString())

	// Set Management Details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()

	detail := apiclient.OrganizationsManagementDetails{
		Name:  &name,
		Value: &value,
	}
	details := []apiclient.OrganizationsManagementDetails{detail}
	organizationsManagement := apiclient.OrganizationsManagement{Details: details}
	createOrganization.SetManagement(organizationsManagement)

	response, d, err := r.client.OrganizationsApi.CreateOrganization(context.Background()).CreateOrganization(createOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Create Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response Header --",
			fmt.Sprintf("%v\n", d.Header),
		)
		return
	}

	// save into the Terraform state.
	data.OrgId = types.StringValue(response.GetId())
	data.Name = types.StringValue(response.GetName())
	data.CloudRegion = types.StringValue(response.Cloud.Region.GetName())
	data.Url = types.StringValue(response.GetUrl())
	data.ApiEnabled = types.BoolValue(response.Api.GetEnabled())
	data.LicensingModel = types.StringValue(response.Licensing.GetModel())

	// Management Details Response
	if len(response.Management.Details) > 0 {

		responseDetails := response.Management.GetDetails()

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

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, d, err := r.client.OrganizationsApi.GetOrganization(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Read Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response Header --",
			fmt.Sprintf("%v\n", d.Header),
		)
		return
	}

	// save into the Terraform state.
	data.OrgId = types.StringValue(response.GetId())
	data.Name = types.StringValue(response.GetName())
	data.CloudRegion = types.StringValue(response.Cloud.Region.GetName())
	data.Url = types.StringValue(response.GetUrl())
	data.ApiEnabled = types.BoolValue(response.Api.GetEnabled())
	data.LicensingModel = types.StringValue(response.Licensing.GetModel())

	// Management Details Response
	if len(response.Management.Details) > 0 {

		responseDetails := response.Management.GetDetails()

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
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId", fmt.Sprintf("%s", data.OrgId.ValueString()))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// create payload
	updateOrganization := apiclient.NewInlineObject166()
	updateOrganization.SetName(data.Name.ValueString())

	// API Enabled
	var enabled = data.ApiEnabled.ValueBool()
	Api := apiclient.OrganizationsOrganizationIdApi{Enabled: &enabled}
	updateOrganization.SetApi(Api)

	// Set Management Details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()
	detail := apiclient.OrganizationsManagementDetails{
		Name:  &name,
		Value: &value,
	}
	details := []apiclient.OrganizationsManagementDetails{detail}
	organizationsManagement := apiclient.OrganizationsManagement{Details: details}
	updateOrganization.SetManagement(organizationsManagement)

	// Initialize provider client
	response, d, err := r.client.OrganizationsApi.UpdateOrganization(context.Background(),
		data.OrgId.ValueString()).UpdateOrganization(*updateOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Update Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response Header --",
			fmt.Sprintf("%v\n", d.Header),
		)
		return
	}

	// save into the Terraform state.
	data.OrgId = types.StringValue(response.GetId())
	data.Name = types.StringValue(response.GetName())
	data.CloudRegion = types.StringValue(response.Cloud.Region.GetName())
	data.Url = types.StringValue(response.GetUrl())
	data.ApiEnabled = types.BoolValue(response.Api.GetEnabled())
	data.LicensingModel = types.StringValue(response.Licensing.GetModel())

	// Management Details Response
	if len(response.Management.Details) > 0 {

		responseDetails := response.Management.GetDetails()

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
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Organization
	response, err := r.client.OrganizationsApi.DeleteOrganization(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Delete Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", response),
		)
		return
	}

	if response.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"-- Delete Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", response),
		)
		return
	} else {
		resp.State.RemoveResource(ctx)
	}

}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
