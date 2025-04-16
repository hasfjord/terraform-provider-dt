# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Notification Rule Acceptance Test"
  project_id   = dt_project.test.id
  device_labels = {
    foo = "bar"
  }
  resolved_notification = true
  schedule = {
    timezone = "Europe/Oslo"
    slots = [
      {
        day_of_week = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
        time_range = [{
          start = {
            hour   = 8
            minute = 0
          }
          end = {
            hour   = 20
            minute = 0
          }
        }]
      }
    ]
  }
  trigger = {
    field = "temperature"
    range = {
      upper = 30
      type  = "OUTSIDE"
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
    }
  ]
}
