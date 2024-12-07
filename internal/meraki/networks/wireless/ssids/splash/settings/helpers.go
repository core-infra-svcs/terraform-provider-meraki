package settings

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

func updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx context.Context, state *resourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	imageObjNull := types.ObjectNull(resourceModelImageAttrTypes)

	if !state.SplashImage.IsNull() {
		// Handle SplashImageState
		var splashImageState resourceModelSplashImage
		err := state.SplashImage.As(ctx, &splashImageState, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		// Handle ImageState
		var imageState resourceModelImage
		err = splashImageState.Image.As(ctx, &imageState, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		if !imageState.Contents.IsNull() && !imageState.Contents.IsUnknown() &&
			!imageState.Format.IsNull() && !imageState.Format.IsUnknown() {
			imageObj, err := types.ObjectValueFrom(ctx, resourceModelImageAttrTypes, imageState)
			if err.HasError() {
				diags.Append(err...)
			}
			return imageObj, diags
		}
	}

	return imageObjNull, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceSplashImageState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *resourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var splashImage resourceModelSplashImage
	splashImageObjNull := types.ObjectNull(resourceModelSplashImageAttrTypes)

	// Md5
	if inlineResp.SplashImage.HasMd5() {
		splashImage.Md5 = types.StringValue(inlineResp.SplashImage.GetMd5())
	} else {
		splashImage.Md5 = types.StringNull()
	}

	// Extension
	if inlineResp.SplashImage.HasExtension() {
		splashImage.Extension = types.StringValue(inlineResp.SplashImage.GetExtension())
	} else {
		splashImage.Extension = types.StringNull()
	}

	// Image
	if !state.SplashImage.IsNull() && !state.SplashImage.IsUnknown() {
		imageObj, err := updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx, state)
		if err.HasError() {
			diags.Append(err...)
		}
		splashImage.Image = imageObj

		splashImageObj, err := types.ObjectValueFrom(ctx, resourceModelSplashImageAttrTypes, splashImage)
		if err.HasError() {
			diags.Append(err...)
		}
		return splashImageObj, diags
	}

	return splashImageObjNull, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceSplashLogoState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *resourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var splashLogo resourceModelSplashLogo
	splashLogoObjNull := types.ObjectNull(resourceModelSplashLogoAttrTypes)

	// Md5
	if inlineResp.SplashLogo.HasMd5() {
		splashLogo.Md5 = types.StringValue(inlineResp.SplashLogo.GetMd5())
	} else {
		splashLogo.Md5 = types.StringNull()
	}

	// Extension
	if inlineResp.SplashLogo.HasExtension() {
		splashLogo.Extension = types.StringValue(inlineResp.SplashLogo.GetExtension())
	} else {
		splashLogo.Extension = types.StringNull()
	}

	// Logo
	if !state.SplashLogo.IsNull() && !state.SplashLogo.IsUnknown() {
		imageObj, err := updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx, state)
		if err.HasError() {
			diags.Append(err...)
		}
		splashLogo.Image = imageObj

		splashLogoObj, err := types.ObjectValueFrom(ctx, resourceModelSplashLogoAttrTypes, splashLogo)
		if err.HasError() {
			diags.Append(err...)
		}
		return splashLogoObj, diags
	}

	return splashLogoObjNull, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceSplashPrepaidFrontState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *resourceModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var splashPrepaidFront resourceModelSplashPrepaidFront
	splashPrepaidFrontObjNull := types.ObjectNull(resourceModelSplashPrepaidFrontAttrTypes)

	// Md5
	if inlineResp.SplashPrepaidFront.HasMd5() {
		splashPrepaidFront.Md5 = types.StringValue(inlineResp.SplashPrepaidFront.GetMd5())
	} else {
		splashPrepaidFront.Md5 = types.StringNull()
	}

	// Extension
	if inlineResp.SplashPrepaidFront.HasExtension() {
		splashPrepaidFront.Extension = types.StringValue(inlineResp.SplashPrepaidFront.GetExtension())
	} else {
		splashPrepaidFront.Extension = types.StringNull()
	}

	// PrepaidFront
	if !state.SplashPrepaidFront.IsNull() && !state.SplashPrepaidFront.IsUnknown() {
		imageObj, err := updateNetworksWirelessSsidsSplashSettingsResourceImageState(ctx, state)
		if err.HasError() {
			diags.Append(err...)
		}
		splashPrepaidFront.Image = imageObj

		splashPrepaidFrontObj, err := types.ObjectValueFrom(ctx, resourceModelSplashPrepaidFrontAttrTypes, splashPrepaidFront)
		if err.HasError() {
			diags.Append(err...)
		}
		return splashPrepaidFrontObj, diags
	}

	return splashPrepaidFrontObjNull, diags
}

