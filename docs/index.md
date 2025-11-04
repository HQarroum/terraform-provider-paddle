---
page_title: "Paddle Provider"
subcategory: ""
description: |-
  Terraform provider for managing Paddle Billing resources.
---

# Paddle Provider

The Paddle provider allows you to manage [Paddle Billing](https://www.paddle.com/) resources such as products, prices, customers, discounts, and webhook notifications using Terraform.

## Example Usage

```terraform
terraform {
  required_providers {
    paddle = {
      source  = "HQarroum/paddle"
      version = "~> 0.1.0"
    }
  }
}

provider "paddle" {
  api_key     = var.paddle_api_key
  environment = "sandbox" # or "production"
}
```

## Authentication

The provider requires a Paddle API key. You can obtain one from your [Paddle Dashboard](https://vendors.paddle.com/authentication).

-> **Note:** It's recommended to use environment variables or Terraform variables to store your API key securely, rather than hardcoding it in your configuration.

```terraform
provider "paddle" {
  api_key     = var.paddle_api_key
  environment = var.paddle_environment
}
```

## Schema

### Required

- `api_key` (String, Sensitive) Paddle API key for authentication. Can also be set via the `PADDLE_API_KEY` environment variable.

### Optional

- `environment` (String) Paddle environment to use. Must be either `sandbox` or `production`. Defaults to `sandbox`. Can also be set via the `PADDLE_ENVIRONMENT` environment variable.
