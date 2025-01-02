package vlan

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"strconv"
	"strings"
)

func (n *NetworksApplianceVLANResourceModelDhcpOption) FromTerraformValue(val tftypes.Value) error {
	// Assuming val is a tftypes.Object with "code", "type", and "value"
	var data map[string]tftypes.Value
	if !val.IsKnown() || val.IsNull() {
		return errors.New("value is unknown or null")
	}

	conversionErr := val.As(&data)
	if conversionErr != nil {
		return conversionErr
	}

	var code, typ, value string
	if conversionErr = data["code"].As(&code); conversionErr != nil {
		return conversionErr
	}
	n.Code = basetypes.NewStringValue(code)

	if conversionErr = data["type"].As(&typ); conversionErr != nil {
		return conversionErr
	}
	n.Type = basetypes.NewStringValue(typ)

	if conversionErr = data["value"].As(&value); conversionErr != nil {
		return conversionErr
	}
	n.Value = basetypes.NewStringValue(value)

	return nil
}

func (n *NetworksApplianceVLANModelReservedIpRange) FromTerraformValue(val tftypes.Value) error {
	// Assuming val is a tftypes.Object with "start", "end", and "comment"
	var data map[string]tftypes.Value
	if !val.IsKnown() || val.IsNull() {
		return errors.New("comment is unknown or null")
	}

	conversionErr := val.As(&data)
	if conversionErr != nil {
		return conversionErr
	}

	var start, end, comment string
	if conversionErr = data["start"].As(&start); conversionErr != nil {
		return conversionErr
	}
	n.Start = basetypes.NewStringValue(start)

	if conversionErr = data["end"].As(&end); conversionErr != nil {
		return conversionErr
	}
	n.End = basetypes.NewStringValue(end)

	if conversionErr = data["comment"].As(&comment); conversionErr != nil {
		return conversionErr
	}
	n.Comment = basetypes.NewStringValue(comment)

	return nil
}

// ToAPIPayload converts the Terraform resource data model into the API payload.
func (m *NetworksApplianceVLANModelIpv6) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6, diag.Diagnostics) {
	tflog.Info(ctx, "[start] NetworksApplianceVLANModelIpv6 ToAPIPayload")

	payload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6{}

	// Convert 'Enabled' field
	payload.Enabled = m.Enabled.ValueBoolPointer()

	var prefixAssignments []Ipv6PrefixAssignment

	// Convert 'PrefixAssignments' field
	err := m.PrefixAssignments.ElementsAs(ctx, &prefixAssignments, false)
	if err != nil {
		return nil, err.Errors()
	}

	for _, prefixAssignment := range prefixAssignments {

		var prefixAssignmentPayload openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner

		prefixAssignmentPayload.SetStaticPrefix(prefixAssignment.StaticPrefix.ValueString())
		prefixAssignmentPayload.SetStaticApplianceIp6(prefixAssignment.StaticApplianceIp6.ValueString())
		prefixAssignmentPayload.SetAutonomous(prefixAssignment.Autonomous.ValueBool())

		var origin Ipv6PrefixAssignmentOrigin
		prefixAssignment.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})

		originPayload, err := origin.ToAPIPayload(ctx)
		if err != nil {
			return nil, err.Errors()
		}

		prefixAssignmentPayload.SetOrigin(*originPayload)

		payload.PrefixAssignments = append(payload.PrefixAssignments, prefixAssignmentPayload)
	}

	tflog.Info(ctx, "[finish] NetworksApplianceVLANModelIpv6 ToAPIPayload")
	return payload, nil
}

// FromAPIResponse transforms an API response into the NetworksApplianceVLANModelIpv6 Terraform structure.
func (m *NetworksApplianceVLANModelIpv6) FromAPIResponse(ctx context.Context, apiResponse *openApiClient.GetNetworkApplianceVlans200ResponseInnerIpv6) diag.Diagnostics {
	tflog.Info(ctx, "[start] NetworksApplianceVLANModelIpv6 FromAPIResponse")
	tflog.Trace(ctx, "NetworksApplianceVLANModelIpv6 FromAPIResponse", map[string]interface{}{
		"apiResponse": apiResponse,
	})
	if apiResponse == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("IPv6 Response Error", "Received nil API response for IPv6")}
	}

	m.Enabled = types.BoolValue(apiResponse.GetEnabled())

	var prefixAssignments []Ipv6PrefixAssignment
	for _, apiPA := range apiResponse.PrefixAssignments {
		var pa Ipv6PrefixAssignment

		diags := pa.FromAPIResponse(ctx, &apiPA)
		if diags.HasError() {
			tflog.Warn(ctx, "failed to extract FromAPIResponse to PrefixAssignments")
			return diags
		}

		prefixAssignments = append(prefixAssignments, pa)
	}

	p, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Ipv6PrefixAssignmentAttrTypes()}, prefixAssignments)

	m.PrefixAssignments = p
	tflog.Info(ctx, "[finish] NetworksApplianceVLANModelIpv6 FromAPIResponse")
	return nil
}

func (pa *Ipv6PrefixAssignment) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner, diag.Diagnostics) {
	apiPA := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
		Autonomous:         pa.Autonomous.ValueBoolPointer(),
		StaticPrefix:       pa.StaticPrefix.ValueStringPointer(),
		StaticApplianceIp6: pa.StaticApplianceIp6.ValueStringPointer(),
	}

	var originObject Ipv6PrefixAssignmentOrigin
	originObjectDiags := pa.Origin.As(ctx, originObject, basetypes.ObjectAsOptions{})
	if originObjectDiags.HasError() {
		return nil, originObjectDiags
	}

	// If 'Origin' is a nested structure, convert it too
	if !pa.Origin.IsNull() {
		originAPIModel, diags := originObject.ToAPIPayload(ctx)
		if diags.HasError() {
			return nil, diags
		}
		apiPA.Origin = originAPIModel
	}

	return apiPA, nil
}

// ToAPIPayload converts the Terraform prefix assignment into the API prefix assignment payload.
func (pa *Ipv6PrefixAssignment) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner, diag.Diagnostics) {
	tflog.Info(ctx, "[start] Ipv6PrefixAssignment ToAPIPayload")

	paPayload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{}

	// Autonomous
	paPayload.Autonomous = pa.Autonomous.ValueBoolPointer()

	// StaticPrefix
	paPayload.StaticPrefix = pa.StaticPrefix.ValueStringPointer()

	// StaticApplianceIp6
	paPayload.StaticApplianceIp6 = pa.StaticApplianceIp6.ValueStringPointer()

	// Origin
	var origin Ipv6PrefixAssignmentOrigin
	diags := pa.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	originPayload, diags := origin.ToAPIPayload(ctx)
	if diags.HasError() {
		return nil, diags
	}

	paPayload.Origin = originPayload

	tflog.Info(ctx, "[finish] Ipv6PrefixAssignment ToAPIPayload")
	return paPayload, nil
}

// ToTerraformObject converts the Ipv6PrefixAssignment instance to a map suitable for ObjectValueFrom.
func (pa *Ipv6PrefixAssignment) ToTerraformObject(ctx context.Context) (map[string]attr.Value, diag.Diagnostics) {
	return map[string]attr.Value{
		"autonomous":           pa.Autonomous,
		"static_prefix":        pa.StaticPrefix,
		"static_appliance_ip6": pa.StaticApplianceIp6,
		"origin":               pa.Origin,
	}, nil
}

