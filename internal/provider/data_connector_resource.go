// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &dataConnectorResource{}
	_ resource.ResourceWithConfigure = &dataConnectorResource{}
)

// NewDataConnectorResource is a helper function to simplify the provider implementation.
func NewDataConnectorResource() resource.Resource {
	return &dataConnectorResource{}
}

// dataConnectorResource is the resource implementation.
type dataConnectorResource struct {
	client *dt.Client
}

// Metadata returns the resource type name.
func (r *dataConnectorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_connector"
}

// Schema defines the schema for the resource.
func (r *dataConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	labelDefault, diags := types.ListValueFrom(ctx, types.StringType, []string{"name"})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The resource name of the data connector. On the form `projects/{project_id}/dataConnectors/{data_connector_id}`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the data connector.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of connector, allowed values: HTTP_PUSH, AZURE_SERVICE_BUS, AZURE_EVENT_HUB, GOOGLE_CLOUD_PUBSUB, AWS_SQS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				Required:    true,
				Description: "The resource name of the project that the data connector belongs to. On the form `projects/{project_id}`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the connector.",
				Default:     stringdefault.StaticString("ACTIVE"),
			},
			"events": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Events to listen on. Empty list is equal to all events.",
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Computed:    true,
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Label keys to include in the event payload.",
				Default:     listdefault.StaticValue(labelDefault),
				Computed:    true,
			},
			"http_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "HTTP configuration for the connector.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:    true,
						Description: "Valid URL using HTTPS.",
					},
					"signature_secret": schema.StringAttribute{
						Optional:           true,
						Description:        "Secret used to sign the payload",
						DeprecationMessage: "The use of signature secret is deprecated, use DT-Asymmetric-Signature to validate the payload instead.",
						Sensitive:          true,
					},
					"headers": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "Headers to include in the HTTP request.",
					},
				},
			},
			"azure_service_bus_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Azure Service Bus configuration for the connector.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:    true,
						Description: "Service bus URL on the form `sb://<namespace>.servicebus.windows.net/<topic>`.",
					},
					"authentication_config": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Authentication configuration for the service bus.",
						Attributes: map[string]schema.Attribute{
							"tenant_id": schema.StringAttribute{
								Required:    true,
								Description: "Azure tenant ID.",
							},
							"client_id": schema.StringAttribute{
								Required:    true,
								Description: "Azure client ID.",
							},
						},
					},
					"broker_properties": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Broker properties for the service bus.",
						Attributes: map[string]schema.Attribute{
							"correlation_id": schema.StringAttribute{
								Optional:    true,
								Description: "Correlation ID for the message.",
							},
						},
					},
				},
			},
			"azure_event_hub_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Azure Event Hub configuration for the connector.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:    true,
						Description: "Event hub URL on the form `sb://<namespace>.servicebus.windows.net/<eventhub>`.",
					},
					"authentication_config": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Authentication configuration for the event hub.",
						Attributes: map[string]schema.Attribute{
							"tenant_id": schema.StringAttribute{
								Required:    true,
								Description: "Azure tenant ID.",
							},
							"client_id": schema.StringAttribute{
								Required:    true,
								Description: "Azure client ID.",
							},
						},
					},
				},
			},
			"pubsub_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Google Cloud Pub/Sub configuration for the connector.",
				Attributes: map[string]schema.Attribute{
					"topic": schema.StringAttribute{
						Required:    true,
						Description: "Pub/Sub topic name.",
					},
					"audience": schema.StringAttribute{
						Required:    true,
						Description: "Audience for the token.",
					},
				},
			},
			"aws_sqs_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "AWS SQS configuration for the connector.",
				Attributes: map[string]schema.Attribute{
					"queue_url": schema.StringAttribute{
						Required:    true,
						Description: "SQS queue URL.",
					},
					"aws_role_arn": schema.StringAttribute{
						Required:    true,
						Description: "AWS role ARN.",
					},
					"audience": schema.StringAttribute{
						Required:    true,
						Description: "Audience for the token.",
					},
				},
			},
		},
	}
}

