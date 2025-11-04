---
page_title: "paddle_discount Data Source - terraform-provider-paddle"
subcategory: ""
description: |-
  Retrieves information about a Paddle discount.
---

# paddle_discount

Retrieves information about an existing Paddle discount.

## Example Usage

```terraform
data "paddle_discount" "example" {
  id = "dsc_01h1vjf8j7v3x9r1f5b4p6n8k2"
}

output "discount_code" {
  value = data.paddle_discount.example.code
}
```

## Schema

### Required

- `id` (String) Paddle discount ID (format: `dsc_...`).

### Read-Only

- `description` (String) Description of the discount.
- `type` (String) Type of discount.
- `amount` (String) Discount amount.
- `code` (String) Discount code.
- `mode` (String) Discount mode.
- `currency_code` (String) Currency code (for flat discounts).
- `recur` (Boolean) Whether discount recurs.
- `maximum_recurring_intervals` (Number) Maximum recurring intervals.
- `restrict_to` (List of String) Restricted product IDs.
- `custom_data` (Map of String) Custom metadata.
- `status` (String) Status of the discount.
- `enabled_for_checkout` (Boolean) Whether enabled for checkout.
- `times_used` (Number) Usage count.
- `created_at` (String) RFC 3339 timestamp when created.
- `updated_at` (String) RFC 3339 timestamp when last updated.
