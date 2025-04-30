# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Cloud connector offline"
  project_id   = data.dt_project.test.id

  trigger = {
    field      = "connectionStatus"
    connection = "CLOUD_CONNECTOR_OFFLINE"
  }
  escalation_levels = [
    {
      display_name = "Escalation Level 1"
      actions = [
        {
          type = "EMAIL"
          email_config = {
            body = "Cloud connector is offline"
            recipients = [
              "this.guy@example.com"
            ]
            subject = "Cloud connector offline alert"
          }
        }
      ]
    }
  ]
}