// FromAPIResponse fills the Ipv6PrefixAssignment with data from the API response.
// This method transforms the OpenAPI response into the format expected by the Terraform provider.
// This method is used in both the vlan_resource and the vlan_datasource
func (pa *Ipv6PrefixAssignment) FromAPIResponse(ctx context.Context, apiPA *openApiClient.GetNetworkApplianceVlans200ResponseInnerIpv6PrefixAssignmentsInner) diag.Diagnostics {
	tflog.Info(ctx, "[start] Ipv6PrefixAssignment FromAPIResponse")
	tflog.Trace(ctx, "NetworksApplianceVLANModelIpv6 FromAPIResponse", map[string]interface{}{
		"apiPA": apiPA,
	})
	if apiPA == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("API Prefix Assignment Error", "Received nil API prefix assignment")}
	}

	pa.Autonomous = types.BoolValue(apiPA.GetAutonomous())

	// staticPrefix
	staticPrefix := apiPA.GetStaticPrefix()

	if staticPrefix == "" {
		// Handle the null scenario
		pa.StaticPrefix = types.StringNull()
	} else {

		pa.StaticPrefix = types.StringValue(apiPA.GetStaticPrefix())

		if pa.StaticPrefix.IsUnknown() {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic("Invalid Assignment", "The 'staticPrefix' field assignment resulted in an unknown value"),
			}
		}
	}

	// staticApplianceIp6
	staticApplianceIp6 := apiPA.GetStaticApplianceIp6()

	if staticApplianceIp6 == "" {
		// Handle the null scenario
		pa.StaticApplianceIp6 = types.StringNull()
	} else {
		pa.StaticApplianceIp6 = types.StringValue(staticApplianceIp6)

		if pa.StaticApplianceIp6.IsUnknown() {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic("Invalid Assignment", "The 'staticApplianceIp6' field assignment resulted in an unknown value"),
			}
		}
	}

	var origin Ipv6PrefixAssignmentOrigin
	originDiags := origin.FromAPIResponse(ctx, apiPA.Origin)
	if originDiags.HasError() {
		return originDiags
	}

	// Use the predefined functions for attribute types and map
	originAttrTypes := Ipv6PrefixAssignmentOriginAttrTypes()
	originAttrMap := Ipv6PrefixAssignmentOriginAttrMap(&origin)

	originTf, diags := types.ObjectValue(originAttrTypes, originAttrMap)
	if diags.HasError() {
		tflog.Warn(ctx, "failed to create object from PrefixAssignment Origin")
		return diags
	}

	pa.Origin = originTf

	tflog.Info(ctx, "[finish] Ipv6PrefixAssignment FromAPIResponse")
	return nil
}

// ToAPIPayload converts the Terraform origin into the API origin payload.
func (o *Ipv6PrefixAssignmentOrigin) ToAPIPayload(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin, diag.Diagnostics) {
	tflog.Info(ctx, "[start] Ipv6PrefixAssignmentOrigin ToAPIPayload")

	originPayload := &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin{}

	originPayload.SetType(o.Type.ValueString())

	var interfaces []string
	// The ElementsAs function usually takes a pointer to the slice
	if diags := o.Interfaces.ElementsAs(ctx, &interfaces, false); diags.HasError() {
		return nil, diags
	}

	// Process interfaces to remove extra quotes if necessary
	for i, iface := range interfaces {
		interfaces[i] = strings.Trim(iface, "\"")
	}

	originPayload.SetInterfaces(interfaces)

	tflog.Info(ctx, "[finish] Ipv6PrefixAssignmentOrigin ToAPIPayload")
	return originPayload, nil
}

// FromAPIResponse fills the Ipv6PrefixAssignmentOrigin with data from the API response.
// This method transforms the OpenAPI response into the format expected by the Terraform provider.
func (o *Ipv6PrefixAssignmentOrigin) FromAPIResponse(ctx context.Context, apiOrigin *openApiClient.CreateNetworkAppliancePrefixesDelegatedStaticRequestOrigin) diag.Diagnostics {
	tflog.Info(ctx, "[start] Ipv6PrefixAssignmentOrigin FromAPIResponse")
	tflog.Trace(ctx, "Ipv6PrefixAssignmentOrigin FromAPIResponse", map[string]interface{}{
		"apiOrigin": apiOrigin,
	})

	// Get the type from API response
	apiType := apiOrigin.GetType()

	// Validate the apiType
	// (Add any specific validation logic here. For example, checking if it's non-empty, or if it matches certain criteria)
	if apiType == "" {
		// Handle the invalid scenario
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("Invalid API Origin Type", "The 'type' field from the API origin is empty"),
		}
	}

	if apiOrigin == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("API Origin Error", "Received nil API origin")}
	}

	// Assuming the value is valid, set the Type
	o.Type = types.StringValue(apiType)

	// Additional check: Verify if the assignment is as expected
	// This part depends on the implementation details of types.StringValue and your specific validation needs
	if o.Type.IsUnknown() || o.Type.IsNull() {
		// Handle the scenario where the assignment didn't work as expected
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("Invalid Assignment", "The 'type' field assignment resulted in an unknown or null value"),
		}
	}

	var interfaces []types.String
	for _, iface := range apiOrigin.Interfaces {
		interfaces = append(interfaces, types.StringValue(iface))
	}

	var diags diag.Diagnostics
	o.Interfaces, diags = types.SetValueFrom(ctx, types.StringType, interfaces)
	if diags.HasError() {
		tflog.Warn(ctx, "failed to create list from Origin interfaces")
		return diags
	}

	tflog.Info(ctx, "[finish] Ipv6PrefixAssignmentOrigin FromAPIResponse")
	return nil
}

