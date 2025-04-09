// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	rangeTypeWithin  = "WITHIN"
	rangeTypeOutside = "OUTSIDE"

	presencePresent = "PRESENT"
	presenceAbsent  = "ABSENT"

	motionDetected   = "MOTION_DETECTED"
	noMotionDetected = "NO_MOTION_DETECTED"

	occupancyOccupied    = "OCCUPIED"
	occupancyNotOccupied = "NOT_OCCUPIED"

	connectionCconOffline   = "CLOUD_CONNECTOR_OFFLINE"
	connectionSensorOffline = "SENSOR_OFFLINE"

	contactOpen  = "OPEN"
	contactClose = "CLOSED"

	notificationActionSMS            = "SMS"
	notificationActionEmail          = "EMAIL"
	notificationActionCorrigo        = "CORRIGO"
	notificationActionServiceChannel = "SERVICE_CHANNEL"
	notificationActionWebhook        = "WEBHOOK"
	notificationActionPhoneCall      = "PHONE_CALL"
	notificationActionSignalTower    = "SIGNAL_TOWER"

	dayMonday    = "Monday"
	dayTuesday   = "Tuesday"
	dayWednesday = "Wednesday"
	dayThursday  = "Thursday"
	dayFriday    = "Friday"
	daySaturday  = "Saturday"
	daySunday    = "Sunday"
)

var (
	validTriggerFields = []string{
		"temperature",
		"co2",
		"relativeHumidity",
		"waterPresent",
		"contact",
		"motion",
		"touch",
		"objectPresent",
		"deskOccupancy",
		"connectionStatus",
	}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &notificationRuleResource{}
	_ resource.ResourceWithConfigure = &notificationRuleResource{}
)

// NewDataConnectorResource is a helper function to simplify the provider implementation.
func NewNotificationRuleResource() resource.Resource {
	return &notificationRuleResource{}
}

// notificationRuleResource is the resource implementation
type notificationRuleResource struct {
	client *dt.Client
}

// Metadata returns the resource type name.
func (r *notificationRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule"
}

