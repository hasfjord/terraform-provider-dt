// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectMemberRoleBindingsResource{}
	_ resource.ResourceWithConfigure   = &projectMemberRoleBindingsResource{}
	_ resource.ResourceWithImportState = &projectMemberRoleBindingsResource{}

	validRoles = []string{
		"roles/project.user",
		"roles/project.admin",
	}
)

// NewMemberResource is a helper function to simplify the provider implementation.
func NewMemberResource() resource.Resource {
	return &projectMemberRoleBindingsResource{}
}

// projectMemberRoleBindingsResource is the resource implementation.
type projectMemberRoleBindingsResource struct {
	client *dt.Client
}

// Metadata returns the resource type name.
func (m *projectMemberRoleBindingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_member_role_bindings"
}

func (r *projectMemberRoleBindingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Schema defines the schema for the resource.
func (m *projectMemberRoleBindingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the project member role binding, in the format `organizations/{organization_id}/roles/{role_id}/members/{member_id}`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"member_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the member, which is the resource name of the project member. Is a number for users, xid for service accounts.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"member_display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the member.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "Resource name of the organization on the format `organizations/{organization_id}`.",
			},
			"projects": schema.SetAttribute{
				Required:    true,
				Description: "List of projects to grant roles to of the format `projects/{project_id}`.",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				// require recreation of the resource any of the projects change
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Email of the project member.",
				// require recreation of the resource if the email changes
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required:    true,
				Description: "Role to assign the member to. ",
				Validators: []validator.String{
					stringvalidator.OneOf(validRoles...),
				},
				// require recreation of the resource if the role changes
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of account the member has. This is either `user` or `serviceAccount`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type membersResourceModel struct {
	Name              types.String `tfsdk:"name"`
	MemberID          types.String `tfsdk:"member_id"`
	MemberDisplayName types.String `tfsdk:"member_display_name"`
	Organization      types.String `tfsdk:"organization"`
	Projects          types.Set    `tfsdk:"projects"`
	Email             types.String `tfsdk:"email"`
	Role              types.String `tfsdk:"role"`
	AccountType       types.String `tfsdk:"account_type"`
}

// Create creates the resource and sets the initial state.
func (m *projectMemberRoleBindingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the data from the request.
	var plan membersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toBeCreated, d := stateToBatchCreateProjectMemberRequest(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := m.client.BatchCreateMemberships(ctx, toBeCreated)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project member",
			"Could not create project member, unexpected error: "+err.Error(),
		)
		return
	}
	state, d := membershipsToState(ctx, plan.Organization.ValueString(), members)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (m *projectMemberRoleBindingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// get the current state
	var state membersResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationID, roleID, memberID, err := decodeID(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error decoding project member ID",
			"Could not decode project member ID, unexpected error: "+err.Error(),
		)
		return
	}
	organization := "organizations/" + organizationID
	role := "roles/" + roleID

	// get the project members for the organization and member ID
	members, err := m.client.ListProjectMemberships(ctx, organization, role, memberID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project member",
			"Could not get project member, unexpected error: "+err.Error(),
		)
		return
	}

	// convert the project member to state
	state, diags = membershipsToState(ctx, organization, members)
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

// Update updates the resource and sets the updated Terraform state on success.
func (m *projectMemberRoleBindingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get the current state
	var plan membersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	memberships, d := stateToMemberships(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := m.client.UpdateMemberships(ctx, memberships, plan.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project member",
			"Could not update project member, unexpected error: "+err.Error(),
		)
		return
	}

	// convert the project member to state
	newState, d := membershipsToState(ctx, plan.Organization.ValueString(), members)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	// set the state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// / Delete deletes the resource and removes the Terraform state on success.
func (m *projectMemberRoleBindingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// get the current state
	var state membersResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toBeDeleted, d := stateToBatchDeleteProjectMembersRequest(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.BatchDeleteMemberships(ctx, toBeDeleted)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project member",
			"Could not delete project member, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (m *projectMemberRoleBindingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	m.client = client
}

func stateToBatchDeleteProjectMembersRequest(ctx context.Context, plan membersResourceModel) (dt.BatchDeleteProjectMembersRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	membersToDelete := make([]string, 0)
	memberID := plan.MemberID.ValueString()

	projects, d := expandStringSet(ctx, plan.Projects)
	diags.Append(d...)
	if diags.HasError() {
		return dt.BatchDeleteProjectMembersRequest{}, diags
	}
	for _, project := range projects {
		if project == "" {
			diags.AddError(
				"Error deleting project member",
				"Project ID cannot be empty",
			)
			return dt.BatchDeleteProjectMembersRequest{}, diags
		}
		membersToDelete = append(membersToDelete, fmt.Sprintf("%s/members/%s", project, memberID))
	}
	return dt.BatchDeleteProjectMembersRequest{
		Names: membersToDelete,
	}, diags
}

func stateToBatchCreateProjectMemberRequest(ctx context.Context, plan membersResourceModel) (dt.BatchCreateProjectsMembersRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	projects, d := expandStringSet(ctx, plan.Projects)
	diags.Append(d...)
	if diags.HasError() {
		return dt.BatchCreateProjectsMembersRequest{}, diags
	}

	members := make([]dt.Members, 0, len(projects))

	for _, project := range projects {
		members = append(members, dt.Members{
			Project: project,
			Email:   plan.Email.ValueString(),
			Roles:   []string{plan.Role.ValueString()},
		})
	}

	return dt.BatchCreateProjectsMembersRequest{
		Members: members,
	}, nil
}

func stateToMemberships(ctx context.Context, plan membersResourceModel) ([]dt.Membership, diag.Diagnostics) {
	var diags diag.Diagnostics
	var memberships []dt.Membership

	projects, d := expandStringSet(ctx, plan.Projects)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	for _, project := range projects {
		if project == "" {
			diags.AddError(
				"Error getting project member",
				"Project ID cannot be empty",
			)
			return nil, diags
		}
		memberships = append(memberships, dt.Membership{
			Name:        fmt.Sprintf("%s/members/%s", project, plan.MemberID.ValueString()),
			DisplayName: plan.MemberDisplayName.ValueString(),
			Email:       plan.Email.ValueString(),
			Roles:       []string{plan.Role.ValueString()},
			AccountType: plan.AccountType.ValueString(),
		})
	}

	return memberships, diags
}

func membershipsToState(ctx context.Context, organization string, memberships []dt.Membership) (membersResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err error
	var projects []string
	var email string
	var role string
	var memberID string
	var accountType string
	var displayName string

	for _, membership := range memberships {
		if email == "" {
			email = membership.Email
		} else if email != membership.Email {
			diags.AddError(
				"Error getting project member",
				"Project memberships must have the same email",
			)
			return membersResourceModel{}, diags
		}

		if role == "" {
			role = membership.Roles[0]
		} else if role != membership.Roles[0] {
			diags.AddError(
				"Error getting project member",
				"Project memberships must have the same role",
			)
			return membersResourceModel{}, diags
		}

		if memberID == "" {
			_, memberID, err = dt.ParseResourceName(membership.Name)
			if err != nil {
				diags.AddError(
					"Error parsing membership name",
					"Could not parse membership name, unexpected error: "+err.Error(),
				)
				return membersResourceModel{}, diags
			}
		} else if !strings.Contains(membership.Name, memberID) {
			diags.AddError(
				"Error getting project member",
				"Project memberships must have the same memberID",
			)
			return membersResourceModel{}, diags
		}

		if accountType == "" {
			accountType = membership.AccountType
		} else if accountType != membership.AccountType {
			diags.AddError(
				"Error getting project member",
				"Project memberships must have the same account type",
			)
			return membersResourceModel{}, diags
		}
		if displayName == "" {
			displayName = membership.DisplayName
		} else if displayName != membership.DisplayName {
			diags.AddError(
				"Error getting project member",
				"Project memberships must have the same display name",
			)
			return membersResourceModel{}, diags
		}
		projectID, err := membership.ProjectID()
		if err != nil {
			diags.AddError(
				"Error getting project ID",
				"Could not get project ID, unexpected error: "+err.Error(),
			)
			return membersResourceModel{}, diags
		}
		projects = append(projects, "projects/"+projectID)
	}

	projectsSet, d := flattenStringSetToAttr(ctx, projects)
	diags.Append(d...)
	if diags.HasError() {
		return membersResourceModel{}, diags
	}

	organizationID := strings.TrimPrefix(organization, "organizations/")
	roleID := strings.TrimPrefix(role, "roles/")

	return membersResourceModel{
		Name:              types.StringValue(fmt.Sprintf("organizations/%s/roles/%s/members/%s", organizationID, roleID, memberID)),
		MemberID:          types.StringValue(memberID),
		MemberDisplayName: types.StringValue(displayName),
		Projects:          projectsSet,
		Organization:      types.StringValue(organization),
		Email:             types.StringValue(email),
		Role:              types.StringValue(role),
		AccountType:       types.StringValue(accountType),
	}, diags
}

func decodeID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 6 {
		return "", "", "", fmt.Errorf("invalid ID format: %s", id)
	}

	if parts[0] != "organizations" || parts[2] != "roles" || parts[4] != "members" {
		return "", "", "", fmt.Errorf("invalid ID format: %s", id)
	}

	return parts[1], parts[3], parts[5], nil
}