func CreateHttpReqPayload(ctx context.Context, data *NetworksApplianceVLANModel) (openApiClient.CreateNetworkApplianceVlanRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	// Log the received request
	tflog.Info(ctx, "[start] Create HTTP Request Payload Call")
	tflog.Trace(ctx, "Create Request Payload", map[string]interface{}{
		"data": data,
	})

	// Initialize the payload
	payload := openApiClient.NewCreateNetworkApplianceVlanRequest(data.Id.ValueString(), data.Name.ValueString())

	// Id
	if !data.VlanId.IsUnknown() && !data.VlanId.IsNull() {

		// API returns this as string, openAPI spec has set as Integer
		payload.SetId(fmt.Sprintf("%v", data.VlanId.ValueInt64()))
	}

	// Name
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	// Subnet
	if !data.Subnet.IsUnknown() && !data.Subnet.IsNull() {
		payload.SetSubnet(data.Subnet.ValueString())
	}

	// ApplianceIp
	if !data.ApplianceIp.IsUnknown() && !data.ApplianceIp.IsNull() {
		payload.SetApplianceIp(data.ApplianceIp.ValueString())
	}

	// GroupPolicyId
	if !data.GroupPolicyId.IsUnknown() && !data.GroupPolicyId.IsNull() {
		payload.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	}

	// TemplateVlanType
	if !data.TemplateVlanType.IsUnknown() && !data.TemplateVlanType.IsNull() {
		payload.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	}

	// Cidr
	if !data.Cidr.IsUnknown() && !data.Cidr.IsNull() {
		payload.SetCidr(data.Cidr.ValueString())
	}

	// Mask
	if !data.Mask.IsUnknown() && !data.Mask.IsNull() {
		payload.SetMask(int32(data.Mask.ValueInt64()))
	}

	// IPV6
	if !data.IPv6.IsUnknown() && !data.IPv6.IsNull() {

		var ipv6 NetworksApplianceVLANModelIpv6
		diags := data.IPv6.As(ctx, &ipv6, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return *payload, diags
		}

		ipv6Payload, ipv6PayloadErr := ipv6.ToAPIPayload(ctx)
		if ipv6PayloadErr.HasError() {
			return *payload, ipv6PayloadErr
		}

		// Enabled
		ipv6Payload.SetEnabled(ipv6.Enabled.ValueBool())

		// Handle Prefix Assignments
		var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
		for _, prefixAssignmentData := range ipv6.PrefixAssignments.Elements() {

			// Convert the prefixAssignmentData (which is of type attr.Value) to your struct
			var pa Ipv6PrefixAssignment

			prefixAssignmentDataDiags := tfsdk.ValueAs(ctx, prefixAssignmentData, &pa)
			if prefixAssignmentDataDiags.HasError() {
				return *payload, prefixAssignmentDataDiags
			}

			// Now create your API client struct
			var prefixAssignment openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
			prefixAssignment.SetAutonomous(pa.Autonomous.ValueBool())
			prefixAssignment.SetStaticPrefix(pa.StaticPrefix.ValueString())
			prefixAssignment.SetStaticApplianceIp6(pa.StaticApplianceIp6.ValueString())

			// Assuming 'Origin' is another struct that you need to convert similarly
			var originData Ipv6PrefixAssignmentOrigin
			originDiags := pa.Origin.As(ctx, &originData, basetypes.ObjectAsOptions{})
			if originDiags.HasError() {
				return *payload, originDiags
			}

			// Populate originData into the prefixAssignment's Origin field
			origin, originDiags := originData.ToAPIPayload(ctx)
			if originDiags.HasError() {
				return *payload, originDiags
			}

			prefixAssignment.SetOrigin(*origin)

			prefixAssignments = append(prefixAssignments, prefixAssignment)
		}

		ipv6Payload.SetPrefixAssignments(prefixAssignments)

		payload.SetIpv6(*ipv6Payload)
	}

	// MandatoryDhcp
	if !data.MandatoryDhcp.IsUnknown() && !data.MandatoryDhcp.IsNull() {
		mandatoryDhcpPayload := openApiClient.NewGetNetworkApplianceVlans200ResponseInnerMandatoryDhcp()
		var mandatoryDhcp NetworksApplianceVLANModelMandatoryDhcp

		diags := data.MandatoryDhcp.As(ctx, &mandatoryDhcp, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		mandatoryDhcpPayload.SetEnabled(mandatoryDhcp.Enabled.ValueBool())

		payload.SetMandatoryDhcp(*mandatoryDhcpPayload)
	}

	tflog.Info(ctx, "[finish] Create HTTP Request Payload")
	tflog.Trace(ctx, "Create Request Payload", map[string]interface{}{
		"payload": payload,
	})

	return *payload, nil
}

func CreateHttpResponse(ctx context.Context, data *NetworksApplianceVLANModel, response *openApiClient.CreateNetworkApplianceVlan201Response) diag.Diagnostics {

	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] CreatePayloadResponse Call")
	tflog.Trace(ctx, "Create Payload Response", map[string]interface{}{
		"response": response,
	})

	// Set to Ids needed for importing resource
	if data.Id.IsUnknown() {
		data.Id = types.StringValue(fmt.Sprintf("%s,%s", data.NetworkId.ValueString(), data.VlanId.String()))
	}

	// VlanId
	if response.HasId() {
		// API returns string, openAPI spec defines int
		idStr := response.GetId()

		// Check if the string is not empty
		if idStr != "" {

			// Convert string to int
			vlanId, err := strconv.Atoi(idStr)
			if err != nil {
				// Handle the error if conversion fails
				resp.AddError("CreateHttpResponse VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s': %v", idStr, err))
			}

			data.VlanId = types.Int64Value(int64(vlanId))
		}
	}

	// InterfaceId
	if response.HasInterfaceId() {
		data.InterfaceId = types.StringValue(response.GetInterfaceId())
	} else {
		if data.InterfaceId.IsUnknown() {
			data.InterfaceId = types.StringNull()
		}
	}

	// Name
	if response.HasName() {
		data.Name = types.StringValue(response.GetName())
	} else {
		if data.Name.IsUnknown() {
			data.Name = types.StringNull()
		}
	}

	// Subnet
	if response.HasSubnet() {
		data.Subnet = types.StringValue(response.GetSubnet())
	} else {
		if data.Subnet.IsUnknown() {
			data.Subnet = types.StringNull()
		}
	}

	// ApplianceIp
	if response.HasApplianceIp() {
		data.ApplianceIp = types.StringValue(response.GetApplianceIp())
	} else {
		if data.ApplianceIp.IsUnknown() {
			data.ApplianceIp = types.StringNull()
		}
	}

	// GroupPolicyId
	if response.HasGroupPolicyId() {
		data.GroupPolicyId = types.StringValue(response.GetGroupPolicyId())
	} else {
		if data.GroupPolicyId.IsUnknown() {
			data.GroupPolicyId = types.StringNull()
		}
	}

	// TemplateVlanType
	if response.HasTemplateVlanType() {
		data.TemplateVlanType = types.StringValue(response.GetTemplateVlanType())
	} else {
		if data.TemplateVlanType.IsUnknown() {
			data.TemplateVlanType = types.StringNull()
		}
	}

	// Cidr
	if response.HasCidr() {
		data.Cidr = types.StringValue(response.GetCidr())
	} else {
		if data.Cidr.IsUnknown() {
			data.Cidr = types.StringNull()
		}
	}

	// Mask
	if response.HasMask() {
		data.Mask = types.Int64Value(int64(response.GetMask()))
	} else {
		if data.Mask.IsUnknown() {
			data.Mask = types.Int64Null()
		}
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANModelMandatoryDhcp{}

		// Enabled
		if response.MandatoryDhcp.HasEnabled() {
			mandatoryDhcp.Enabled = types.BoolValue(response.MandatoryDhcp.GetEnabled())
		}

		mandatoryDhcpAttributes := map[string]attr.Type{
			"enabled": types.BoolType,
		}

		objectVal, diags := types.ObjectValueFrom(ctx, mandatoryDhcpAttributes, mandatoryDhcp)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.MandatoryDhcp = objectVal
	} else {
		if data.MandatoryDhcp.IsUnknown() {
			data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
				"enabled": types.BoolType,
			})
		}
	}

	if response.HasIpv6() {
		ipv6Instance := NetworksApplianceVLANModelIpv6{}
		diags := ipv6Instance.FromAPIResponse(ctx, response.Ipv6)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Object, diags := types.ObjectValueFrom(ctx, Ipv6AttrTypes(), ipv6Instance)
		if diags.HasError() {
			resp.Append(diags...)
		}
		tflog.Warn(ctx, fmt.Sprintf("CreateHttpResponse: %v", ipv6Object.String()))

		data.IPv6 = ipv6Object

		tflog.Warn(ctx, fmt.Sprintf("CREATE PAYLOAD: %v", data.IPv6.String()))

	} else {
		if data.IPv6.IsUnknown() {
			ipv6Instance := NetworksApplianceVLANModelIpv6{}
			ipv6Prefixes := []Ipv6PrefixAssignment{}

			ipv6PrefixesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Ipv6PrefixAssignmentAttrTypes()}, ipv6Prefixes)
			if diags.HasError() {
				resp.Append(diags...)
			}

			ipv6Instance.PrefixAssignments = ipv6PrefixesList

			ipv6Object, diags := types.ObjectValueFrom(ctx, Ipv6AttrTypes(), ipv6Instance)
			if diags.HasError() {
				resp.Append(diags...)
			}

			data.IPv6 = ipv6Object
		}
	}

	tflog.Info(ctx, "[finish] CreateResponsePayloadResponse Call")

	return resp
}

