CREATE TABLE IF NOT EXISTS defi.stellar_dex_liquidity_actions (
    event_id varchar NOT NULL,
    pool_contract_id varchar NOT NULL,
    user_type varchar NOT NULL,
    user_account varchar NOT NULL,
    action_type varchar NOT NULL,
    token_a_type varchar,
    token_a_account varchar,
    token_b_type varchar,
    token_b_account varchar,
    amount_a varchar,
    amount_b varchar,
    liquidity varchar,
    new_reserve_a varchar,
    new_reserve_b varchar,
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

CREATE INDEX on defi.stellar_dex_liquidity_actions (block_date);
CREATE INDEX on defi.stellar_dex_liquidity_actions (pool_contract_id);
CREATE INDEX on defi.stellar_dex_liquidity_actions (token_a_account);
CREATE INDEX on defi.stellar_dex_liquidity_actions (action_type);
CREATE INDEX on defi.stellar_dex_liquidity_actions (token_b_account);
