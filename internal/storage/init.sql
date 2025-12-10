CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS wallets (
    id_wallet UUID PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0, -- в копейках
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallet_operations (
    id_operation UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_wallet UUID NOT NULL REFERENCES wallets(id_wallet) ON DELETE CASCADE,
    operation_type VARCHAR(16) NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_wallet_operations_wallet_id_created_at
    ON wallet_operations (id_wallet, created_at);