variable "paddle_api_key" {
  description = "Paddle API key"
  type        = string
  sensitive   = true
}

variable "paddle_environment" {
  description = "Paddle environment (sandbox or production)"
  type        = string
  default     = "sandbox"
}

variable "webhook_url" {
  description = "Your webhook endpoint URL that will receive Paddle events"
  type        = string

  validation {
    condition     = can(regex("^https://", var.webhook_url))
    error_message = "Webhook URL must use HTTPS protocol."
  }
}
