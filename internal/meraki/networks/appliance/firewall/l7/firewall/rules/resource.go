package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksApplianceFirewallL7FirewallRulesResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceFirewallL7FirewallRulesResource{}

func NewResource() resource.Resource {
	return &NetworksApplianceFirewallL7FirewallRulesResource{}
}

// NetworksApplianceFirewallL7FirewallRulesResource defines the resource implementation.
type NetworksApplianceFirewallL7FirewallRulesResource struct {
	client *openApiClient.APIClient
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_firewall_l7_firewall_rules"
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkApplianceFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner

	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}

	updateNetworkApplianceFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallL7FirewallRulesRequest(updateNetworkApplianceFirewallL7FirewallRules).Execute()
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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).Execute()
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
		return
	} else {
		resp.Diagnostics.Append()
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkApplianceFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner

	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}

	updateNetworkApplianceFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallL7FirewallRulesRequest(updateNetworkApplianceFirewallL7FirewallRules).Execute()
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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rules := []openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner{}
	updateNetworkApplianceFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkApplianceFirewallL7FirewallRulesRequest()
	updateNetworkApplianceFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallL7FirewallRulesRequest(updateNetworkApplianceFirewallL7FirewallRules).Execute()
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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
