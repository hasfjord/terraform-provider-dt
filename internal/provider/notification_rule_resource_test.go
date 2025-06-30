// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Setup separate project for the test.
// There can only be 10 data connectors per project.
var notificationRuleProviderConfig = providerConfig + `
data "dt_project" "test" {
	name = "projects/d0919uq3tjjs739bf18g"
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
					resource.TestCheckNoResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.filter.product_equivalent_temperature.%"),
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
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/with_schedule.tf"),
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
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
			// Import testing
			{
				ResourceName:                         "dt_notification_rule.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_notification_rule.test"].Primary.Attributes["name"], nil
				},
			},
			// Update and read testing
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/email_sms_escalation.tf"),
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
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		// Test case for the ccon offline trigger
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/ccon_offline_trigger.tf"),
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
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		// Test case for the sensor offline trigger
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/sensor_offline_trigger.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Sensor offline"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "connectionStatus"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.connection", "SENSOR_OFFLINE"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger_delay", "900s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.body", "Sensor $name is offline"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.subject", "Sensor offline alert"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.email_config.recipients.0", "this.guy@example.com"),
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test case for the disabled rule
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/disabled_rule.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "display_name", "Disabled notification rule"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "enabled", "false"),
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
					resource.TestCheckNoResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test case for reminder notifications
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/reminder_notification.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "With reminder notification"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "relativeHumidity"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.lower", "30"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.upper", "70"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "reminder_notification", "true"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "SMS"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.sms_config.body", "Relative humidity $relativeHumidity% is out of range"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.sms_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.sms_config.recipients.0", "+4798765432"),
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test case for "all" escalation types
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/all_escalations.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "All escalation types"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "temperature"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.lower", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.range.type", "OUTSIDE"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.#", "6"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "corrigo"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.escalate_after", "3600s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "2"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "CORRIGO"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.asset_id", "asset-id-1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.client_id", "client-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.client_secret", "super-secret"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.company_name", "company-name"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.contact_address", "contact-address"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.contact_name", "contact-name"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.customer_id", "customer-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.sub_type_id", "sub-type-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.corrigo_config.task_id", "task-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.type", "CORRIGO"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.asset_id", "asset-id-2"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.client_id", "client-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.client_secret", "super-secret"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.company_name", "company-name"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.contact_address", "contact-address"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.contact_name", "contact-name"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.customer_id", "customer-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.sub_type_id", "sub-type-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.1.corrigo_config.task_id", "task-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.display_name", "email"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.escalate_after", "3600s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.email_config.body", "Temperature $celsius is over the limit"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.email_config.subject", "Temperature Alert"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.1.actions.0.email_config.recipients.0", "someone@example.com"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.display_name", "phone call"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.escalate_after", "3600s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.actions.0.type", "PHONE_CALL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.actions.0.phone_call_config.introduction", "This is an automated call from Disruptive Technologies"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.actions.0.phone_call_config.message", "Temperature $celsius is over the limit for device $name"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.actions.0.phone_call_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.2.actions.0.phone_call_config.recipients.0", "+4798765432"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.display_name", "service channel"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.escalate_after", "3600s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.actions.0.type", "SERVICE_CHANNEL"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.actions.0.service_channel_config.asset_tag_id", "asset-tag-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.actions.0.service_channel_config.description", "Temperature $celsius is over the limit"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.actions.0.service_channel_config.store_id", "store-id"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.3.actions.0.service_channel_config.trade", "REFRIGERATION"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.display_name", "SMS"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.escalate_after", "3600s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.actions.0.type", "SMS"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.actions.0.sms_config.body", "Temperature $celsius is over the limit"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.actions.0.sms_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.4.actions.0.sms_config.recipients.0", "+4798765432"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.display_name", "webhook"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.actions.0.type", "WEBHOOK"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.actions.0.webhook_config.url", "https://example.com/webhook"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.actions.0.webhook_config.headers.%", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.actions.0.webhook_config.headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.5.actions.0.webhook_config.signature_secret", "super-secret"),
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// test case for signal tower
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/signal_tower.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Signal Tower"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "connectionStatus"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.connection", "CLOUD_CONNECTOR_OFFLINE"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.display_name", "signal tower"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.escalate_after", "3600s"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.type", "SIGNAL_TOWER"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "escalation_levels.0.actions.0.signal_tower_config.cloud_connector_name", "projects/d0919uq3tjjs739bf18g/devices/emud091aassh1nc738nel0g"),
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// test case for inverse schedule
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/inverse_schedule.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.test", "display_name", "Off Hours Schedule"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.timezone", "Europe/Oslo"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.#", "2"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.day_of_week.#", "5"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.day_of_week.0", "Monday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.day_of_week.1", "Tuesday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.day_of_week.2", "Wednesday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.day_of_week.3", "Thursday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.day_of_week.4", "Friday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.time_range.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.time_range.0.start.hour", "8"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.time_range.0.start.minute", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.time_range.0.end.hour", "20"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.0.time_range.0.end.minute", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.day_of_week.#", "2"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.day_of_week.0", "Saturday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.day_of_week.1", "Sunday"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.time_range.0.start.hour", "10"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.time_range.0.start.minute", "30"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.time_range.0.end.hour", "18"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.slots.1.time_range.0.end.minute", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "schedule.inverse", "true"),
					resource.TestCheckResourceAttr("dt_notification_rule.test", "trigger.field", "temperature"),
					resource.TestCheckNoResourceAttr("dt_notification_rule.test", "trigger.range.filter.product_equivalent_temperature.%"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test case for the disabled rule
			{
				Config: notificationRuleProviderConfig + readTestFile(t, "../../testdata/notification_rule/pet_filter.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "display_name", "Range with PET filter"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "enabled", "true"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.field", "temperature"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.lower", "-10"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.upper", "6"),
					// check that the PET filter is set
					resource.TestCheckResourceAttrSet("dt_notification_rule.my_notification_rule", "trigger.range.filter.product_equivalent_temperature.%"),
					// check that the PET filter is empty
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "trigger.range.filter.product_equivalent_temperature.#", "0"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.display_name", "Escalation Level 1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.type", "EMAIL"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.body", "The glycol buffer $name is out of range with a temperature of $celsius°C"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.subject", "Temperature Alert $name"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.recipients.#", "1"),
					resource.TestCheckResourceAttr("dt_notification_rule.my_notification_rule", "escalation_levels.0.actions.0.email_config.recipients.0", "someone@example.com"),
				),
			},
		},
	})
}
