# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "my_notification_rule" {
  display_name = "Range with PET filter"
  project_id   = data.dt_project.test.id
  trigger = {
    field = "temperature"
    range = {
      lower = -10
      upper = 6
      filter = {
        product_equivalent_temperature = {
          // This is an empty object, as the ProductEquivalentTemperature filter does not have any attributes.
        }
      }
    }
  }
  escalation_levels = [
    {
      display_name = "Escalation Level 1"
      actions = [
        {
          type = "EMAIL"
          email_config = {
            body = "The glycol buffer $name is out of range with a temperature of $celsius°C"
            recipients = [
              "someone@example.com"
            ]
            subject = "Temperature Alert $name"
          }
        }
      ]
    }
  ]
}
