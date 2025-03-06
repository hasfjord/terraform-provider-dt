// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

func NewDeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

type deviceDataSource struct {
	client dt.Client
}

// Metadata returns the data source type name.
func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the data source.
func (d deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID of the device.",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID of the project.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The resource name of the device. On the form `projects/{project_id}/devices/{device_id}`.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the device.",
			},
			"labels": schema.MapAttribute{
				Computed:    true,
				Description: "The labels of the device.",
				ElementType: types.StringType,
			},
		},
	}
}

type DeviceDataSourceModel struct {
	DeviceID  types.String `tfsdk:"device_id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Labels    types.Map    `tfsdk:"labels"`
}

func (d deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// retrieve data source configuration
	var config DeviceDataSourceModel
	diag := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	device, err := d.client.GetDevice(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get device", err.Error())
		return
	}

	deviceID, projectID, err := idFromName(device.Name)
	if err != nil {
		resp.Diagnostics.AddError("failed to get device ID and project ID", err.Error())
		return
	}
	labels, diag := types.MapValueFrom(ctx, types.StringType, device.Labels)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	state := DeviceDataSourceModel{
		DeviceID:  types.StringValue(deviceID),
		ProjectID: types.StringValue(projectID),
		Name:      types.StringValue(device.Name),
		Type:      types.StringValue(device.Type),
		Labels:    labels,
	}

	diag = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
}

func (d *deviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dt.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"invalid provider data",
			fmt.Sprintf("Expected *dt.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = *client
}

func idFromName(name string) (string, string, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 4 {
		return name, "", fmt.Errorf("invalid device name: %s", name)
	}
	return parts[len(parts)-1], parts[len(parts)-3], nil
}
