# Paddle Subscriptions Example

This example demonstrates how to create a SaaS product with recurring subscription pricing (monthly and annual plans) using the Paddle Terraform provider.

## What This Creates

- ✅ A SaaS product
- ✅ Monthly subscription plan ($29/month) with 14-day trial
- ✅ Annual subscription plan ($290/year - save $58) with 14-day trial
- ✅ A discount code (20% off for early adopters)

## Prerequisites

- Terraform >= 1.0
- Paddle API key (sandbox or production)

## Setup

1. **Set your Paddle API key:**

```bash
export PADDLE_API_KEY="your-paddle-api-key"
export PADDLE_ENVIRONMENT="sandbox"
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

## Outputs

After deployment, you'll see:
- Product ID
- Monthly price ID  
- Annual price ID
- Discount code

## Testing

Use the Paddle Checkout to test your subscription:

```
https://sandbox-checkout.paddle.com/checkout/custom/[PRICE_ID]
```

## Cleanup

To destroy all resources:

```bash
terraform destroy
```
