package organizations

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsSnmpResource{}
var _ resource.ResourceWithImportState = &OrganizationsSnmpResource{}

func NewOrganizationsSnmpResource() resource.Resource {
	return &OrganizationsSnmpResource{}
}

// OrganizationsSnmpResource defines the resource implementation.
type OrganizationsSnmpResource struct {
	client *openApiClient.APIClient
}

// OrganizationsSnmpResourceModel describes the resource data model.
type OrganizationsSnmpResourceModel struct {
	Id             jsontypes.String   `tfsdk:"id"`
	OrganizationId jsontypes.String   `tfsdk:"organization_id" json:"organizationId"`
	V2cEnabled     jsontypes.Bool     `tfsdk:"v2c_enabled" json:"v2cEnabled,omitempty"`
	V3Enabled      jsontypes.Bool     `tfsdk:"v3_enabled" json:"v3Enabled,omitempty"`
	V3AuthMode     jsontypes.String   `tfsdk:"v3_auth_mode" json:"v3AuthMode,omitempty"`
	V3AuthPass     jsontypes.String   `tfsdk:"v3_auth_pass" json:"v3AuthPass,omitempty"`
	V3PrivMode     jsontypes.String   `tfsdk:"v3_priv_mode" json:"v3PrivMode,omitempty"`
	V3PrivPass     jsontypes.String   `tfsdk:"v3_priv_pass" json:"v3PrivPass,omitempty"`
	PeerIps        []jsontypes.String `tfsdk:"peer_ips" json:"peerIps,omitempty"`
}

func (r *OrganizationsSnmpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_snmp"
}

func (r *OrganizationsSnmpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				Description: "The ID of the organization",
				Required:    true,
				CustomType:  jsontypes.StringType,
			},
			"v2c_enabled": schema.BoolAttribute{
				Description: "Boolean indicating whether SNMP version 2c is enabled for the organization.",
				Required:    true,
				CustomType:  jsontypes.BoolType,
			},
			"v3_enabled": schema.BoolAttribute{
				Description: "Boolean indicating whether SNMP version 3 is enabled for the organization.",
				Required:    true,
				CustomType:  jsontypes.BoolType,
			},
			"v3_auth_mode": schema.StringAttribute{
				Description: "The SNMP version 3 authentication mode. Can be either 'MD5' or 'SHA'.",
				Optional:    true,
				Computed:    true,
				CustomType:  jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("MD5", "SHA"),
				},
			},
			"v3_auth_pass": schema.StringAttribute{
				Description: "The SNMP version 3 authentication password.",
				Optional:    true,
				Sensitive:   true,
				CustomType:  jsontypes.StringType,
			},
			"v3_priv_mode": schema.StringAttribute{
				Description: "The SNMP version 3 privacy mode. Can be either 'DES' or 'AES128'.",
				Optional:    true,
				Computed:    true,
				CustomType:  jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("DES", "AES128"),
				},
			},
			"v3_priv_pass": schema.StringAttribute{
				Description: "The SNMP version 3 privacy password.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				CustomType:  jsontypes.StringType,
			},
			"peer_ips": schema.SetAttribute{
				Description: "The list of IPv4 addresses that are allowed to access the SNMP server.",
				ElementType: jsontypes.StringType,
				Optional:    true,
				Computed:    true,
				CustomType:  jsontypes.SetType[jsontypes.String](),
			},
		},
	}
}

func (r *OrganizationsSnmpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrganizationsSnmpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationsSnmpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateSnmpRequest := openApiClient.NewUpdateOrganizationSnmpRequest()
	updateSnmpRequest.V2cEnabled = data.V2cEnabled.ValueBoolPointer()
	updateSnmpRequest.V3Enabled = data.V3Enabled.ValueBoolPointer()
	updateSnmpRequest.V3AuthMode = data.V3AuthMode.ValueStringPointer()
	updateSnmpRequest.V3AuthPass = data.V3AuthPass.ValueStringPointer()
	updateSnmpRequest.V3PrivMode = data.V3PrivMode.ValueStringPointer()
	updateSnmpRequest.V3PrivPass = data.V3PrivPass.ValueStringPointer()

	for _, peer := range data.PeerIps {
		updateSnmpRequest.PeerIps = append(updateSnmpRequest.PeerIps, peer.ValueString())
	}

	_, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationSnmp(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSnmpRequest(*updateSnmpRequest).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error creating group policy",
			fmt.Sprintf("Could not create group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = data.OrganizationId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsSnmpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationsSnmpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.OrganizationsApi.GetOrganizationSnmp(context.Background(), data.OrganizationId.ValueString()).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error reading group policy",
			fmt.Sprintf("Could not read group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsSnmpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationsSnmpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateSnmpRequest := openApiClient.NewUpdateOrganizationSnmpRequest()
	updateSnmpRequest.V2cEnabled = data.V2cEnabled.ValueBoolPointer()
	updateSnmpRequest.V3Enabled = data.V3Enabled.ValueBoolPointer()
	updateSnmpRequest.V3AuthMode = data.V3AuthMode.ValueStringPointer()
	updateSnmpRequest.V3AuthPass = data.V3AuthPass.ValueStringPointer()
	updateSnmpRequest.V3PrivMode = data.V3PrivMode.ValueStringPointer()
	updateSnmpRequest.V3PrivPass = data.V3PrivPass.ValueStringPointer()

	for _, peer := range data.PeerIps {
		updateSnmpRequest.PeerIps = append(updateSnmpRequest.PeerIps, peer.ValueString())
	}

	_, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationSnmp(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSnmpRequest(*updateSnmpRequest).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error updating group policy",
			fmt.Sprintf("Could not update group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = data.OrganizationId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsSnmpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationsSnmpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	defaultState := false

	updateSnmpRequest := openApiClient.NewUpdateOrganizationSnmpRequest()
	updateSnmpRequest.V2cEnabled = &defaultState
	updateSnmpRequest.V3Enabled = &defaultState

	_, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationSnmp(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationSnmpRequest(*updateSnmpRequest).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error deleting group policy",
			fmt.Sprintf("Could not delete group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Trace(ctx, "removed resource")
}

func (r *OrganizationsSnmpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
