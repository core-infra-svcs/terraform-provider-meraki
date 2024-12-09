package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &ManagementInterfaceResource{}
	_ resource.ResourceWithConfigure   = &ManagementInterfaceResource{}
	_ resource.ResourceWithImportState = &ManagementInterfaceResource{}
)

func NewResource() resource.Resource {
	return &ManagementInterfaceResource{}
}

// ManagementInterfaceResource defines the resource implementation.
type ManagementInterfaceResource struct {
	client *openApiClient.APIClient
}

func (r *ManagementInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_management_interface"
}

func (r *ManagementInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GetResourceSchema
}

func (r *ManagementInterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ManagementInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()

	if !data.Wan1.IsNull() && !data.Wan1.IsUnknown() {
		var wan1Plan wanModel
		data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})
		var staticDNS1 []string
		wan1Plan.StaticDns.ElementsAs(ctx, &staticDNS1, false)

		wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{}
		wan1Enabled := wan1Plan.WanEnabled.ValueString()
		if wan1Plan.WanEnabled.IsNull() || wan1Enabled == "" {
			wan1Enabled = "not configured"
		}
		wan1.SetWanEnabled(wan1Enabled)
		wan1.SetStaticDns(staticDNS1)
		wan1.SetStaticGatewayIp(wan1Plan.StaticGatewayIp.ValueString())
		wan1.SetStaticSubnetMask(wan1Plan.StaticSubnetMask.ValueString())
		wan1.SetStaticIp(wan1Plan.StaticIp.ValueString())
		wan1.SetUsingStaticIp(wan1Plan.UsingStaticIp.ValueBool())
		wan1.SetVlan(int32(wan1Plan.Vlan.ValueInt64()))

		if !wan1Plan.Vlan.IsNull() {
			vlan := int32(wan1Plan.Vlan.ValueInt64())
			wan1.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan1 payload before API call", map[string]interface{}{
			"wan1": wan1,
		})

		payload.Wan1 = &wan1
	}

	if !data.Wan2.IsNull() && !data.Wan2.IsUnknown() {
		var wan2Plan wanModel
		data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})
		var staticDNS2 []string
		wan2Plan.StaticDns.ElementsAs(ctx, &staticDNS2, false)

		wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{}
		wan2Enabled := wan2Plan.WanEnabled.ValueString()
		if wan2Plan.WanEnabled.IsNull() || wan2Enabled == "" {
			wan2Enabled = "not configured"
		}
		wan2.SetWanEnabled(wan2Enabled)
		wan2.SetStaticDns(staticDNS2)
		wan2.SetStaticGatewayIp(wan2Plan.StaticGatewayIp.ValueString())
		wan2.SetStaticSubnetMask(wan2Plan.StaticSubnetMask.ValueString())
		wan2.SetStaticIp(wan2Plan.StaticIp.ValueString())
		wan2.SetUsingStaticIp(wan2Plan.UsingStaticIp.ValueBool())
		wan2.SetVlan(int32(wan2Plan.Vlan.ValueInt64()))

		if !wan2Plan.Vlan.IsNull() {
			vlan := int32(wan2Plan.Vlan.ValueInt64())
			wan2.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan2 payload", map[string]interface{}{
			"wan2": wan2,
		})

		payload.Wan2 = &wan2
	}

	if data.Serial.IsNull() || data.Serial.IsUnknown() {
		resp.Diagnostics.AddError(
			"Serial Number Not Found",
			"The serial number must be provided to create the device management interface.",
		)
		return
	}

	serial := data.Serial.ValueString()
	if serial == "" {
		resp.Diagnostics.AddError(
			"Serial Number Empty",
			"The serial number provided is empty. Ensure the serial number is set correctly.",
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
		return inline, respHttp, err
	})
	if err != nil {
		tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"No Management interface information found in API",
			utils.HttpDiagnostics(httpResp),
		)
	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
			return
		}
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}
	}

	// Extract the wan_enabled value directly from Wan1
	var wan1Plan wanModel
	data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})

	// Assign wan_enabled to a variable
	wan1EnabledPlan := wan1Plan.WanEnabled.ValueString()

	var wan2Plan wanModel
	data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})

	// Assign wan_enabled to a variable
	wan2EnabledPlan := wan2Plan.WanEnabled.ValueString()

	diags = updateResourceState(ctx, &data, inlineResp, httpResp, wan1EnabledPlan, wan2EnabledPlan)
	data.Id = types.StringValue(data.Serial.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "create resource")
}