// Schema defines the schema for the resource.
func (r *notificationRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
				Description: `The resource name of the rule. The resource name has the following format: "projects/{project_id}/rules/{rule_id}". 
							 	The name is ignored when creating a new rule.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether or not the rule is enabled.",
				Default:     booldefault.StaticBool(true),
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "he display name of the rule that is visible in Studio.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The DT project ID of the rule.",
			},
			"devices": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: `An optional list of device resource names that this rule applies to.
								If the list is empty, the rule applies to all devices in the project,
								or those matching all the labels in device_labels (if present).`,
				// Defaults to empty list
				Default:    listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.List{listvalidator.SizeAtLeast(1)},
			},
			"device_labels": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: `An optional map of labels to use as a filter for which devices this rule applies to.
								This applies regardless of whether or not the devices field is set. The map can contain
								both label key/value pairs, or just label keys. If multiple labels are specified, the
								device must match all of them to be included.`,
				// Defaults to empty list
				Default:    mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				Validators: []validator.Map{mapvalidator.SizeAtLeast(1)},
			},
			"trigger": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The condition that that needs to be met before the actions are executed.",
				Attributes: map[string]schema.Attribute{
					"field": schema.StringAttribute{
						Required:    true,
						Description: "The data field to use for the criteria.",
						Validators:  []validator.String{stringvalidator.OneOf(validTriggerFields...)},
					},
					"range": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The range of values that the data field has to be within or without for the criteria to be met.",
						Attributes: map[string]schema.Attribute{
							"lower": schema.Float64Attribute{
								Optional:    true,
								Description: "The minimum value of the range.",
							},
							"upper": schema.Float64Attribute{
								Optional:    true,
								Description: "The maximum value of the range.",
							},
							"type": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Description: fmt.Sprintf(
									"The type of range to use. Must be one of `%s` or `%s`. Default is `%s`.",
									rangeTypeWithin,
									rangeTypeOutside,
									rangeTypeWithin),
								Validators: []validator.String{stringvalidator.OneOf([]string{rangeTypeWithin, rangeTypeOutside}...)},
								Default:    stringdefault.StaticString(rangeTypeWithin),
							},
						},
					},
					"presence": schema.StringAttribute{
						Optional: true,
						Description: fmt.Sprintf(
							"The presence of a object or person detected by a sensor. Must be one of `%s` or `%s`. Default is %s",
							presencePresent,
							presenceAbsent,
							presencePresent),
						Validators: []validator.String{stringvalidator.OneOf([]string{presencePresent, presenceAbsent}...)},
					},
					"motion": schema.StringAttribute{
						Optional: true,
						Description: fmt.Sprintf(
							"The motion detected by a sensor. Must be one of `%s` or `%s`. Default is %s",
							motionDetected,
							noMotionDetected,
							motionDetected),
						Validators: []validator.String{stringvalidator.OneOf([]string{motionDetected, noMotionDetected}...)},
					},
					"occupancy": schema.StringAttribute{
						Optional: true,
						Description: fmt.Sprintf(
							"Occupancy detected by a sensor. Must be one of `%s` or `%s`. Default is %s",
							occupancyOccupied,
							occupancyNotOccupied,
							occupancyOccupied,
						),
						Validators: []validator.String{stringvalidator.OneOf([]string{occupancyOccupied, occupancyNotOccupied}...)},
					},
					"connection": schema.StringAttribute{
						Optional: true,
						Description: fmt.Sprintf(
							"The connection status of a device can be set to either `%s`for cloud connector offline or `%s` for sensor offline. Default is `%s`",
							connectionCconOffline,
							connectionSensorOffline,
							connectionCconOffline,
						),
						Validators: []validator.String{stringvalidator.OneOf([]string{connectionCconOffline, connectionSensorOffline}...)},
					},
					"contact": schema.StringAttribute{
						Optional: true,
						Description: fmt.Sprintf(
							"The open status of a contact(door and window) sensor. Must be one of `%s` or `%s`. Default is %s",
							contactOpen,
							contactClose,
							contactOpen,
						),
						Validators: []validator.String{stringvalidator.OneOf([]string{contactOpen, contactClose}...)},
					},
					"trigger_count": schema.Int32Attribute{
						Computed: true,
						Description: `The number of times a device has to meet the trigger criteria (enter triggering mode) 
											before a notification is published. The value has to be greater than 1 for this feature
											to be "enabled". A value of 0 or 1 is considered equivalent.
											Note that this feature can't be used with trigger delay, reminder notifications or 
											resolved notifications.`,
						Default: int32default.StaticInt32(1),
					},
				},
			},
			"escalation_levels": schema.ListNestedAttribute{
				Optional: true,
				Description: ` A list of escalation levels that will be used throughout the lifecycle of alerts
    							that are created by this rule. The first escalation level will be used when the
    							alert is triggered, and the next escalation level will be used when the alert is
    							escalated, and so on. Each escalation level needs at least one action, and there
    							needs to be at least one escalation level.`,
				Validators: []validator.List{listvalidator.SizeAtLeast(1)},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Required:    true,
							Description: "he display namse of the escalation level that is visible in Studio.",
						},
						"actions": schema.ListNestedAttribute{
							Required:     true,
							Description:  "The list of actions that will be executed once this escalation level is reached.",
							NestedObject: notificationAction,
							Validators:   []validator.List{listvalidator.SizeAtLeast(1)},
						},
						"escalate_after": schema.StringAttribute{
							Optional: true,
							Description: ` 	The amount of time to wait before escalating to the next escalation level. Only
    										relevant if there are more escalation levels that follows this one.`,
							Validators: []validator.String{durationValidator},
						},
					},
				},
			},
			"schedule": schema.SingleNestedAttribute{
				Optional: true,
				Description: `A schedule limits at what times the rule will be evaluated, and events will be processed. 
								When an event is received outside the schedule, the device will never be put on the delay
								queue, and a trigger counter (if enabled) will not be incremented.`,
				Attributes: map[string]schema.Attribute{
					"timezone": schema.StringAttribute{
						Required: true,
						Description: `The timezone for which the schedule applies. This will automatically handle DST if the correct zones are used.
										E.g. "Europe/Oslo", "America/Los_Angeles", "UTC"
										See https://en.wikipedia.org/wiki/List_of_tz_database_time_zones`,
					},
					"slots": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Slots of time where the rule should be active or inactive.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"day_of_week": schema.ListAttribute{
									Required: true,
									Description: fmt.Sprintf(
										"Day of week for the slot. Must be one of %s,%s,%s,%s,%s,%s,%s, empty means every day",
										dayMonday,
										dayTuesday,
										dayWednesday,
										dayThursday,
										dayFriday,
										daySaturday,
										daySunday,
									),
									ElementType: types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(
											stringvalidator.OneOf(
												dayMonday,
												dayTuesday,
												dayWednesday,
												dayThursday,
												dayFriday,
												daySaturday,
												daySunday,
											),
										),
									},
								},
								"time_range": schema.ListNestedAttribute{
									Optional:    true,
									Description: "Ranges of time where the rule should be active or inactive.",
									Validators:  []validator.List{listvalidator.SizeAtLeast(1)},
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"start": schema.SingleNestedAttribute{
												Required:    true,
												Description: "The start time of the slot.",
												Attributes: map[string]schema.Attribute{
													"hour": schema.Int32Attribute{
														Required:    true,
														Description: "The hour of the slot.",
													},
													"minute": schema.Int32Attribute{
														Required:    true,
														Description: "The minute of the slot.",
													},
												},
											},
											"end": schema.SingleNestedAttribute{
												Required:    true,
												Description: "The end time of the slot in HH:MM format.",
												Attributes: map[string]schema.Attribute{
													"hour": schema.Int32Attribute{
														Required:    true,
														Description: "The hour of the slot.",
													},
													"minute": schema.Int32Attribute{
														Required:    true,
														Description: "The minute of the slot.",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"trigger_delay": schema.StringAttribute{
				Optional: true,
				Description: `The amount of time to wait before executing the actions. This is useful to avoid
    							sending notifications for short-lived conditions, or when a notification is only desired
    							if a condition has been met for an extended period of time (eg. fridge temp above a certain
    							value for 1h, or a door open for 30 minutes). If this is not set, the actions will be
    							executed immediately.`,
				Validators: []validator.String{durationValidator},
			},
			"reminder_notification": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether or not to send a reminder notifications",
				Default:     booldefault.StaticBool(false),
			},
			"resolved_notification": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether or not to send a resolved notifications",
				Default:     booldefault.StaticBool(false),
			},
			"unacknowledge_after": schema.StringAttribute{
				Optional: true,
				Description: ` The amount of time before an acknowledged alert created by this rule is unacknowledged.
								When an alert is unacknowledged, the notifications for its current escalation level
								will be sent out once again, and the escalation timer will start. The default value
								is 4 hours.`,
				Validators: []validator.String{durationValidator},
			},
			"actions": schema.ListNestedAttribute{
				Optional:           true,
				DeprecationMessage: "Use the `escalation_levels` attribute instead.",
				Description:        "The list of actions that will be executed when the trigger is met.",
				NestedObject:       notificationAction,
			},
		},
	}
}

