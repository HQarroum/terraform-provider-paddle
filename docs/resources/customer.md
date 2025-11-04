---
page_title: "paddle_customer Resource - terraform-provider-paddle"
subcategory: ""
description: |-
  Manages a Paddle customer.
---

# paddle_customer (Resource)

Manages a Paddle customer. Customers are people or businesses that buy from you.

## Example Usage

```terraform
resource "paddle_customer" "example" {
  email  = "customer@example.com"
  name   = "John Doe"
  locale = "en"

  custom_data = {
    company     = "Acme Corp"
    employee_id = "EMP-001"
  }
}
```

## Schema

### Required

- `email` (String) Email address of the customer.

### Optional

- `name` (String) Full name of the customer.
- `locale` (String) Valid IETF BCP 47 locale tag (e.g., `en`, `en-US`). Defaults to `en`.
- `custom_data` (Map of String) Custom metadata as key-value pairs.

### Read-Only

- `id` (String) Paddle customer ID (format: `ctm_...`).
- `status` (String) Status of the customer (`active` or `archived`).
- `created_at` (String) RFC 3339 timestamp when the customer was created.
- `updated_at` (String) RFC 3339 timestamp when the customer was last updated.

## Import

Customers can be imported using the Paddle customer ID:

```shell
terraform import paddle_customer.example ctm_01h1vjf0j84dfq3fh0trr7nqxb
```
