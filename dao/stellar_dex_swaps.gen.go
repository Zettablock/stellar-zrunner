// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"github.com/lib/pq"
	"time"
)

const TableNameStellarDexSwap = "stellar_dex_swaps"

// StellarDexSwap mapped from table <stellar_dex_swaps>
type StellarDexSwap struct {
	EventID           string         `gorm:"column:event_id;primaryKey" json:"event_id"`
	PoolContractID    string         `gorm:"column:pool_contract_id;not null" json:"pool_contract_id"`
	UserType          string         `gorm:"column:user_type;not null" json:"user_type"`
	UserAccount       string         `gorm:"column:user_account;not null" json:"user_account"`
	TokenAType        string         `gorm:"column:token_a_type" json:"token_a_type"`
	TokenAAccount     string         `gorm:"column:token_a_account" json:"token_a_account"`
	TokenBType        string         `gorm:"column:token_b_type" json:"token_b_type"`
	TokenBAccount     string         `gorm:"column:token_b_account" json:"token_b_account"`
	AmountAIn         string         `gorm:"column:amount_a_in" json:"amount_a_in"`
	AmountBIn         string         `gorm:"column:amount_b_in" json:"amount_b_in"`
	AmountAOut        string         `gorm:"column:amount_a_out" json:"amount_a_out"`
	AmountBOut        string         `gorm:"column:amount_b_out" json:"amount_b_out"`
	SpreadAmount      int32          `gorm:"column:spread_amount" json:"spread_amount"`
	ReferralFeeAmount int32          `gorm:"column:referral_fee_amount" json:"referral_fee_amount"`
	ParsedJSON        string         `gorm:"column:parsed_json;type:jsonb" json:"parsed_json"`
	Ledger            int64          `gorm:"column:ledger" json:"ledger"`
	LedgerClosedAt    time.Time      `gorm:"column:ledger_closed_at" json:"ledger_closed_at"`
	Topic             pq.StringArray `gorm:"column:topic;type:text[]" json:"topic"`
	Value             string         `gorm:"column:value" json:"value"`
	TransactionHash   string         `gorm:"column:transaction_hash" json:"transaction_hash"`
	ProcessTime       time.Time      `gorm:"column:process_time" json:"process_time"`
	BlockDate         time.Time      `gorm:"column:block_date" json:"block_date"`
}

// TableName StellarDexSwap's table name
func (*StellarDexSwap) TableName() string {
	return TableNameStellarDexSwap
}