var notificationAction = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
			Description: fmt.Sprintf(
				"The type of notification action. Must be one of %s,%s,%s,%s,%s,%s,%s, default is %s",
				notificationActionSMS,
				notificationActionEmail,
				notificationActionCorrigo,
				notificationActionServiceChannel,
				notificationActionWebhook,
				notificationActionPhoneCall,
				notificationActionSignalTower,
				notificationActionEmail,
			),
		},
		"sms_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for sending an SMS notification.",
			Attributes: map[string]schema.Attribute{
				"recipients": schema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
					Description: "The phone numbers to send the SMS to.",
				},
				"body": schema.StringAttribute{
					Required:    true,
					Description: "The message to send in the SMS.",
				},
			},
		},
		"email_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for sending an email notification.",
			Attributes: map[string]schema.Attribute{
				"recipients": schema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
					Description: "The email addresses to send the email to.",
				},
				"subject": schema.StringAttribute{
					Required:    true,
					Description: "The subject of the email.",
				},
				"body": schema.StringAttribute{
					Required:    true,
					Description: "The body of the email.",
				},
			},
		},
		"corrigo_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for creating a Corrigo work order.",
			Attributes: map[string]schema.Attribute{
				"asset_id": schema.StringAttribute{
					Required:    true,
					Description: "The asset ID of the device.",
				},
				"task_id": schema.StringAttribute{
					Required:    true,
					Description: "The task ID of the device.",
				},
				"customer_id": schema.StringAttribute{
					Required:    true,
					Description: "The customer ID of the device.",
				},
				"client_secret": schema.StringAttribute{
					Required:    true,
					Description: "The client secret of the device.",
					Sensitive:   true,
				},
				"company_name": schema.StringAttribute{
					Required:    true,
					Description: "The company name of the device.",
				},
				"sub_type_id": schema.StringAttribute{
					Required:    true,
					Description: "The sub type ID of the device.",
				},
				"contact_name": schema.StringAttribute{
					Required:    true,
					Description: "The contact name of the device.",
				},
				"contact_address": schema.StringAttribute{
					Required:    true,
					Description: "The contact address of the device.",
				},
				"work_order_description": schema.StringAttribute{
					Optional: true,
					Description: `Optional field that will populate the description of the work order. 
    								If this is not specified, a default value will be used.`,
				},
				"studio_dashboard_url": schema.StringAttribute{
					Optional: true,
					Description: `Optional field to allow users to set the Studio dashboard link that
    								should be included in the Corrigo Work Order. If this is not specified,
    								the defaultx (initial) dashboard will be used in the link.`,
				},
			},
		},
		"service_channel_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for sending a notification to a service channel.",
			Attributes: map[string]schema.Attribute{
				"store_id": schema.StringAttribute{
					Optional: true,
					Description: `A 4-digit number that identifies the store the asset is in.
    								If this is not set, it will be derived from the "service_channel_store_id"
    								label in Studio.`,
				},
				"asset_tag_id": schema.StringAttribute{
					Optional: true,
					Description: ` tag that is set on the asset in ServiceChannel. Will be used to find the 
									asset, and link it to work orders. 
									When specified, the trade of the asset will be used for the work order instead 
									of the provided asset. If this is not set, or if an asset couldn't be found 
									based on this tag, the provided trade will be used instead to derive fields 
									dependent on the trade. If this is not set, it will be derived from the 
									"service_channel_asset_tag_id" label in Studio.`,
				},
				"trade": schema.StringAttribute{
					Optional: true,
					Description: `The trade to use if the asset tag id is either not specified or no matching
									asset could be found. If this is not set, it will be derived from the 
									"service_channel_trade" label in Studio.
									Examples of a trade could be "REFRIGERATION" or "HOT FOOD".`,
				},
				"description": schema.StringAttribute{
					Optional: true,
					Description: `The description that will appear on the work order. If this is not set, 
									it will be derived from the "service_channel_work_order_description"
									label in Studio.`,
				},
			},
		},
		"webhook_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for sending a webhook notification.",
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					Required:    true,
					Description: "Valid URL using HTTPS.",
				},
				"signature_secret": schema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					Description: "Use a custom secret to sign the data.",
				},
				"headers": schema.MapAttribute{
					Optional:    true,
					Description: "The headers to include in the webhook request.",
					ElementType: types.StringType,
				},
			},
		},
		"phone_call_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for sending a phone call notification.",
			Attributes: map[string]schema.Attribute{
				"recipients": schema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
					Description: "A list of the phone numbers to call. Must be in E.164 format.",
				},
				"introduction": schema.StringAttribute{
					Required: true,
					Description: `Used to introduce the call to the callee.
    								Example: "This is an automated call from Disruptive Technologies.`,
				},
				"message": schema.StringAttribute{
					Required: true,
					Description: `The message that should be read to the callee.
    								Example: "The temperature in the fridge is above the threshold, currently at $celsius.`,
				},
			},
		},
		"signal_tower_config": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The configuration for sending a signal tower notification.",
			Attributes: map[string]schema.Attribute{
				"cloud_connector_name": schema.StringAttribute{
					Required: true,
					Description: `The resource name of the Cloud Connector that has the signal tower
									connected to it. We currently only support one signal tower per
									Cloud Connector.`,
				},
			},
		},
	},
}

