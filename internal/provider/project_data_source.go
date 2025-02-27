package provider

import (
	"context"
	"fmt"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

// NewProjectDataSource is a helper function to simplify the provider implementation.
func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

// projectDataSource is the data source implementation.
type projectDataSource struct {
	client dt.Client
}

// Metadata returns the data source type name.
func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source.
func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID of the project.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The resource name of the project. On the form `projects/{project_id}`.",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the project.",
			},
			"inventory": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the project is an inventory project.",
			},
			"organization": schema.StringAttribute{
				Computed:    true,
				Description: "The reource name of the organization that the project belongs to. on the form `organizations/{organization_id}`.",
			},
			"organization_display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the organization that the project belongs to.",
			},
			"sensor_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of sensors in the project.",
			},
			"cloud_connector_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of cloud connectors in the project.",
			},
			"location": schema.ObjectAttribute{
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"latitude":      types.Float64Type,
					"longitude":     types.Float64Type,
					"time_location": types.StringType,
				},
			},
		},
	}
}

// projectModel is the data model for the data source.
type projectDataSourceModel struct {
	ID                      types.String                    `tfsdk:"id"`
	Name                    types.String                    `tfsdk:"name"`
	DisplayName             types.String                    `tfsdk:"display_name"`
	Inventory               types.Bool                      `tfsdk:"inventory"`
	Organization            types.String                    `tfsdk:"organization"`
	OrganizationDisplayName types.String                    `tfsdk:"organization_display_name"`
	SensorCount             types.Int32                     `tfsdk:"sensor_count"`
	CloudConnectorCount     types.Int32                     `tfsdk:"cloud_connector_count"`
	Location                *projectLocationDataSourceModel `tfsdk:"location"`
}

type projectLocationDataSourceModel struct {
	Latitude     types.Float64 `tfsdk:"latitude"`
	Longitude    types.Float64 `tfsdk:"longitude"`
	TimeLocation types.String  `tfsdk:"time_location"`
}

// Read refreshes the Terraform state with the latest data.
func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// retrieve data source configuration
	var config projectDataSourceModel
	diag := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	project, err := d.client.GetProject(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get project", err.Error())
		return
	}

	state, diags := projectToState(project)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *projectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
