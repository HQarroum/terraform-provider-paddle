# Paddle Notifications Example

This example demonstrates how to configure webhook notifications for Paddle events using the Paddle Terraform provider.

## What This Creates

- ✅ A SaaS product
- ✅ A monthly subscription price
- ✅ A webhook notification endpoint configured to receive transaction and subscription events

## Prerequisites

- Terraform >= 1.0
- Paddle API key (sandbox or production)
- A webhook endpoint URL (your server that will receive Paddle events)

## Setup

1. **Set your Paddle API key and webhook URL:**

```bash
export PADDLE_API_KEY="your-paddle-api-key"
export PADDLE_ENVIRONMENT="sandbox"
export TF_VAR_webhook_url="https://your-api.com/paddle/webhook"
```

2. **Initialize Terraform:**

```bash
terraform init
```

3. **Review the plan:**

```bash
terraform plan
```

4. **Deploy:**

```bash
terraform apply
```

## Webhook Events

The webhook is configured to receive the following events:

- `transaction.completed` - When a transaction completes successfully
- `transaction.updated` - When a transaction is updated
- `subscription.created` - When a new subscription is created
- `subscription.updated` - When a subscription is updated
- `subscription.canceled` - When a subscription is canceled
- `customer.created` - When a new customer is created

## Testing Webhooks

1. Deploy this configuration
2. Make a test purchase using the checkout URL from outputs
3. Check your webhook endpoint for incoming events

**Tip:** Use tools like [webhook.site](https://webhook.site) or [ngrok](https://ngrok.com) for testing during development.

## Outputs

After deployment, you'll see:
- Product ID
- Price ID
- Webhook notification setting ID
- Checkout URL for testing

## Cleanup

To destroy all resources:

```bash
terraform destroy
```
