// Copyright (c) HashiCorp, Inc.

package dt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type BatchCreateProjectsMembersRequest struct {
	Members []Members `json:"members"`
}

type Members struct {
	Project string   `json:"project"`
	Email   string   `json:"email"`
	Roles   []string `json:"roles"`
}

type BatchDeleteProjectMembersRequest struct {
	Names []string `json:"names"`
}

type CreateProjectMemberRequest struct {
	Roles []string `json:"roles"`
	Email string   `json:"email"`
}

type MembershipResponse struct {
	Memberships []Membership `json:"members"`
}

type ListProjectMembersResponse struct {
	Members       []Membership `json:"members"`
	NextPageToken string       `json:"nextPageToken"`
}

type Membership struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Roles       []string `json:"roles"`
	Email       string   `json:"email"`
	AccountType string   `json:"accountType"`
}

func (m Membership) ProjectID() (string, error) {
	projectID, _, err := ParseResourceName(m.Name)
	if err != nil {
		return "", fmt.Errorf("dt: failed to parse resource name: %w", err)
	}
	return projectID, nil
}

func (m Membership) ID() (string, error) {
	_, memberID, err := ParseResourceName(m.Name)
	if err != nil {
		return "", fmt.Errorf("dt: failed to parse resource name: %w", err)
	}
	return memberID, nil
}

// ListProjectMemberships lists all memberships for a given organization and member.
func (c *Client) ListProjectMemberships(ctx context.Context, organization, role, memberID string) ([]Membership, error) {
	var members []Membership

	params := map[string]string{
		"memberId":     memberID,
		"organization": organization,
		"pageSize":     "100", // Default page size
		"pageToken":    "",
	}

	// use the project wildcard to list all memberships across all projects in the organization:
	url := c.URL + "/v2/projects/-/members"

	for {
		responseBody, err := c.DoRequest(ctx, http.MethodGet, url, nil, params)
		if err != nil {
			return nil, err
		}

		var member ListProjectMembersResponse
		if err := json.Unmarshal(responseBody, &member); err != nil {
			return nil, err
		}
		members = append(members, member.Members...)

		if member.NextPageToken == "" {
			break
		}

		// If there is a next page token, set it for the next request:
		params["pageToken"] = member.NextPageToken
	}
	tflog.Debug(ctx, fmt.Sprintf("dt: found %d memberships for member %s in organization %s", len(members), memberID, organization))

	// Filter out memberships that do not match the specified role:
	filteredMembers := make([]Membership, 0, len(members))
	for _, member := range members {
		if len(member.Roles) != 1 {
			return nil, fmt.Errorf("dt: expected exactly one role for member %s in organization %s, got %d roles", memberID, organization, len(member.Roles))
		}
		if member.Roles[0] == role {
			filteredMembers = append(filteredMembers, member)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("dt: found %d memberships with role %s for member %s in organization %s", len(filteredMembers), role, memberID, organization))

	return filteredMembers, nil
}

// BatchCreateMemberships creates multiple project memberships in a single request.
// This is to avoid sending out multiple emails for each project membership.
func (c *Client) BatchCreateMemberships(ctx context.Context, req BatchCreateProjectsMembersRequest) ([]Membership, error) {
	url := c.URL + "/v2/projects/-/members:batchCreate"

	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to marshal create members request: %w", err)
	}

	responseBody, err := c.DoRequest(ctx, "POST", url, requestBody, nil)
	if err != nil {
		return nil, err
	}

	var response MembershipResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	return response.Memberships, nil
}

type UpdateProjectMemberRequest struct {
	Roles []string `json:"roles"`
}

// UpdateMemberships updates a project member in place.
func (c *Client) UpdateMemberships(ctx context.Context, memberships []Membership, role string) ([]Membership, error) {
	updatedMembers := make([]Membership, 0, len(memberships))
	for _, member := range memberships {
		projectID, memberID, err := ParseResourceName(member.Name)
		if err != nil {
			return nil, fmt.Errorf("dt: failed to parse resource name: %w", err)
		}

		// This api only allows a single role to be set for a member:
		member.Roles = []string{role}

		url := c.URL + "/v2/projects/" + projectID + "/members/" + memberID
		requestBody, err := json.Marshal(member)
		if err != nil {
			return nil, fmt.Errorf("dt: failed to marshal memberships: %w", err)
		}

		responseBody, err := c.DoRequest(ctx, "PATCH", url, requestBody, nil)
		if err != nil {
			return nil, fmt.Errorf("dt: failed to update memberships: %w", err)
		}

		var updatedMember Membership
		if err := json.Unmarshal(responseBody, &updatedMember); err != nil {
			return nil, fmt.Errorf("dt: failed to unmarshal updated memberships: %w", err)
		}

		updatedMembers = append(updatedMembers, updatedMember)
	}

	return updatedMembers, nil
}

func (c *Client) BatchDeleteMemberships(ctx context.Context, req BatchDeleteProjectMembersRequest) error {
	url := c.URL + "/v2/projects/-/members:batchDelete"

	requestBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("dt: failed to marshal batch delete memberships request: %w", err)
	}

	_, err = c.DoRequest(ctx, "POST", url, requestBody, nil)
	if err != nil {
		return fmt.Errorf("dt: failed to delete memberships: %w", err)
	}

	return nil
}