// ReadHttpResponse - used by READ, UPDATE & DELETE funcs
func ReadHttpResponse(ctx context.Context, data *NetworksApplianceVLANModel, response *openApiClient.GetNetworkApplianceVlans200ResponseInner) diag.Diagnostics {

	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] ReadHttpResponse Call")
	tflog.Trace(ctx, "Read Response Payload ", map[string]interface{}{
		"response": response,
	})

	// Set to Ids needed for importing resource
	if data.Id.IsUnknown() {
		data.Id = types.StringValue(fmt.Sprintf("%s,%s", data.NetworkId.ValueString(), data.VlanId.String()))
	}

	// VlanId
	if response.HasId() {
		// API returns string, openAPI spec defines int
		idStr := response.GetId()

		// Check if the string is not empty
		if idStr != "" {

			// Convert string to int
			vlanId, err := strconv.Atoi(idStr)
			if err != nil {
				// Handle the error if conversion fails
				resp.AddError("CreateHttpResponse VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s': %v", idStr, err))
			}

			data.VlanId = types.Int64Value(int64(vlanId))
		}
	}

	// InterfaceId
	if response.HasInterfaceId() {
		data.InterfaceId = types.StringValue(response.GetInterfaceId())
	} else {
		if data.InterfaceId.IsUnknown() {
			data.InterfaceId = types.StringNull()
		}
	}

	// Name
	if response.HasName() {
		data.Name = types.StringValue(response.GetName())
	} else {
		if data.Name.IsUnknown() {
			data.Name = types.StringNull()
		}
	}

	// Subnet
	if response.HasSubnet() {
		data.Subnet = types.StringValue(response.GetSubnet())
	} else {
		if data.Subnet.IsUnknown() {
			data.Subnet = types.StringNull()
		}
	}

	// ApplianceIp
	if response.HasApplianceIp() {
		data.ApplianceIp = types.StringValue(response.GetApplianceIp())
	} else {
		if data.ApplianceIp.IsUnknown() {
			data.ApplianceIp = types.StringNull()
		}
	}

	// GroupPolicyId
	if response.HasGroupPolicyId() {
		data.GroupPolicyId = types.StringValue(response.GetGroupPolicyId())
	} else {
		if data.GroupPolicyId.IsUnknown() {
			data.GroupPolicyId = types.StringNull()
		}
	}

	// TemplateVlanType
	if response.HasTemplateVlanType() {
		data.TemplateVlanType = types.StringValue(response.GetTemplateVlanType())
	} else {
		if data.TemplateVlanType.IsUnknown() {
			data.TemplateVlanType = types.StringNull()
		}
	}

	// Cidr
	if response.HasCidr() {
		data.Cidr = types.StringValue(response.GetCidr())
	} else {
		if data.Cidr.IsUnknown() {
			data.Cidr = types.StringNull()
		}
	}

	// Mask
	if response.HasMask() {
		data.Mask = types.Int64Value(int64(response.GetMask()))
	} else {
		if data.Mask.IsUnknown() {
			data.Mask = types.Int64Null()
		}
	}

	// DhcpRelayServerIps
	if response.HasDhcpRelayServerIps() {
		var dhcpRelayServerIps []attr.Value
		for _, dhcpRelayServerIp := range response.GetDhcpRelayServerIps() {
			dhcpRelayServerIps = append(dhcpRelayServerIps, types.StringValue(dhcpRelayServerIp))

		}
		dhcpRelayServerIpsRespData, dhcpRelayServerIpsDiags := basetypes.NewListValue(types.StringType, dhcpRelayServerIps)
		if dhcpRelayServerIpsDiags.HasError() {
			resp.Append(dhcpRelayServerIpsDiags...)
		}

		data.DhcpRelayServerIps = dhcpRelayServerIpsRespData
	} else {
		if data.DhcpRelayServerIps.IsUnknown() {
			data.DhcpRelayServerIps = basetypes.NewListNull(types.StringType)
		}
	}

	// DhcpHandling
	if response.HasDhcpHandling() {
		data.DhcpHandling = types.StringValue(response.GetDhcpHandling())

	} else {
		if data.DhcpHandling.IsUnknown() {
			data.DhcpHandling = types.StringNull()
		}
	}

	// DhcpLeaseTime
	if response.HasDhcpLeaseTime() {
		data.DhcpLeaseTime = types.StringValue(response.GetDhcpLeaseTime())

	} else {
		if data.DhcpLeaseTime.IsUnknown() {
			data.DhcpLeaseTime = types.StringNull()
		}
	}

	// DhcpBootOptionsEnabled
	if response.HasDhcpBootOptionsEnabled() {
		data.DhcpBootOptionsEnabled = types.BoolValue(response.GetDhcpBootOptionsEnabled())
	} else {
		if data.DhcpBootOptionsEnabled.IsUnknown() {
			data.DhcpBootOptionsEnabled = types.BoolNull()
		}
	}

	// DhcpBootNextServer
	if response.HasDhcpBootNextServer() {
		data.DhcpBootNextServer = types.StringValue(response.GetDhcpBootNextServer())
	} else {
		if data.DhcpBootNextServer.IsUnknown() {
			data.DhcpBootNextServer = types.StringNull()
		}
	}

	// DhcpBootFilename
	if response.HasDhcpBootFilename() {
		data.DhcpBootFilename = types.StringValue(response.GetDhcpBootFilename())

	} else {
		if data.DhcpBootFilename.IsUnknown() {
			data.DhcpBootFilename = types.StringNull()
		}
	}

	// FixedIpAssignments
	if response.HasFixedIpAssignments() {
		fixedIpAssignmentsMap := map[string]attr.Value{}

		fixedIpAssignmentAttrTypes := map[string]attr.Type{
			"ip":   types.StringType,
			"name": types.StringType,
		}

		for macAddress, assignmentInterface := range response.GetFixedIpAssignments() {

			fixedIpAssignmentObject := map[string]attr.Value{
				"ip":   types.StringValue(assignmentInterface.GetIp()),
				"name": types.StringValue(assignmentInterface.GetName()),
			}

			fixedIpAssignmentValue, fixedIpAssignmentsDiags := types.ObjectValue(fixedIpAssignmentAttrTypes, fixedIpAssignmentObject)
			if fixedIpAssignmentsDiags.HasError() {
				resp.Append(fixedIpAssignmentsDiags...)
				continue
			}

			fixedIpAssignmentsMap[macAddress] = fixedIpAssignmentValue
		}

		fixedIpAssignmentsValue, fixedIpAssignmentsDiags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: fixedIpAssignmentAttrTypes}, fixedIpAssignmentsMap)
		if fixedIpAssignmentsDiags.HasError() {
			resp.Append(fixedIpAssignmentsDiags...)
		}

		data.FixedIpAssignments = fixedIpAssignmentsValue

	} else {
		if data.FixedIpAssignments.IsUnknown() {
			fixedIpAssignmentsAttrTypes := map[string]attr.Type{
				"ip":   types.StringType,
				"name": types.StringType,
			}

			data.FixedIpAssignments = types.MapNull(
				types.ObjectType{AttrTypes: fixedIpAssignmentsAttrTypes},
			)
		}
	}

	if response.HasReservedIpRanges() {

		reservedIpRangeAttrTypes := map[string]attr.Type{
			"comment": types.StringType,
			"end":     types.StringType,
			"start":   types.StringType,
		}

		// Define the ObjectType
		objectType := types.ObjectType{AttrTypes: reservedIpRangeAttrTypes}

		// Create a slice to hold the ObjectValues
		var objectValues []attr.Value

		for _, reservedIpRange := range response.GetReservedIpRanges() {

			// Construct the map for the current ObjectValue
			valuesMap, valuesMapErrs := basetypes.NewObjectValue(reservedIpRangeAttrTypes, map[string]attr.Value{
				"start":   basetypes.NewStringValue(reservedIpRange.GetStart()),
				"end":     basetypes.NewStringValue(reservedIpRange.GetEnd()),
				"comment": basetypes.NewStringValue(reservedIpRange.GetComment()),
			})

			if valuesMapErrs.HasError() {
				for _, valuesMapErr := range valuesMapErrs.Errors() {
					tflog.Error(ctx, valuesMapErr.Summary()+valuesMapErr.Detail())
				}
				resp.Append(valuesMapErrs...)

				continue
			}

			// Create the ObjectValue for the current DHCP option
			reservedIpRangeValue, reservedIpRangeDiags := types.ObjectValueFrom(ctx, reservedIpRangeAttrTypes, valuesMap)
			if reservedIpRangeDiags.HasError() {
				for _, dhcpOptionReservedIpRangeDiag := range reservedIpRangeDiags.Errors() {
					tflog.Error(ctx, dhcpOptionReservedIpRangeDiag.Summary()+dhcpOptionReservedIpRangeDiag.Detail())
				}
				resp.Append(reservedIpRangeDiags...)

				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, reservedIpRangeValue)

		}

		// Create a ListValue from the slice of ObjectValue
		reservedIpRangesValue, reservedIpRangesListDiags := types.ListValueFrom(ctx, objectType, objectValues)
		if reservedIpRangesListDiags.HasError() {
			resp.Append(reservedIpRangesListDiags...)
		}

		data.ReservedIpRanges = reservedIpRangesValue

	} else {
		if data.ReservedIpRanges.IsUnknown() {
			data.ReservedIpRanges = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"comment": types.StringType,
					"end":     types.StringType,
					"start":   types.StringType,
				},
			})
		}
	}

	// DnsNameservers
	if response.HasDnsNameservers() {
		data.DnsNameservers = types.StringValue(response.GetDnsNameservers())
	} else {
		if data.DnsNameservers.IsUnknown() {
			data.DnsNameservers = types.StringNull()
		}
	}

	// VpnNatSubnet
	if response.HasVpnNatSubnet() {
		data.VpnNatSubnet = types.StringValue(response.GetVpnNatSubnet())

	} else {
		if data.VpnNatSubnet.IsUnknown() {
			data.VpnNatSubnet = types.StringNull()
		}
	}

	// DhcpOptions
	if response.HasDhcpOptions() {

		// Define the structure of each object in the list
		dhcpOptionsAttrTypes := map[string]attr.Type{
			"code":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		}

		// Define the ObjectType
		objectType := types.ObjectType{AttrTypes: dhcpOptionsAttrTypes}

		// Create a slice to hold the ObjectValues
		var objectValues []attr.Value

		for _, dhcpOption := range response.GetDhcpOptions() {
			// Construct the map for the current ObjectValue
			valuesMap, valuesMapErrs := basetypes.NewObjectValue(dhcpOptionsAttrTypes, map[string]attr.Value{
				"code":  basetypes.NewStringValue(dhcpOption.GetCode()),
				"type":  basetypes.NewStringValue(dhcpOption.GetType()),
				"value": basetypes.NewStringValue(dhcpOption.GetValue()),
			})

			if valuesMapErrs.HasError() {
				for _, valuesMapErr := range valuesMapErrs.Errors() {
					tflog.Error(ctx, valuesMapErr.Summary()+valuesMapErr.Detail())
				}
				resp.Append(valuesMapErrs...)

				continue
			}

			// Create the ObjectValue for the current DHCP option
			dhcpOptionValue, dhcpOptionsDiags := types.ObjectValueFrom(ctx, dhcpOptionsAttrTypes, valuesMap)
			if dhcpOptionsDiags.HasError() {
				for _, dhcpOptionsDiag := range dhcpOptionsDiags.Errors() {
					tflog.Error(ctx, dhcpOptionsDiag.Summary()+dhcpOptionsDiag.Detail())
				}
				resp.Append(dhcpOptionsDiags...)

				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, dhcpOptionValue)
		}

		// Create a ListValue from the slice of ObjectValue
		dhcpOptionsListValue, dhcpOptionsListDiags := types.ListValueFrom(ctx, objectType, objectValues)
		if dhcpOptionsListDiags.HasError() {
			resp.Append(dhcpOptionsListDiags...)
		}

		data.DhcpOptions = dhcpOptionsListValue
	} else {
		if data.DhcpOptions.IsUnknown() {
			data.DhcpOptions = types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"code":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			}})
		}
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANModelMandatoryDhcp{}

		// Enabled
		if response.MandatoryDhcp.HasEnabled() {
			mandatoryDhcp.Enabled = types.BoolValue(response.MandatoryDhcp.GetEnabled())
		}

		mandatoryDhcpAttributes := map[string]attr.Type{
			"enabled": types.BoolType,
		}

		objectVal, diags := types.ObjectValueFrom(ctx, mandatoryDhcpAttributes, mandatoryDhcp)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.MandatoryDhcp = objectVal
	} else {
		if data.MandatoryDhcp.IsUnknown() {
			data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
				"enabled": types.BoolType,
			})
		}
	}

	// IPv6
	// Assuming response is a structure containing your data
	if response.HasIpv6() {
		ipv6Instance := NetworksApplianceVLANModelIpv6{}
		diags := ipv6Instance.FromAPIResponse(ctx, response.Ipv6)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Object, diags := types.ObjectValueFrom(ctx, Ipv6AttrTypes(), ipv6Instance)

		if diags.HasError() {
			resp.Append(diags...)
		}

		data.IPv6 = ipv6Object
	} else {
		if data.IPv6.IsUnknown() {
			ipv6Instance := NetworksApplianceVLANModelIpv6{}
			ipv6Prefixes := []Ipv6PrefixAssignment{}

			ipv6PrefixesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Ipv6PrefixAssignmentAttrTypes()}, ipv6Prefixes)
			if diags.HasError() {
				resp.Append(diags...)
			}

			ipv6Instance.PrefixAssignments = ipv6PrefixesList

			ipv6Object, diags := types.ObjectValueFrom(ctx, Ipv6AttrTypes(), ipv6Instance)
			if diags.HasError() {
				resp.Append(diags...)
			}

			data.IPv6 = ipv6Object
		}
	}

	tflog.Info(ctx, "[finish] ReadHttpResponse Call")

	return resp
}

