CREATE TABLE IF NOT EXISTS defi.stellar_dex_pools (
    pool_contract_id varchar NOT NULL,
    token_a_type varchar,
    token_a_account varchar,
    token_b_type varchar,
    token_b_account varchar,
    factory_contract_id varchar,
    parsed_json jsonb,
    event_id varchar NOT NULL,
    ledger int8,
    ledger_closed_at timestamp,
    topic text[],
    value text,
    transaction_hash text,
    process_time timestamp,
    block_date date,
    PRIMARY KEY (pool_contract_id)
);

CREATE INDEX on defi.stellar_dex_pools (block_date);
CREATE INDEX on defi.stellar_dex_pools (token_a_account);
CREATE INDEX on defi.stellar_dex_pools (token_b_account);
