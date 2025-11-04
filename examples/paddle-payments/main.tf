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

# Ebook product
resource "paddle_product" "ebook" {
  name         = "The Complete Guide to Terraform"
  description  = "Master Infrastructure as Code with this comprehensive ebook - UPDATED"
  tax_category = "standard"
  image_url    = "https://example.com/terraform-ebook-cover.png"

  custom_data = {
    format       = "pdf-epub-mobi"
    pages        = "350"
    publish_year = "2025"
  }
}

# Ebook price - $39 one-time
resource "paddle_price" "ebook_price" {
  product_id  = paddle_product.ebook.id
  name        = "Ebook Purchase"
  description = "One-time purchase - instant download"

  unit_price = {
    amount        = "3900"
    currency_code = "USD"
  }
}

# Video course product
resource "paddle_product" "video_course" {
  name         = "Advanced Terraform Patterns"
  description  = "Video course with real-world examples and best practices"
  tax_category = "standard"
  image_url    = "https://example.com/terraform-course-cover.png"

  custom_data = {
    duration_hours = "12"
    video_count    = "45"
    level          = "advanced"
  }
}

# Video course price - $149 one-time
resource "paddle_price" "video_course_price" {
  product_id  = paddle_product.video_course.id
  name        = "Course Purchase"
  description = "One-time purchase - lifetime access"

  unit_price = {
    amount        = "14900"
    currency_code = "USD"
  }
}

# Bundle discount - 15% off
resource "paddle_discount" "bundle_discount" {
  description = "Bundle Discount - 15% off when buying both"
  type        = "percentage"
  mode        = "standard"
  amount      = "15"
  code        = "BUNDLE2025V6"

  restrict_to = [
    paddle_product.ebook.id,
    paddle_product.video_course.id
  ]

  custom_data = {
    campaign = "bundle-offer"
    offer    = "ebook-plus-course"
  }
}

# Outputs

output "ebook_product_id" {
  description = "The ID of the ebook product"
  value       = paddle_product.ebook.id
}

output "video_course_product_id" {
  description = "The ID of the video course product"
  value       = paddle_product.video_course.id
}

output "bundle_discount_code" {
  description = "The bundle discount code"
  value       = paddle_discount.bundle_discount.code
}

output "checkout_url_ebook" {
  description = "Checkout URL for ebook"
  value       = "https://sandbox-checkout.paddle.com/checkout/custom/${paddle_price.ebook_price.id}"
}

output "checkout_url_course" {
  description = "Checkout URL for video course"
  value       = "https://sandbox-checkout.paddle.com/checkout/custom/${paddle_price.video_course_price.id}"
}
