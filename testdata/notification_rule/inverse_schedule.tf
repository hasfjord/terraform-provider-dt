# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Off Hours Schedule"
  project_id   = data.dt_project.test.id
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
      },
      {
        day_of_week = ["Saturday", "Sunday"]
        time_range = [{
          start = {
            hour   = 10
            minute = 30
          }
          end = {
            hour   = 18
            minute = 0
          }
        }]
      }
    ]
    inverse = true
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
