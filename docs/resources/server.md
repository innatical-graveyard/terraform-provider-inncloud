---
page_title: "server Resource - terraform-provider-inncloud"
subcategory: ""
description: |-
  The server resource allows you to configure an Innatical Cloud server.
---

# Resource `inncloud_server`

The server resource allows you to configure an Innatical Cloud server.

## Example Usage

```terraform
resource "inncloud_server" "Example" {
  name = "Example"
  model = "starter"
  image = "ubuntu-20.04"
  region = "LA1"
  cycle = "month"
}
```

## Argument Reference

- `name` - (Required) The name of the server
- `model` - (Required) The model of the server. Can be starter, pro, business, or enterprise.
- `image` - (Required) The image to build the server on. Can be ubuntu-20.0.4.
- `region` - (Required) The region of the server. Can be LA1 or HEL1.
- `cycle` - (Required) The billing cycle of the server. Can be month or year.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

### Server

- `ip` - The server's IP.
