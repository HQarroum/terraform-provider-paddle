---
page_title: "paddle_discount Resource - terraform-provider-paddle"
subcategory: ""
description: |-
  Manages a Paddle discount code.
---

# paddle_discount (Resource)

Manages a Paddle discount. Discounts let you offer money off transactions.

## Example Usage

```terraform
resource "paddle_discount" "launch_offer" {
  description = "Early Adopter Discount - 20% off"
  type        = "percentage"
  mode        = "standard"
  amount      = "20"
  code        = "LAUNCH2025"

  restrict_to = [
    paddle_product.saas_platform.id
  ]

  custom_data = {
    campaign    = "launch-2025"
    valid_until = "2025-12-31"
  }
}
```

## Schema

### Required

- `description` (String) Description of the discount.
- `type` (String) Type of discount. One of: `percentage`, `flat`, `flat_per_seat`.
- `amount` (String) Amount to discount. For `percentage`: 0.01-100. For `flat`/`flat_per_seat`: amount in lowest denomination.

### Optional

- `code` (String) Unique code that customers use to redeem. Not case-sensitive. Omit for seller-applied discounts.
- `mode` (String) Discount mode. One of: `standard`, `custom`. Defaults to `standard`. Cannot be changed after creation.
- `currency_code` (String) Three-letter ISO 4217 currency code. Required for `flat` and `flat_per_seat` types.
- `recur` (Boolean) Whether discount applies for multiple billing periods. Defaults to `false`.
- `maximum_recurring_intervals` (Number) Number of billing periods discount recurs for. Requires `recur = true`.
- `restrict_to` (List of String) List of product IDs this discount is restricted to.
- `custom_data` (Map of String) Custom metadata as key-value pairs.

### Read-Only

- `id` (String) Paddle discount ID (format: `dsc_...`).
- `status` (String) Status of the discount (`active` or `archived`).
- `enabled_for_checkout` (Boolean) Whether customers can redeem at checkout.
- `times_used` (Number) Number of times discount has been used.
- `created_at` (String) RFC 3339 timestamp when the discount was created.
- `updated_at` (String) RFC 3339 timestamp when the discount was last updated.

## Import

Discounts can be imported using the Paddle discount ID:

```shell
terraform import paddle_discount.example dsc_01h1vjf8j7v3x9r1f5b4p6n8k2
```
