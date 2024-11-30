package dt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt/oidc"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	URL        string
	httpClient http.Client
	oidc       *oidc.Client
}

type Config struct {
	Oidc oidc.Config
	URL  string
}

func NewClient(cfg Config) *Client {
	return &Client{
		URL:        cfg.URL,
		httpClient: *http.DefaultClient,
		oidc:       oidc.NewClient(cfg.Oidc),
	}
}

func (c *Client) WithHttpClient(httpClient http.Client) *Client {
	c.httpClient = httpClient
	return c
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

type Location struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TimeLocation string  `json:"timeLocation"`
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d: %s", e.StatusCode, e.Body)
}

func (c *Client) GetProject(ctx context.Context, project string) (Project, error) {
	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), idFromProject(project))

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
	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), idFromProject(project.Name))
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
	url := fmt.Sprintf("%s/projects", strings.TrimSuffix(c.URL, "/"))

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
	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), idFromProject(project))

	// Send a DELETE request to the API
	_, err := c.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DoRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "method", method)
	ctx = tflog.SetField(ctx, "url", url)

	tflog.Debug(ctx, "sending request to DT API")

	// Get an OIDC token and set it as a Bearer token in the request
	token, err := c.oidc.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to get OIDC token: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to send request: %w", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to read response body: %w, status: %d", err, response.StatusCode)
	}
	if response.StatusCode != http.StatusOK {
		ctx = tflog.SetField(ctx, "status_code", response.StatusCode)
		ctx = tflog.SetField(ctx, "body", string(bodyBytes))
		tflog.Debug(ctx, "received non-200 status code from DT API")
		return nil, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	return bodyBytes, nil
}

func idFromProject(project string) string {
	return strings.Split(project, "/")[1]
}