// Data model
type notificationRuleModel struct {
	Name                 types.String              `tfsdk:"name"`
	Enabled              types.Bool                `tfsdk:"enabled"`
	DisplayName          types.String              `tfsdk:"display_name"`
	ProjectID            types.String              `tfsdk:"project_id"`
	Devices              types.List                `tfsdk:"devices"`
	DeviceLabels         types.Map                 `tfsdk:"device_labels"`
	Trigger              triggerModel              `tfsdk:"trigger"`
	EscalationLevels     []escalationLevelModel    `tfsdk:"escalation_levels"`
	Schedule             *scheduleModel            `tfsdk:"schedule"`
	TriggerDelay         types.String              `tfsdk:"trigger_delay"`
	ReminderNotification types.Bool                `tfsdk:"reminder_notification"`
	ResolvedNotification types.Bool                `tfsdk:"resolved_notification"`
	UnacknowledgeAfter   types.String              `tfsdk:"unacknowledge_after"`
	Actions              []notificationActionModel `tfsdk:"actions"`
}

type escalationLevelModel struct {
	DisplayName   types.String              `tfsdk:"display_name"`
	Actions       []notificationActionModel `tfsdk:"actions"`
	EscalateAfter types.String              `tfsdk:"escalate_after"`
}

type notificationActionModel struct {
	Type                 types.String               `tfsdk:"type"`
	SMSConfig            *smsConfigModel            `tfsdk:"sms_config"`
	EmailConfig          *emailConfigModel          `tfsdk:"email_config"`
	CorrigoConfig        *corrigoConfigModel        `tfsdk:"corrigo_config"`
	ServiceChannelConfig *serviceChannelConfigModel `tfsdk:"service_channel_config"`
	WebhookConfig        *webhookConfigModel        `tfsdk:"webhook_config"`
	PhoneCallConfig      *phoneCallConfigModel      `tfsdk:"phone_call_config"`
	SignalTowerConfig    *signalTowerConfigModel    `tfsdk:"signal_tower_config"`
}

type smsConfigModel struct {
	Recipients types.List   `tfsdk:"recipients"`
	Body       types.String `tfsdk:"body"`
}

type emailConfigModel struct {
	Recipients types.List   `tfsdk:"recipients"`
	Subject    types.String `tfsdk:"subject"`
	Body       types.String `tfsdk:"body"`
}

type corrigoConfigModel struct {
	AssetID              types.String `tfsdk:"asset_id"`
	TaskID               types.String `tfsdk:"task_id"`
	CustomerID           types.String `tfsdk:"customer_id"`
	ClientSecret         types.String `tfsdk:"client_secret"`
	CompanyName          types.String `tfsdk:"company_name"`
	SubTypeID            types.String `tfsdk:"sub_type_id"`
	ContactName          types.String `tfsdk:"contact_name"`
	ContactAddress       types.String `tfsdk:"contact_address"`
	WorkOrderDescription types.String `tfsdk:"work_order_description"`
	StudioDashboardURL   types.String `tfsdk:"studio_dashboard_url"`
}

type serviceChannelConfigModel struct {
	StoreID     types.String `tfsdk:"store_id"`
	AssetTagID  types.String `tfsdk:"asset_tag_id"`
	Trade       types.String `tfsdk:"trade"`
	Description types.String `tfsdk:"description"`
}

type webhookConfigModel struct {
	URL             types.String `tfsdk:"url"`
	SignatureSecret types.String `tfsdk:"signature_secret"`
	Headers         types.Map    `tfsdk:"headers"`
}

type phoneCallConfigModel struct {
	Recipients   types.List   `tfsdk:"recipients"`
	Introduction types.String `tfsdk:"introduction"`
	Message      types.String `tfsdk:"message"`
}

type signalTowerConfigModel struct {
	CloudConnectorName types.String `tfsdk:"cloud_connector_name"`
}

type triggerModel struct {
	Field        types.String `tfsdk:"field"`
	Range        *rangeModel  `tfsdk:"range"`
	Presence     types.String `tfsdk:"presence"`
	Motion       types.String `tfsdk:"motion"`
	Occupancy    types.String `tfsdk:"occupancy"`
	Connection   types.String `tfsdk:"connection"`
	Contact      types.String `tfsdk:"contact"`
	TriggerCount types.Int32  `tfsdk:"trigger_count"`
}

type rangeModel struct {
	Lower types.Float64 `tfsdk:"lower"`
	Upper types.Float64 `tfsdk:"upper"`
	Type  types.String  `tfsdk:"type"`
}

type scheduleModel struct {
	Timezone types.String `tfsdk:"timezone"`
	Slots    []slotsModel `tfsdk:"slots"`
}

type slotsModel struct {
	DayOfWeek []types.String   `tfsdk:"day_of_week"`
	TimeRange []timeRangeModel `tfsdk:"time_range"`
}

type timeRangeModel struct {
	Start timeOfDayModel `tfsdk:"start"`
	End   timeOfDayModel `tfsdk:"end"`
}

type timeOfDayModel struct {
	Hour   types.Int32 `tfsdk:"hour"`
	Minute types.Int32 `tfsdk:"minute"`
}

