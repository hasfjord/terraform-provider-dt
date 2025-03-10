// Copyright (c) HashiCorp, Inc.

package dt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Project struct {
	Name                    string   `json:"name"`
	DisplayName             string   `json:"displayName"`
	Inventory               bool     `json:"inventory"`
	Organization            string   `json:"organization"`
	OrganizationDisplayName string   `json:"organizationDisplayName"`
	SensorCount             int      `json:"sensorCount"`
	CloudConnectorCount     int      `json:"cloudConnectorCount"`
	Location                Location `json:"location"`
}

func (p Project) ID() (string, error) {
	return idFromProject(p.Name)
}

type Location struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TimeLocation string  `json:"timeLocation"`
}

func (c *Client) GetProject(ctx context.Context, projectName string) (Project, error) {
	// Get the project ID from the project name
	projectID, err := idFromProject(projectName)
	if err != nil {
		return Project{}, fmt.Errorf("failed to get project ID: %w", err)
	}

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), projectID)

	// Send a GET request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Project{}, err
	}

	var p Project
	err = json.Unmarshal(responseBody, &p)
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

func (c *Client) UpdateProject(ctx context.Context, project Project) (Project, error) {
	// Get the project ID from the project name
	projectID, err := idFromProject(project.Name)
	if err != nil {
		return Project{}, fmt.Errorf("failed to get project ID: %w", err)
	}

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), projectID)
	body, err := json.Marshal(project)
	if err != nil {
		return Project{}, err
	}

	// Send a PUT request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return Project{}, err
	}

	// Unmarshal the response body to a Project struct
	var p Project
	err = json.Unmarshal(responseBody, &p)
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

type createProjectRequest struct {
	DisplayName  string `json:"displayName"`
	Organization string `json:"organization"`
	Location     struct {
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		TimeLocation string  `json:"timeLocation"`
	} `json:"location"`
}

func (c *Client) CreateProject(ctx context.Context, project Project) (Project, error) {
	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects
	url := fmt.Sprintf("%s/v2/projects", strings.TrimSuffix(c.URL, "/"))

	createProjectRequest := createProjectRequest{
		DisplayName:  project.DisplayName,
		Organization: project.Organization,
		Location:     project.Location,
	}

	// Marshal the project to JSON
	body, err := json.Marshal(createProjectRequest)
	if err != nil {
		return Project{}, err
	}

	// Send a POST request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Project{}, err
	}

	// Unmarshal the response body to a Project struct
	var p Project
	err = json.Unmarshal(responseBody, &p)
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

func (c *Client) DeleteProject(ctx context.Context, project string) error {
	// Get the project ID from the project name
	projectID, err := idFromProject(project)
	if err != nil {
		return fmt.Errorf("failed to get project ID: %w", err)
	}

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/v2/projects/%s", strings.TrimSuffix(c.URL, "/"), projectID)

	// Send a DELETE request to the API
	_, err = c.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return nil
}

func idFromProject(project string) (string, error) {
	parts := strings.Split(project, "/")
	if len(parts) != 2 {
		return project, fmt.Errorf("invalid project name: %s", project)
	}
	return parts[1], nil
}
