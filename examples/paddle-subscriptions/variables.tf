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
