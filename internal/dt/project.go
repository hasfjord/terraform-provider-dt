// Copyright (c) HashiCorp, Inc.

package dt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type ListProjectResponse struct {
	Projects []Project `json:"projects"`
}

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

type projectCache struct {
	projects map[string]Project

	mu sync.RWMutex
}

func (c *projectCache) getProject(name string) (Project, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if project, ok := c.projects[name]; ok {
		return project, true
	}
	return Project{}, false
}

func (c *projectCache) setProject(project Project) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.projects[project.Name] = project
}

func (c *Client) GetProject(ctx context.Context, projectName string) (Project, error) {
	// first check if the project is in the cache
	if project, ok := c.projectCache.getProject(projectName); ok {
		return project, nil
	}

	// call the API to get all the projects in the org and populate the cache
	projects, err := c.listProjects(ctx)
	if err != nil {
		return Project{}, fmt.Errorf("failed to list projects: %w", err)
	}

	// populate the cache with the projects
	for _, project := range projects.Projects {
		c.projectCache.setProject(project)
	}
	// Now that the cache is populated, we can get the project by name
	project, ok := c.projectCache.getProject(projectName)
	if !ok {
		return Project{}, fmt.Errorf("project not found: %s", projectName)
	}

	return project, nil
}

func (c *Client) listProjects(ctx context.Context) (ListProjectResponse, error) {
	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects
	url := fmt.Sprintf("%s/v2/projects", strings.TrimSuffix(c.URL, "/"))

	// Send a GET request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ListProjectResponse{}, err
	}

	var projects ListProjectResponse
	err = json.Unmarshal(responseBody, &projects)
	if err != nil {
		return ListProjectResponse{}, err
	}

	return projects, nil
}

func (c *Client) UpdateProject(ctx context.Context, project Project) (Project, error) {
	// Get the project ID from the project name
	projectID, err := idFromProject(project.Name)
	if err != nil {
		return Project{}, fmt.Errorf("failed to get project ID: %w", err)
	}

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/v2/projects/%s", strings.TrimSuffix(c.URL, "/"), projectID)
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
