---
page_title: "paddle_product Data Source - terraform-provider-paddle"
subcategory: ""
description: |-
  Retrieves information about a Paddle product.
---

# paddle_product (Data Source)

Retrieves information about an existing Paddle product.

## Example Usage

```terraform
data "paddle_product" "example" {
  id = "pro_01h1vjes1y163xfj1rh1tkfb65"
}

output "product_name" {
  value = data.paddle_product.example.name
}
```

## Schema

### Required

- `id` (String) Paddle product ID (format: `pro_...`).

### Read-Only

- `name` (String) Name of the product.
- `description` (String) Description of the product.
- `tax_category` (String) Tax category for the product.
- `image_url` (String) URL of the product image.
- `custom_data` (Map of String) Custom metadata.
- `status` (String) Status of the product.
- `created_at` (String) RFC 3339 timestamp when created.
- `updated_at` (String) RFC 3339 timestamp when last updated.
