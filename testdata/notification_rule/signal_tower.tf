# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Signal Tower"
  project_id   = data.dt_project.test.id

  trigger = {
    field      = "connectionStatus"
    connection = "CLOUD_CONNECTOR_OFFLINE"
  }
  resolved_notification = true
  escalation_levels = [
    {
      display_name   = "signal tower"
      escalate_after = "3600s"
      actions = [{
        type = "SIGNAL_TOWER"
        signal_tower_config = {
          cloud_connector_name = "projects/${data.dt_project.test.id}/devices/emud091aassh1nc738nel0g"
        }
      }]
    }
  ]
}
