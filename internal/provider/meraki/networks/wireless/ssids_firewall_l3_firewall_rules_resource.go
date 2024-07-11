package wireless

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &NetworksWirelessSsidsFirewallL3FirewallRulesResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsFirewallL3FirewallRulesResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsFirewallL3FirewallRulesResource{} // Interface for resources with import state functionality
)

func NewNetworksWirelessSsidsFirewallL3FirewallRulesResource() resource.Resource {
	return &NetworksWirelessSsidsFirewallL3FirewallRulesResource{}
}

type NetworksWirelessSsidsFirewallL3FirewallRulesResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The NetworksWirelessSsidsFirewallL3FirewallRulesResourceModel structure describes the data model.
type NetworksWirelessSsidsFirewallL3FirewallRulesResourceModel struct {
	Id             jsontypes.String                                                `tfsdk:"id"`
	NetworkId      jsontypes.String                                                `tfsdk:"network_id" json:"network_id"`
	Number         jsontypes.String                                                `tfsdk:"number"`
	AllowLanAccess jsontypes.Bool                                                  `tfsdk:"allow_lan_access"`
	Rules          []NetworksWirelessSsidsFirewallL3FirewallRulesResourceModelRule `tfsdk:"rules" json:"rules"`
}

type NetworksWirelessSsidsFirewallL3FirewallRulesResourceModelRule struct {
	Comment  jsontypes.String `tfsdk:"comment"`
	DestCidr jsontypes.String `tfsdk:"dest_cidr"`
	DestPort jsontypes.String `tfsdk:"dest_port"`
	Policy   jsontypes.String `tfsdk:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids_firewall_l3_firewall_rules"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		// It should provide a clear and concise description of the resource.
		MarkdownDescription: "NetworksWirelessSsidsFirewallL3FirewallRules for Updating Networks Wireless Ssids Firewall L3FirewallRules",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"number": schema.StringAttribute{
				MarkdownDescription: "SsIds SsidNumber",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"allow_lan_access": schema.BoolAttribute{
				MarkdownDescription: "Allow wireless client access to local LAN (boolean value - true allows access and false denies access) (optional)",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"rules": schema.SetNestedAttribute{
				MarkdownDescription: "An ordered array of the firewall rules for this SSID (not including the local LAN access rule or the default rule)",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"comment": schema.StringAttribute{
							MarkdownDescription: "Description of the rule (optional)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_cidr": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'Any'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_port": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"policy": schema.StringAttribute{
							MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'Any')",
							Required:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"tcp", "udp", "icmp", "icmp6", "Any"}...),
							},
						},
					}}},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// This allows the resource to use the configured provider for any API calls it needs to make.
	r.client = client
}

// Create method is responsible for creating a new resource.
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsFirewallL3FirewallRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidFirewallL3FirewallRules := *openApiClient.NewUpdateNetworkWirelessSsidFirewallL3FirewallRulesRequest()
	updateNetworkWirelessSsidFirewallL3FirewallRules.SetAllowLanAccess(data.AllowLanAccess.ValueBool())
	var rules []openApiClient.UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequestRulesInner
	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequestRulesInner
			if attribute.Comment != jsontypes.StringValue("Default rule") {
				if attribute.Comment != jsontypes.StringValue("Wireless clients accessing LAN") {
					rule.SetComment(attribute.Comment.ValueString())
					rule.SetDestCidr(attribute.DestCidr.ValueString())
					rule.SetDestPort(attribute.DestPort.ValueString())
					rule.SetPolicy(attribute.Policy.ValueString())
					rule.SetProtocol(attribute.Protocol.ValueString())
					rules = append(rules, rule)
				}
			}
		}
	}
	updateNetworkWirelessSsidFirewallL3FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.FirewallApi.UpdateNetworkWirelessSsidFirewallL3FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequest(updateNetworkWirelessSsidFirewallL3FirewallRules).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksWirelessSsidsFirewallL3FirewallRulesResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.FirewallApi.GetNetworkWirelessSsidFirewallL3FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksWirelessSsidsFirewallL3FirewallRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidFirewallL3FirewallRules := *openApiClient.NewUpdateNetworkWirelessSsidFirewallL3FirewallRulesRequest()
	updateNetworkWirelessSsidFirewallL3FirewallRules.SetAllowLanAccess(data.AllowLanAccess.ValueBool())
	var rules []openApiClient.UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequestRulesInner
	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequestRulesInner
			if attribute.Comment != jsontypes.StringValue("Default rule") {
				if attribute.Comment != jsontypes.StringValue("Wireless clients accessing LAN") {
					rule.SetComment(attribute.Comment.ValueString())
					rule.SetDestCidr(attribute.DestCidr.ValueString())
					rule.SetDestPort(attribute.DestPort.ValueString())
					rule.SetPolicy(attribute.Policy.ValueString())
					rule.SetProtocol(attribute.Protocol.ValueString())
					rules = append(rules, rule)
				}
			}
		}
	}
	updateNetworkWirelessSsidFirewallL3FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.FirewallApi.UpdateNetworkWirelessSsidFirewallL3FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequest(updateNetworkWirelessSsidFirewallL3FirewallRules).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksWirelessSsidsFirewallL3FirewallRulesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	updateNetworkWirelessSsidFirewallL3FirewallRules := *openApiClient.NewUpdateNetworkWirelessSsidFirewallL3FirewallRulesRequest()
	updateNetworkWirelessSsidFirewallL3FirewallRules.SetRules(nil)
	updateNetworkWirelessSsidFirewallL3FirewallRules.SetAllowLanAccess(true)
	_, httpResp, err := r.client.FirewallApi.UpdateNetworkWirelessSsidFirewallL3FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL3FirewallRulesRequest(updateNetworkWirelessSsidFirewallL3FirewallRules).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL3FirewallRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, serial number. Got: %q", req.ID),
		)
		return
	}

	// Set the attributes required for making a Read API call in the state.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), idParts[1])...)

	// If there were any errors setting the attributes, return early.
	if resp.Diagnostics.HasError() {
		return
	}
}
