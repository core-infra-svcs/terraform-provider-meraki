package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
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

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *openApiClient.APIClient
}

type resourceModel struct {
	Id                jsontypes.String    `tfsdk:"id"`
	OrganizationId    jsontypes.String    `tfsdk:"organization_id" json:"organizationId"`
	SyslogDefaultRule jsontypes.Bool      `tfsdk:"syslog_default_rule"`
	Rules             []resourceModelRule `tfsdk:"rules" json:"rules"`
}

type resourceModelRule struct {
	Comment       jsontypes.String `tfsdk:"comment"`
	DestCidr      jsontypes.String `tfsdk:"dest_cidr"`
	DestPort      jsontypes.String `tfsdk:"dest_port"`
	Policy        jsontypes.String `tfsdk:"policy"`
	Protocol      jsontypes.String `tfsdk:"protocol"`
	SrcPort       jsontypes.String `tfsdk:"src_port"`
	SrcCidr       jsontypes.String `tfsdk:"src_cidr"`
	SysLogEnabled jsontypes.Bool   `tfsdk:"syslog_enabled"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_appliance_vpn_vpn_firewall_rules"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "OrganizationsApplianceVpnVpnFirewallRules resource for updating Organizations Appliance Vpn Vpn Firewall Rules.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"syslog_default_rule": schema.BoolAttribute{
				MarkdownDescription: "Log the special default rule (boolean value - enable only if you've configured a syslog server) (optional)",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"rules": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"comment": schema.StringAttribute{
							MarkdownDescription: "Description of the rule (optional)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_cidr": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'any'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_port": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination port(s) (integer in the range 1-65535), or 'any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_cidr": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of source IP address(es) (in IP or CIDR notation), or 'any' (note: FQDN not supported for source addresses)",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_port": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of source port(s) (integer in the range 1-65535), or 'any'",
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
							MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'any')",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"syslog_enabled": schema.BoolAttribute{
							MarkdownDescription: "Log this rule to syslog (true or false, boolean value) - only applicable if a syslog has been configured (optional)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
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

	organizationsApplianceVpnVpnFirewallRules := *openApiClient.NewUpdateOrganizationApplianceVpnVpnFirewallRulesRequest()
	var rules []openApiClient.UpdateOrganizationApplianceVpnVpnFirewallRulesRequestRulesInner

	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateOrganizationApplianceVpnVpnFirewallRulesRequestRulesInner
			if attribute.Comment != jsontypes.StringValue("Default rule") {
				rule.SetComment(attribute.Comment.ValueString())
				rule.SetDestCidr(attribute.DestCidr.ValueString())
				rule.SetDestPort(attribute.DestPort.ValueString())
				rule.SetSrcCidr(attribute.SrcCidr.ValueString())
				rule.SetSrcPort(attribute.SrcPort.ValueString())
				rule.SetPolicy(attribute.Policy.ValueString())
				rule.SetProtocol(attribute.Protocol.ValueString())
				rule.SetSyslogEnabled(attribute.SysLogEnabled.ValueBool())
				rules = append(rules, rule)
			}
		}
	}

	organizationsApplianceVpnVpnFirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateOrganizationApplianceVpnVpnFirewallRules(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationApplianceVpnVpnFirewallRulesRequest(organizationsApplianceVpnVpnFirewallRules).Execute()
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

	data.Id = jsontypes.StringValue(data.OrganizationId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetOrganizationApplianceVpnVpnFirewallRules(context.Background(), data.OrganizationId.ValueString()).Execute()
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

	organizationsApplianceVpnVpnFirewallRules := *openApiClient.NewUpdateOrganizationApplianceVpnVpnFirewallRulesRequest()
	var rules []openApiClient.UpdateOrganizationApplianceVpnVpnFirewallRulesRequestRulesInner

	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateOrganizationApplianceVpnVpnFirewallRulesRequestRulesInner
			if attribute.Comment != jsontypes.StringValue("Default rule") {
				rule.SetComment(attribute.Comment.ValueString())
				rule.SetDestCidr(attribute.DestCidr.ValueString())
				rule.SetDestPort(attribute.DestPort.ValueString())
				rule.SetSrcCidr(attribute.SrcCidr.ValueString())
				rule.SetSrcPort(attribute.SrcPort.ValueString())
				rule.SetPolicy(attribute.Policy.ValueString())
				rule.SetProtocol(attribute.Protocol.ValueString())
				rule.SetSyslogEnabled(attribute.SysLogEnabled.ValueBool())
				rules = append(rules, rule)
			}
		}
	}

	organizationsApplianceVpnVpnFirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateOrganizationApplianceVpnVpnFirewallRules(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationApplianceVpnVpnFirewallRulesRequest(organizationsApplianceVpnVpnFirewallRules).Execute()
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

	organizationsApplianceVpnVpnFirewallRules := *openApiClient.NewUpdateOrganizationApplianceVpnVpnFirewallRulesRequest()
	var rules []openApiClient.UpdateOrganizationApplianceVpnVpnFirewallRulesRequestRulesInner
	organizationsApplianceVpnVpnFirewallRules.SetRules(rules)
	organizationsApplianceVpnVpnFirewallRules.SetSyslogDefaultRule(false)

	_, httpResp, err := r.client.ApplianceApi.UpdateOrganizationApplianceVpnVpnFirewallRules(context.Background(), data.OrganizationId.ValueString()).UpdateOrganizationApplianceVpnVpnFirewallRulesRequest(organizationsApplianceVpnVpnFirewallRules).Execute()
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

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
