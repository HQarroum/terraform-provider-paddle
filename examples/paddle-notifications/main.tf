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

# Create a product
resource "paddle_product" "api_service" {
  name         = "API Service Platform"
  description  = "RESTful API service with webhook integration"
  tax_category = "standard"
  image_url    = "https://example.com/api-service-logo.png"

  custom_data = {
    service_type = "api"
    features     = "webhooks,rate-limiting,analytics"
  }
}

# Create a monthly price
resource "paddle_price" "monthly" {
  product_id  = paddle_product.api_service.id
  name        = "Monthly API Access"
  description = "Billed monthly"

  unit_price = {
    amount        = "4900"
    currency_code = "USD"
  }

  billing_cycle = {
    interval  = "month"
    frequency = 1
  }

  trial_period = {
    interval  = "day"
    frequency = 7
  }
}

# Configure webhook notification endpoint
resource "paddle_notification_setting" "webhook" {
  description              = "Production webhook for transactions and subscriptions"
  destination              = var.webhook_url
  active                   = true
  include_sensitive_fields = false

  # Subscribe to transaction events
  subscribed_events = [
    "transaction.completed",
    "transaction.updated",
    "subscription.created",
    "subscription.updated",
    "subscription.canceled",
    "customer.created"
  ]
}

# Outputs

output "product_id" {
  description = "The ID of the product"
  value       = paddle_product.api_service.id
}

output "price_id" {
  description = "The ID of the monthly price"
  value       = paddle_price.monthly.id
}

output "webhook_url" {
  description = "The configured webhook URL"
  value       = paddle_notification_setting.webhook.destination
}

output "checkout_url" {
  description = "Checkout URL for testing (make a purchase to trigger webhook events)"
  value       = "https://sandbox-checkout.paddle.com/checkout/custom/${paddle_price.monthly.id}"
}
