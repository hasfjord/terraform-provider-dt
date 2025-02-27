package dt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DataConnector struct {
	Name                  string                 `json:"name"`
	Type                  string                 `json:"type"`
	DisplayName           string                 `json:"displayName"`
	Status                string                 `json:"status"`
	Events                []string               `json:"events"`
	Labels                []string               `json:"labels"`
	HTTPConfig            *HTTPConfig            `json:"httpConfig"`
	AzureServiceBusConfig *AzureServiceBusConfig `json:"azureServiceBusConfig"`
	AzureEventHubConfig   *AzureEventHubConfig   `json:"azureEventHubConfig"`
	PubsubConfig          *PubsubConfig          `json:"pubsubConfig"`
	AWSSQSConfig          *AWSSQSConfig          `json:"awsSqsConfig"`
}

type HTTPConfig struct {
	Url             string            `json:"url"`
	SignatureSecret string            `json:"signatureSecret"`
	Headers         map[string]string `json:"headers"`
}

type AzureServiceBusConfig struct {
	URL                  string               `json:"url"`
	AuthenticationConfig AuthenticationConfig `json:"authenticationConfig"`
	BrokerProperties     BrokerProperties     `json:"brokerProperties"`
}

type BrokerProperties struct {
	CorrelationID string `json:"correlationId"`
}

type AzureEventHubConfig struct {
	URL                  string               `json:"url"`
	AuthenticationConfig AuthenticationConfig `json:"authenticationConfig"`
}

type AuthenticationConfig struct {
	TenantID string `json:"tenantId"`
	ClientID string `json:"clientId"`
}

type PubsubConfig struct {
	Topic    string `json:"topic"`
	Audience string `json:"audience"`
}

type AWSSQSConfig struct {
	QueueUrl   string `json:"queueUrl"`
	AwsRoleArn string `json:"awsRoleArn"`
	Audience   string `json:"audience"`
}

// GetDatConnector retrieves a data connector by name.
func (c *Client) GetDataConnector(ctx context.Context, dataConnector string) (DataConnector, error) {
	projectID, dataConnectorID, err := parseDataConnectorResourceName(dataConnector)
	if err != nil {
		return DataConnector{}, err
	}
	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}/dataconnectors/{data_connector_id}
	url := fmt.Sprintf("%s/projects/%s/dataconnectors/%s", strings.TrimSuffix(c.URL, "/"), projectID, dataConnectorID)

	// Send a GET request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return DataConnector{}, err
	}

	var dc DataConnector
	err = json.Unmarshal(responseBody, &dc)
	if err != nil {
		return DataConnector{}, err
	}

	return dc, nil
}

type CreateDataConnectorRequest struct {
	DisplayName           string                 `json:"displayName"`
	Type                  string                 `json:"type"`
	Status                string                 `json:"status"`
	Events                []string               `json:"events"`
	Labels                []string               `json:"labels"`
	HTTPConfig            *HTTPConfig            `json:"httpConfig"`
	AzureServiceBusConfig *AzureServiceBusConfig `json:"azureServiceBusConfig"`
	AzureEventHubConfig   *AzureEventHubConfig   `json:"azureEventHubConfig"`
	PubsubConfig          *PubsubConfig          `json:"pubsubConfig"`
	AWSSQSConfig          *AWSSQSConfig          `json:"awsSqsConfig"`
}

func dataConnectorToCreateDataConnectorRequest(dc DataConnector) CreateDataConnectorRequest {
	return CreateDataConnectorRequest{
		DisplayName:           dc.DisplayName,
		Type:                  dc.Type,
		Status:                dc.Status,
		Events:                dc.Events,
		Labels:                dc.Labels,
		HTTPConfig:            dc.HTTPConfig,
		AzureServiceBusConfig: dc.AzureServiceBusConfig,
		AzureEventHubConfig:   dc.AzureEventHubConfig,
		PubsubConfig:          dc.PubsubConfig,
		AWSSQSConfig:          dc.AWSSQSConfig,
	}
}

// CreateDataConnector creates a new data connector.
func (c *Client) CreateDataConnector(ctx context.Context, projectID string, dataConnector DataConnector) (DataConnector, error) {
	// Convert the DataConnector struct to a CreateDataConnectorRequest struct to remove the generated Name field and add the ProjectID field.
	request := dataConnectorToCreateDataConnectorRequest(dataConnector)

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}/dataconnectors
	url := fmt.Sprintf("%s/projects/%s/dataconnectors", strings.TrimSuffix(c.URL, "/"), projectID)
	body, err := json.Marshal(request)
	if err != nil {
		return DataConnector{}, err
	}

	// Send a POST request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return DataConnector{}, err
	}

	// Unmarshal the response body to a DataConnector struct
	var newDC DataConnector
	err = json.Unmarshal(responseBody, &newDC)
	if err != nil {
		return DataConnector{}, err
	}

	return newDC, nil
}

// UpdateDataConnector updates an existing data connector.
func (c *Client) UpdateDataConnector(ctx context.Context, dc DataConnector) (DataConnector, error) {
	projectID, dataConnectorID, err := parseDataConnectorResourceName(dc.Name)
	if err != nil {
		return DataConnector{}, err
	}

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}/dataconnectors/{data_connector_id}
	url := fmt.Sprintf("%s/projects/%s/dataconnectors/%s", strings.TrimSuffix(c.URL, "/"), projectID, dataConnectorID)
	body, err := json.Marshal(dc)
	if err != nil {
		return DataConnector{}, err
	}

	// Send a PATCH request to the API
	responseBody, err := c.DoRequest(ctx, http.MethodPatch, url, strings.NewReader(string(body)))
	if err != nil {
		return DataConnector{}, err
	}

	// Unmarshal the response body to a DataConnector struct
	var updatedDC DataConnector
	err = json.Unmarshal(responseBody, &updatedDC)
	if err != nil {
		return DataConnector{}, err
	}

	return updatedDC, nil
}

// DeleteDataConnector deletes a data connector.
func (c *Client) DeleteDataConnector(ctx context.Context, dataConnector string) error {
	projectID, dataConnectorID, err := parseDataConnectorResourceName(dataConnector)
	if err != nil {
		return err
	}

	// Create the URL for the API request: https://api.disruptive-technologies.com/v2/projects/{project_id}/dataconnectors/{data_connector_id}
	url := fmt.Sprintf("%s/projects/%s/dataconnectors/%s", strings.TrimSuffix(c.URL, "/"), projectID, dataConnectorID)

	// Send a DELETE request to the API
	_, err = c.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (d DataConnector) ProjectID() string {
	projectID, _, _ := parseDataConnectorResourceName(d.Name)
	return projectID
}

func (d DataConnector) DataConnectorID() string {
	_, dataConnectorID, _ := parseDataConnectorResourceName(d.Name)
	return dataConnectorID
}

// parseDataConnectorResourceName parses a data connector resource name into projectID and dataConnectorID.
func parseDataConnectorResourceName(name string) (projectID string, dataConnectorID string, err error) {
	parts := strings.Split(name, "/")
	if len(parts) != 4 {
		return "", "", fmt.Errorf("invalid data connector name: %s", name)
	}
	return parts[1], parts[3], nil
}
