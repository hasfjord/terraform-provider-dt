// Copyright (c) HashiCorp, Inc.

package dt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Emulator struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

func (c *Client) GetEmulator(ctx context.Context, name string) (Emulator, error) {
	projectID, deviceID, err := ParseResourceName(name)
	if err != nil {
		return Emulator{}, err
	}

	url := c.EmulatorURL + "/v2/projects/" + projectID + "/devices/" + deviceID
	responseBody, err := c.DoRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		return Emulator{}, err
	}

	var emulator Emulator
	if err := json.Unmarshal(responseBody, &emulator); err != nil {
		return Emulator{}, err
	}
	return emulator, nil
}

func (c *Client) CreateEmulator(ctx context.Context, projectID string, emulatorToBeCreated Emulator) (Emulator, error) {
	body, err := json.Marshal(emulatorToBeCreated)
	if err != nil {
		return Emulator{}, err
	}

	url := c.EmulatorURL + "/v2/projects/" + projectID + "/devices"
	responseBody, err := c.DoRequest(ctx, "POST", url, body, nil)
	if err != nil {
		return Emulator{}, err
	}

	var createdEmulator Emulator
	if err := json.Unmarshal(responseBody, &createdEmulator); err != nil {
		return Emulator{}, err
	}
	return createdEmulator, nil
}

func (c *Client) DeleteEmulator(ctx context.Context, name string) error {
	projectID, deviceID, err := ParseResourceName(name)
	if err != nil {
		return err
	}

	url := c.EmulatorURL + "/v2/projects/" + projectID + "/devices/" + deviceID
	_, err = c.DoRequest(ctx, "DELETE", url, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateEmulator(ctx context.Context, emulator Emulator) (Emulator, error) {
	projectID, deviceID, err := ParseResourceName(emulator.Name)
	if err != nil {
		return Emulator{}, err
	}

	url := c.EmulatorURL + "/v2/projects/" + projectID + "/devices/" + deviceID

	body, err := json.Marshal(emulator)
	if err != nil {
		return Emulator{}, err
	}

	responseBody, err := c.DoRequest(ctx, "PUT", url, body, nil)
	if err != nil {
		return Emulator{}, err
	}

	var updatedEmulator Emulator
	if err := json.Unmarshal(responseBody, &updatedEmulator); err != nil {
		return Emulator{}, err
	}
	return updatedEmulator, nil
}

func (e *Emulator) ProjectID() string {
	projectID, _, _ := parseEmulatorResourceName(e.Name)
	return projectID
}

func (e *Emulator) DeviceID() string {
	_, deviceID, _ := parseEmulatorResourceName(e.Name)
	return deviceID
}

func parseEmulatorResourceName(name string) (string, string, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 4 {
		return "", "", fmt.Errorf("invalid resource name: %s", name)
	}
	return parts[1], parts[3], nil
}