func updateNetworksWirelessSsidsSplashSettingsResourceState(ctx context.Context, inlineResp openApiClient.GetNetworkWirelessSsidSplashSettings200Response, state *resourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set ID for the new resource.
	if !state.NetworkId.IsNull() && !state.NetworkId.IsUnknown() && !state.Number.IsNull() && !state.Number.IsUnknown() {
		state.Id = types.StringValue(fmt.Sprintf("%s,%s", state.NetworkId.ValueString(), state.Number.ValueString()))
	} else {
		state.Id = types.StringNull()
	}

	// Simple Attributes
	state.SplashUrl = types.StringValue(inlineResp.GetSplashUrl())
	state.UseSplashUrl = types.BoolValue(inlineResp.GetUseSplashUrl())
	state.SplashTimeout = types.Int64Value(int64(inlineResp.GetSplashTimeout()))
	state.WelcomeMessage = types.StringValue(inlineResp.GetWelcomeMessage())
	state.RedirectUrl = types.StringValue(inlineResp.GetRedirectUrl())
	state.AllowSimultaneousLogins = types.BoolValue(inlineResp.GetAllowSimultaneousLogins())
	state.BlockAllTrafficBeforeSignOn = types.BoolValue(inlineResp.GetBlockAllTrafficBeforeSignOn())
	state.ControllerDisconnectionBehavior = types.StringValue(inlineResp.GetControllerDisconnectionBehavior())
	state.UseRedirectUrl = types.BoolValue(inlineResp.GetUseRedirectUrl())

	// Handle null values for Billing
	if inlineResp.Billing != nil {
		var freeAccess resourceModelFreeAccess
		if inlineResp.Billing.FreeAccess != nil {
			freeAccess.DurationInMinutes = types.Int64Value(int64(inlineResp.Billing.FreeAccess.GetDurationInMinutes()))
			freeAccess.Enabled = types.BoolValue(inlineResp.Billing.FreeAccess.GetEnabled())
		}
		freeAccessObject, err := types.ObjectValueFrom(ctx, resourceModelFreeAccessAttrTypes, freeAccess)
		if err.HasError() {
			diags.Append(err...)
		}

		billing := resourceModelBilling{
			ReplyToEmailAddress:           types.StringValue(inlineResp.Billing.GetReplyToEmailAddress()),
			PrepaidAccessFastLoginEnabled: types.BoolValue(inlineResp.Billing.GetPrepaidAccessFastLoginEnabled()),
			FreeAccess:                    freeAccessObject,
		}

		billingObj, err := types.ObjectValueFrom(ctx, resourceModelBillingAttrTypes, billing)
		if err.HasError() {
			diags.Append(err...)
		}

		state.Billing = billingObj
	}

	// Handle null values for GuestSponsorship
	if inlineResp.GuestSponsorship != nil {
		var guestSponsorship resourceModelGuestSponsorship
		guestSponsorship.DurationInMinutes = types.Int64Value(int64(inlineResp.GuestSponsorship.GetDurationInMinutes()))
		guestSponsorship.GuestCanRequestTimeframe = types.BoolValue(inlineResp.GuestSponsorship.GetGuestCanRequestTimeframe())

		guestSponsorshipObj, err := types.ObjectValueFrom(ctx, resourceModelGuestSponsorshipAttrTypes, guestSponsorship)
		if err.HasError() {
			diags.Append(err...)
		}
		state.GuestSponsorship = guestSponsorshipObj
	}

	// Handle null values for SentryEnrollment
	if inlineResp.SentryEnrollment != nil {
		var sentryEnrollment resourceModelSentryEnrollment

		// enforcedSystems obj
		enforcedSystemsObj, err := types.ListValueFrom(ctx, types.StringType, inlineResp.SentryEnrollment.GetEnforcedSystems())
		if err.HasError() {
			diags.Append(err...)
		}
		sentryEnrollment.EnforcedSystems = enforcedSystemsObj

		sentryEnrollment.Strength = types.StringValue(inlineResp.SentryEnrollment.GetStrength())

		// systemsManagerNetwork
		if inlineResp.SentryEnrollment.SystemsManagerNetwork != nil {
			var systemsManagerNetwork resourceModelSystemsManagerNetwork
			systemsManagerNetwork.Id = types.StringValue(inlineResp.SentryEnrollment.SystemsManagerNetwork.GetId())
			systemsManagerNetworkObj, err := types.ObjectValueFrom(ctx, resourceModelSystemsManagerNetworkAttrTypes, systemsManagerNetwork)
			if err.HasError() {
				diags.Append(err...)
			}

			sentryEnrollment.SystemsManagerNetwork = systemsManagerNetworkObj
		}

		sentryEnrollmentObj, err := types.ObjectValueFrom(ctx, resourceModelSentryEnrollmentAttrTypes, sentryEnrollment)
		if err.HasError() {
			diags.Append(err...)
		}
		state.SentryEnrollment = sentryEnrollmentObj
	}

	// Handle null values for SplashImage
	splashImage, err := updateNetworksWirelessSsidsSplashSettingsResourceSplashImageState(ctx, inlineResp, state)
	if err.HasError() {
		diags.Append(err...)
	}
	state.SplashImage = splashImage

	// Handle null values for SplashLogo
	splashLogo, err := updateNetworksWirelessSsidsSplashSettingsResourceSplashLogoState(ctx, inlineResp, state)
	if err.HasError() {
		diags.Append(err...)
	}
	state.SplashLogo = splashLogo

	// Handle null values for SplashPrepaidFront
	splashPrepaidFront, err := updateNetworksWirelessSsidsSplashSettingsResourceSplashPrepaidFrontState(ctx, inlineResp, state)
	if err.HasError() {
		diags.Append(err...)
	}
	state.SplashPrepaidFront = splashPrepaidFront

	return diags
}

