# Terraform Provider for LINE LIFF

[![Terraform Registry](https://img.shields.io/badge/terraform-registry-blueviolet.svg)](https://registry.terraform.io/providers/sugarshin/liff/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sugarshin/terraform-provider-liff)](https://goreportcard.com/report/github.com/sugarshin/terraform-provider-liff)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

A Terraform provider for managing [LINE LIFF (LINE Front-end Framework)](https://developers.line.biz/en/docs/liff/overview/) applications.

LIFF is a platform provided by LINE for developing web apps that run within the LINE app. This provider allows you to manage LIFF apps as infrastructure-as-code using [LIFF Server API](https://developers.line.biz/en/reference/liff-server/).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider plugin)

## Installation

```hcl
terraform {
  required_providers {
    liff = {
      source  = "sugarshin/liff"
      version = "~> 0.0"
    }
  }
}
```

## Authentication

The provider supports two authentication methods. All credentials can be configured via the provider block or environment variables.

### Option 1: Channel ID + Channel Secret (Recommended)

The provider automatically issues a [Stateless Channel Access Token](https://developers.line.biz/en/docs/messaging-api/generate-json-web-token/) (valid for 15 minutes, no issuance limit). This is the recommended approach for CI/CD pipelines.

```hcl
provider "liff" {
  channel_id     = var.line_channel_id
  channel_secret = var.line_channel_secret
}
```

Or via environment variables:

```sh
export LIFF_CHANNEL_ID="your-channel-id"
export LIFF_CHANNEL_SECRET="your-channel-secret"
```

### Option 2: Channel Access Token (Direct)

If you already have a Channel Access Token, you can specify it directly. This takes precedence over Option 1.

```hcl
provider "liff" {
  channel_access_token = var.line_channel_access_token
}
```

Or via environment variable:

```sh
export LIFF_CHANNEL_ACCESS_TOKEN="your-access-token"
```

### Authentication Priority

1. `channel_access_token` (direct specification)
2. `channel_id` + `channel_secret` (Stateless Token auto-issuance)
3. Environment variable fallback for each attribute

### Environment Variables

| Variable | Description |
|----------|-------------|
| `LIFF_CHANNEL_ACCESS_TOKEN` | Channel Access Token (direct) |
| `LIFF_CHANNEL_ID` | LINE Login Channel ID |
| `LIFF_CHANNEL_SECRET` | LINE Login Channel Secret |

## Usage

### Managing a LIFF App

```hcl
resource "liff_app" "my_app" {
  description            = "My LIFF App"
  permanent_link_pattern = "concat"
  bot_prompt             = "normal"
  scope                  = ["openid", "profile"]

  view {
    type        = "full"
    url         = "https://example.com"
    module_mode = false
  }

  features {
    qr_code = true
  }
}
```

### Importing an Existing LIFF App

```sh
terraform import liff_app.my_app "1234567890-AbCdEfGh"
```

### Listing All LIFF Apps

```hcl
data "liff_apps" "all" {}

output "app_ids" {
  value = data.liff_apps.all.apps[*].liff_id
}
```

### Looking Up a Specific LIFF App

```hcl
data "liff_app" "existing" {
  liff_id = "1234567890-AbCdEfGh"
}

output "app_url" {
  value = data.liff_app.existing.view.url
}
```

## Resources

### `liff_app`

Manages the full lifecycle (create, read, update, delete, import) of a LIFF application.

#### Argument Reference

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `description` | String | Optional | - | Name of the LIFF app. Cannot include "LINE" or similar strings. |
| `permanent_link_pattern` | String | Optional | - | How additional information in LIFF URLs is handled. Only `"concat"` is supported. |
| `bot_prompt` | String | Optional | `"none"` | Bot link feature setting. One of `"normal"`, `"aggressive"`, `"none"`. |
| `scope` | List(String) | Optional | - | Array of scopes: `"openid"`, `"email"`, `"profile"`, `"chat_message.write"`. |
| `view` | Block | **Required** | - | LIFF app view settings. See [view](#view). |
| `features` | Block | Optional | - | LIFF app feature settings. See [features](#features). |

##### `view`

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | String | **Required** | - | Size of the LIFF app view. One of `"compact"`, `"tall"`, `"full"`. |
| `url` | String | **Required** | - | Endpoint URL. Must be HTTPS. |
| `module_mode` | Bool | Optional | `false` | Whether to use the LIFF app in modular mode. |

##### `features`

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `ble` | Bool | Computed | - | Whether the LIFF app supports BLE (read-only). |
| `qr_code` | Bool | Optional | `false` | Whether to use the 2D code reader. |

#### Attribute Reference

| Attribute | Description |
|-----------|-------------|
| `liff_id` | The LIFF app ID (e.g., `1234567890-AbCdEfGh`). |

## Data Sources

### `liff_apps`

Fetches all LIFF apps in the channel.

#### Attribute Reference

| Attribute | Description |
|-----------|-------------|
| `apps` | List of LIFF apps. Each element has the same attributes as the `liff_app` resource. |

### `liff_app`

Fetches a specific LIFF app by its ID.

#### Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `liff_id` | String | **Required** | The LIFF app ID to look up. |

#### Attribute Reference

All attributes from the `liff_app` resource are available as computed attributes.

## Limitations

- **No individual GET API**: The LIFF Server API only provides a "list all" endpoint. The provider fetches all apps and filters by `liff_id`. With a maximum of 30 LIFF apps per channel, this has negligible performance impact.
- **`features.ble` is read-only**: BLE (Bluetooth Low Energy) support cannot be set via the API. The `ble` attribute is computed only.
- **Maximum 30 LIFF apps per channel**: This is a LINE platform limitation.

## Developing

### Building

```sh
go install
```

### Running Tests

Unit tests:

```sh
make test
```

Acceptance tests (creates real resources):

```sh
export LIFF_CHANNEL_ID="your-channel-id"
export LIFF_CHANNEL_SECRET="your-channel-secret"
make testacc
```

### Generating Documentation

```sh
make generate
```

## License

[Mozilla Public License v2.0](./LICENSE)
