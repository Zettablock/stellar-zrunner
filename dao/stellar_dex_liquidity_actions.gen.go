// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"time"
	"github.com/lib/pq"
)

const TableNameStellarDexLiquidityAction = "stellar_dex_liquidity_actions"

// StellarDexLiquidityAction mapped from table <stellar_dex_liquidity_actions>
type StellarDexLiquidityAction struct {
	EventID         string    `gorm:"column:event_id;primaryKey" json:"event_id"`
	PoolContractID  string    `gorm:"column:pool_contract_id;not null" json:"pool_contract_id"`
	UserType        string    `gorm:"column:user_type;not null" json:"user_type"`
	UserAccount     string    `gorm:"column:user_account;not null" json:"user_account"`
	ActionType      string    `gorm:"column:action_type;not null" json:"action_type"`
	TokenAType      string    `gorm:"column:token_a_type" json:"token_a_type"`
	TokenAAccount   string    `gorm:"column:token_a_account" json:"token_a_account"`
	TokenBType      string    `gorm:"column:token_b_type" json:"token_b_type"`
	TokenBAccount   string    `gorm:"column:token_b_account" json:"token_b_account"`
	Amount0         string    `gorm:"column:amount_0" json:"amount_0"`
	Amount1         string    `gorm:"column:amount_1" json:"amount_1"`
	Liquidity       string    `gorm:"column:liquidity" json:"liquidity"`
	NewReserve0     string    `gorm:"column:new_reserve_0" json:"new_reserve_0"`
	NewReserve1     string    `gorm:"column:new_reserve_1" json:"new_reserve_1"`
	ParsedJSON      string    `gorm:"column:parsed_json;type:jsonb" json:"parsed_json"`
	Ledger          int64     `gorm:"column:ledger" json:"ledger"`
	LedgerClosedAt  time.Time `gorm:"column:ledger_closed_at" json:"ledger_closed_at"`
	Topic           pq.StringArray    `gorm:"column:topic;type:text[]" json:"topic"`
	Value           string    `gorm:"column:value" json:"value"`
	TransactionHash string    `gorm:"column:transaction_hash" json:"transaction_hash"`
	ProcessTime     time.Time `gorm:"column:process_time" json:"process_time"`
	BlockDate       time.Time `gorm:"column:block_date" json:"block_date"`
}

// TableName StellarDexLiquidityAction's table name
func (*StellarDexLiquidityAction) TableName() string {
	return TableNameStellarDexLiquidityAction
}
