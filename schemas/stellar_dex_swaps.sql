CREATE TABLE IF NOT EXISTS defi.stellar_dex_swaps (
    event_id varchar NOT NULL,
    pool_contract_id varchar NOT NULL,
    user_type varchar NOT NULL,
    user_account varchar NOT NULL,
    token_a_type varchar,
    token_a_account varchar,
    token_b_type varchar,
    token_b_account varchar,
    amount_0_in varchar,
    amount_1_in varchar,
    amount_0_out varchar,
    amount_1_out varchar,
    spread_amount int4,
    referral_fee_amount int4,
    parsed_json jsonb,
    ledger int8,
    ledger_closed_at timestamp,
    topic text[],
    value text,
    transaction_hash text,
    process_time timestamp,
    block_date date,
    PRIMARY KEY (event_id)
);

CREATE INDEX on defi.stellar_dex_swaps (block_date);
CREATE INDEX on defi.stellar_dex_swaps (pool_contract_id);
CREATE INDEX on defi.stellar_dex_swaps (token_a_account);
CREATE INDEX on defi.stellar_dex_swaps (token_b_account);
