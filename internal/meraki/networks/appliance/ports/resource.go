package ports

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_ports"
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAppliancePortRequest()

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() && data.Vlan != jsontypes.Int64Value(0) {
		var vlan = int32(data.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	if !data.Type.IsUnknown() && !data.Type.IsNull() && data.Type != jsontypes.StringValue("") {
		payload.Type = data.Type.ValueStringPointer()
	}

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.Enabled = data.Enabled.ValueBoolPointer()
	}

	if !data.Accesspolicy.IsUnknown() && !data.Accesspolicy.IsNull() && data.Accesspolicy != jsontypes.StringValue("") {
		payload.AccessPolicy = data.Accesspolicy.ValueStringPointer()
	}
	if !data.Allowedvlans.IsUnknown() && !data.Allowedvlans.IsNull() && data.Allowedvlans != jsontypes.StringValue("") {
		payload.AllowedVlans = data.Allowedvlans.ValueStringPointer()
	}
	if !data.Dropuntaggedtraffic.IsUnknown() && !data.Dropuntaggedtraffic.IsNull() {
		payload.DropUntaggedTraffic = data.Dropuntaggedtraffic.ValueBoolPointer()
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).UpdateNetworkAppliancePortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, httpResp, err := r.client.ApplianceApi.GetNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAppliancePortRequest()

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() && data.Vlan != jsontypes.Int64Value(0) {
		var vlan = int32(data.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	if !data.Type.IsUnknown() && !data.Type.IsNull() && data.Type != jsontypes.StringValue("") {
		payload.Type = data.Type.ValueStringPointer()
	}

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.Enabled = data.Enabled.ValueBoolPointer()
	}

	if !data.Accesspolicy.IsUnknown() && !data.Accesspolicy.IsNull() && data.Accesspolicy != jsontypes.StringValue("") {
		payload.AccessPolicy = data.Accesspolicy.ValueStringPointer()
	}
	if !data.Allowedvlans.IsUnknown() && !data.Allowedvlans.IsNull() && data.Allowedvlans != jsontypes.StringValue("") {
		payload.AllowedVlans = data.Allowedvlans.ValueStringPointer()
	}
	if !data.Dropuntaggedtraffic.IsUnknown() && !data.Dropuntaggedtraffic.IsNull() {
		payload.DropUntaggedTraffic = data.Dropuntaggedtraffic.ValueBoolPointer()
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).UpdateNetworkAppliancePortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAppliancePortRequest()

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() && data.Vlan != jsontypes.Int64Value(0) {
		var vlan = int32(data.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	if !data.Type.IsUnknown() && !data.Type.IsNull() && data.Type != jsontypes.StringValue("") {
		payload.Type = data.Type.ValueStringPointer()
	}

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.Enabled = data.Enabled.ValueBoolPointer()
	}

	if !data.Accesspolicy.IsUnknown() && !data.Accesspolicy.IsNull() && data.Accesspolicy != jsontypes.StringValue("") {
		payload.AccessPolicy = data.Accesspolicy.ValueStringPointer()
	}
	if !data.Allowedvlans.IsUnknown() && !data.Allowedvlans.IsNull() && data.Allowedvlans != jsontypes.StringValue("") {
		payload.AllowedVlans = data.Allowedvlans.ValueStringPointer()
	}
	if !data.Dropuntaggedtraffic.IsUnknown() && !data.Dropuntaggedtraffic.IsNull() {
		payload.DropUntaggedTraffic = data.Dropuntaggedtraffic.ValueBoolPointer()
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).UpdateNetworkAppliancePortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}