type dataConnectorResourceModel struct {
	Name                  types.String           `tfsdk:"name"`
	DisplayName           types.String           `tfsdk:"display_name"`
	Project               types.String           `tfsdk:"project"`
	Type                  types.String           `tfsdk:"type"`
	Status                types.String           `tfsdk:"status"`
	Events                types.List             `tfsdk:"events"`
	Labels                types.List             `tfsdk:"labels"`
	HTTPConfig            *httpConfig            `tfsdk:"http_config"`
	AzureServiceBusConfig *azureServiceBusConfig `tfsdk:"azure_service_bus_config"`
	AzureEventHubConfig   *azureEventHubConfig   `tfsdk:"azure_event_hub_config"`
	PubsubConfig          *pubsubConfig          `tfsdk:"pubsub_config"`
	AWSSQSConfig          *awsSQSConfig          `tfsdk:"aws_sqs_config"`
}

type httpConfig struct {
	URL             types.String `tfsdk:"url"`
	SignatureSecret types.String `tfsdk:"signature_secret"`
	Headers         types.Map    `tfsdk:"headers"`
}

type azureServiceBusConfig struct {
	URL                  types.String          `tfsdk:"url"`
	AuthenticationConfig *authenticationConfig `tfsdk:"authentication_config"`
	BrokerProperties     *brokerProperties     `tfsdk:"broker_properties"`
}

type brokerProperties struct {
	CorrelationID types.String `tfsdk:"correlation_id"`
}

type azureEventHubConfig struct {
	URL                  types.String          `tfsdk:"url"`
	AuthenticationConfig *authenticationConfig `tfsdk:"authentication_config"`
}

type authenticationConfig struct {
	TenantID types.String `tfsdk:"tenant_id"`
	ClientID types.String `tfsdk:"client_id"`
}

type pubsubConfig struct {
	Topic    types.String `tfsdk:"topic"`
	Audience types.String `tfsdk:"audience"`
}

