---
page_title: "paddle_price Data Source - terraform-provider-paddle"
subcategory: ""
description: |-
  Retrieves information about a Paddle price.
---

# paddle_price (Data Source)

Retrieves information about an existing Paddle price.

## Example Usage

```terraform
data "paddle_price" "example" {
  id = "pri_01h1vjeh4bt4y75vf5c69azkm9"
}

output "price_amount" {
  value = data.paddle_price.example.unit_price.amount
}
```

## Schema

### Required

- `id` (String) Paddle price ID (format: `pri_...`).

### Read-Only

- `product_id` (String) Associated product ID.
- `name` (String) Name of the price.
- `description` (String) Description of the price.
- `tax_mode` (String) Tax calculation mode.
- `unit_price` (Block) Price per unit.
- `billing_cycle` (Block) Billing cycle information.
- `trial_period` (Block) Trial period information.
- `quantity` (Block) Quantity limits.
- `status` (String) Status of the price.
- `created_at` (String) RFC 3339 timestamp when created.
- `updated_at` (String) RFC 3339 timestamp when last updated.
