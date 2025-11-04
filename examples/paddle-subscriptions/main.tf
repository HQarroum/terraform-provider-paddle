terraform {
  required_version = ">= 1.0"

  required_providers {
    paddle = {
      source  = "HQarroum/paddle"
      version = "~> 0.1.0"
    }
  }
}

provider "paddle" {
  api_key     = var.paddle_api_key
  environment = var.paddle_environment
}

# Create a SaaS product
resource "paddle_product" "saas_platform" {
  name         = "CloudFlow SaaS Platform"
  description  = "Professional project management and collaboration platform"
  tax_category = "standard"
  image_url    = "https://example.com/cloudflow-logo.png"

  custom_data = {
    product_tier = "standard"
    category     = "productivity"
  }
}

# Monthly subscription - $29/month
resource "paddle_price" "monthly" {
  product_id  = paddle_product.saas_platform.id
  name        = "Monthly Plan"
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

# Annual subscription - $290/year (saves $58)
resource "paddle_price" "annual" {
  product_id  = paddle_product.saas_platform.id
  name        = "Annual Plan"
  description = "Billed annually - Save $58 per year"

  unit_price = {
    amount        = "29000"
    currency_code = "USD"
  }

  billing_cycle = {
    interval  = "year"
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

# Early adopter discount - 20% off
resource "paddle_discount" "launch_discount" {
  description = "Early Adopter Discount - 20% off"
  type        = "percentage"
  mode        = "standard"
  amount      = "20"
  code        = "LAUNCH_2025"

  restrict_to = [
    paddle_product.saas_platform.id
  ]

  custom_data = {
    campaign    = "launch-2025"
    valid_until = "2025-12-31"
  }
}

# Outputs
output "product_id" {
  description = "The ID of the SaaS product"
  value       = paddle_product.saas_platform.id
}

output "discount_code" {
  description = "The discount code for early adopters"
  value       = paddle_discount.launch_discount.code
}

output "checkout_url_monthly" {
  description = "Checkout URL for monthly plan"
  value       = "https://sandbox-checkout.paddle.com/checkout/custom/${paddle_price.monthly.id}"
}

output "checkout_url_annual" {
  description = "Checkout URL for annual plan"
  value       = "https://sandbox-checkout.paddle.com/checkout/custom/${paddle_price.annual.id}"
}
