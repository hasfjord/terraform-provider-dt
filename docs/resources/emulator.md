---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dt_emulator Resource - dt"
subcategory: ""
description: |-
  
---

# dt_emulator (Resource)



## Example Usage

```terraform
# Copyright (c) HashiCorp, Inc.

resource "dt_emulator" "my_emulator" {
  display_name = "Terraform created emulator"
  project_id   = "d0ito5m62hus73ae3lr0"
  type         = "touch"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the emulator.
- `project_id` (String) The project ID to create the emulator in.
- `type` (String) The type of emulator valid types are: touch, temperature, proximity, touchCounter, proximityCounter, humidity, waterDetector, co2, motion, contact, deskOccupancy, ccon

### Optional

- `labels` (Map of String) A map of labels to assign to the emulator.

### Read-Only

- `name` (String) The resource name of the emulator on the form: `projects/{project_id}/devices/{device_id}`
- `system_labels` (Map of String) A map of system labels assigned to the emulator. Read only
