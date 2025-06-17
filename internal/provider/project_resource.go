// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

// NewProjectResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// projectResource is the resource implementation.
type projectResource struct {
	client *dt.Client
}

// Metadata returns the resource type name.
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Schema defines the schema for the resource.
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The project ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The resource name of the project. On the form `projects/{project_id}`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the project.",
			},
			"inventory": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the project is an inventory project.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The reource name of the organization that the project belongs to. on the form `organizations/{organization_id}`.",
			},
			"organization_display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the organization that the project belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sensor_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of sensors in the project.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"cloud_connector_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of cloud connectors in the project.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The location of the project.",
				Attributes: map[string]schema.Attribute{
					"latitude": schema.Float64Attribute{
						Optional:    true,
						Description: "The latitude of the project in Degrees Decimal. This is used to determine the time zone of the project.",
					},
					"longitude": schema.Float64Attribute{
						Optional:    true,
						Description: "The longitude of the project in Degrees Decimal. This is used to determine the time zone of the project.",
					},
					"time_location": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The time location of the project. This is used to determine the time zone of the project. For example, `Europe/Oslo`.",
						Default:     stringdefault.StaticString("UTC"),
					},
				},
			},
		},
	}
}

// projectModel is the data model for the data source.
type projectResourceModel struct {
	ID                      types.String                  `tfsdk:"id"`
	Name                    types.String                  `tfsdk:"name"`
	DisplayName             types.String                  `tfsdk:"display_name"`
	Inventory               types.Bool                    `tfsdk:"inventory"`
	Organization            types.String                  `tfsdk:"organization"`
	OrganizationDisplayName types.String                  `tfsdk:"organization_display_name"`
	SensorCount             types.Int32                   `tfsdk:"sensor_count"`
	CloudConnectorCount     types.Int32                   `tfsdk:"cloud_connector_count"`
	Location                *projectLocationResourceModel `tfsdk:"location"`
}

type projectLocationResourceModel struct {
	Latitude     types.Float64 `tfsdk:"latitude"`
	Longitude    types.Float64 `tfsdk:"longitude"`
	TimeLocation types.String  `tfsdk:"time_location"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the data from the request.
	var plan projectResourceModel
	dias := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(dias...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := stateToProject(plan)

	// Create the project.
	project, err := r.client.CreateProject(ctx, project)
	if err != nil {
		resp.Diagnostics.AddError("failed to create project", err.Error())
		return
	}

	plan, diags := projectToState(project)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the Terraform state.
	dias = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(dias...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// get the current state
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get the project from the API
	project, err := r.client.GetProject(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get project", err.Error())
		return
	}

	state, diags = projectToState(project)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// retrieve values from plan
	var state projectResourceModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate the api request from the plan
	toBeUpdated := stateToUpdateProjectRequest(state)
	project, err := r.client.UpdateProject(ctx, toBeUpdated)
	if err != nil {
		resp.Diagnostics.AddError("failed to update project", err.Error())
		return
	}

	// set the updated state
	updateProjectState(project, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// retrieve the current state
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete the project
	err := r.client.DeleteProject(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete project", err.Error())
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dt.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dt.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func projectToState(project dt.Project) (projectResourceModel, diag.Diagnostics) {
	id, err := project.ID()
	if err != nil {
		diags := diag.NewErrorDiagnostic("ID", "failed to get project ID")
		return projectResourceModel{}, diag.Diagnostics{diags}
	}
	var latitude types.Float64
	if project.Location.Latitude != nil {
		latitude = types.Float64PointerValue(project.Location.Latitude)
	} else {
		latitude = types.Float64Null()
	}
	var longitude types.Float64
	if project.Location.Longitude != nil {
		longitude = types.Float64PointerValue(project.Location.Longitude)
	} else {
		longitude = types.Float64Null()
	}
	return projectResourceModel{
		ID:                      types.StringValue(id),
		Name:                    types.StringValue(project.Name),
		DisplayName:             types.StringValue(project.DisplayName),
		Inventory:               types.BoolValue(project.Inventory),
		Organization:            types.StringValue(project.Organization),
		OrganizationDisplayName: types.StringValue(project.OrganizationDisplayName),
		SensorCount:             types.Int32Value(int32(project.SensorCount)),
		CloudConnectorCount:     types.Int32Value(int32(project.CloudConnectorCount)),
		Location: &projectLocationResourceModel{
			Latitude:     latitude,
			Longitude:    longitude,
			TimeLocation: types.StringValue(project.Location.TimeLocation),
		},
	}, nil
}

func stateToProject(state projectResourceModel) dt.Project {
	project := dt.Project{
		Name:                    state.Name.ValueString(),
		DisplayName:             state.DisplayName.ValueString(),
		Inventory:               state.Inventory.ValueBool(),
		Organization:            state.Organization.ValueString(),
		OrganizationDisplayName: state.OrganizationDisplayName.ValueString(),
		SensorCount:             int(state.SensorCount.ValueInt32()),
		CloudConnectorCount:     int(state.CloudConnectorCount.ValueInt32()),
	}
	if state.Location != nil {
		project.Location = dt.Location{
			Latitude:     state.Location.Latitude.ValueFloat64Pointer(),
			Longitude:    state.Location.Longitude.ValueFloat64Pointer(),
			TimeLocation: state.Location.TimeLocation.ValueString(),
		}
	}
	return project
}

func stateToUpdateProjectRequest(state projectResourceModel) dt.EditableProject {
	updateRequest := dt.EditableProject{
		Name:         state.Name.ValueString(),
		DisplayName:  state.DisplayName.ValueString(),
		Organization: state.Organization.ValueString(),
	}
	if state.Location != nil {
		updateRequest.Location = dt.Location{
			Latitude:     state.Location.Latitude.ValueFloat64Pointer(),
			Longitude:    state.Location.Longitude.ValueFloat64Pointer(),
			TimeLocation: state.Location.TimeLocation.ValueString(),
		}
	}
	return updateRequest
}

func updateProjectState(project dt.EditableProject, state *projectResourceModel) {
	state.Name = types.StringValue(project.Name)
	state.DisplayName = types.StringValue(project.DisplayName)
	state.Organization = types.StringValue(project.Organization)
	state.Location = &projectLocationResourceModel{
		Latitude:     types.Float64PointerValue(project.Location.Latitude),
		Longitude:    types.Float64PointerValue(project.Location.Longitude),
		TimeLocation: types.StringValue(project.Location.TimeLocation),
	}
}