// Create creates the resource and sets the initial Terraform state.s
func (r *notificationRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the data from the request
	var plan notificationRuleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the data to the dt.NotificationRule
	toBeCreated, diags := stateToNotificationRule(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the notification rule
	created, err := r.client.CreateNotificationRule(ctx, plan.ProjectID.ValueString(), toBeCreated)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification rule",
			fmt.Sprintf("Could not create notification rule: %s", err),
		)
		return
	}

	// Convert the created notification rule to the state model
	state, diags := notificationRuleToState(ctx, created)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *notificationRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve the data from the request
	var state notificationRuleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the notification rule
	notificationRule, err := r.client.GetNotificationRule(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading notification rule",
			fmt.Sprintf("Could not read notification rule: %s", err),
		)
		return
	}

	// Convert the notification rule to the state model
	state, diags = notificationRuleToState(ctx, notificationRule)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *notificationRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the data from the request
	var state notificationRuleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the notification rule
	err := r.client.DeleteNotificationRule(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting notification rule",
			fmt.Sprintf("Could not delete notification rule: %s", err),
		)
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *notificationRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the data from the request
	var plan notificationRuleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the data to the dt.NotificationRule
	toBeUpdated, diags := stateToNotificationRule(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the notification rule
	updated, err := r.client.UpdateNotificationRule(ctx, toBeUpdated)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating notification rule",
			fmt.Sprintf("Could not update notification rule: %s", err),
		)
		return
	}

	// Convert the updated notification rule to the state model
	state, diags := notificationRuleToState(ctx, updated)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *notificationRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dt.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Type",
			fmt.Sprintf("Expected *dt.Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

// NotificationRuleResource converts the dt.NotificationRule to the state model.
func notificationRuleToState(ctx context.Context, notificationRule dt.NotificationRule) (notificationRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	devicesList, d := types.ListValueFrom(ctx, types.StringType, notificationRule.Devices)
	diags = append(diags, d...)

	deviceLabelsMap, d := types.MapValueFrom(ctx, types.StringType, notificationRule.DeviceLabels)
	diags = append(diags, d...)

	escalationLevels, d := escalationLevelToState(ctx, notificationRule.EscalationLevels)
	diags = append(diags, d...)

	actions, d := notificationActionToState(ctx, notificationRule.Actions) // nolint: staticcheck // Actions is deprecated, but included for completeness.
	diags = append(diags, d...)

	// Convert the notification rule to the state model.
	var state notificationRuleModel

	state.Name = types.StringValue(notificationRule.Name)

	projectID, _, err := dt.ParseResourceName(notificationRule.Name)
	if err != nil {
		diags.AddError(
			"Error parsing notification rule name",
			fmt.Sprintf("Could not parse notification rule name: %s", err),
		)
		return state, diags
	}
	state.ProjectID = types.StringValue(projectID)
	state.Enabled = types.BoolValue(notificationRule.Enabled)
	state.DisplayName = types.StringValue(notificationRule.DisplayName)
	state.Devices = devicesList
	state.DeviceLabels = deviceLabelsMap
	state.Trigger = triggerToState(notificationRule.Trigger)
	state.EscalationLevels = escalationLevels
	state.Schedule = scheduleToState(notificationRule.Schedule)

	// Set the trigger delay to null if it is not set.
	state.TriggerDelay = types.StringNull()
	if notificationRule.TriggerDelay != nil {
		state.TriggerDelay = types.StringValue(*notificationRule.TriggerDelay)
	}

	state.ReminderNotification = types.BoolValue(notificationRule.ReminderNotification)
	state.ResolvedNotification = types.BoolValue(notificationRule.ResolvedNotification)

	// Set the unacknowledge after to null if it is not set.
	state.UnacknowledgeAfter = types.StringNull()
	if notificationRule.UnacknowledgesAfter != nil {
		state.UnacknowledgeAfter = types.StringValue(*notificationRule.UnacknowledgesAfter)
	}
	state.Actions = actions
	return state, diags
}

// escalationLevelToState converts the dt.EscalationLevel to the state model.
func escalationLevelToState(ctx context.Context, dtEscalationLevel []dt.EscalationLevel) ([]escalationLevelModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	escalationLevels := make([]escalationLevelModel, 0, len(dtEscalationLevel))
	for _, level := range dtEscalationLevel {
		notificationActions, d := notificationActionToState(ctx, level.Actions)
		diags = append(diags, d...)

		// Set the escalate after to null if it is not set.
		escalateAfter := types.StringNull()
		if level.EscalateAfter != nil {
			escalateAfter = types.StringValue(*level.EscalateAfter)
		}

		escalationLevels = append(escalationLevels, escalationLevelModel{
			DisplayName:   types.StringValue(level.DisplayName),
			Actions:       notificationActions,
			EscalateAfter: escalateAfter,
		})

	}

	return escalationLevels, diags
}

// notificationActionToState converts the dt.Action to the state model.
func notificationActionToState(ctx context.Context, dtNotificationActions []dt.NotificationAction) ([]notificationActionModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(dtNotificationActions) == 0 {
		return nil, diags
	}

	notificationActions := make([]notificationActionModel, 0, len(dtNotificationActions))

	for _, action := range dtNotificationActions {
		smsConfig, d := smsConfigToState(ctx, action.SMSConfig)
		diags = append(diags, d...)

		emailConfig, d := emailConfigToState(ctx, action.EmailConfig)
		diags = append(diags, d...)

		webhookConfig, d := webhookConfigToState(ctx, action.WebhookConfig)
		diags = append(diags, d...)

		phoneCallConfig, d := phoneCallConfigToState(ctx, action.PhoneCallConfig)
		diags = append(diags, d...)

		notificationActions = append(notificationActions, notificationActionModel{
			Type:                 types.StringValue(action.Type),
			SMSConfig:            smsConfig,
			EmailConfig:          emailConfig,
			CorrigoConfig:        corrigoConfigToState(action.CorrigoConfig),
			ServiceChannelConfig: serviceChannelConfigToState(action.ServiceChannelConfig),
			WebhookConfig:        webhookConfig,
			PhoneCallConfig:      phoneCallConfig,
			SignalTowerConfig:    signalTowerConfigToState(action.SignalTowerConfig),
		})
	}

	return notificationActions, diags
}

func smsConfigToState(ctx context.Context, smsConfig *dt.SMSConfig) (*smsConfigModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if smsConfig == nil {
		return nil, diags
	}

	recipientsList, d := types.ListValueFrom(ctx, types.StringType, smsConfig.Recipients)
	diags = append(diags, d...)

	return &smsConfigModel{
		Recipients: recipientsList,
		Body:       types.StringValue(smsConfig.Body),
	}, diags
}

func emailConfigToState(ctx context.Context, emailConfig *dt.EmailConfig) (*emailConfigModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if emailConfig == nil {
		return nil, diags
	}

	recipientsList, d := types.ListValueFrom(ctx, types.StringType, emailConfig.Recipients)
	diags = append(diags, d...)

	return &emailConfigModel{
		Recipients: recipientsList,
		Subject:    types.StringValue(emailConfig.Subject),
		Body:       types.StringValue(emailConfig.Body),
	}, diags
}

func corrigoConfigToState(corrigoConfig *dt.CorrigoConfig) *corrigoConfigModel {
	if corrigoConfig == nil {
		return nil
	}

	return &corrigoConfigModel{
		AssetID:              types.StringValue(corrigoConfig.AssetID),
		TaskID:               types.StringValue(corrigoConfig.TaskID),
		CustomerID:           types.StringValue(corrigoConfig.CustomerID),
		ClientSecret:         types.StringValue(corrigoConfig.ClientSecret),
		CompanyName:          types.StringValue(corrigoConfig.CompanyName),
		SubTypeID:            types.StringValue(corrigoConfig.SubTypeID),
		ContactName:          types.StringValue(corrigoConfig.ContactName),
		ContactAddress:       types.StringValue(corrigoConfig.ContactAddress),
		WorkOrderDescription: types.StringValue(corrigoConfig.WorkOrderDescription),
		StudioDashboardURL:   types.StringValue(corrigoConfig.StudioDashboardURL),
	}
}

func serviceChannelConfigToState(serviceChannelConfig *dt.ServiceChannelConfig) *serviceChannelConfigModel {
	if serviceChannelConfig == nil {
		return nil
	}

	return &serviceChannelConfigModel{
		StoreID:     types.StringValue(serviceChannelConfig.StoreID),
		AssetTagID:  types.StringValue(serviceChannelConfig.AssetTagID),
		Trade:       types.StringValue(serviceChannelConfig.Trade),
		Description: types.StringValue(serviceChannelConfig.Description),
	}
}

func webhookConfigToState(ctx context.Context, webhookConfig *dt.WebhookConfig) (*webhookConfigModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if webhookConfig == nil {
		return nil, diags
	}

	headersMap, d := types.MapValueFrom(ctx, types.StringType, webhookConfig.Headers)
	diags = append(diags, d...)

	return &webhookConfigModel{
		URL:             types.StringValue(webhookConfig.URL),
		SignatureSecret: types.StringValue(webhookConfig.SignatureSecret),
		Headers:         headersMap,
	}, diags
}

func phoneCallConfigToState(ctx context.Context, phoneCallConfig *dt.PhoneCallConfig) (*phoneCallConfigModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if phoneCallConfig == nil {
		return nil, diags
	}

	recipientsList, d := types.ListValueFrom(ctx, types.StringType, phoneCallConfig.Recipients)
	diags = append(diags, d...)

	return &phoneCallConfigModel{
		Recipients:   recipientsList,
		Introduction: types.StringValue(phoneCallConfig.Introduction),
		Message:      types.StringValue(phoneCallConfig.Message),
	}, diags
}

func signalTowerConfigToState(signalTowerConfig *dt.SignalTowerConfig) *signalTowerConfigModel {
	if signalTowerConfig == nil {
		return nil
	}

	return &signalTowerConfigModel{
		CloudConnectorName: types.StringValue(signalTowerConfig.CloudConnectorName),
	}
}

func scheduleToState(schedule *dt.Schedule) *scheduleModel {
	if schedule == nil {
		return nil
	}
	scheduleModel := scheduleModel{}
	scheduleModel.Timezone = types.StringValue(schedule.Timezone)

	scheduleModel.Slots = make([]slotsModel, len(schedule.Slots))
	for slotIndex, slot := range schedule.Slots {
		scheduleModel.Slots[slotIndex].DayOfWeek = make([]types.String, len(slot.DaysOfWeek))
		for dayIndex, day := range slot.DaysOfWeek {
			scheduleModel.Slots[slotIndex].DayOfWeek[dayIndex] = types.StringValue(day)
		}

		scheduleModel.Slots[slotIndex].TimeRange = make([]timeRangeModel, len(slot.TimeRange))
		for timeRangeIndex, timeRange := range slot.TimeRange {
			scheduleModel.Slots[slotIndex].TimeRange[timeRangeIndex] = timeRangeModel{
				Start: timeOfDayModel{
					Hour:   types.Int32Value(timeRange.Start.Hour),
					Minute: types.Int32Value(timeRange.Start.Minute),
				},
				End: timeOfDayModel{
					Hour:   types.Int32Value(timeRange.End.Hour),
					Minute: types.Int32Value(timeRange.End.Minute),
				},
			}
		}
	}

	return &scheduleModel
}

func triggerToState(trigger dt.Trigger) triggerModel {

	model := triggerModel{
		Field: types.StringValue(trigger.Field),
	}

	if trigger.Range != nil {
		model.Range = &rangeModel{
			Lower: types.Float64PointerValue(trigger.Range.Lower),
			Upper: types.Float64PointerValue(trigger.Range.Upper),
			Type:  types.StringValue(trigger.Range.Type),
		}
	}

	model.Presence = types.StringNull()
	if trigger.Presence != nil {
		model.Presence = types.StringValue(*trigger.Presence)
	}
	model.Motion = types.StringNull()
	if trigger.Motion != nil {
		model.Motion = types.StringValue(*trigger.Motion)
	}
	model.Occupancy = types.StringNull()
	if trigger.Occupancy != nil {
		model.Occupancy = types.StringValue(*trigger.Occupancy)
	}
	model.Connection = types.StringNull()
	if trigger.Connection != nil {
		model.Connection = types.StringValue(*trigger.Connection)
	}
	model.Contact = types.StringNull()
	if trigger.Contact != nil {
		model.Contact = types.StringValue(*trigger.Contact)
	}
	model.TriggerCount = types.Int32Value(trigger.TriggerCount)
	if trigger.TriggerCount == 0 {
		model.TriggerCount = types.Int32Null()
	}

	return model
}

func stateToNotificationRule(ctx context.Context, state notificationRuleModel) (dt.NotificationRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	devices, d := expandStringList(ctx, state.Devices)
	diags = append(diags, d...)

	deviceLabels := make(map[string]string)
	d = state.DeviceLabels.ElementsAs(ctx, &deviceLabels, false)
	diags = append(diags, d...)

	escalationLevels, d := stateToEscalationLevels(ctx, state.EscalationLevels)
	diags = append(diags, d...)

	actions, d := stateToNotificationAction(ctx, state.Actions)
	diags = append(diags, d...)

	triggerDelay := state.TriggerDelay.ValueStringPointer()
	unacknowledgeAfter := state.UnacknowledgeAfter.ValueStringPointer()
	schedule := stateToSchedule(state.Schedule)

	return dt.NotificationRule{
		Name:                 state.Name.ValueString(),
		Enabled:              state.Enabled.ValueBool(),
		DisplayName:          state.DisplayName.ValueString(),
		Devices:              devices,
		DeviceLabels:         deviceLabels,
		Trigger:              stateToTrigger(state.Trigger),
		EscalationLevels:     escalationLevels,
		Schedule:             schedule,
		TriggerDelay:         triggerDelay,
		ReminderNotification: state.ReminderNotification.ValueBool(),
		ResolvedNotification: state.ResolvedNotification.ValueBool(),
		UnacknowledgesAfter:  unacknowledgeAfter,
		Actions:              actions,
	}, diags
}

func stateToTrigger(state triggerModel) dt.Trigger {
	trigger := dt.Trigger{
		Field:        state.Field.ValueString(),
		Presence:     state.Presence.ValueStringPointer(),
		Motion:       state.Motion.ValueStringPointer(),
		Occupancy:    state.Occupancy.ValueStringPointer(),
		Connection:   state.Connection.ValueStringPointer(),
		Contact:      state.Contact.ValueStringPointer(),
		TriggerCount: state.TriggerCount.ValueInt32(),
	}
	if state.Range == nil {
		return trigger
	}

	trigger.Range = &dt.Range{
		Lower: state.Range.Lower.ValueFloat64Pointer(),
		Upper: state.Range.Upper.ValueFloat64Pointer(),
		Type:  state.Range.Type.ValueString(),
	}
	return trigger
}

func stateToEscalationLevels(ctx context.Context, state []escalationLevelModel) ([]dt.EscalationLevel, diag.Diagnostics) {
	var diags diag.Diagnostics
	levels := []dt.EscalationLevel{}
	for _, level := range state {
		actions, d := stateToNotificationAction(ctx, level.Actions)
		diags = append(diags, d...)

		escalateAfter := level.EscalateAfter.ValueStringPointer()

		levels = append(levels, dt.EscalationLevel{
			DisplayName:   level.DisplayName.ValueString(),
			Actions:       actions,
			EscalateAfter: escalateAfter,
		})
	}

	return levels, diags
}

func stateToNotificationAction(ctx context.Context, state []notificationActionModel) ([]dt.NotificationAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(state) == 0 {
		return nil, diags
	}

	actions := make([]dt.NotificationAction, 0, len(state))
	for _, action := range state {
		smsConfig, d := stateToSMSConfig(ctx, action.SMSConfig)
		diags = append(diags, d...)

		emailConfig, d := stateToEmailConfig(ctx, action.EmailConfig)
		diags = append(diags, d...)

		webhookConfig, d := stateToWebhookConfig(ctx, action.WebhookConfig)
		diags = append(diags, d...)

		phoneCallConfig, d := stateToPhoneCallConfig(ctx, action.PhoneCallConfig)
		diags = append(diags, d...)

		actions = append(actions, dt.NotificationAction{
			Type:                 action.Type.ValueString(),
			SMSConfig:            smsConfig,
			EmailConfig:          emailConfig,
			CorrigoConfig:        stateToCorrigoConfig(action.CorrigoConfig),
			ServiceChannelConfig: stateToServiceChannelConfig(action.ServiceChannelConfig),
			WebhookConfig:        webhookConfig,
			PhoneCallConfig:      phoneCallConfig,
			SignalTowerConfig:    stateToSignalTowerConfig(action.SignalTowerConfig),
		})
	}

	return actions, diags
}

func stateToSMSConfig(ctx context.Context, state *smsConfigModel) (*dt.SMSConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if state == nil {
		return nil, nil
	}

	recipients, d := expandStringList(ctx, state.Recipients)
	diags = append(diags, d...)

	return &dt.SMSConfig{
		Recipients: recipients,
		Body:       state.Body.ValueString(),
	}, diags
}

func stateToEmailConfig(ctx context.Context, state *emailConfigModel) (*dt.EmailConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	if state == nil {
		return nil, nil
	}

	recipients, d := expandStringList(ctx, state.Recipients)
	diags = append(diags, d...)

	return &dt.EmailConfig{
		Recipients: recipients,
		Subject:    state.Subject.ValueString(),
		Body:       state.Body.ValueString(),
	}, diags
}

func stateToCorrigoConfig(state *corrigoConfigModel) *dt.CorrigoConfig {
	if state == nil {
		return nil
	}

	return &dt.CorrigoConfig{
		AssetID:              state.AssetID.ValueString(),
		TaskID:               state.TaskID.ValueString(),
		CustomerID:           state.CustomerID.ValueString(),
		ClientSecret:         state.ClientSecret.ValueString(),
		CompanyName:          state.CompanyName.ValueString(),
		SubTypeID:            state.SubTypeID.ValueString(),
		ContactName:          state.ContactName.ValueString(),
		ContactAddress:       state.ContactAddress.ValueString(),
		WorkOrderDescription: state.WorkOrderDescription.ValueString(),
		StudioDashboardURL:   state.StudioDashboardURL.ValueString(),
	}
}

func stateToServiceChannelConfig(state *serviceChannelConfigModel) *dt.ServiceChannelConfig {
	if state == nil {
		return nil
	}

	return &dt.ServiceChannelConfig{
		StoreID:     state.StoreID.ValueString(),
		AssetTagID:  state.AssetTagID.ValueString(),
		Trade:       state.Trade.ValueString(),
		Description: state.Description.ValueString(),
	}
}

func stateToWebhookConfig(ctx context.Context, state *webhookConfigModel) (*dt.WebhookConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if state == nil {
		return nil, nil
	}

	headers := make(map[string]string)
	d := state.Headers.ElementsAs(ctx, &headers, false)
	diags = append(diags, d...)
	if diags.HasError() {
		return nil, diags
	}

	return &dt.WebhookConfig{
		URL:             state.URL.ValueString(),
		SignatureSecret: state.SignatureSecret.ValueString(),
		Headers:         headers,
	}, nil
}

func stateToPhoneCallConfig(ctx context.Context, state *phoneCallConfigModel) (*dt.PhoneCallConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if state == nil {
		return nil, nil
	}

	recipients, d := expandStringList(ctx, state.Recipients)
	diags = append(diags, d...)

	return &dt.PhoneCallConfig{
		Recipients:   recipients,
		Introduction: state.Introduction.ValueString(),
		Message:      state.Message.ValueString(),
	}, diags
}

func stateToSignalTowerConfig(state *signalTowerConfigModel) *dt.SignalTowerConfig {
	if state == nil {
		return nil
	}

	return &dt.SignalTowerConfig{
		CloudConnectorName: state.CloudConnectorName.ValueString(),
	}
}

func stateToSchedule(state *scheduleModel) *dt.Schedule {
	if state == nil {
		return nil
	}
	schedule := dt.Schedule{
		Timezone: state.Timezone.ValueString(),
		Slots:    make([]dt.Slot, len(state.Slots)),
	}

	for slotIndex, slot := range state.Slots {
		daysOfWeek := make([]string, len(slot.DayOfWeek))
		for dayIndex, day := range slot.DayOfWeek {
			daysOfWeek[dayIndex] = day.ValueString()
		}

		timeRanges := make([]dt.TimeRange, len(slot.TimeRange))
		for timeRangeIndex, timeRange := range slot.TimeRange {
			timeRanges[timeRangeIndex] = dt.TimeRange{
				Start: dt.TimeOfDay{
					Hour:   timeRange.Start.Hour.ValueInt32(),
					Minute: timeRange.Start.Minute.ValueInt32(),
				},
				End: dt.TimeOfDay{
					Hour:   timeRange.End.Hour.ValueInt32(),
					Minute: timeRange.End.Minute.ValueInt32(),
				},
			}
		}

		schedule.Slots[slotIndex] = dt.Slot{
			DaysOfWeek: daysOfWeek,
			TimeRange:  timeRanges,
		}
	}

	return &schedule
}
