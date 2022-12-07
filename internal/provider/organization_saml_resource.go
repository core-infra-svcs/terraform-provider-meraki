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
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationSamlResource{}
var _ resource.ResourceWithImportState = &OrganizationSamlResource{}

func NewOrganizationSamlResource() resource.Resource {
	return &OrganizationSamlResource{}
}

// OrganizationSamlResource defines the resource implementation.
type OrganizationSamlResource struct {
	client *apiclient.APIClient
}

// OrganizationSamlResourceModel describes the resource data model.
type OrganizationSamlResourceModel struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (r *OrganizationSamlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_saml"
}

func (r *OrganizationSamlResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "OrganizationSaml resource - Enable or disable a SAML in your organization",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				Required:            false,
				Optional:            false,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"organization_id": {
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Type:                types.StringType,
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"enabled": {
				Description:         "Boolean for updating SAML SSO enabled settings.",
				MarkdownDescription: "Boolean for updating SAML SSO enabled settings.",
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
		},
	}, nil
}

func (r *OrganizationSamlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationSamlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization Id", fmt.Sprintf("Id: %s", data.OrganizationId.ValueString()))
		return
	}

	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing Enabled Status", fmt.Sprintf("Enabled: %v", data.Enabled.ValueBool()))
		return
	}

	// Create HTTP request body
	enableOrganizationSaml := *apiclient.NewInlineObject213()
	enableOrganizationSaml.SetEnabled(data.Enabled.ValueBool())

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.UpdateOrganizationSaml(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSaml(enableOrganizationSaml).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
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

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.Enabled = types.BoolValue(inlineResp.GetEnabled())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *OrganizationSamlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization Id", fmt.Sprintf("Id: %s", data.OrganizationId.ValueString()))
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.GetOrganizationSaml(context.Background(), data.OrganizationId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
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

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.Enabled = types.BoolValue(inlineResp.GetEnabled())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationSamlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization Id", fmt.Sprintf("Id: %s", data.OrganizationId.ValueString()))
		return
	}

	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing Enabled Status", fmt.Sprintf("Enabled: %v", data.Enabled.ValueBool()))
		return
	}

	// Create HTTP request body
	enableOrganizationSaml := *apiclient.NewInlineObject213()
	enableOrganizationSaml.SetEnabled(data.Enabled.ValueBool())

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.UpdateOrganizationSaml(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSaml(enableOrganizationSaml).Execute()
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

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.Enabled = types.BoolValue(inlineResp.GetEnabled())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationSamlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationSamlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
