# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "provider_test_notification_rule" {
  display_name = "Terraform created notification rule"
  project_id   = "myProjectID"
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
            body = "Temperature $celsius°c is out of range"
            recipients = [
              "this.guy@example.com"
            ]
            subject = "Temperature alert"
          }
        }
      ]
    }
  ]
}
