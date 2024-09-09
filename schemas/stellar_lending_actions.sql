CREATE TABLE IF NOT EXISTS defi.stellar_lending_actions (
    event_id varchar NOT NULL,
    pool_contract_id varchar NOT NULL,
    user_type varchar NOT NULL,
    user_account varchar NOT NULL,
    action_type varchar NOT NULL,
    token_type varchar,
    token_account varchar,
    request_amount varchar,
    btoken_amount varchar,
    btoken_type varchar,
    btoken_account varchar,
    dtoken_amount varchar,
    dtoken_type varchar,
    dtoken_account varchar,
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

CREATE INDEX on defi.stellar_lending_actions (block_date);
CREATE INDEX on defi.stellar_lending_actions (pool_contract_id);
CREATE INDEX on defi.stellar_lending_actions (token_account);
CREATE INDEX on defi.stellar_lending_actions (action_type);
