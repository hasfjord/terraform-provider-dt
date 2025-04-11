// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Setup separate project for the test.
// There can only be 10 data connectors per project.
var notificationRuleProviderConfig = providerConfig + `
resource "dt_project" "test" {
	display_name = "Notification Rule Acceptance Test Project"
	organization = "organizations/cvinmt9aq9sc738g6eog"
	location = {
		time_location = "Europe/Oslo"
	}
}

`

func TestAccNotificationRulesResourceExamples(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: providerConfig + readTestFile(t, "../../examples/resources/dt_notification_rule/resource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "display_name", "Terraform created notification rule"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.field", "temperature"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.lower", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.upper", "30"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.body", "Temperature $celsius°C is out of range"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.subject", "Temperature Alert"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.recipients.0", "someone@example.com"),
				),
			},
		},
	})
}

func TestAccNotificationRuleResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../test/testdata/notification_rule/with_schedule.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Notification Rule Acceptance Test"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "device_labels.%", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "device_labels.foo", "bar"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "resolved_notification", "true"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "temperature"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.type", "OUTSIDE"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.upper", "30"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.body", "Temperature $celsius is out of range"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.subject", "Temperature Alert"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.0", "this.guy@example.com"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.subject", "Temperature Alert"),
				),
			},
			// Update and read testing
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../test/testdata/notification_rule/email_sms_escalation.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Notification Rule Acceptance Test Updated"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "temperature"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.lower", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.upper", "35"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.#", "2"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.body", "Temperature $celsius is out of range"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.subject", "Temperature Alert"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.0", "this.guy@example.com"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.display_name", "Escalation Level 2"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.type", "SMS"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.sms_config.body", "Temperature $celsius is out of range"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.sms_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.sms_config.recipients.0", "+4798765432"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../test/testdata/notification_rule/offline_trigger.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Cloud connector offline"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "connectionStatus"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.connection", "CLOUD_CONNECTOR_OFFLINE"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.body", "Cloud connector is offline"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.subject", "Cloud connector offline alert"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.0", "this.guy@example.com"),
				),
			},
		},
	})
}
