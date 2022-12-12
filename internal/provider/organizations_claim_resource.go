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

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsClaimResource{}
var _ resource.ResourceWithImportState = &OrganizationsClaimResource{}

func NewOrganizationsClaimResource() resource.Resource {
	return &OrganizationsClaimResource{}
}

// OrganizationsClaimResource defines the resource implementation.
type OrganizationsClaimResource struct {
	client *apiclient.APIClient
}

// OrganizationsClaimResourceModel describes the resource data model.
type OrganizationsClaimResourceModel struct {
	// TODO - ALL resources and DATA SOURCES must have an "id" filed for the acceptance tests to run.
	// TODO - DO NOT DELETE OR USE for randomly generated postgres Ids such as "organization_id", "network_id", etc...
	Id types.String `tfsdk:"id"`

	// TODO - Check that all names are in SnakeCase
	// TODO - Check that each item is typed correctly and the names match the tfsdk.Schema attributes.
	// TODO - Tip: "types.Object" should always be modified.

}

func (r *OrganizationsClaimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_claim"
}

func (r *OrganizationsClaimResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "OrganizationsClaim resource - ",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},

			// TODO - Inspect each Attribute and ensure each matches the API endpoint.
			// TODO - typically Optional + Computed is set to allow the user and the API to modify these values.
			// TODO - Required cannot be true while optional or computed is true
			// TODO - Even if an attribute is required to even make the API call do not mark it as required.
			// TODO - Often not all of the API calls use it. For now we add separate logic in each call for validation.
			// TODO - Check the attributes MarkdownDescription.
			// TODO - If the value is a password or token set the Sensitive flag to true.

		},
	}, nil
}

