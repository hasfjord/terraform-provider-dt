# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "my_notification_rule" {
  display_name = "Disabled notification rule"
  enabled      = false
  project_id   = data.dt_project.test.id
  trigger = {
    field = "temperature"
    range = {
      lower = 0
      upper = 30
    }
  }
  escalation_levels = [
    {
      display_name = "Escalation Level 1"
      actions = [
        {
          type = "EMAIL"
          email_config = {
            body = "Temperature $celsius°C is out of range"
            recipients = [
              "someone@example.com"
            ]
            subject = "Temperature Alert"
          }
        }
      ]
    }
  ]
}
