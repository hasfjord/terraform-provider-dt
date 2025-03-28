# Copyright (c) HashiCorp, Inc.

resource "dt_project" "my_project" {
  display_name = "Notification Rule Acceptance Test Project"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  location = {
    time_location = "Europe/Oslo"
  }
}

resource "dt_notification_rule" "my_notification_rule" {
  display_name = "Terraform created notification rule"
  project_id   = dt_project.my_project.id
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