func (r *ManagementInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resourceModel

	// Read Terraform prior state data into the model

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//inlineResp, httpResp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()
		return inline, respHttp, err
	})
	if err != nil {
		tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"No Management interface information found in API",
			utils.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

	}

	// Extract the wan_enabled value directly from Wan1
	var wan1Plan wanModel
	data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})

	// Assign wan_enabled to a variable
	wan1EnabledPlan := wan1Plan.WanEnabled.ValueString()

	var wan2Plan wanModel
	data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})

	// Assign wan_enabled to a variable
	wan2EnabledPlan := wan2Plan.WanEnabled.ValueString()

	diags = updateResourceState(ctx, &data, inlineResp, httpResp, wan1EnabledPlan, wan2EnabledPlan)

	data.Id = types.StringValue(data.Serial.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *ManagementInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()

	if !data.Wan1.IsNull() && !data.Wan1.IsUnknown() {
		var wan1Plan wanModel
		data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})
		var staticDNS1 []string
		wan1Plan.StaticDns.ElementsAs(ctx, &staticDNS1, false)

		wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{}
		wan1Enabled := wan1Plan.WanEnabled.ValueString()
		if wan1Plan.WanEnabled.IsNull() || wan1Enabled == "" {
			wan1Enabled = "not configured"
		}
		wan1.SetWanEnabled(wan1Enabled)
		wan1.SetStaticDns(staticDNS1)
		wan1.SetStaticGatewayIp(wan1Plan.StaticGatewayIp.ValueString())
		wan1.SetStaticSubnetMask(wan1Plan.StaticSubnetMask.ValueString())
		wan1.SetStaticIp(wan1Plan.StaticIp.ValueString())
		wan1.SetUsingStaticIp(wan1Plan.UsingStaticIp.ValueBool())
		wan1.SetVlan(int32(wan1Plan.Vlan.ValueInt64()))

		if !wan1Plan.Vlan.IsNull() {
			vlan := int32(wan1Plan.Vlan.ValueInt64())
			wan1.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan1 payload before API call", map[string]interface{}{
			"wan1": wan1,
		})

		payload.Wan1 = &wan1
	}

	if !data.Wan2.IsNull() && !data.Wan2.IsUnknown() {
		var wan2Plan wanModel
		data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})
		var staticDNS2 []string
		wan2Plan.StaticDns.ElementsAs(ctx, &staticDNS2, false)

		wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{}
		wan2Enabled := wan2Plan.WanEnabled.ValueString()
		if wan2Plan.WanEnabled.IsNull() || wan2Enabled == "" {
			wan2Enabled = "not configured"
		}
		wan2.SetWanEnabled(wan2Enabled)
		wan2.SetStaticDns(staticDNS2)
		wan2.SetStaticGatewayIp(wan2Plan.StaticGatewayIp.ValueString())
		wan2.SetStaticSubnetMask(wan2Plan.StaticSubnetMask.ValueString())
		wan2.SetStaticIp(wan2Plan.StaticIp.ValueString())
		wan2.SetUsingStaticIp(wan2Plan.UsingStaticIp.ValueBool())
		wan2.SetVlan(int32(wan2Plan.Vlan.ValueInt64()))

		if !wan2Plan.Vlan.IsNull() {
			vlan := int32(wan2Plan.Vlan.ValueInt64())
			wan2.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan2 payload", map[string]interface{}{
			"wan2": wan2,
		})

		payload.Wan2 = &wan2
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
		return inline, respHttp, err
	})
	if err != nil {
		tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"No Management interface information found in API",
			utils.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

	}

	// Extract the wan_enabled value directly from Wan1
	var wan1Plan wanModel
	data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})

	// Assign wan_enabled to a variable
	wan1EnabledPlan := wan1Plan.WanEnabled.ValueString()

	var wan2Plan wanModel
	data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})

	// Assign wan_enabled to a variable
	wan2EnabledPlan := wan2Plan.WanEnabled.ValueString()

	diags = updateResourceState(ctx, &data, inlineResp, httpResp, wan1EnabledPlan, wan2EnabledPlan)

	data.Id = types.StringValue(data.Serial.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	tflog.Debug(ctx, "Updated state after API call", map[string]interface{}{
		"data": data,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "updated resource")
}

func (r *ManagementInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()
	wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{}
	wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{}

	payload.Wan1 = &wan1
	payload.Wan2 = &wan2

	_, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"No Management interface information found in API",
			utils.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *ManagementInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
