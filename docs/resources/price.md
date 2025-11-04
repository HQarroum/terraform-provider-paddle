---
page_title: "paddle_price Resource - terraform-provider-paddle"
subcategory: ""
description: |-
  Manages a Paddle price for a product.
---

# paddle_price (Resource)

Manages a Paddle price. Prices determine how much and how often you charge for a product.

## Example Usage

### One-Time Price

```terraform
resource "paddle_price" "ebook" {
  product_id  = paddle_product.my_product.id
  name        = "Ebook Purchase"
  description = "One-time purchase - instant download"

  unit_price = {
    amount        = "3900"
    currency_code = "USD"
  }
}
```

### Recurring Subscription Price

```terraform
resource "paddle_price" "monthly" {
  product_id  = paddle_product.saas_platform.id
  name        = "Monthly Subscription"
  description = "Billed monthly"

  unit_price = {
    amount        = "2900"
    currency_code = "USD"
  }

  billing_cycle = {
    interval  = "month"
    frequency = 1
  }

  trial_period = {
    interval  = "day"
    frequency = 14
  }

  quantity = {
    minimum = 1
    maximum = 100
  }
}
```

## Schema

### Required

- `product_id` (String) Paddle product ID that this price is for. Cannot be changed after creation.
- `description` (String) Internal description for the price, not shown to customers.

### Optional

- `name` (String) Name of this price, shown to customers.
- `tax_mode` (String) How tax is calculated. One of: `account_setting`, `internal`, `external`. Defaults to `account_setting`.
- `unit_price` (Block) Price per unit.
  - `amount` (String, Required) Amount in the lowest denomination (e.g., cents).
  - `currency_code` (String, Required) Three-letter ISO 4217 currency code.
- `billing_cycle` (Block) Billing cycle for recurring prices.
  - `interval` (String, Required) Billing interval: `day`, `week`, `month`, or `year`.
  - `frequency` (Number, Required) Number of intervals between billings.
- `trial_period` (Block) Trial period for subscriptions.
  - `interval` (String, Required) Trial interval: `day`, `week`, `month`, or `year`.
  - `frequency` (Number, Required) Number of intervals for the trial.
- `quantity` (Block) Quantity limits.
  - `minimum` (Number, Required) Minimum quantity.
  - `maximum` (Number, Required) Maximum quantity.

### Read-Only

- `id` (String) Paddle price ID (format: `pri_...`).
- `status` (String) Status of the price (`active` or `archived`).
- `created_at` (String) RFC 3339 timestamp when the price was created.
- `updated_at` (String) RFC 3339 timestamp when the price was last updated.

## Import

Prices can be imported using the Paddle price ID:

```shell
terraform import paddle_price.example pri_01h1vjeh4bt4y75vf5c69azkm9
```
