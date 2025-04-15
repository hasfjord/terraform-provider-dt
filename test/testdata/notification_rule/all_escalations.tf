# Copyright (c) HashiCorp, Inc.

resource "dt_notification_rule" "test" {
  display_name = "Cloud connector offline"
  project_id   = dt_project.test.id

  trigger = {
    field = "temperature"
    range = {
      lower = 0
      type  = "OUTSIDE"
    }
  }
  escalation_levels = [
    # TODO: Fix empty sensitive parameters
    # {
    #   display_name   = "corrigo"
    #   escalate_after = "3600s"
    #   actions = [{
    #     type = "CORRIGO"
    #     corrigo_config = {
    #       asset_id        = "asset-id"
    #       client_id       = "client-id"
    #       client_secret   = "super-secret"
    #       company_name    = "company-name"
    #       contact_address = "contact-address"
    #       contact_name    = "contact-name"
    #       customer_id     = "customer-id"
    #       sub_type_id     = "sub-type-id"
    #       task_id         = "task-id"
    #     }
    #   }]
    # },
    {
      display_name   = "email"
      escalate_after = "3600s"
      actions = [{
        type = "EMAIL"
        email_config = {
          body = "Temperature $celsius is over the limit"
          recipients = [
            "this.guy@example.com"
          ]
          subject = "Cloud connector offline alert"
        }
      }]
    },
    {
      display_name   = "phone call"
      escalate_after = "3600s"
      actions = [{
        type = "PHONE_CALL"
        phone_call_config = {
          introduction = "This is an automated call from Disruptive Technologies"
          message      = "Temperature $celsius is over the limit for device $name"
          recipients = [
            "+4798765432"
          ]
        }
      }]
    },
    {
      display_name   = "service channel"
      escalate_after = "3600s"
      actions = [{
        type = "SERVICE_CHANNEL"
        service_channel_config = {
          asset_tag_id = "asset-tag-id"
          description  = "Temperature $celsius is over the limit"
          store_id     = "store-id"
          trade        = "REFRIGERATION"
        }
      }]
    },
    # TODO: This requires an existing cloud connector device
    # This can be achieved once the provider supports creating emulators.
    # {
    #   display_name = "signal tower"
    #   actions = [{
    #     type = "SIGNAL_TOWER"
    #     signal_tower_config = {
    #       cloud_connector_name = "projects/${dt_project.test.id}/devices/123456789"
    #     }
    #   }]
    # },
    {
      display_name   = "SMS"
      escalate_after = "3600s"
      actions = [
        {
          type = "SMS"
          sms_config = {
            body = "Temperature $celsius is over the limit"
            recipients = [
              "+4798765432"
            ]
          }
        }
      ]
    },
    {
      display_name   = "webhook"
      escalate_after = "3600s"
      actions = [{
        type = "WEBHOOK"
        webhook_config = {
          url = "https://example.com/webhook"
          headers = {
            "Content-Type" : "application/json"
          }
          signature_secret = "super-secret"
        }
      }]
    }
  ]
}
