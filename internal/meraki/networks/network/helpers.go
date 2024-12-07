package network

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"strings"
)

func updateNetworksNetworksResourceCreatePayload(plan *resourceModel) (openApiClient.CreateOrganizationNetworkRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// name
	name := plan.Name.ValueString()

	// ProductTypes
	var productTypes []string
	if !plan.ProductTypes.IsNull() && !plan.ProductTypes.IsUnknown() {
		for _, product := range plan.ProductTypes.Elements() {
			pt := fmt.Sprint(strings.Trim(product.String(), "\""))
			productTypes = append(productTypes, pt)
		}
	}

	// Create HTTP request body
	payload := openApiClient.NewCreateOrganizationNetworkRequest(name, productTypes)

	// Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		for _, tag := range plan.Tags.Elements() {
			t := fmt.Sprint(strings.Trim(tag.String(), "\""))
			tags = append(tags, t)
		}
		payload.SetTags(tags)
	}

	//    TimeZone
	if !plan.Timezone.IsNull() && !plan.Timezone.IsUnknown() {
		payload.SetTimeZone(plan.Timezone.ValueString())
	}

	// CopyFromNetworkId
	if !plan.CopyFromNetworkId.IsNull() && !plan.CopyFromNetworkId.IsUnknown() {
		payload.SetCopyFromNetworkId(plan.CopyFromNetworkId.ValueString())
	}

	// Notes
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		payload.SetNotes(plan.Notes.ValueString())
	}

	return *payload, diags

}

func updateNetworksNetworksResourceUpdatePayload(plan *resourceModel) (openApiClient.UpdateNetworkRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload := openApiClient.NewUpdateNetworkRequest()

	//   Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())
	}

	//    TimeZone
	if !plan.Timezone.IsNull() && !plan.Timezone.IsUnknown() {
		payload.SetTimeZone(plan.Timezone.ValueString())
	}

	//    Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		for _, tag := range plan.Tags.Elements() {
			tags = append(tags, tag.String())
		}
		payload.SetTags(tags)
	}

	//    EnrollmentString
	if !plan.EnrollmentString.IsNull() && !plan.EnrollmentString.IsUnknown() {
		payload.SetEnrollmentString(plan.EnrollmentString.ValueString())
	}

	//    Notes
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		payload.SetNotes(plan.Notes.ValueString())
	}

	return *payload, diags

}

func createNetworksNetworksResourceState(ctx context.Context, state *resourceModel, inlineResp *openApiClient.GetNetwork200Response, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	//  Id (NetworkId)
	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		state.NetworkId = types.StringValue(inlineResp.GetId())
	}

	orgId := fmt.Sprint(strings.Trim(inlineResp.GetOrganizationId(), "\""))
	state.OrganizationId = types.StringValue(orgId)

	//  Id (Terraform Resource)
	if !state.NetworkId.IsNull() || !state.NetworkId.IsUnknown() && !state.OrganizationId.IsNull() || !state.OrganizationId.IsUnknown() {
		importId := state.OrganizationId.String() + "," + inlineResp.GetId()
		state.Id = types.StringValue(importId)
	} else {
		state.Id = types.StringNull()
	}

	//    Name
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name = types.StringValue(inlineResp.GetName())
	}

	//    ProductTypes
	if state.ProductTypes.IsNull() || state.ProductTypes.IsUnknown() {

		var productTypesList []string

		productTypesList = append(productTypesList, inlineResp.ProductTypes...)

		productTypesListObj, err := types.SetValueFrom(ctx, types.StringType, productTypesList)
		if err.HasError() {
			diags.Append(err...)
		}

		state.ProductTypes = productTypesListObj

	}

	//  	Timezone
	if state.Timezone.IsNull() || state.Timezone.IsUnknown() {
		state.Timezone = types.StringValue(inlineResp.GetTimeZone())
	}

	//    Tags
	if state.Tags.IsNull() || state.Tags.IsUnknown() {

		// Tags
		var tagsList []string
		for _, tag := range inlineResp.Tags {
			// Strip any extra quotes from the tags
			tagsList = append(tagsList, strings.Trim(tag, `"`))
		}
		tagsListObj, err := types.SetValueFrom(ctx, types.StringType, tagsList)
		if err.HasError() {
			diags.Append(err...)
		}
		state.Tags = tagsListObj

	}

	//    EnrollmentString
	if state.EnrollmentString.IsNull() || state.EnrollmentString.IsUnknown() {
		if inlineResp.GetEnrollmentString() == "" {
			state.EnrollmentString = types.StringNull()
		} else {
			state.EnrollmentString = types.StringValue(inlineResp.GetEnrollmentString())
		}

	}

	//    Url
	if state.Url.IsNull() || state.Url.IsUnknown() {
		state.Url = types.StringValue(inlineResp.GetUrl())
	}

	//    Notes
	if state.Notes.IsNull() || state.Notes.IsUnknown() {
		state.Notes = types.StringValue(inlineResp.GetNotes())
	}

	//    IsBoundToConfigTemplate
	if state.IsBoundToConfigTemplate.IsNull() || state.IsBoundToConfigTemplate.IsUnknown() {
		state.IsBoundToConfigTemplate = types.BoolValue(inlineResp.GetIsBoundToConfigTemplate())
	}

	// CopyFromNetworkId
	if state.CopyFromNetworkId.IsNull() || state.CopyFromNetworkId.IsUnknown() {
		state.CopyFromNetworkId = types.StringNull()
	}

	return diags
}
