package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"

	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationSamlResource{}
	_ resource.ResourceWithConfigure   = &OrganizationSamlResource{}
	_ resource.ResourceWithImportState = &OrganizationSamlResource{}
)

func NewOrganizationSamlResource() resource.Resource {
	return &OrganizationSamlResource{}
}

// OrganizationSamlResource defines the resource implementation.
type OrganizationSamlResource struct {
	client *openApiClient.APIClient
}

// OrganizationSamlResourceModel describes the resource data model.
type OrganizationSamlResourceModel struct {
	Id             types.String     `tfsdk:"id"`
	OrganizationId jsontypes.String `tfsdk:"organization_id"`
	Enabled        jsontypes.Bool   `tfsdk:"enabled"`
}

func (r *OrganizationSamlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_saml"
}

func (r *OrganizationSamlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the SAML SSO enabled settings for an organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Toggle depicting if SAML SSO settings are enabled",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
		},
	}
}

func (r *OrganizationSamlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationSamlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	enableOrganizationSaml := *openApiClient.NewInlineObject213()
	enableOrganizationSaml.SetEnabled(data.Enabled.ValueBool())

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.UpdateOrganizationSaml(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSaml(enableOrganizationSaml).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
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
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.Enabled = jsontypes.BoolValue(inlineResp.GetEnabled())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *OrganizationSamlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.GetOrganizationSaml(context.Background(), data.OrganizationId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
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

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.Enabled = jsontypes.BoolValue(inlineResp.GetEnabled())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationSamlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationSamlResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Create HTTP request body
	enableOrganizationSaml := *openApiClient.NewInlineObject213()
	enableOrganizationSaml.SetEnabled(data.Enabled.ValueBool())

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.UpdateOrganizationSaml(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSaml(enableOrganizationSaml).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
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
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.Enabled = jsontypes.BoolValue(inlineResp.GetEnabled())

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
