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

// DISCLAIMER: The Notification Rule API is not released yet and is subject to change.

// NotificationRule represents a notification rule in the Disruptive Technologies platform.
type NotificationRule struct {
	Name                 string            `json:"name"`
	Enabled              bool              `json:"enabled"`
	DisplayName          string            `json:"displayName"`
	Devices              []string          `json:"devices"`
	DeviceLables         map[string]string `json:"deviceLabels"`
	Trigger              Trigger           `json:"trigger"`
	EscalationLevels     []EscalationLevel `json:"escalationLevels"`
	Schedule             *Schedule         `json:"schedule"`
	TriggerDelay         *string           `json:"triggerDelay"`
	ReminderNotification bool              `json:"reminderNotifications"`
	ResolvedNotification bool              `json:"resolvedNotifications"`
	UnacknowledgesAfter  *string           `json:"unacknowledgesAfter"`
	// Deprecated: Use EscalationLevels instead, included for completeness.
	Actions []NotificationAction `json:"actions"`
}

// EscalationLevel represents an escalation level in a notification rule.
type EscalationLevel struct {
	DisplayName   string               `json:"displayName"`
	Actions       []NotificationAction `json:"actions"`
	EscalateAfter *string              `json:"escalateAfter"`
}

// Note: some of these notification types are not available for all customers.
type NotificationAction struct {
	Type                 string                `json:"type"`
	SMSConfig            *SMSConfig            `json:"sms"`
	EmailConfig          *EmailConfig          `json:"email"`
	CorrigoConfig        *CorrigoConfig        `json:"corrigo"`
	ServiceChannelConfig *ServiceChannelConfig `json:"serviceChannel"`
	WebhookConfig        *WebhookConfig        `json:"webhook"`
	PhoneCallConfig      *PhoneCallConfig      `json:"phoneCall"`
	SignalTowerConfig    *SignalTowerConfig    `json:"signalTower"`
}

type SMSConfig struct {
	Recipients []string `json:"recipients"`
	Message    string   `json:"message"`
}

type EmailConfig struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
}

type CorrigoConfig struct {
	AssetID              string `json:"assetId"`
	TaskID               string `json:"taskId"`
	CustomerID           string `json:"customerId"`
	ClientSecret         string `json:"clientSecret"`
	CompanyName          string `json:"companyName"`
	SubTypeID            string `json:"subTypeId"`
	ContactName          string `json:"contactName"`
	ContactAddress       string `json:"contactAddress"`
	WorkOrderDescription string `json:"workOrderDescription"`
	StudioDashboardURL   string `json:"studioDashboardUrl"`
}

type ServiceChannelConfig struct {
	StoreID     string `json:"storeId"`
	AssetTagID  string `json:"assetTagId"`
	Trade       string `json:"trade"`
	Description string `json:"description"`
}

type WebhookConfig struct {
	URL             string            `json:"url"`
	SignatureSecret string            `json:"signatureSecret"`
	Headers         map[string]string `json:"headers"`
}

type PhoneCallConfig struct {
	Recipients   []string `json:"recipients"`
	Introduction string   `json:"introduction"`
	Message      string   `json:"message"`
}

type SignalTowerConfig struct {
	CloudConnectorName string `json:"cloudConnectorName"`
}

type Trigger struct {
	Field        string  `json:"field"`
	Range        *Range  `json:"range"`
	Presence     *string `json:"presence"`
	Motion       *string `json:"motion"`
	Occupancy    *string `json:"occupancy"`
	Connection   *string `json:"connection"`
	Contact      *string `json:"contact"`
	TriggerCount int32   `json:"triggerCount"`
}

type Range struct {
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
	Type  string  `json:"type"`
}

type Schedule struct {
	Timezone string `json:"timezone"`
	Slots    []Slot `json:"slots"`
	Inverse  bool   `json:"inverse"`
}

type Slot struct {
	DaysOfWeek []string    `json:"day"`
	TimeRange  []TimeRange `json:"times"`
}

type TimeRange struct {
	Start TimeOfDay `json:"start"`
	End   TimeOfDay `json:"end"`
}

type TimeOfDay struct {
	Hour   int32 `json:"hour"`
	Minute int32 `json:"minute"`
}

// GetNotificationRule returns a notification rule by resource name.
func (c *Client) GetNotificationRule(ctx context.Context, name string) (NotificationRule, error) {
	projectID, ruleID, err := ParseResourceName(name)
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to parse resource name: %w", err)
	}

	url := fmt.Sprintf("%s/v2alpha/projects/%s/rules/%s", strings.TrimSuffix(c.URL, "/"), projectID, ruleID)
	responseBody, err := c.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to get notification rule: %w", err)
	}

	var rule NotificationRule
	if err := json.Unmarshal(responseBody, &rule); err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to unmarshal notification rule: %w", err)
	}

	return rule, nil
}

// CreateNotificationRule creates a new notification rule.
func (c *Client) CreateNotificationRule(ctx context.Context, projectID string, rule NotificationRule) (NotificationRule, error) {
	url := fmt.Sprintf("%s/v2alpha/projects/%s/rules", strings.TrimSuffix(c.URL, "/"), projectID)

	body, err := json.Marshal(rule)
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to marshal notification rule: %w", err)
	}

	responseBody, err := c.DoRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to create notification rule: %w", err)
	}

	var createdRule NotificationRule
	if err := json.Unmarshal(responseBody, &createdRule); err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to unmarshal created notification rule: %w", err)
	}

	return createdRule, nil
}

// UpdateNotificationRule updates an existing notification rule.
func (c *Client) UpdateNotificationRule(ctx context.Context, rule NotificationRule) (NotificationRule, error) {
	projectID, ruleID, err := ParseResourceName(rule.Name)
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to parse resource name: %w", err)
	}

	url := fmt.Sprintf("%s/v2alpha/projects/%s/rules/%s", strings.TrimSuffix(c.URL, "/"), projectID, ruleID)

	body, err := json.Marshal(rule)
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to marshal notification rule: %w", err)
	}

	responseBody, err := c.DoRequest(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to update notification rule: %w", err)
	}

	var updatedRule NotificationRule
	if err := json.Unmarshal(responseBody, &updatedRule); err != nil {
		return NotificationRule{}, fmt.Errorf("dt: failed to unmarshal updated notification rule: %w", err)
	}

	return updatedRule, nil
}

// DeleteNotificationRule deletes a notification rule.
func (c *Client) DeleteNotificationRule(ctx context.Context, name string) error {
	projectID, ruleID, err := ParseResourceName(name)
	if err != nil {
		return fmt.Errorf("dt: failed to parse resource name: %w", err)
	}

	url := fmt.Sprintf("%s/v2alpha/projects/%s/rules/%s", strings.TrimSuffix(c.URL, "/"), projectID, ruleID)
	_, err = c.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("dt: failed to delete notification rule: %w", err)
	}

	return nil
}

// ParseResourceName is a helper function to parse the resource name projects/{projectID}/rules/{ruleID}
// into projectID and notificationRuleID.
func ParseResourceName(name string) (string, string, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 4 {
		return "", "", fmt.Errorf("dt: invalid resource name: %s", name)
	}
	return parts[1], parts[3], nil
}
