// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"strings"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &emulatorResource{}
	_ resource.ResourceWithConfigure   = &emulatorResource{}
	_ resource.ResourceWithImportState = &emulatorResource{}
)

var (
	validEmulatorTypes = []string{
		"touch",
		"temperature",
		"proximity",
		"touchCounter",
		"proximityCounter",
		"humidity",
		"waterDetector",
		"co2",
		"motion",
		"contact",
		"deskOccupancy",
		"ccon",
	}
)

// NewEmulatorResource is a helper function to simplify the provider implementation.
func NewEmulatorResource() resource.Resource {
	return &emulatorResource{}
}

// emulatorResource is the resource implementation.
type emulatorResource struct {
	client *dt.Client
}

func (r *emulatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_emulator"
}

func (r *emulatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *emulatorResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	labelDefault, diags := basetypes.NewMapValueFrom(ctx, types.StringType, map[string]string{})
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The resource name of the emulator on the form: `projects/{project_id}/devices/{device_id}`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the emulator.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The project ID to create the emulator in.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: `The type of emulator valid types are: ` + strings.Join(validEmulatorTypes, ", "),
				Validators: []validator.String{
					stringvalidator.OneOf(validEmulatorTypes...),
				},
			},
			"system_labels": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of system labels assigned to the emulator. Read only",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A map of labels to assign to the emulator.",
				Default:     mapdefault.StaticValue(labelDefault),
			},
		},
	}
}

type emulatorResourceModel struct {
	Name         types.String `tfsdk:"name"`
	DisplayName  types.String `tfsdk:"display_name"`
	ProjectID    types.String `tfsdk:"project_id"`
	Type         types.String `tfsdk:"type"`
	SystemLabels types.Map    `tfsdk:"system_labels"`
	Labels       types.Map    `tfsdk:"labels"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *emulatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the data from the request
	var plan emulatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	toBeCreated, diags := stateToEmulator(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Create the emulator
	created, err := r.client.CreateEmulator(ctx, plan.ProjectID.ValueString(), toBeCreated)
	if err != nil {
		resp.Diagnostics.AddError("Error creating emulator", err.Error())
		return
	}

	state, diags := emulatorToState(ctx, created)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the Terraform state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *emulatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// retrieve the current state
	var state emulatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Get the emulator
	emulator, err := r.client.GetEmulator(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading emulator", err.Error())
		return
	}

	state, diags = emulatorToState(ctx, emulator)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the Terraform state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

// Update updates the resource and refreshes the Terraform state.
func (r *emulatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// retrieve the current state
	var state emulatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// retrieve the plan
	var plan emulatorResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	toBeUpdated, diags := stateToEmulator(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Update the emulator
	updated, err := r.client.UpdateEmulator(ctx, toBeUpdated)
	if err != nil {
		resp.Diagnostics.AddError("Error updating emulator", err.Error())
		return
	}

	state, diags = emulatorToState(ctx, updated)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Set the Terraform state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *emulatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// retrieve the current state
	var state emulatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Delete the emulator
	err := r.client.DeleteEmulator(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting emulator", err.Error())
		return
	}
}

func (r *emulatorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func stateToEmulator(ctx context.Context, state emulatorResourceModel) (dt.Emulator, diag.Diagnostics) {
	var diags diag.Diagnostics

	labelsMap := make(map[string]string)
	d := state.Labels.ElementsAs(ctx, &labelsMap, false)
	diags.Append(d...)
	if d.HasError() {
		return dt.Emulator{}, diags
	}

	// add the system labels
	labelsMap["name"] = state.DisplayName.ValueString()
	labelsMap["virtual-sensor"] = ""

	return dt.Emulator{
		Name:   state.Name.ValueString(),
		Type:   state.Type.ValueString(),
		Labels: labelsMap,
	}, diags
}

func emulatorToState(ctx context.Context, emulator dt.Emulator) (emulatorResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var labels = make(map[string]string)
	var systemLabels = make(map[string]string)
	for key, value := range emulator.Labels {
		if key == "name" || key == "virtual-sensor" {
			systemLabels[key] = value
		} else {
			labels[key] = value
		}
	}
	displayName := systemLabels["name"]
	if displayName == "" {
		displayName = emulator.DeviceID()
	}

	labelsMap, d := basetypes.NewMapValueFrom(ctx, types.StringType, labels)
	diags.Append(d...)

	systemLabelsMap, d := basetypes.NewMapValueFrom(ctx, types.StringType, systemLabels)
	diags.Append(d...)

	return emulatorResourceModel{
		Name:         types.StringValue(emulator.Name),
		DisplayName:  types.StringValue(displayName),
		Type:         types.StringValue(emulator.Type),
		ProjectID:    types.StringValue(emulator.ProjectID()),
		SystemLabels: systemLabelsMap,
		Labels:       labelsMap,
	}, diags
}
