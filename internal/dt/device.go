// Copyright (c) HashiCorp, Inc.

package dt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Device struct {
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Labels        map[string]string `json:"labels"`
	ProductNumber string            `json:"productNumber"`
}

func (c *Client) GetDevice(ctx context.Context, deviceName string) (*Device, error) {
	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(c.URL, "/"), deviceName)
	responseBody, err := c.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("dt: failed to get device: %w", err)
	}

	var device Device
	if err := json.Unmarshal(responseBody, &device); err != nil {
		return nil, fmt.Errorf("dt: failed to unmarshal device: %w", err)
	}

	return &device, nil
}
