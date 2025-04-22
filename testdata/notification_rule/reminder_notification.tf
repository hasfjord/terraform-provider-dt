# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "With reminder notification"
  project_id   = dt_project.test.id
  trigger = {
    field = "relativeHumidity"
    range = {
      lower = 30
      upper = 70
    }
  }
  reminder_notification = true
  escalation_levels = [
    {
      display_name = "Escalation Level 1"
      actions = [
        {
          type = "SMS"
          sms_config = {
            body = "Relative humidity $relativeHumidity% is out of range"
            recipients = [
              "+4798765432"
            ]
          }
        }
      ]
    }
  ]
}
