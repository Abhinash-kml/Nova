CREATE TABLE IF NOT EXISTS currencies(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    code VARCHAR UNIQUE NOT NULL,
    kind VARCHAR NOT NULL,
    is_purchasable BOOLEAN NOT NULL DEFAULT FALSE,
    is_tradable BOOLEAN NOT NULL DEFAULT FALSE,
    expires_after_days INT NOT NULL,
    max_balance_minor BIGINT NOT NULL,
    sort_order SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets(
    user_id UUID NOT NULL,
    currency_id INT NOT NULL,
    account_kind VARCHAR NOT NULL,
    balance_minor BIGINT NOT NULL,

    PRIMARY KEY(user_id, currency_id),

    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,

    FOREIGN KEY(currency_id) REFERENCES currencies(id)
    ON DELETE CASCADE,

    CONSTRAINT unique_user_currency UNIQUE(user_id, currency_id)
);

CREATE TYPE IF NOT EXISTS WALLET_TRANSACTION_SOURCE AS ENUM (
    'match_reward',
    'store_purchase',
    'admin_grant',
    'refund'
);

CREATE TABLE IF NOT EXISTS wallet_ledger (
    user_id UUID,
    currency_id INT,
    transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('credit', 'debit')),
    transaction_source WALLET_TRANSACTION_SOURCE NOT NULL,
    amount_minor BIGINT NOT NULL,
    reference_id UUID NOT NULL UNIQUE, -- Idempotency key from game engine (e.g., match_id, purchase_id)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(user_id, currency_id, created_at),

    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE SET NULL,

    FOREIGN KEY(currency_id) REFERENCES currencies(id)
    ON DELETE SET NULL
);

-- Crucial index for summing up logs or pulling history quickly
CREATE INDEX idx_wallet_tx_wallet_id ON wallet_ledger(wallet_id);