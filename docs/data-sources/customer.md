---
page_title: "paddle_customer Data Source - terraform-provider-paddle"
subcategory: ""
description: |-
  Retrieves information about a Paddle customer.
---

# paddle_customer (Data Source)

Retrieves information about an existing Paddle customer.

## Example Usage

```terraform
data "paddle_customer" "example" {
  id = "ctm_01h1vjf0j84dfq3fh0trr7nqxb"
}

output "customer_email" {
  value = data.paddle_customer.example.email
}
```

## Schema

### Required

- `id` (String) Paddle customer ID (format: `ctm_...`).

### Read-Only

- `email` (String) Email address of the customer.
- `name` (String) Full name of the customer.
- `locale` (String) Locale tag.
- `custom_data` (Map of String) Custom metadata.
- `status` (String) Status of the customer.
- `created_at` (String) RFC 3339 timestamp when created.
- `updated_at` (String) RFC 3339 timestamp when last updated.
