CREATE TABLE IF NOT EXISTS currencies(
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    code VARCHAR UNIQUE NOT NULL,
    kind VARCHAR NOT NULL,
    is_purchasable BOOLEAN NOT NULL DEFAULT FALSE,
    is_tradable BOOLEAN NOT NULL DEFAULT FALSE,
    expires_after_days INT NOT NULL,
    max_balance_minor INT NOT NULL,
    sort_order SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    currency_id INT NOT NULL,
    account_kind VARCHAR NOT NULL,
    balance_minor BIGINT NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(currency_id) REFERENCES currencies(id),
    CONSTRAINT unique_user_currency UNIQUE(user_id, currency_id)
);

CREATE TABLE IF NOT EXISTS wallet_ledger (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    
    -- "credit" (adding money) or "debit" (spending money)
    transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('credit', 'debit')),
    
    -- Always store the absolute change amount as a positive number
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    
    -- System tracking metadata
    reference_id UUID NOT NULL UNIQUE, -- Idempotency key from game engine (e.g., match_id, purchase_id)
    reason VARCHAR(50) NOT NULL,       -- 'match_reward', 'store_purchase', 'admin_grant'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(wallet_id) REFERENCES wallets(id)
);

-- Crucial index for summing up logs or pulling history quickly
CREATE INDEX idx_wallet_tx_wallet_id ON wallet_ledger(wallet_id);