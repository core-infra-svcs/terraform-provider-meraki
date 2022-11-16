package provider

import (
	"context"
    "encoding/json"
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
var _ resource.Resource = &OrganizationsSamlIdpsResource{}
var _ resource.ResourceWithImportState = &OrganizationsSamlIdpsResource{}

func NewOrganizationsSamlIdpsResource() resource.Resource {
	return &OrganizationsSamlIdpsResource{}
}

// OrganizationsSamlIdpsResource defines the resource implementation.
type OrganizationsSamlIdpsResource struct {
	client *apiclient.APIClient
}

// OrganizationsSamlIdpsResourceModel describes the resource data model.
type OrganizationsSamlIdpsResourceModel struct {
        Id  					types.String `tfsdk:"id"`
		Consumerurl             types.String `tfsdk:"consumer_url"`
		Idpid                   types.String `tfsdk:"idp_id"`
		OrganizationId          types.String `tfsdk:"organization_id"`
		Slologouturl          	types.String `tfsdk:"slo_logout_url"`
		X509certsha1fingerprint types.String `tfsdk:"x_509cert_sha1_fingerprint"`
}

func (r *OrganizationsSamlIdpsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_idps"
}

func (r *OrganizationsSamlIdpsResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OrganizationsSamlIdps resource - Get a SAML IdP from your organization",

		Attributes: map[string]tfsdk.Attribute{
            // TODO - Inspect each Attribute and ensure each matches the API endpoint.
            // TODO - Check the attributes MarkdownDescription.
            // TODO - Ensure either Required OR Optional + Computed is set.
            // TODO - If the value is a password or token set the Sensitive flag to true.

			"id": {
				Description:    "Example identifier needed for terraform",
				MarkdownDescription:    "Example identifier needed for terraform",
				Type:     types.StringType,
				Required:   false,
				Optional:   false,
				Computed: true,
				Sensitive:  false,
				Attributes: nil,
				DeprecationMessage: "",
				Validators: nil,
				PlanModifiers:  nil,
			},
			"consumer_url": {
				Description:    "URL that is consuming SAML Identity Provider (IdP)",
				MarkdownDescription:    "URL that is consuming SAML Identity Provider (IdP)",
				Type:     types.StringType,
				Required:   false,
				Optional:   true,
				Computed: false,
				Sensitive:  false,
				Attributes: nil,
				DeprecationMessage: "",
				Validators: nil,
				PlanModifiers:  nil,
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
			"idp_id": {
				Description:    "ID associated with the SAML Identity Provider (IdP)",
				MarkdownDescription:    "ID associated with the SAML Identity Provider (IdP)",
				Type:     types.StringType,
				Required:   true,
				Optional:   false,
				Computed: false,
				Sensitive:  false,
				Attributes: nil,
				DeprecationMessage: "",
				Validators: nil,
				PlanModifiers:  nil,
			},
			"slo_logout_url": {
				Description:    "Dashboard will redirect users to this URL when they sign out.",
				MarkdownDescription:    "Dashboard will redirect users to this URL when they sign out.",
				Type:     types.StringType,
				Required:   false,
				Optional:   true,
				Computed: false,
				Sensitive:  false,
				Attributes: nil,
				DeprecationMessage: "",
				Validators: nil,
				PlanModifiers:  nil,
			},
			"x_509cert_sha1_fingerprint": {
				Description:    "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
				MarkdownDescription:    "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
				Type:     types.StringType,
				Required:   false,
				Optional:   true,
				Computed: false,
				Sensitive:  false,
				Attributes: nil,
				DeprecationMessage: "",
				Validators: nil,
				PlanModifiers:  nil,
			},
        },
    }, nil
}

func (r *OrganizationsSamlIdpsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsSamlIdpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsSamlIdpsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO - POST/PUT call to API endpoint using HTTP Client.

	createOrganizationsSamlIdp := *apiclient.NewInlineObject166(data.OrganizationId.ValueString())
	inlineResp, httpResp, err := r.client.IdpsApi.CreateOrganizationSamlIdp(context.Background()).CreateOrganizationSamlIdp(createOrganizationsSamlIdp).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"-- Create Error --",
				fmt.Sprintf("%v\n", err.Error()),
			)
			resp.Diagnostics.AddError(
				"-- Response --",
				fmt.Sprintf("%v\n", httpResp),
			)
			return
		}

    // TODO - if missing a strongly typed response object in HTTP client:
        // 1) cast response data into an empty interface
        // 2) unmarshal into custom struct.
        /*
        responseData, _ := json.Marshal(response)
        var results apiclient.InlineResponse20064
        json.Unmarshal(responseData, &results)
        */

    // TODO - save response to Terraform state.

	data.Id = types.StringValue("example-id")
	data.OrganizationId = types.StringValue(inlineResp.GetId) // check if it is GetOrganizationId in inlineresp
	data.Consumerurl = types.StringValue(inlineResp.GetConsumerurl)
	data.Idpid = types.StringValue(inlineResp.GetIdpid)
	data.Slologouturl = types.StringValue(inlineResp.GetSlologouturl)
	data.X509certsha1fingerprint = types.StringValue(inlineResp.zgetX509certsha1fingerprint)


	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsSamlIdpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *Organizations{Organizationid}SamlIdps{Idpid}ResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO - GET call to API endpoint using HTTP Client.
	    /*
	    response, d, err := r.client.OrganizationsApi.GetOrganization(context.Background(), data.Id.Value).Execute()
        	if err != nil {
        		resp.Diagnostics.AddError(
        			"-- Read Error --",
        			fmt.Sprintf("%v\n", err.Error()),
        		)
        		resp.Diagnostics.AddError(
        			"-- Response --",
        			fmt.Sprintf("%v\n", d),
        		)
        		return
        	}
	    */

	// TODO - save response to Terraform state.
            /*
            data.Id = types.String{Value: response.GetId()}
            data.Name = types.String{Value: response.GetName()}
            */

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsSamlIdpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *Organizations{Organizationid}SamlIdps{Idpid}ResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO - create payload
        /*
        payload := apiclient.NewInlineObject167()
        payload.SetName(data.Name.Value)
        */

	// TODO - PUT/POST call to API endpoint using HTTP Client.
    	    /*
    	    response, d, err := r.provider.client.OrganizationsApi.UpdateOrganization(context.Background(),
            		data.Id.Value).UpdateOrganization(*payload).Execute()
            	if err != nil {
            		resp.Diagnostics.AddError(
            			"-- Update Error --",
            			fmt.Sprintf("%v\n", err.Error()),
            		)
            		resp.Diagnostics.AddError(
            			"-- Response --",
            			fmt.Sprintf("%v\n", d),
            		)
            		return
            	}
    	    */

    	// TODO - save response to Terraform state.
               /*
                data.Id = types.String{Value: response.GetId()}
                data.Name = types.String{Value: response.GetName()}
                */

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsSamlIdpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *Organizations{Organizationid}SamlIdps{Idpid}ResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO - Delete call to API endpoint using HTTP Client.
        /*
        response, err := r.client.OrganizationsApi.DeleteOrganization(context.Background(), data.Id.Value).Execute()
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
        */

    resp.State.RemoveResource(ctx)

}

func (r *OrganizationsSamlIdpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