type awsSQSConfig struct {
	QueueURL   types.String `tfsdk:"queue_url"`
	AWSRoleArn types.String `tfsdk:"aws_role_arn"`
	Audience   types.String `tfsdk:"audience"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *dataConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the data from the request
	var plan dataConnectorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toBeCreated, diags := stateToDataConnector(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the data connector
	created, err := r.client.CreateDataConnector(ctx, plan.Project.ValueString(), toBeCreated)
	if err != nil {
		resp.Diagnostics.AddError("failed to create data connector", err.Error())
		return
	}

	state, diags := dataConnectorToState(ctx, created)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the Terraform state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dataConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get the current state
	var state dataConnectorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the data connector from the API
	dataConnector, err := r.client.GetDataConnector(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get data connector", err.Error())
		return
	}

	state, diags = dataConnectorToState(ctx, dataConnector)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dataConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get the current state
	var state dataConnectorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the data connector
	err := r.client.DeleteDataConnector(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete data connector", err.Error())
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *dataConnectorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dt.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"invalid provider data",
			"Provider data is not of the expected type",
		)
		return
	}

	r.client = client
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dataConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get the current state
	var plan dataConnectorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataConnector, diags := stateToDataConnector(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the data connector
	dataConnector, err := r.client.UpdateDataConnector(ctx, dataConnector)
	if err != nil {
		resp.Diagnostics.AddError("failed to update data connector", err.Error())
		return
	}

	state, diag := dataConnectorToState(ctx, dataConnector)
	resp.Diagnostics.Append(diag...)

	// Set the Terraform state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// stateToDataConnector converts the resource model to the API model.
func stateToDataConnector(ctx context.Context, plan dataConnectorResourceModel) (dt.DataConnector, diag.Diagnostics) {
	var diags diag.Diagnostics

	events, d := expandStringList(ctx, plan.Events)
	diags = append(diags, d...)

	labels, d := expandStringList(ctx, plan.Labels)
	diags = append(diags, d...)

	var httpPushConfig *dt.HTTPConfig
	var azureServiceBusConfig *dt.AzureServiceBusConfig
	var azureEventHubConfig *dt.AzureEventHubConfig
	var pubsubConfig *dt.PubsubConfig
	var awsSQSConfig *dt.AWSSQSConfig

	diags = append(diags, validateTypeConfig(plan)...)

	dataConnector := dt.DataConnector{
		Name:        plan.Name.ValueString(),
		DisplayName: plan.DisplayName.ValueString(),
		Type:        plan.Type.ValueString(),
		Status:      plan.Status.ValueString(),
		Events:      events,
		Labels:      labels,
	}
	switch plan.Type.ValueString() {
	case "HTTP_PUSH":
		headersMap := make(map[string]string)
		for key, value := range plan.HTTPConfig.Headers.Elements() {
			headersMap[key] = value.String()
		}
		httpPushConfig = &dt.HTTPConfig{
			Url:             plan.HTTPConfig.URL.ValueString(),
			SignatureSecret: plan.HTTPConfig.SignatureSecret.ValueString(),
			Headers:         headersMap,
		}
		dataConnector.HTTPConfig = httpPushConfig
	case "AZURE_SERVICE_BUS":
		azureServiceBusConfig = &dt.AzureServiceBusConfig{
			URL: plan.AzureServiceBusConfig.URL.ValueString(),
			AuthenticationConfig: dt.AuthenticationConfig{
				TenantID: plan.AzureServiceBusConfig.AuthenticationConfig.TenantID.ValueString(),
				ClientID: plan.AzureServiceBusConfig.AuthenticationConfig.ClientID.ValueString(),
			},
			BrokerProperties: dt.BrokerProperties{
				CorrelationID: plan.AzureServiceBusConfig.BrokerProperties.CorrelationID.ValueString(),
			},
		}
		dataConnector.AzureServiceBusConfig = azureServiceBusConfig
	case "AZURE_EVENT_HUB":
		azureEventHubConfig = &dt.AzureEventHubConfig{
			URL: plan.AzureEventHubConfig.URL.ValueString(),
			AuthenticationConfig: dt.AuthenticationConfig{
				TenantID: plan.AzureEventHubConfig.AuthenticationConfig.TenantID.ValueString(),
				ClientID: plan.AzureEventHubConfig.AuthenticationConfig.ClientID.ValueString(),
			},
		}
		dataConnector.AzureEventHubConfig = azureEventHubConfig
	case "GOOGLE_CLOUD_PUBSUB":
		pubsubConfig = &dt.PubsubConfig{
			Topic:    plan.PubsubConfig.Topic.ValueString(),
			Audience: plan.PubsubConfig.Audience.ValueString(),
		}
		dataConnector.PubsubConfig = pubsubConfig
	case "AWS_SQS":
		awsSQSConfig = &dt.AWSSQSConfig{
			QueueUrl:   plan.AWSSQSConfig.QueueURL.ValueString(),
			AwsRoleArn: plan.AWSSQSConfig.AWSRoleArn.ValueString(),
			Audience:   plan.AWSSQSConfig.Audience.ValueString(),
		}
		dataConnector.AWSSQSConfig = awsSQSConfig
	default:
		diags = append(diags,
			diag.NewAttributeErrorDiagnostic(
				path.Root("type"),
				"Invalid data connector type",
				fmt.Sprintf("Type %q is not a valid data connector type", plan.Type.ValueString()),
			),
		)
	}

	return dataConnector, diags
}

// dataConnectorToState converts the API model to the resource model.
func dataConnectorToState(ctx context.Context, dataConnector dt.DataConnector) (dataConnectorResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	labelsList, d := flattenStringListToAttr(ctx, dataConnector.Labels)
	diags = append(diags, d...)

	eventsList, d := flattenStringListToAttr(ctx, dataConnector.Events)
	diags = append(diags, d...)

	resourceModel := dataConnectorResourceModel{
		Name:        types.StringValue(dataConnector.Name),
		DisplayName: types.StringValue(dataConnector.DisplayName),
		Project:     types.StringValue(dataConnector.ProjectID()),
		Type:        types.StringValue(dataConnector.Type),
		Status:      types.StringValue(dataConnector.Status),
		Events:      eventsList,
		Labels:      labelsList,
	}

	switch dataConnector.Type {
	case "HTTP_PUSH":
		// Convert the headers map to a Terraform map
		headersMap, d := basetypes.NewMapValueFrom(ctx, types.StringType, dataConnector.HTTPConfig.Headers)
		diags = append(diags, d...)

		httpConfigModel := &httpConfig{
			URL:             types.StringValue(dataConnector.HTTPConfig.Url),
			SignatureSecret: types.StringValue(dataConnector.HTTPConfig.SignatureSecret),
			Headers:         headersMap,
		}
		resourceModel.HTTPConfig = httpConfigModel
	case "AZURE_SERVICE_BUS":
		azureServiceBusModel := &azureServiceBusConfig{
			URL: types.StringValue(dataConnector.AzureServiceBusConfig.URL),
			AuthenticationConfig: &authenticationConfig{
				TenantID: types.StringValue(dataConnector.AzureServiceBusConfig.AuthenticationConfig.TenantID),
				ClientID: types.StringValue(dataConnector.AzureServiceBusConfig.AuthenticationConfig.ClientID),
			},
			BrokerProperties: &brokerProperties{
				CorrelationID: types.StringValue(dataConnector.AzureServiceBusConfig.BrokerProperties.CorrelationID),
			},
		}
		resourceModel.AzureServiceBusConfig = azureServiceBusModel
	case "AZURE_EVENT_HUB":
		azureEventHubModel := &azureEventHubConfig{
			URL: types.StringValue(dataConnector.AzureEventHubConfig.URL),
			AuthenticationConfig: &authenticationConfig{
				TenantID: types.StringValue(dataConnector.AzureEventHubConfig.AuthenticationConfig.TenantID),
				ClientID: types.StringValue(dataConnector.AzureEventHubConfig.AuthenticationConfig.ClientID),
			},
		}
		resourceModel.AzureEventHubConfig = azureEventHubModel
	case "GOOGLE_CLOUD_PUBSUB":
		pubsubConfig := &pubsubConfig{
			Topic:    types.StringValue(dataConnector.PubsubConfig.Topic),
			Audience: types.StringValue(dataConnector.PubsubConfig.Audience),
		}
		resourceModel.PubsubConfig = pubsubConfig
	case "AWS_SQS":
		awsSQSConfig := &awsSQSConfig{
			QueueURL:   types.StringValue(dataConnector.AWSSQSConfig.QueueUrl),
			AWSRoleArn: types.StringValue(dataConnector.AWSSQSConfig.AwsRoleArn),
			Audience:   types.StringValue(dataConnector.AWSSQSConfig.Audience),
		}
		resourceModel.AWSSQSConfig = awsSQSConfig
	default:
		diags = append(diags,
			diag.NewAttributeErrorDiagnostic(
				path.Root("type"),
				"Invalid data connector type",
				fmt.Sprintf("Type %q is not a valid data connector type", dataConnector.Type),
			),
		)
	}

	return resourceModel, diags
}

func validateTypeConfig(plan dataConnectorResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	// Exclusive fields
	if plan.Type.ValueString() != "HTTP_PUSH" && plan.HTTPConfig != nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("http_config"), "http_config", "HTTP configuration is only allowed for HTTP_PUSH data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() != "AZURE_SERVICE_BUS" && plan.AzureServiceBusConfig != nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("azure_service_bus_config"), "azure_service_bus_config", "Azure Service Bus configuration is only allowed for AZURE_SERVICE_BUS data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() != "AZURE_EVENT_HUB" && plan.AzureEventHubConfig != nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("azure_event_hub_config"), "azure_event_hub_config", "Azure Event Hub configuration is only allowed for AZURE_EVENT_HUB data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() != "GOOGLE_CLOUD_PUBSUB" && plan.PubsubConfig != nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("pubsub_config"), "pubsub_config", "Google Cloud Pub/Sub configuration is only allowed for GOOGLE_CLOUD_PUBSUB data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() != "AWS_SQS" && plan.AWSSQSConfig != nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("aws_sqs_config"), "aws_sqs_config", "AWS SQS configuration is only allowed for AWS_SQS data connectors")
		diags = append(diags, diag)
	}

	// Required fields
	if plan.Type.ValueString() == "HTTP_PUSH" && plan.HTTPConfig == nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("http_config"), "http_config", "HTTP configuration is required for HTTP_PUSH data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() == "AZURE_SERVICE_BUS" && plan.AzureServiceBusConfig == nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("azure_service_bus_config"), "azure_service_bus_config", "Azure Service Bus configuration is required for AZURE_SERVICE_BUS data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() == "AZURE_EVENT_HUB" && plan.AzureEventHubConfig == nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("azure_event_hub_config"), "azure_event_hub_config", "Azure Event Hub configuration is required for AZURE_EVENT_HUB data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() == "GOOGLE_CLOUD_PUBSUB" && plan.PubsubConfig == nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("pubsub_config"), "pubsub_config", "Google Cloud Pub/Sub configuration is required for GOOGLE_CLOUD_PUBSUB data connectors")
		diags = append(diags, diag)
	}
	if plan.Type.ValueString() == "AWS_SQS" && plan.AWSSQSConfig == nil {
		diag := diag.NewAttributeErrorDiagnostic(path.Root("aws_sqs_config"), "aws_sqs_config", "AWS SQS configuration is required for AWS_SQS data connectors")
		diags = append(diags, diag)
	}

	return diags
}
