---
page_title: "paddle_product Resource - terraform-provider-paddle"
subcategory: ""
description: |-
  Manages a Paddle product.
---

# paddle_product

Manages a Paddle product. Products are the items or services that you sell.

## Example Usage

```terraform
resource "paddle_product" "saas_platform" {
  name         = "CloudFlow SaaS Platform"
  description  = "Professional project management and collaboration platform"
  tax_category = "saas"
  image_url    = "https://example.com/product-image.png"

  custom_data = {
    product_tier = "standard"
    category     = "productivity"
  }
}
```

## Schema

### Required

- `name` (String) Name of the product.
- `tax_category` (String) Tax category for the product. Must be one of: `standard`, `digital-goods`, `ebooks`, `implementation-services`, `professional-services`, `saas`, `software-programming-services`, `training-services`.

### Optional

- `description` (String) Description of the product.
- `image_url` (String) URL of the product image.
- `custom_data` (Map of String) Custom metadata as key-value pairs.

### Read-Only

- `id` (String) Paddle product ID (format: `pro_...`).
- `status` (String) Status of the product (`active` or `archived`).
- `created_at` (String) RFC 3339 timestamp when the product was created.
- `updated_at` (String) RFC 3339 timestamp when the product was last updated.

## Import

Products can be imported using the Paddle product ID:

```shell
terraform import paddle_product.example pro_01h1vjes1y163xfj1rh1tkfb65
```