func NetworksWirelessSsidsSplashSettingsResourcePayload(ctx context.Context, data *resourceModel) (openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload := *openApiClient.NewUpdateNetworkWirelessSsidSplashSettingsRequest()
	payload.SetSplashUrl(data.SplashUrl.ValueString())
	payload.SetUseSplashUrl(data.UseSplashUrl.ValueBool())
	payload.SetSplashTimeout(int32(data.SplashTimeout.ValueInt64()))
	payload.SetWelcomeMessage(data.WelcomeMessage.ValueString())
	payload.SetRedirectUrl(data.RedirectUrl.ValueString())
	payload.SetUseRedirectUrl(data.UseRedirectUrl.ValueBool())
	payload.SetAllowSimultaneousLogins(data.AllowSimultaneousLogins.ValueBool())
	payload.SetControllerDisconnectionBehavior(data.ControllerDisconnectionBehavior.ValueString())
	payload.SetBlockAllTrafficBeforeSignOn(data.BlockAllTrafficBeforeSignOn.ValueBool())

	// Billing
	if !data.Billing.IsUnknown() && !data.Billing.IsNull() {
		var billing openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestBilling
		var billingData resourceModelBilling

		err := data.Billing.As(ctx, &billingData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		billing.SetReplyToEmailAddress(billingData.ReplyToEmailAddress.ValueString())
		billing.SetPrepaidAccessFastLoginEnabled(billingData.PrepaidAccessFastLoginEnabled.ValueBool())

		if !billingData.FreeAccess.IsUnknown() && !billingData.FreeAccess.IsNull() {
			var freeAccess openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestBillingFreeAccess
			var freeAccessData resourceModelFreeAccess

			err := billingData.FreeAccess.As(ctx, &freeAccessData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			freeAccess.SetDurationInMinutes(int32(freeAccessData.DurationInMinutes.ValueInt64()))
			freeAccess.SetEnabled(freeAccessData.Enabled.ValueBool())
			billing.SetFreeAccess(freeAccess)
		}
		payload.SetBilling(billing)
	}

	// Guest Sponsorship
	if !data.GuestSponsorship.IsUnknown() && !data.GuestSponsorship.IsNull() {
		var guestSponsorship openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestGuestSponsorship
		var guestSponsorshipData resourceModelGuestSponsorship

		err := data.GuestSponsorship.As(ctx, &guestSponsorshipData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		guestSponsorship.SetDurationInMinutes(int32(guestSponsorshipData.DurationInMinutes.ValueInt64()))
		guestSponsorship.SetGuestCanRequestTimeframe(guestSponsorshipData.GuestCanRequestTimeframe.ValueBool())
		payload.SetGuestSponsorship(guestSponsorship)
	}

	// Sentry Enrollment
	if !data.SentryEnrollment.IsUnknown() && !data.SentryEnrollment.IsNull() {
		var sentryEnrollment openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSentryEnrollment
		var sentryEnrollmentData resourceModelSentryEnrollment

		err := data.SentryEnrollment.As(ctx, &sentryEnrollmentData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		sentryEnrollment.SetStrength(sentryEnrollmentData.Strength.ValueString())

		var enforcedSystems []string
		for _, v := range sentryEnrollmentData.EnforcedSystems.Elements() {
			enforcedSystems = append(enforcedSystems, v.(types.String).ValueString())
		}
		sentryEnrollment.SetEnforcedSystems(enforcedSystems)

		if !sentryEnrollmentData.SystemsManagerNetwork.IsUnknown() && !sentryEnrollmentData.SystemsManagerNetwork.IsNull() {
			var systemsManagerNetwork openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSentryEnrollmentSystemsManagerNetwork
			var systemsManagerNetworkData resourceModelSystemsManagerNetwork

			err := sentryEnrollmentData.SystemsManagerNetwork.As(ctx, &systemsManagerNetworkData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			systemsManagerNetwork.SetId(systemsManagerNetworkData.Id.ValueString())
			sentryEnrollment.SetSystemsManagerNetwork(systemsManagerNetwork)
		}
		payload.SetSentryEnrollment(sentryEnrollment)
	}

	// Splash Image
	if !data.SplashImage.IsUnknown() && !data.SplashImage.IsNull() {
		var splashImage openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashImage
		var splashImageData resourceModelSplashImage

		err := data.SplashImage.As(ctx, &splashImageData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		splashImage.SetExtension(splashImageData.Extension.ValueString())
		splashImage.SetMd5(splashImageData.Md5.ValueString())

		if !splashImageData.Image.IsUnknown() && !splashImageData.Image.IsNull() {
			var image openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashImageImage
			var imageData resourceModelImage

			err := splashImageData.Image.As(ctx, &imageData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			image.SetContents(imageData.Contents.ValueString())
			image.SetFormat(imageData.Format.ValueString())
			splashImage.SetImage(image)
		}
		payload.SetSplashImage(splashImage)
	}

	// Splash Logo
	if !data.SplashLogo.IsUnknown() && !data.SplashLogo.IsNull() {
		var splashLogo openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashLogo
		var splashLogoData resourceModelSplashLogo

		err := data.SplashLogo.As(ctx, &splashLogoData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		splashLogo.SetExtension(splashLogoData.Extension.ValueString())
		splashLogo.SetMd5(splashLogoData.Md5.ValueString())

		if !splashLogoData.Image.IsUnknown() && !splashLogoData.Image.IsNull() {
			var imageLogo openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashLogoImage
			var imageLogoData resourceModelImage

			err := splashLogoData.Image.As(ctx, &imageLogoData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			imageLogo.SetContents(imageLogoData.Contents.ValueString())
			imageLogo.SetFormat(imageLogoData.Format.ValueString())
			splashLogo.SetImage(imageLogo)
		}

		payload.SetSplashLogo(splashLogo)
	}

	// Splash Prepaid Front
	if !data.SplashPrepaidFront.IsUnknown() && !data.SplashPrepaidFront.IsNull() {
		var splashPrepaidFront openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashPrepaidFront
		var splashPrepaidFrontData resourceModelSplashPrepaidFront

		err := data.SplashPrepaidFront.As(ctx, &splashPrepaidFrontData, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags.Append(err...)
		}

		splashPrepaidFront.SetExtension(splashPrepaidFrontData.Extension.ValueString())
		splashPrepaidFront.SetMd5(splashPrepaidFrontData.Md5.ValueString())

		if !splashPrepaidFrontData.Image.IsUnknown() && !splashPrepaidFrontData.Image.IsNull() {
			var imagePrepaidFront openApiClient.UpdateNetworkWirelessSsidSplashSettingsRequestSplashPrepaidFrontImage
			var imagePrepaidFrontData resourceModelImage

			err := splashPrepaidFrontData.Image.As(ctx, &imagePrepaidFrontData, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags.Append(err...)
			}

			imagePrepaidFront.SetContents(imagePrepaidFrontData.Contents.ValueString())
			imagePrepaidFront.SetFormat(imagePrepaidFrontData.Format.ValueString())
			splashPrepaidFront.SetImage(imagePrepaidFront)
		}
		payload.SetSplashPrepaidFront(splashPrepaidFront)
	}

	return payload, diags
}
