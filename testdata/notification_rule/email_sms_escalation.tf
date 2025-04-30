# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Notification Rule Acceptance Test Updated"
  project_id   = data.dt_project.test.id
  trigger = {
    field = "temperature"
    range = {
      lower = 0
      upper = 35
    }
  }
  escalation_levels = [
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
