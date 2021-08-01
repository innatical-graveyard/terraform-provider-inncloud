---
page_title: "Provider: inncloud"
subcategory: ""
description: |-
  Terraform provider for interacting with the Innatical Cloud API.
---

# inncloud Provider

The inncloud provider allows you to manage your Innatical Cloud resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "inncloud" {
  token = "owo"
  project_id = "b2e677d3-d340-4462-86fe-bc7efcbd2d95"
}
```

## Schema

### Optional

- **token** (String, Optional) Token to authenticate to Innatical Cloud API
- **project_id** (String, Optional) Project to use with the Innatical Cloud API
