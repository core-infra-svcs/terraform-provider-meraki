package provider

import (
	"context"
	"encoding/json"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-golang/client"

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
	Api_enabled     types.Bool   `tfsdk:"api_enabled"`
	Cloud_region    types.String `tfsdk:"cloud_region"`
	Id              types.String `tfsdk:"id"`
	Licensing_model types.String `tfsdk:"licensing_model"`
	Name            types.String `tfsdk:"name"`
	Url             types.String `tfsdk:"url"`
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Organization resource - Manage the organizations that the user has privileges on",

		Attributes: map[string]tfsdk.Attribute{

			// TODO - As a developer I must manually inspect the required/optional status of each attribute
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
				Description:         "Organization ID",
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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	createOrganization := *apiclient.NewInlineObject166(data.Name.Value)
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
		resp.Diagnostics.AddError(
			"-- Request URL --",
			fmt.Sprintf("%v\n", d.Request.URL),
		)
		return
	}

	// extract map string interface data into struct
	responseData, _ := json.Marshal(response)
	var results apiclient.InlineResponse20064
	json.Unmarshal(responseData, &results)

	// save into the Terraform state.
	data.Id = types.String{Value: results.GetId()}
	data.Name = types.String{Value: results.GetName()}
	data.Cloud_region = types.String{Value: results.Cloud.Region.GetName()}
	data.Url = types.String{Value: results.GetUrl()}
	data.Api_enabled = types.Bool{Value: results.Api.GetEnabled()}
	data.Licensing_model = types.String{Value: results.Licensing.GetModel()}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	response, d, err := r.client.OrganizationsApi.GetOrganization(context.Background(), data.Id.Value).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Read Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response Header --",
			fmt.Sprintf("%v\n", d.Header),
		)
		resp.Diagnostics.AddError(
			"-- Request URL --",
			fmt.Sprintf("%v\n", d.Request.URL),
		)
		return
	}

	data.Id = types.String{Value: response.GetId()}
	data.Name = types.String{Value: response.GetName()}
	data.Cloud_region = types.String{Value: response.Cloud.Region.GetName()}
	data.Url = types.String{Value: response.GetUrl()}
	data.Api_enabled = types.Bool{Value: response.Api.GetEnabled()}
	data.Licensing_model = types.String{Value: response.Licensing.GetModel()}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.

	// create payload
	updateOrganization := apiclient.NewInlineObject167()
	updateOrganization.SetName(data.Name.Value)

	// nested params are strongly typed
	var enabled apiclient.OrganizationsOrganizationIdApi
	enabled.SetEnabled(data.Api_enabled.Value)
	updateOrganization.SetApi(enabled)

	// Initialize provider client
	response, d, err := r.client.OrganizationsApi.UpdateOrganization(context.Background(),
		data.Id.Value).UpdateOrganization(*updateOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Update Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response Header --",
			fmt.Sprintf("%v\n", d.Header),
		)
		resp.Diagnostics.AddError(
			"-- Request URL --",
			fmt.Sprintf("%v\n", d.Request.URL),
		)
		return
	}

	// extract map string interface data into struct
	responseData, _ := json.Marshal(response)
	var results apiclient.InlineResponse20064
	json.Unmarshal(responseData, &results)

	// save into the Terraform state.
	data.Id = types.String{Value: results.GetId()}
	data.Name = types.String{Value: results.GetName()}
	data.Cloud_region = types.String{Value: results.Cloud.Region.GetName()}
	data.Url = types.String{Value: results.GetUrl()}
	data.Api_enabled = types.Bool{Value: results.Api.GetEnabled()}
	data.Licensing_model = types.String{Value: results.Licensing.GetModel()}

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// Delete Organization
	response, err := r.client.OrganizationsApi.DeleteOrganization(context.Background(), data.Id.Value).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Delete Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response Header --",
			fmt.Sprintf("%v\n", response.Header),
		)
		resp.Diagnostics.AddError(
			"-- Request URL --",
			fmt.Sprintf("%v\n", response.Request.URL),
		)
		return
	}

}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
