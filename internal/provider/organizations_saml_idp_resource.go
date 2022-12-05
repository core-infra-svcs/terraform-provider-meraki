package provider

import (
	"context"
	"encoding/json"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsSamlIdpResource{}
var _ resource.ResourceWithImportState = &OrganizationsSamlIdpResource{}

func NewOrganizationsSamlIdpResource() resource.Resource {
	return &OrganizationsSamlIdpResource{}
}

// OrganizationsSamlIdpResource defines the resource implementation.
type OrganizationsSamlIdpResource struct {
	client *apiclient.APIClient
}

// OrganizationsSamlIdpResourceModel describes the resource data model.
type OrganizationsSamlIdpResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	OrganizationId          types.String `tfsdk:"organization_id"`
	ConsumerUrl             types.String `tfsdk:"consumer_url"`
	IdpId                   types.String `tfsdk:"idp_id"`
	SloLogoutUrl            types.String `tfsdk:"slo_logout_url"`
	X509CertSha1Fingerprint types.String `tfsdk:"x_509cert_sha1_fingerprint"`
}

func (r *OrganizationsSamlIdpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_idp"
}

func (r *OrganizationsSamlIdpResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// TODO - This description is used by the documentation generator and the language server.
		MarkdownDescription: "OrganizationsSamlIdp resource - ",
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
			"consumer_url": {
				Description:         "URL that is consuming SAML Identity Provider (IdP)",
				MarkdownDescription: "URL that is consuming SAML Identity Provider (IdP)",
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
			"idp_id": {
				Description:         "ID associated with the SAML Identity Provider (IdP)",
				MarkdownDescription: "ID associated with the SAML Identity Provider (IdP)",
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
			"slo_logout_url": {
				Description:         "Dashboard will redirect users to this URL when they sign out.",
				MarkdownDescription: "Dashboard will redirect users to this URL when they sign out.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"x_509cert_sha1_fingerprint": {
				Description:         "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
				MarkdownDescription: "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
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

func (r *OrganizationsSamlIdpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsSamlIdpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganizationsSamlIdp := *apiclient.NewInlineObject214(data.X509CertSha1Fingerprint.ValueString())
	createOrganizationsSamlIdp.SetSloLogoutUrl(data.SloLogoutUrl.ValueString())

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SamlApi.CreateOrganizationSamlIdp(context.Background(), data.OrganizationId.ValueString()).CreateOrganizationSamlIdp(createOrganizationsSamlIdp).Execute()
	if err != nil {
		// BUG - HTTP Client is unable to unmarshal data into typed response []client.InlineResponse20095, returns empty
	}

	// unmarshal http body into inlineResp object
	var inlineResp *apiclient.InlineResponse20095
	body, _ := io.ReadAll(httpResp.Body)
	if err := json.Unmarshal(body, &inlineResp); err != nil {
		fmt.Println(err)
		resp.Diagnostics.AddError(
			"Failed to unmarshal JSON into typed response",
			fmt.Sprintf("%v", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
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

	resp.Diagnostics.Append()

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.IdpId = types.StringValue(inlineResp.GetIdpId())
	data.ConsumerUrl = types.StringValue(inlineResp.GetConsumerUrl())
	data.SloLogoutUrl = types.StringValue(inlineResp.GetSloLogoutUrl())
	data.X509CertSha1Fingerprint = types.StringValue(inlineResp.GetX509certSha1Fingerprint())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *OrganizationsSamlIdpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("organization_id: %s", data.OrganizationId.ValueString()))
		return
	}

	// check for required parameters
	if len(data.IdpId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing saml idp id", fmt.Sprintf("idp_id: %s", data.IdpId.ValueString()))
		return
	}

	//ToDo:- Check if this is really needed as it may not be in the required list, but is needed for this api call
	//if len(data.IdpId.ValueString()) < 1 {
	//	resp.Diagnostics.AddError("Missing idp_Id", fmt.Sprintf("idp_id: %s", data.IdpId.ValueString()))
	//	return
	//}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.GetOrganizationSamlIdp(context.Background(), data.OrganizationId.ValueString(), data.IdpId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success inlineResp code
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

	data.Id = types.StringValue("example-id")
	data.IdpId = types.StringValue(inlineResp.GetIdpId())
	data.ConsumerUrl = types.StringValue(inlineResp.GetConsumerUrl())
	data.SloLogoutUrl = types.StringValue(inlineResp.GetSloLogoutUrl())
	data.X509CertSha1Fingerprint = types.StringValue(inlineResp.GetX509certSha1Fingerprint())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsSamlIdpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("organization_id: %s", data.OrganizationId.ValueString()))
		return
	}

	// Create HTTP request body
	updateOrganizationsSamlIdp := apiclient.NewInlineObject215()
	updateOrganizationsSamlIdp.SetX509certSha1Fingerprint(data.X509CertSha1Fingerprint.ValueString())
	updateOrganizationsSamlIdp.SetSloLogoutUrl(data.SloLogoutUrl.ValueString())

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.SamlApi.UpdateOrganizationSamlIdp(context.Background(),
		data.OrganizationId.ValueString(), data.IdpId.ValueString()).UpdateOrganizationSamlIdp(*updateOrganizationsSamlIdp).Execute()
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

	data.Id = types.StringValue("example-id")
	data.IdpId = types.StringValue(inlineResp[0].GetIdpId())
	data.ConsumerUrl = types.StringValue(inlineResp[0].GetConsumerUrl())
	data.SloLogoutUrl = types.StringValue(inlineResp[0].GetSloLogoutUrl())
	data.X509CertSha1Fingerprint = types.StringValue(inlineResp[0].GetX509certSha1Fingerprint())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsSamlIdpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsSamlIdpResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// check for required parameters
	if len(data.OrganizationId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organization_Id", fmt.Sprintf("organization_id: %s", data.OrganizationId.ValueString()))
		return
	}

	// Initialize provider client and make API call
	httpResp, err := r.client.SamlApi.DeleteOrganizationSamlIdp(context.Background(), data.OrganizationId.ValueString(), data.IdpId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
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

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationsSamlIdpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
