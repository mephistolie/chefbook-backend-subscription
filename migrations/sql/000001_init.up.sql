CREATE TYPE plan AS ENUM ('premium', 'maximum');

CREATE TYPE source AS ENUM ('google', 'yookassa', 'firebase', 'promocode');

CREATE TABLE subscriptions
(
    user_id              uuid PRIMARY KEY         NOT NULL,
    plan                 plan                     NOT NULL,
    source               source                   NOT NULL,
    start_timestamp      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now():: timestamp,
    expiration_timestamp TIMESTAMP WITH TIME ZONE          DEFAULT NULL,
    auto_renew           BOOLEAN                           DEFAULT false,

    last_check           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now():: timestamp,
    UNIQUE (user_id, plan, source)
);

CREATE INDEX subscriptions_renew_key ON subscriptions (expiration_timestamp, last_check);

CREATE TABLE google
(
    user_id uuid REFERENCES subscriptions (user_id) ON DELETE CASCADE NOT NULL,
    purchaseToken                                                     NOT NULL UNIQUE,
);

CREATE TABLE inbox
(
    message_id uuid PRIMARY KEY         NOT NULL UNIQUE,
    timestamp  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now():: timestamp
);