---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dt_project Data Source - dt"
subcategory: ""
description: |-
  
---

# dt_project (Data Source)



## Example Usage

```terraform
# Copyright (c) HashiCorp, Inc.

data "dt_project" "test_project" {
  provider = disruptive-technologies
  name     = "projects/your-project-id"
}

output "project" {
  value = data.dt_project.test_project
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The resource name of the project. On the form `projects/{project_id}`.

### Read-Only

- `cloud_connector_count` (Number) The number of cloud connectors in the project.
- `display_name` (String) The display name of the project.
- `id` (String) The resource ID of the project.
- `inventory` (Boolean) Whether the project is an inventory project.
- `location` (Object) (see [below for nested schema](#nestedatt--location))
- `organization` (String) The reource name of the organization that the project belongs to. on the form `organizations/{organization_id}`.
- `organization_display_name` (String) The display name of the organization that the project belongs to.
- `sensor_count` (Number) The number of sensors in the project.

<a id="nestedatt--location"></a>
### Nested Schema for `location`

Read-Only:

- `latitude` (Number)
- `longitude` (Number)
- `time_location` (String)
