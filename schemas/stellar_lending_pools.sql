CREATE TABLE IF NOT EXISTS defi.stellar_lending_pools (
    pool_contract_id varchar NOT NULL,
    pool_name varchar,
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


CREATE INDEX on defi.stellar_lending_pools (block_date);
CREATE INDEX on defi.stellar_lending_pools (pool_name);
CREATE INDEX on defi.stellar_lending_pools (factory_contract_id);
