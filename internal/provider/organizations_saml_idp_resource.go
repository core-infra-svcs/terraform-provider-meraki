package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"

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
	_ resource.Resource                = &OrganizationsSamlIdpResource{}
	_ resource.ResourceWithConfigure   = &OrganizationsSamlIdpResource{}
	_ resource.ResourceWithImportState = &OrganizationsSamlIdpResource{}
)

func NewOrganizationsSamlIdpResource() resource.Resource {
	return &OrganizationsSamlIdpResource{}
}

// OrganizationsSamlIdpResource defines the resource implementation.
type OrganizationsSamlIdpResource struct {
	client *openApiClient.APIClient
}

// OrganizationsSamlIdpResourceModel describes the resource data model.
type OrganizationsSamlIdpResourceModel struct {
	Id                      types.String     `tfsdk:"id"`
	OrganizationId          jsontypes.String `tfsdk:"organization_id"`
	ConsumerUrl             jsontypes.String `tfsdk:"consumer_url"`
	IdpId                   jsontypes.String `tfsdk:"idp_id"`
	SloLogoutUrl            jsontypes.String `tfsdk:"slo_logout_url"`
	X509CertSha1Fingerprint jsontypes.String `tfsdk:"x_509_cert_sha1_fingerprint"`
}

func (r *OrganizationsSamlIdpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_idp"
}

func (r *OrganizationsSamlIdpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the SAML IdPs in your organization",

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
			"consumer_url": schema.StringAttribute{
				Description: "URL that is consuming SAML Identity Provider (IdP)",
				Optional:    true,
				Computed:    true,
				CustomType:  jsontypes.StringType,
			},
			"idp_id": schema.StringAttribute{
				MarkdownDescription: "ID associated with the SAML Identity Provider (IdP)",
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
			"slo_logout_url": schema.StringAttribute{
				MarkdownDescription: "Dashboard will redirect users to this URL when they sign out.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"x_509_cert_sha1_fingerprint": schema.StringAttribute{
				MarkdownDescription: "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
		},
	}
}

func (r *OrganizationsSamlIdpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsSamlIdpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganizationsSamlIdp := *openApiClient.NewInlineObject214(data.X509CertSha1Fingerprint.ValueString())
	createOrganizationsSamlIdp.SetSloLogoutUrl(data.SloLogoutUrl.ValueString())

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SamlApi.CreateOrganizationSamlIdp(context.Background(), data.OrganizationId.ValueString()).CreateOrganizationSamlIdp(createOrganizationsSamlIdp).Execute()
	//if err != nil {
	// BUG - HTTP Client is unable to unmarshal data into typed response []client.InlineResponse20095, returns empty
	//}
	if httpResp == nil {
		resp.Diagnostics.AddError(
			"Failed to get http response",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	// unmarshal http body into inlineResp object
	var inlineResp *openApiClient.InlineResponse20095
	body, _ := io.ReadAll(httpResp.Body)
	if err = json.Unmarshal(body, &inlineResp); err != nil {
		resp.Diagnostics.AddError(
			"Failed to unmarshal JSON into typed response",
			fmt.Sprintf("%v", err.Error()),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

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

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.IdpId = jsontypes.StringValue(inlineResp.GetIdpId())
	data.ConsumerUrl = jsontypes.StringValue(inlineResp.GetConsumerUrl())
	data.SloLogoutUrl = jsontypes.StringValue(inlineResp.GetSloLogoutUrl())
	data.X509CertSha1Fingerprint = jsontypes.StringValue(inlineResp.GetX509certSha1Fingerprint())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *OrganizationsSamlIdpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.GetOrganizationSamlIdp(context.Background(), data.OrganizationId.ValueString(), data.IdpId.ValueString()).Execute()
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
	}

	data.Id = types.StringValue("example-id")
	data.IdpId = jsontypes.StringValue(inlineResp.GetIdpId())
	data.ConsumerUrl = jsontypes.StringValue(inlineResp.GetConsumerUrl())
	data.SloLogoutUrl = jsontypes.StringValue(inlineResp.GetSloLogoutUrl())
	data.X509CertSha1Fingerprint = jsontypes.StringValue(inlineResp.GetX509certSha1Fingerprint())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsSamlIdpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Create HTTP request body
	updateOrganizationsSamlIdp := openApiClient.NewInlineObject215()
	updateOrganizationsSamlIdp.SetX509certSha1Fingerprint(data.X509CertSha1Fingerprint.ValueString())
	updateOrganizationsSamlIdp.SetSloLogoutUrl(data.SloLogoutUrl.ValueString())

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SamlApi.UpdateOrganizationSamlIdp(context.Background(),
		data.OrganizationId.ValueString(), data.IdpId.ValueString()).UpdateOrganizationSamlIdp(*updateOrganizationsSamlIdp).Execute()
	//if err != nil {
	// BUG - HTTP Client is unable to unmarshal data into typed response []client.InlineResponse20095, returns empty
	//}
	if httpResp == nil {
		resp.Diagnostics.AddError(
			"Failed to get http response",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	// unmarshal http body into inlineResp object
	var inlineResp *openApiClient.InlineResponse20095
	body, _ := io.ReadAll(httpResp.Body)
	if err = json.Unmarshal(body, &inlineResp); err != nil {
		resp.Diagnostics.AddError(
			"Failed to unmarshal JSON into typed response",
			fmt.Sprintf("%v", err.Error()),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

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

	data.Id = types.StringValue("example-id")
	data.IdpId = jsontypes.StringValue(inlineResp.GetIdpId())
	data.ConsumerUrl = jsontypes.StringValue(inlineResp.GetConsumerUrl())
	data.SloLogoutUrl = jsontypes.StringValue(inlineResp.GetSloLogoutUrl())
	data.X509CertSha1Fingerprint = jsontypes.StringValue(inlineResp.GetX509certSha1Fingerprint())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsSamlIdpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	httpResp, err := r.client.SamlApi.DeleteOrganizationSamlIdp(context.Background(), data.OrganizationId.ValueString(), data.IdpId.ValueString()).Execute()
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
	}

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationsSamlIdpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
