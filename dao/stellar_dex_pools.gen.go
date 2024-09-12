package dao

import (
	"github.com/lib/pq"
	"time"
)

const TableNameStellarDexPool = "stellar_dex_pools"

// StellarDexPool mapped from table <stellar_dex_pools>
type StellarDexPool struct {
	PoolContractID    string         `gorm:"column:pool_contract_id;primaryKey" json:"pool_contract_id"`
	TokenAType        string         `gorm:"column:token_a_type" json:"token_a_type"`
	TokenAAccount     string         `gorm:"column:token_a_account" json:"token_a_account"`
	TokenBType        string         `gorm:"column:token_b_type" json:"token_b_type"`
	TokenBAccount     string         `gorm:"column:token_b_account" json:"token_b_account"`
	FactoryContractID string         `gorm:"column:factory_contract_id" json:"factory_contract_id"`
	ParsedJSON        string         `gorm:"column:parsed_json;type:jsonb" json:"parsed_json"`
	EventID           string         `gorm:"column:event_id;not null" json:"event_id"`
	Ledger            int64          `gorm:"column:ledger" json:"ledger"`
	LedgerClosedAt    time.Time      `gorm:"column:ledger_closed_at" json:"ledger_closed_at"`
	Topic             pq.StringArray `gorm:"column:topic;type:text[]" json:"topic"`
	Value             string         `gorm:"column:value" json:"value"`
	TransactionHash   string         `gorm:"column:transaction_hash" json:"transaction_hash"`
	ProcessTime       time.Time      `gorm:"column:process_time" json:"process_time"`
	BlockDate         time.Time      `gorm:"column:block_date" json:"block_date"`
}

// TableName StellarDexPool's table name
func (*StellarDexPool) TableName() string {
	return TableNameStellarDexPool
}
