---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dt Provider"
subcategory: ""
description: |-
  
---

# dt Provider



## Example Usage

```terraform
# Copyright (c) HashiCorp, Inc.

provider "dt" {
  url            = "https://api.disruptive-technologies.com"
  emulator_url   = "https://emulator.disruptive-technologies.com"
  token_endpoint = "https://identity.disruptive-technologies.com/oauth2/token"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `email` (String) The email address used to authenticate with the OIDC provider.
- `emulator_url` (String) The URL of the emulator server.
- `key_id` (String) The key ID from the service account.
- `key_secret` (String, Sensitive) The key secret from the service account.
- `token_endpoint` (String) The token endpoint for the OIDC provider.
- `url` (String) The URL of the API server.
