# Unofficial Disruptive Technologies Terraform Provider
## Disclaimer
This provider is not officially supported by Disruptive Technologies or Hashicorp, it is a hobby project and should be used at your own risk.

## Features
The provider currently supports the following resources and data sources:

- [x] Device data source
- [x] Data Connector resource
- [ ] Data Connector data source
- [ ] Labels Resource
- [ ] Labels Data Source
- [ ] Organization Data Source
- [x] Project data source
- [x] Project resource
- [x] Rules Resource
- [ ] Rules Data Source

## Usage

The provider requires a DT service account. Se how to setup a service account [here](https://disruptive.gitbook.io/docs/service-accounts/creating-a-service-account).

The provider requires the following variables to be set:
- `DT_API_KEY_ID` - The ID for the DT Service Account key
- `DT_API_KEY_SECRET` - The secret for the DT Service Account key
- `DT_OIDC_EMAIL` - The email for the DT Service Account

These variables are sensitive and should not be committed to version control.

Here is an example of how to configure the provider:

```hcl
terraform {
  required_providers {
    disruptive-technologies = {
      source = "registry.terraform.io/hasfjord/dt"
    }
  }
}

provider "disruptive-technologies" {
  url            = "https://api.disruptive-technologies.com"
  token_endpoint = "https://identity.disruptive-technologies.com/oauth2/token"
}
```

See the [examples](examples) directory for example usage.
