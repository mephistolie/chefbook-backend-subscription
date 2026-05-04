# ChefBook Backend Subscription Service

The subscription service owns subscription state, subscription source tracking, Google Play subscription confirmation, renewal/cancellation events, and subscription-related MQ processing.

## Responsibilities

- Return profile subscriptions and current subscription.
- Confirm Google subscription purchases.
- Acknowledge Google subscription purchases after successful claim.
- Process Google real-time developer notifications.
- Track subscription source, expiration, and auto-renew state.
- Import legacy premium state from Firebase profile import messages.
- Delete profile subscription state when profile deletion messages arrive.

## Main RPC Families

- `GetProfileSubscriptions`
- `GetProfileCurrentSubscription`
- `ConfirmGoogleSubscription`

## Dependencies

- Calls `auth` for user email when sending subscription-related mail.
- Uses Google Android Publisher API for purchase validation and acknowledgement.
- Uses Google Pub/Sub for Google real-time developer notifications.
- Consumes profile lifecycle messages from MQ when configured.
- Owns its PostgreSQL schema and migrations.

## Database Ownership

Owns:

- `subscriptions` - user plan, source, start, expiration, and auto-renew state.
- `google` - Google purchase token binding.
- `inbox` - idempotent MQ message processing.

```mermaid
erDiagram
    SUBSCRIPTIONS {
        uuid user_id PK
        plan plan
        source source
        timestamptz start_timestamp
        timestamptz expiration_timestamp
        boolean auto_renew
    }

    GOOGLE {
        uuid user_id FK
        varchar purchaseToken UK
    }

    SUBSCRIPTION_INBOX {
        uuid message_id PK
        timestamptz timestamp
    }

    SUBSCRIPTIONS ||--o| GOOGLE : has_google_purchase
```

Important constraints:

- Subscriptions are unique by user, plan, and source.
- `google.purchaseToken` is globally unique.
- `user_id` is a logical cross-service reference.
