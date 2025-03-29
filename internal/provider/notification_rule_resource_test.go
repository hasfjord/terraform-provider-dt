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
				Config: providerConfig + readExampleFile(t, "../../examples/resources/dt_notification_rule/resource.tf"),
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
				Config: notificationRuleProviderConfig + `
				resource "dt_notification_rule" "test" {
					display_name = "Notification Rule Acceptance Test"
					project_id = dt_project.test.id
					device_labels = {
						foo = "bar"
					}
					trigger = {
						field = "temperature"
						range = {
							lower = 0
							upper = 30
							}
						}
					escalation_levels =[
						{
							display_name = "Escalation Level 1"
							actions = [
								{
									type = "EMAIL"
									email_config = {
										body = "Temperature $celsius is out of range"
										recipients = [
											"this.guy@example.com"
										]
										subject = "Temperature Alert"
									}
								}
							]
						}
					]
				}

				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Notification Rule Acceptance Test"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "device_labels.%", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "device_labels.foo", "bar"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "temperature"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.lower", "0"),
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
				Config: notificationRuleProviderConfig + `
				resource "dt_notification_rule" "test" {
					display_name = "Notification Rule Acceptance Test Updated"
					project_id = dt_project.test.id
					trigger = {
						field = "temperature"
						range = {
							lower = 0
							upper = 35
						}
					}
					escalation_levels =[
						{
							display_name = "Escalation Level 1"
							actions = [
								{
									type = "EMAIL"
									email_config = {
										body = "Temperature $celsius is out of range"
										recipients = [
											"this.guy@example.com"
										]
										subject = "Temperature Alert"
									}
								}
							]
							escalate_after = "3600s"
						},
						{
							display_name = "Escalation Level 2"
							actions = [
								{
									type = "SMS"
									sms_config = {
										body = "Temperature $celsius is out of range"
										recipients = [
											"+4798765432"
										]
									}
								}
							]
						}
					]
				}

				`,
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
}
