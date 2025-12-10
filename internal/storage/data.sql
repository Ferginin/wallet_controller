INSERT INTO wallets (id_wallet, balance)
VALUES
    ('11111111-1111-1111-1111-111111111111', 0),
    ('22222222-2222-2222-2222-222222222222', 50000),
    ('33333333-3333-3333-3333-333333333333', 100000)
ON CONFLICT (id_wallet) DO NOTHING;

-- INSERT INTO wallet_operations (id_operation, id_wallet, operation_type, amount)
-- VALUES
--     (uuid_generate_v4(), '22222222-2222-2222-2222-222222222222', 'DEPOSIT', 20000),
--     (uuid_generate_v4(), '22222222-2222-2222-2222-222222222222', 'WITHDRAW', 10000),
--     (uuid_generate_v4(), '33333333-3333-3333-3333-333333333333', 'DEPOSIT', 50000)
-- ON CONFLICT (id_operation) DO NOTHING;
