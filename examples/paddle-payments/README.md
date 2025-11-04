# Paddle One-Time Payments Example

This example demonstrates how to create products with one-time payment pricing using the Paddle Terraform provider.

## What This Creates

- ✅ Ebook product ($39 one-time payment)
- ✅ Video course product ($149 one-time payment)
- ✅ A discount code (15% off when buying both)

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
- Ebook product ID and price ID
- Video course product ID and price ID
- Bundle discount code

## Testing

Use the Paddle Checkout to test purchases:

```
https://sandbox-checkout.paddle.com/checkout/custom/[PRICE_ID]
```

Apply discount code: `BUNDLE15`

## Cleanup

To destroy all resources:

```bash
terraform destroy
```