func UpdateHttpReqPayload(ctx context.Context, data *NetworksApplianceVLANModel) (*openApiClient.UpdateNetworkApplianceVlanRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] UpdateHttpReqPayload Call")

	payload := openApiClient.NewUpdateNetworkApplianceVlanRequest()

	// Name
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	// Subnet
	if !data.Subnet.IsUnknown() && !data.Subnet.IsNull() {
		payload.SetSubnet(data.Subnet.ValueString())
	}

	// ApplianceIp
	if !data.ApplianceIp.IsUnknown() && !data.ApplianceIp.IsNull() {
		payload.SetApplianceIp(data.ApplianceIp.ValueString())
	}

	// GroupPolicyId
	if !data.GroupPolicyId.IsUnknown() && !data.GroupPolicyId.IsNull() {
		payload.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	}

	// VpnNatSubnet
	if !data.VpnNatSubnet.IsUnknown() && !data.VpnNatSubnet.IsNull() {
		payload.SetVpnNatSubnet(data.VpnNatSubnet.ValueString())
	}

	// DhcpHandling
	if !data.DhcpHandling.IsUnknown() && !data.DhcpHandling.IsNull() {
		payload.SetDhcpHandling(data.DhcpHandling.ValueString())
	}

	// DhcpRelayServerIps
	if !data.DhcpRelayServerIps.IsUnknown() && !data.DhcpRelayServerIps.IsNull() {
		var dhcpRelayServerIps []string

		for _, dhcpRelayServerIp := range data.DhcpRelayServerIps.Elements() {
			dhcpRelayServerIps = append(dhcpRelayServerIps, dhcpRelayServerIp.String())
		}

		payload.SetDhcpRelayServerIps(dhcpRelayServerIps)
	}

	// DhcpLeaseTime
	if !data.DhcpLeaseTime.IsUnknown() && !data.DhcpLeaseTime.IsNull() {
		payload.SetDhcpLeaseTime(data.DhcpLeaseTime.ValueString())
	}

	// DhcpBootOptionsEnabled
	if !data.DhcpBootOptionsEnabled.IsUnknown() && !data.DhcpBootOptionsEnabled.IsNull() {
		payload.SetDhcpBootOptionsEnabled(data.DhcpBootOptionsEnabled.ValueBool())
	}

	// DhcpBootNextServer
	if !data.DhcpBootNextServer.IsUnknown() && !data.DhcpBootNextServer.IsNull() {
		payload.SetDhcpBootNextServer(data.DhcpBootNextServer.ValueString())
	}

	// DhcpBootFilename
	if !data.DhcpBootFilename.IsUnknown() && !data.DhcpBootFilename.IsNull() {
		payload.SetDhcpBootFilename(data.DhcpBootFilename.ValueString())
	}

	// FixedIpAssignments
	if !data.FixedIpAssignments.IsUnknown() && !data.FixedIpAssignments.IsNull() {

		var fixedIpAssignments map[string]FixedIpAssignmentTerraform
		fixedIpAssignmentsDiags := data.FixedIpAssignments.ElementsAs(ctx, &fixedIpAssignments, true)
		if fixedIpAssignmentsDiags.HasError() {
			resp.AddError(
				"Create Payload Failure, FixedIpAssignments", fmt.Sprintf("%v", fixedIpAssignmentsDiags),
			)
		}

		fixedIpAssignmentsIf := map[string]interface{}{}
		for key, val := range fixedIpAssignments {
			apiAttrs := FixedIpAssignment{}
			if !val.IP.IsUnknown() {
				val := val.IP.ValueString()
				apiAttrs.IP = val
			}

			if !val.Name.IsUnknown() {
				apiAttrs.Name = val.Name.ValueString()
			}

			fixedIpAssignmentsIf[key] = apiAttrs
		}

		payload.SetFixedIpAssignments(fixedIpAssignmentsIf)
	}

	// ReservedIpRanges
	if !data.ReservedIpRanges.IsUnknown() && !data.ReservedIpRanges.IsNull() {

		var reservedIpRanges []openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner

		for _, reservedIpRange := range data.ReservedIpRanges.Elements() {

			var reservedIpRangePayload openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner
			var reservedIpRangeData NetworksApplianceVLANModelReservedIpRange

			// Convert dhcpOption (types.Object) to dhcpOptionData (struct)
			reservedIpRangeValue, reservedIpRangeDiags := reservedIpRange.ToTerraformValue(ctx)
			if reservedIpRangeDiags != nil {
				// Handle errors during conversion
				resp.AddError(
					"Error converting reservedIpRange",
					reservedIpRangeDiags.Error(),
				)
				continue
			}

			err := reservedIpRangeData.FromTerraformValue(reservedIpRangeValue)
			if err != nil {
				// Handle errors during conversion
				resp.AddError("Error converting reservedIpRange Value", err.Error())
				tflog.Warn(ctx, err.Error())
				continue
			}

			// Set Payload
			reservedIpRangePayload.SetStart(reservedIpRangeData.Start.ValueString())
			reservedIpRangePayload.SetEnd(reservedIpRangeData.End.ValueString())
			reservedIpRangePayload.SetComment(reservedIpRangeData.Comment.ValueString())

			reservedIpRanges = append(reservedIpRanges, reservedIpRangePayload)
		}

		payload.SetReservedIpRanges(reservedIpRanges)
	}

	// DnsNameservers
	if !data.DnsNameservers.IsUnknown() && !data.DnsNameservers.IsNull() {
		payload.SetDnsNameservers(data.DnsNameservers.ValueString())
	}

	// DhcpOptions
	if !data.DhcpOptions.IsUnknown() && !data.DhcpOptions.IsNull() {

		var dhcpOptionsPayload []openApiClient.GetNetworkApplianceVlans200ResponseInnerDhcpOptionsInner

		for _, dhcpOption := range data.DhcpOptions.Elements() {
			var dhcpOptionPayload openApiClient.GetNetworkApplianceVlans200ResponseInnerDhcpOptionsInner
			var dhcpOptionData NetworksApplianceVLANResourceModelDhcpOption

			// Convert dhcpOption (types.Object) to dhcpOptionData (struct)
			dhcpOptionValue, dhcpOptionDiags := dhcpOption.ToTerraformValue(ctx)
			if dhcpOptionDiags != nil {
				// Handle errors during conversion
				resp.AddError(
					"Error converting DHCP option",
					dhcpOptionDiags.Error(),
				)
				continue
			}

			err := dhcpOptionData.FromTerraformValue(dhcpOptionValue)
			if err != nil {
				// Handle errors during conversion
				resp.AddError("Error converting DHCP option", err.Error())
				tflog.Warn(ctx, err.Error())
				continue
			}

			// Set Payload
			dhcpOptionPayload.SetCode(dhcpOptionData.Code.ValueString())
			dhcpOptionPayload.SetType(dhcpOptionData.Type.ValueString())
			dhcpOptionPayload.SetValue(dhcpOptionData.Value.ValueString())

			dhcpOptionsPayload = append(dhcpOptionsPayload, dhcpOptionPayload)

		}

		payload.SetDhcpOptions(dhcpOptionsPayload)
	}

	// TemplateVlanType
	if !data.TemplateVlanType.IsUnknown() && !data.TemplateVlanType.IsNull() {
		payload.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	}

	// Cidr
	if !data.Cidr.IsUnknown() && !data.Cidr.IsNull() {
		payload.SetCidr(data.Cidr.ValueString())
	}

	// Mask
	if !data.Mask.IsUnknown() && !data.Mask.IsNull() {
		payload.SetMask(int32(data.Mask.ValueInt64()))
	}

	// MandatoryDhcp
	if !data.MandatoryDhcp.IsUnknown() && !data.MandatoryDhcp.IsNull() {
		mandatoryDhcpPayload := openApiClient.NewGetNetworkApplianceVlans200ResponseInnerMandatoryDhcp()
		var mandatoryDhcp NetworksApplianceVLANModelMandatoryDhcp

		diags := data.MandatoryDhcp.As(ctx, &mandatoryDhcp, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		mandatoryDhcpPayload.SetEnabled(mandatoryDhcp.Enabled.ValueBool())

		payload.SetMandatoryDhcp(*mandatoryDhcpPayload)
	}

	// IPV6
	if !data.IPv6.IsUnknown() && !data.IPv6.IsNull() {

		var ipv6 NetworksApplianceVLANModelIpv6
		diags := data.IPv6.As(ctx, &ipv6, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		ipv6Payload, ipv6PayloadErr := ipv6.ToAPIPayload(ctx)
		if ipv6PayloadErr.HasError() {
			return nil, ipv6PayloadErr
		}

		// Enabled
		ipv6Payload.SetEnabled(ipv6.Enabled.ValueBool())

		// Handle Prefix Assignments
		var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
		for _, prefixAssignmentData := range ipv6.PrefixAssignments.Elements() {
			// Convert the prefixAssignmentData (which is of type attr.Value) to your struct
			var pa Ipv6PrefixAssignment

			prefixAssignmentDataDiags := tfsdk.ValueAs(ctx, prefixAssignmentData, &pa)
			if prefixAssignmentDataDiags.HasError() {
				return nil, prefixAssignmentDataDiags
			}

			// Now create your API client struct
			var prefixAssignment openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
			prefixAssignment.SetAutonomous(pa.Autonomous.ValueBool())
			prefixAssignment.SetStaticPrefix(pa.StaticPrefix.ValueString())
			prefixAssignment.SetStaticApplianceIp6(pa.StaticApplianceIp6.ValueString())

			// Assuming 'Origin' is another struct that you need to convert similarly
			var originData Ipv6PrefixAssignmentOrigin
			originDiags := pa.Origin.As(ctx, &originData, basetypes.ObjectAsOptions{})
			if originDiags.HasError() {
				return nil, originDiags
			}

			// Populate originData into the prefixAssignment's Origin field
			origin, originDiags := originData.ToAPIPayload(ctx)
			if originDiags.HasError() {
				return nil, originDiags
			}

			prefixAssignment.SetOrigin(*origin)

			prefixAssignments = append(prefixAssignments, prefixAssignment)
		}

		ipv6Payload.SetPrefixAssignments(prefixAssignments)

		payload.SetIpv6(*ipv6Payload)
	}

	tflog.Info(ctx, "[finish] UpdateHttpReqPayload Call")
	tflog.Trace(ctx, "Update Request Payload", map[string]interface{}{
		"payload": payload,
	})

	return payload, nil
}

func DatasourceReadHttpResponse(ctx context.Context, data *NetworksApplianceVLANModel, response *openApiClient.GetNetworkApplianceVlans200ResponseInner) diag.Diagnostics {

	resp := diag.Diagnostics{}

	// Id field only returns "", this is a bug in the HTTP client

	// InterfaceId
	if response.HasInterfaceId() {
		data.InterfaceId = types.StringValue(response.GetInterfaceId())
	} else {
		data.InterfaceId = types.StringNull()
	}

	// Name
	if response.HasName() {
		data.Name = types.StringValue(response.GetName())
	} else {
		data.Name = types.StringNull()
	}

	// Subnet
	if response.HasSubnet() {
		data.Subnet = types.StringValue(response.GetSubnet())
	} else {
		data.Subnet = types.StringNull()
	}

	// ApplianceIp
	if response.HasApplianceIp() {
		data.ApplianceIp = types.StringValue(response.GetApplianceIp())
	} else {
		data.ApplianceIp = types.StringNull()
	}

	// GroupPolicyId
	if response.HasGroupPolicyId() {
		data.GroupPolicyId = types.StringValue(response.GetGroupPolicyId())
	} else {
		data.GroupPolicyId = types.StringNull()
	}

	// TemplateVlanType
	if response.HasTemplateVlanType() {
		data.TemplateVlanType = types.StringValue(response.GetTemplateVlanType())
	} else {
		data.TemplateVlanType = types.StringNull()
	}

	// Cidr
	if response.HasCidr() {
		data.Cidr = types.StringValue(response.GetCidr())
	} else {
		data.Cidr = types.StringNull()
	}

	// Mask
	if response.HasMask() {
		data.Mask = types.Int64Value(int64(response.GetMask()))
	} else {
		data.Mask = types.Int64Null()
	}

	// DhcpRelayServerIps
	data.DhcpRelayServerIps = basetypes.NewListNull(types.StringType)
	if response.HasDhcpRelayServerIps() {
		var dhcpRelayServerIps []attr.Value
		for _, dhcpRelayServerIp := range response.GetDhcpRelayServerIps() {
			dhcpRelayServerIps = append(dhcpRelayServerIps, types.StringValue(dhcpRelayServerIp))

		}
		dhcpRelayServerIpsRespData, dhcpRelayServerIpsDiags := basetypes.NewListValue(types.StringType, dhcpRelayServerIps)
		if dhcpRelayServerIpsDiags.HasError() {
			resp.Append(dhcpRelayServerIpsDiags...)
		}

		data.DhcpRelayServerIps = dhcpRelayServerIpsRespData
	}

	// DhcpHandling
	if response.HasDhcpHandling() {
		data.DhcpHandling = types.StringValue(response.GetDhcpHandling())

	} else {
		data.DhcpHandling = types.StringNull()
	}

	// DhcpLeaseTime
	if response.HasDhcpLeaseTime() {
		data.DhcpLeaseTime = types.StringValue(response.GetDhcpLeaseTime())

	} else {
		data.DhcpLeaseTime = types.StringNull()
	}

	// DhcpBootOptionsEnabled
	if response.HasDhcpBootOptionsEnabled() {
		data.DhcpBootOptionsEnabled = types.BoolValue(response.GetDhcpBootOptionsEnabled())
	} else {
		data.DhcpBootOptionsEnabled = types.BoolValue(false)
	}

	// DhcpBootNextServer
	if response.HasDhcpBootNextServer() {
		data.DhcpBootNextServer = types.StringValue(response.GetDhcpBootNextServer())
	} else {
		data.DhcpBootNextServer = types.StringNull()
	}

	// DhcpBootFilename
	if response.HasDhcpBootFilename() {
		data.DhcpBootFilename = types.StringValue(response.GetDhcpBootFilename())

	} else {
		data.DhcpBootFilename = types.StringNull()
	}

	// FixedIpAssignments
	if response.HasFixedIpAssignments() {
		fixedIpAssignmentsMap := map[string]attr.Value{}

		fixedIpAssignmentAttrTypes := map[string]attr.Type{
			"ip":   types.StringType,
			"name": types.StringType,
		}

		for macAddress, assignmentInterface := range response.GetFixedIpAssignments() {

			fixedIpAssignmentObject := map[string]attr.Value{
				"ip":   types.StringValue(assignmentInterface.GetIp()),
				"name": types.StringValue(assignmentInterface.GetName()),
			}

			fixedIpAssignmentValue, fixedIpAssignmentsDiags := types.ObjectValue(fixedIpAssignmentAttrTypes, fixedIpAssignmentObject)
			if fixedIpAssignmentsDiags.HasError() {
				resp.Append(fixedIpAssignmentsDiags...)
				continue
			}

			fixedIpAssignmentsMap[macAddress] = fixedIpAssignmentValue
		}

		fixedIpAssignmentsValue, fixedIpAssignmentsDiags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: fixedIpAssignmentAttrTypes}, fixedIpAssignmentsMap)
		if fixedIpAssignmentsDiags.HasError() {
			resp.Append(fixedIpAssignmentsDiags...)
		}

		data.FixedIpAssignments = fixedIpAssignmentsValue

	} else {
		fixedIpAssignmentsAttrTypes := map[string]attr.Type{
			"ip":   types.StringType,
			"name": types.StringType,
		}

		data.FixedIpAssignments = types.MapNull(
			types.ObjectType{AttrTypes: fixedIpAssignmentsAttrTypes},
		)
	}

	if response.HasReservedIpRanges() {

		reservedIpRangeAttrTypes := map[string]attr.Type{
			"comment": types.StringType,
			"end":     types.StringType,
			"start":   types.StringType,
		}

		// Define the ObjectType
		objectType := types.ObjectType{AttrTypes: reservedIpRangeAttrTypes}

		// Create a slice to hold the ObjectValues
		var objectValues []attr.Value

		for _, reservedIpRange := range response.GetReservedIpRanges() {

			// Construct the map for the current ObjectValue
			valuesMap, valuesMapErrs := basetypes.NewObjectValue(reservedIpRangeAttrTypes, map[string]attr.Value{
				"start":   basetypes.NewStringValue(reservedIpRange.GetStart()),
				"end":     basetypes.NewStringValue(reservedIpRange.GetEnd()),
				"comment": basetypes.NewStringValue(reservedIpRange.GetComment()),
			})

			if valuesMapErrs.HasError() {
				for _, valuesMapErr := range valuesMapErrs.Errors() {
					tflog.Error(ctx, valuesMapErr.Summary()+valuesMapErr.Detail())
				}
				resp.Append(valuesMapErrs...)

				continue
			}

			// Create the ObjectValue for the current DHCP option
			reservedIpRangeValue, reservedIpRangeDiags := types.ObjectValueFrom(ctx, reservedIpRangeAttrTypes, valuesMap)
			if reservedIpRangeDiags.HasError() {
				for _, dhcpOptionReservedIpRangeDiag := range reservedIpRangeDiags.Errors() {
					tflog.Error(ctx, dhcpOptionReservedIpRangeDiag.Summary()+dhcpOptionReservedIpRangeDiag.Detail())
				}
				resp.Append(reservedIpRangeDiags...)

				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, reservedIpRangeValue)

		}

		// Create a ListValue from the slice of ObjectValue
		reservedIpRangesValue, reservedIpRangesListDiags := types.ListValueFrom(ctx, objectType, objectValues)
		if reservedIpRangesListDiags.HasError() {
			resp.Append(reservedIpRangesListDiags...)
		}

		data.ReservedIpRanges = reservedIpRangesValue

	} else {
		data.ReservedIpRanges = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"comment": types.StringType,
				"end":     types.StringType,
				"start":   types.StringType,
			},
		})
	}

	// DnsNameservers
	if response.HasDnsNameservers() {
		data.DnsNameservers = types.StringValue(response.GetDnsNameservers())
	} else {
		data.DnsNameservers = types.StringNull()
	}

	// VpnNatSubnet
	if response.HasVpnNatSubnet() {
		data.VpnNatSubnet = types.StringValue(response.GetVpnNatSubnet())

	} else {
		data.VpnNatSubnet = types.StringNull()
	}

	// DhcpOptions
	if response.HasDhcpOptions() {
		// Define the structure of each object in the list
		dhcpOptionsAttrTypes := map[string]attr.Type{
			"code":  types.StringType,
			"type":  types.StringType,
			"value": types.StringType,
		}

		// Define the ObjectType
		objectType := types.ObjectType{AttrTypes: dhcpOptionsAttrTypes}

		// Create a slice to hold the ObjectValues
		var objectValues []attr.Value

		for _, dhcpOption := range response.GetDhcpOptions() {
			valuesMap := map[string]attr.Value{
				"code":  types.StringValue(dhcpOption.GetCode()),
				"type":  types.StringValue(dhcpOption.GetType()),
				"value": types.StringValue(dhcpOption.GetValue()),
			}

			// Create the ObjectValue for the current DHCP option
			dhcpOptionValue, dhcpOptionsDiags := types.ObjectValue(dhcpOptionsAttrTypes, valuesMap)
			if dhcpOptionsDiags.HasError() {
				resp.Append(dhcpOptionsDiags...)
				continue
			}

			// Add the ObjectValue to the slice
			objectValues = append(objectValues, dhcpOptionValue)
		}

		// Create a ListValue from the slice of ObjectValues
		dhcpOptionsListValue, dhcpOptionsListDiags := types.ListValue(objectType, objectValues)
		if dhcpOptionsListDiags.HasError() {
			resp.Append(dhcpOptionsListDiags...)
		}

		data.DhcpOptions = dhcpOptionsListValue
	} else {
		data.DhcpOptions = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"code":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
		})
	}

	// Mandatory DHCP
	if response.HasMandatoryDhcp() {
		mandatoryDhcp := NetworksApplianceVLANModelMandatoryDhcp{}

		// Enabled
		if response.MandatoryDhcp.HasEnabled() {
			mandatoryDhcp.Enabled = types.BoolValue(response.MandatoryDhcp.GetEnabled())
		}

		mandatoryDhcpAttributes := map[string]attr.Type{
			"enabled": types.BoolType,
		}

		objectVal, diags := types.ObjectValueFrom(ctx, mandatoryDhcpAttributes, mandatoryDhcp)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.MandatoryDhcp = objectVal
	} else {
		data.MandatoryDhcp = types.ObjectNull(map[string]attr.Type{
			"enabled": types.BoolType,
		})
	}

	// IPv6
	// Assuming response is a structure containing your data
	if response.HasIpv6() {
		ipv6Instance := NetworksApplianceVLANModelIpv6{}
		diags := ipv6Instance.FromAPIResponse(ctx, response.Ipv6)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Object, diags := types.ObjectValueFrom(ctx, Ipv6AttrTypes(), ipv6Instance)

		if diags.HasError() {
			resp.Append(diags...)
		}

		data.IPv6 = ipv6Object
	} else {
		ipv6Instance := NetworksApplianceVLANModelIpv6{}
		ipv6Prefixes := []Ipv6PrefixAssignment{}

		ipv6PrefixesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Ipv6PrefixAssignmentAttrTypes()}, ipv6Prefixes)
		if diags.HasError() {
			resp.Append(diags...)
		}

		ipv6Instance.PrefixAssignments = ipv6PrefixesList

		ipv6Object, diags := types.ObjectValueFrom(ctx, Ipv6AttrTypes(), ipv6Instance)
		if diags.HasError() {
			resp.Append(diags...)
		}

		data.IPv6 = ipv6Object
	}

	return resp
}
