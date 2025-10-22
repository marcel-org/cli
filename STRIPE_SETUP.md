# Stripe Premium Subscription Setup

This guide will help you set up Stripe for the Marcel Premium subscription feature.

## Prerequisites

- A Stripe account ([Sign up here](https://dashboard.stripe.com/register))
- Access to your Stripe Dashboard

## Setup Steps

### 1. Get Your Stripe API Keys

1. Log in to your [Stripe Dashboard](https://dashboard.stripe.com)
2. Go to **Developers** > **API keys**
3. Copy your **Publishable key** (starts with `pk_test_` or `pk_live_`)
4. Copy your **Secret key** (starts with `sk_test_` or `sk_live_`)

### 2. Create a Subscription Product

1. In your Stripe Dashboard, go to **Products** > **Add product**
2. Set the following:
   - **Name**: Marcel Premium
   - **Description**: Premium subscription for Marcel productivity app
   - **Pricing model**: Standard pricing
   - **Price**: $4.99 USD
   - **Billing period**: Monthly
3. Click **Save product**
4. Copy the **Price ID** (starts with `price_`) - you'll need this

### 3. Configure Environment Variables

#### Backend (.env)

Add these variables to your `backend/.env` file:

```bash
# Stripe Configuration
STRIPE_SECRET_KEY=sk_test_your_secret_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_publishable_key_here
```

### 4. Update the Price ID in Frontend

Edit `/frontend/src/routes/(app)/premium/+page.svelte` and replace the price ID on line 37:

```typescript
priceId: 'YOUR_PRICE_ID_HERE', // Replace with your actual Price ID from step 2
```

### 5. Set Up Stripe Webhooks

Webhooks allow Stripe to notify your backend when subscription events occur (payments, cancellations, etc.).

#### For Development (using Stripe CLI):

1. Install the Stripe CLI: https://stripe.com/docs/stripe-cli
2. Log in: `stripe login`
3. Forward webhooks to your local backend:
   ```bash
   stripe listen --forward-to localhost:8080/stripe/webhook
   ```
4. Copy the webhook signing secret (starts with `whsec_`) and add it to your `.env` file as `STRIPE_WEBHOOK_SECRET`

#### For Production:

1. In your Stripe Dashboard, go to **Developers** > **Webhooks**
2. Click **Add endpoint**
3. Set the **Endpoint URL** to: `https://your-backend-url.com/stripe/webhook`
4. Select the following events to listen to:
   - `checkout.session.completed`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
   - `invoice.payment_succeeded`
   - `invoice.payment_failed`
5. Click **Add endpoint**
6. Copy the **Signing secret** and add it to your production `.env` file as `STRIPE_WEBHOOK_SECRET`

### 6. Test the Integration

#### Test Cards (for test mode):

- **Success**: `4242 4242 4242 4242`
- **Decline**: `4000 0000 0000 0002`
- **Requires authentication**: `4000 0025 0000 3155`

Use any future expiration date and any 3-digit CVC.

#### Test Flow:

1. Start your backend: `cd backend && bun run src/index.ts`
2. Start your frontend: `cd frontend && bun run dev`
3. Log in to Marcel
4. Go to Settings > Premium or navigate to `/premium`
5. Click "Subscribe Now"
6. Use a test card to complete the payment
7. You should be redirected back with a success message
8. Your profile should now show the premium badge

### 7. Verify Webhook Events

1. In the Stripe Dashboard, go to **Developers** > **Webhooks**
2. Click on your webhook endpoint
3. Check the **Recent events** section to see if events are being received

### 8. Go Live

When you're ready to accept real payments:

1. Replace test keys with live keys in your `.env` files
2. Create a new product in live mode with the same pricing
3. Update the frontend with the live price ID
4. Set up production webhooks (step 5)
5. Test with a real card (recommended: use a small amount first)

## Security Notes

- **Never commit** your `.env` files to git
- Store webhook secrets securely
- Use HTTPS in production
- Validate webhook signatures (already implemented in the code)
- Test thoroughly before going live

## Features Implemented

### Backend (`/backend/src/routes/stripe.ts`)

- ✅ Create Stripe Checkout Session
- ✅ Create Stripe Customer Portal Session (for managing subscriptions)
- ✅ Webhook handler for subscription events
- ✅ Database sync for subscription status
- ✅ Automatic premium status management

### Frontend

- ✅ Premium subscription page (`/premium`)
- ✅ Premium badge on profile pages
- ✅ Premium section in settings
- ✅ Stripe checkout integration
- ✅ Subscription status display

### Database

- ✅ `Subscription` table with Stripe data
- ✅ `PaymentHistory` table for payment tracking
- ✅ `isPremium` field on User model

## Troubleshooting

### Webhook not receiving events

- Check that the webhook URL is correct and accessible
- Verify the webhook signing secret matches
- Check server logs for errors
- Use Stripe CLI for local testing

### Payment succeeds but user not marked as premium

- Check webhook events in Stripe Dashboard
- Check backend logs for webhook processing errors
- Verify database migrations ran successfully
- Check that `userId` is being passed correctly in checkout metadata

### Checkout session creation fails

- Verify Stripe secret key is correct
- Check that price ID is valid
- Ensure user has a valid email address
- Check backend logs for detailed error messages

## Support

For Stripe-specific issues, refer to:
- [Stripe Documentation](https://stripe.com/docs)
- [Stripe API Reference](https://stripe.com/docs/api)
- [Stripe Support](https://support.stripe.com/)

For Marcel-specific issues, check the main README or contact the development team.
