# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Sensor offline"
  project_id   = dt_project.test.id

  trigger = {
    field      = "connectionStatus"
    connection = "SENSOR_OFFLINE"
  }
  trigger_delay = "900s" # 15 minutes
  escalation_levels = [
    {
      display_name = "Escalation Level 1"
      actions = [
        {
          type = "EMAIL"
          email_config = {
            body = "Sensor $name is offline"
            recipients = [
              "this.guy@example.com"
            ]
            subject = "Sensor offline alert"
          }
        }
      ]
    }
  ]
}
