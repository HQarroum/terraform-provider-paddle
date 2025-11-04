---
page_title: "paddle_notification_setting Resource - terraform-provider-paddle"
subcategory: ""
description: |-
  Manages a Paddle notification setting (webhook).
---

# paddle_notification_setting

Manages a Paddle notification setting. Notification settings allow you to subscribe to events and receive webhooks.

## Example Usage

```terraform
resource "paddle_notification_setting" "webhook" {
  description              = "Production webhook endpoint"
  destination              = "https://api.example.com/paddle/webhook"
  active                   = true
  include_sensitive_fields = false

  subscribed_events = [
    "transaction.completed",
    "transaction.updated",
    "subscription.created",
    "subscription.updated",
    "subscription.canceled",
    "customer.created"
  ]
}
```

## Schema

### Required

- `description` (String) Short description for this notification destination.
- `destination` (String) Webhook endpoint URL or email address. URLs must be HTTPS.
- `subscribed_events` (List of String) List of event types to subscribe to.

### Optional

- `type` (String) Type of notification destination. One of: `url` (webhook), `email`. Defaults to `url`. Cannot be changed after creation.
- `active` (Boolean) Whether the notification destination is active. Defaults to `true`.
- `include_sensitive_fields` (Boolean) Whether to include sensitive fields in webhook payloads.
- `traffic_source` (String) Filter events by source. One of: `platform`, `api`. Omit to receive all events.

### Read-Only

- `id` (String) Paddle notification setting ID (format: `ntfset_...`).
- `endpoint_secret_key` (String, Sensitive) Webhook secret key for signature verification.
- `api_version` (Number) API version for event payloads.

## Import

Notification settings can be imported using the Paddle notification setting ID:

```shell
terraform import paddle_notification_setting.example ntfset_01h1vjfbk9m2q4r7x3w5t8n6p0
```