func (r *OrganizationsClaimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsClaimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	/*
			// TODO - Check for required parameters
		    	if len(data.OrgId.ValueString()) < 1 {
		    		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		    		return
		    	}
	*/

	if resp.Diagnostics.HasError() {
		return
	}

	/*
			// TODO - Create and Validate Payload
		    createOrganizationAdmin := *apiclient.NewInlineObject176(
		        data.Email.ValueString(),
		        data.Name.ValueString(),
		        data.OrgAccess.ValueString())

		    	// Tags
		    	if len(data.Tags) < 0 {
		    		var tags []apiclient.OrganizationsOrganizationIdAdminsTags
		    		for _, attribute := range data.Tags {
		    			var tag apiclient.OrganizationsOrganizationIdAdminsTags
		    			tag.Tag = attribute.Tag.ValueString()
		    			tag.Access = attribute.Tag.ValueString()
		    			tags = append(tags, tag)
		    		}
		    		createOrganizationAdmin.SetTags(tags)
		    	}

		    	// Networks
		    	if len(data.Networks) < 0 {
		    		var networks []apiclient.OrganizationsOrganizationIdAdminsNetworks
		    		for _, attribute := range data.Networks {
		    			var network apiclient.OrganizationsOrganizationIdAdminsNetworks
		    			network.Id = attribute.Id.ValueString()
		    			network.Access = attribute.Access.ValueString()
		    			networks = append(networks, network)
		    		}
		    		createOrganizationAdmin.SetNetworks(networks)
		    	}

		    	if data.AuthenticationMethod.IsNull() != true {
		    		createOrganizationAdmin.SetAuthenticationMethod(data.AuthenticationMethod.ValueString())
		    	}
	*/

	/*
			// TODO - Check the API client /docs for examples of each API call.
			inlineResp, httpResp, err := r.client.AdminsApi.CreateOrganizationAdmin(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdmin(createOrganizationAdmin).Execute()
		    	if err != nil {
		    		resp.Diagnostics.AddError(
		    			"Failed to create resource",
		    			fmt.Sprintf("%v\n", err.Error()),
		    		)
		    	}
	*/

	/*
	   // TODO - Check Postman or the API Spec to ensure the response code matches
	   // Check for API success response code
	   if httpResp.StatusCode != 201 {
	       resp.Diagnostics.AddError(
	           "Unexpected HTTP Response Status Code",
	           fmt.Sprintf("%v", httpResp.StatusCode),
	       )
	   	}
	*/

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	/*
				// TODO - If the API call returns a typed response access like below. Otherwise use existing resources for reference in handling untyped responses.
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
	*/

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsClaimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	/*
			// TODO - Check for required parameters
		    if len(data.OrgId.ValueString()) < 1 {
		        resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		        return
		    }


			if len(data.AdminId.ValueString()) < 1 {
				resp.Diagnostics.AddError("Missing AdminId", fmt.Sprintf("Value: %s", data.AdminId.ValueString()))
				return
			}
	*/

	if resp.Diagnostics.HasError() {
		return
	}

	/*
			TODO -
			inlineResp, httpResp, err := r.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.OrgId.ValueString()).Execute()
		    	if err != nil {
		    		resp.Diagnostics.AddError(
		    			"Failed to read resource",
		    			fmt.Sprintf("%v\n", err.Error()),
		    		)
		    	}
	*/

	/*
			// TODO - Check for API success inlineResp code
		    if httpResp.StatusCode != 200 {
		        resp.Diagnostics.AddError(
		            "Unexpected HTTP Response Status Code",
		            fmt.Sprintf("%v", httpResp.StatusCode),
		        )
		    }
	*/

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO - Save data into Terraform state
	/*
			data.Id = types.StringValue("example-id")
		    data.OrgId = types.StringValue(inlineResp.GetId())
		    data.Name = types.StringValue(inlineResp.GetName())
		    data.CloudRegion = types.StringValue(inlineResp.Cloud.Region.GetName())
		    data.Url = types.StringValue(inlineResp.GetUrl())
		    data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
		    data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())
	*/

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsClaimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// TODO - For update it is often required to look at both plan and state data to find dynamically generated postgres Ids.
	var data *OrganizationsClaimResourceModel
	var stateData *OrganizationsClaimResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	/*
			// TODO - Check for required parameters
		    if len(data.OrgId.ValueString()) < 1 {
		        resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		        return
		    }


			// Check state for required attribute
			if len(data.AdminId.ValueString()) < 1 {
				data.AdminId = stateData.AdminId
			}

			if len(data.AdminId.ValueString()) < 1 {
				resp.Diagnostics.AddError("Missing AdminId", fmt.Sprintf("AdminId: %s", data.AdminId.ValueString()))
				return
			}
	*/

	if resp.Diagnostics.HasError() {
		return
	}

	/*
			// TODO - Create Payload
		    updateOrganizationAdmin := *apiclient.NewInlineObject177()
		    updateOrganizationAdmin.SetName(data.Name.ValueString())
		    updateOrganizationAdmin.SetOrgAccess(data.OrgAccess.ValueString())

		    // Tags
		    if len(data.Tags) < 0 {
		        var tags []apiclient.OrganizationsOrganizationIdAdminsTags
		        for _, attribute := range data.Tags {
		            var tag apiclient.OrganizationsOrganizationIdAdminsTags
		            tag.Tag = attribute.Tag.ValueString()
		            tag.Access = attribute.Tag.ValueString()
		            tags = append(tags, tag)
		        }
		        updateOrganizationAdmin.SetTags(tags)
		    }

		    // Networks
		    if len(data.Networks) < 0 {
		        var networks []apiclient.OrganizationsOrganizationIdAdminsNetworks
		        for _, attribute := range data.Networks {
		            var network apiclient.OrganizationsOrganizationIdAdminsNetworks
		            network.Id = attribute.Id.ValueString()
		            network.Access = attribute.Access.ValueString()
		            networks = append(networks, network)
		        }
		        updateOrganizationAdmin.SetNetworks(networks)
		    }
	*/

	/*
			// TODO - Create API Call
			inlineResp, httpResp, err := r.client.AdminsApi.UpdateOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).UpdateOrganizationAdmin(updateOrganizationAdmin).Execute()
		    	if err != nil {
		    		resp.Diagnostics.AddError(
		    			"Failed to update resource",
		    			fmt.Sprintf("%v\n", err.Error()),
		    		)
		    	}
	*/

	/*
			// TODO - Check for API success response code
		    if httpResp.StatusCode != 200 {
		        resp.Diagnostics.AddError(
		            "Unexpected HTTP Response Status Code",
		            fmt.Sprintf("%v", httpResp.StatusCode),
		        )
		    }
	*/

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// TODO -  Save data into Terraform state
	/*
			data.Id = types.StringValue("example-id")
		    data.OrgId = types.StringValue(inlineResp.GetId())
		    data.Name = types.StringValue(inlineResp.GetName())
		    data.CloudRegion = types.StringValue(inlineResp.Cloud.Region.GetName())
		    data.Url = types.StringValue(inlineResp.GetUrl())
		    data.ApiEnabled = types.BoolValue(inlineResp.Api.GetEnabled())
		    data.LicensingModel = types.StringValue(inlineResp.Licensing.GetModel())
	*/

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsClaimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	/*
			// TODO - Check for required parameters
		    if len(data.OrgId.ValueString()) < 1 {
		        resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		        return
		    }

		    if len(data.AdminId.ValueString()) < 1 {
		        resp.Diagnostics.AddError("Missing AdminId", fmt.Sprintf("Value: %s", data.AdminId.ValueString()))
		        return
		    }
	*/

	if resp.Diagnostics.HasError() {
		return
	}

	/*
			TODO - DELETE API CALL
			httpResp, err := r.client.AdminsApi.DeleteOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).Execute()
		    	if err != nil {
		    		resp.Diagnostics.AddError(
		    			"Failed to delete resource",
		    			fmt.Sprintf("%v\n", err.Error()),
		    		)
		    	}
	*/

	/*
			// TODO - Check for API success response code
		    if httpResp.StatusCode != 204 {
		        resp.Diagnostics.AddError(
		            "Unexpected HTTP Response Status Code",
		            fmt.Sprintf("%v", httpResp.StatusCode),
		        )
		    }
	*/

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationsClaimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
