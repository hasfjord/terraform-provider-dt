package dt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	username   string
	password   string
	URL        string
	httpClient http.Client
}

func NewClient(URL string) *Client {
	return &Client{
		URL:        URL,
		httpClient: *http.DefaultClient,
	}
}

func (c *Client) WithBasicAuth(username, password string) *Client {
	c.username = username
	c.password = password
	return c
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
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), idFromProject(project))
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Project{}, err
	}

	ctx = tflog.SetField(ctx, "project", project)
	ctx = tflog.SetField(ctx, "url", url)
	ctx = tflog.SetField(ctx, "username", c.username)
	ctx = tflog.SetField(ctx, "password_len", len(c.password))

	tflog.Debug(ctx, "sending request")

	request.SetBasicAuth(c.username, c.password)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Project{}, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Project{}, err
	}
	if response.StatusCode != http.StatusOK {
		return Project{}, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var p Project
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

func (c *Client) UpdateProject(ctx context.Context, project Project) (Project, error) {
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), idFromProject(project.Name))
	body, err := json.Marshal(project)
	if err != nil {
		return Project{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return Project{}, err
	}

	request.SetBasicAuth(c.username, c.password)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Project{}, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Project{}, err
	}
	if response.StatusCode != http.StatusOK {
		return Project{}, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var p Project
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

func (c *Client) CreateProject(ctx context.Context, project Project) (Project, error) {
	url := fmt.Sprintf("%s/projects", strings.TrimSuffix(c.URL, "/"))
	body, err := json.Marshal(project)
	if err != nil {
		return Project{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Project{}, err
	}

	request.SetBasicAuth(c.username, c.password)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Project{}, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Project{}, err
	}
	if response.StatusCode != http.StatusOK {
		return Project{}, &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var p Project
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil {
		return Project{}, err
	}

	return p, nil
}

func (c *Client) DeleteProject(ctx context.Context, project string) error {
	url := fmt.Sprintf("%s/projects/%s", strings.TrimSuffix(c.URL, "/"), idFromProject(project))
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	request.SetBasicAuth(c.username, c.password)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return &HTTPError{
			StatusCode: response.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	return nil
}

func idFromProject(project string) string {
	return strings.Split(project, "/")[1]
}
