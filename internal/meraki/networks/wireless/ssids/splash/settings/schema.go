package settings

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksWirelessSsidsSplashSettings",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"number": schema.StringAttribute{
				MarkdownDescription: "SsIds Number",
				Required:            true,
			},
			"splash_url": schema.StringAttribute{
				MarkdownDescription: "The custom splash URL of the click-through splash page. Note that the URL can be configured without necessarily being used. In order to enable the custom URL, see 'useSplashUrl'",
				Optional:            true,
				Computed:            true,
			},
			"use_splash_url": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.",
				Optional:            true,
				Computed:            true,
			},
			"splash_timeout": schema.Int64Attribute{
				MarkdownDescription: "Splash timeout in minutes. This will determine how often users will see the splash page.",
				Optional:            true,
				Computed:            true,
			},
			"welcome_message": schema.StringAttribute{
				MarkdownDescription: "The welcome message for the users on the splash page.",
				Optional:            true,
				Computed:            true,
			},
			"redirect_url": schema.StringAttribute{
				MarkdownDescription: "The custom redirect URL where the users will go after the splash page.",
				Optional:            true,
				Computed:            true,
			},
			"use_redirect_url": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether the users will be redirected to the custom splash url. A custom splash URL must be set if this is true. Note that depending on your SSID's access control settings, it may not be possible to use the custom splash URL.",
				Optional:            true,
				Computed:            true,
			},
			"block_all_traffic_before_sign_on": schema.BoolAttribute{
				MarkdownDescription: "How restricted allowing traffic should be. If true, all traffic types are blocked until the splash page is acknowledged. If false, all non-HTTP traffic is allowed before the splash page is acknowledged.",
				Optional:            true,
				Computed:            true,
			},
			"controller_disconnection_behavior": schema.StringAttribute{
				MarkdownDescription: "How login attempts should be handled when the controller is unreachable. Can be either 'open', 'restricted', or 'default'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("open", "restricted", "default"),
				},
			},
			"allow_simultaneous_logins": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to allow simultaneous logins from different devices.",
				Optional:            true,
				Computed:            true,
			},
			"billing": schema.SingleNestedAttribute{
				MarkdownDescription: "Details associated with billing splash.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"reply_to_email_address": schema.StringAttribute{
						MarkdownDescription: "The email address that receives replies from clients.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"prepaid_access_fast_login_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not billing uses the fast login prepaid access option.",
						Optional:            true,
						Computed:            true,
					},
					"free_access": schema.SingleNestedAttribute{
						MarkdownDescription: "Details associated with a free access plan with limits.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"duration_in_minutes": schema.Int64Attribute{
								MarkdownDescription: "How long a device can use a network for free..",
								Optional:            true,
								Computed:            true,
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether or not free access is enabled.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
				},
			},
			"guest_sponsorship": schema.SingleNestedAttribute{
				MarkdownDescription: "Details associated with guest sponsored splash.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"duration_in_minutes": schema.Int64Attribute{
						MarkdownDescription: "Duration in minutes of sponsored guest authorization. Must be between 1 and 60480 (6 weeks).",
						Optional:            true,
						Computed:            true,
					},
					"guest_can_request_time_frame": schema.BoolAttribute{
						MarkdownDescription: "Whether or not guests can specify how much time they are requesting.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
			"sentry_enrollment": schema.SingleNestedAttribute{
				MarkdownDescription: "Systems Manager sentry enrollment splash settings.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"strength": schema.StringAttribute{
						MarkdownDescription: "The strength of the enforcement of selected system types. Must be one of: 'focused', 'click-through', and 'strict'.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("click-through", "focused", "strict"),
						},
					},
					"enforced_systems": schema.ListAttribute{
						MarkdownDescription: "The system types that the Sentry enforces. Must be included in: 'iOS, 'Android', 'macOS', and 'Windows'.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
					"systems_manager_network": schema.SingleNestedAttribute{
						MarkdownDescription: "Systems Manager network targeted for sentry enrollment..",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "The network ID of the Systems Manager network.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"splash_image": schema.SingleNestedAttribute{
				MarkdownDescription: "The image used in the splash page.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new image.",
								Optional:            true,
								Computed:            true,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("gif", "jpg", "png"),
								},
							},
						},
					},
				},
			},
			"splash_logo": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new image.",
								Optional:            true,
								Computed:            true,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("gif", "jpg", "png"),
								},
							},
						},
					},
				},
			},
			"splash_prepaid_front": schema.SingleNestedAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension": schema.StringAttribute{
						MarkdownDescription: "The extension of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"md5": schema.StringAttribute{
						MarkdownDescription: "The MD5 value of the image file.",
						Optional:            true,
						Computed:            true,
					},
					"image": schema.SingleNestedAttribute{
						MarkdownDescription: "Properties for setting a new image.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The file contents (a base 64 encoded string) of your new image.",
								Optional:            true,
								Computed:            true,
							},
							"format": schema.StringAttribute{
								MarkdownDescription: "The format of the encoded contents. Supported formats are 'png', 'gif', and jpg'.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("gif", "jpg", "png"),
								},
							},
						},
					},
				},
			},
		},
	}
}
