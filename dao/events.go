package dao

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Event mapped from table <events>
type Event struct {
	Id				string			`gorm:"column:id;primaryKey" json:"id"`
	Type			string			`gorm:"column:type;not null" json:"type"`
	ContractId		string			`gorm:"column:contract_id;not null" json:"contract_id"`
	LedgerClosedAt	time.Time		`gorm:"column:ledger_closed_at;type:timestamp" json:"ledger_closed_at"`
	Ledger			int64			`gorm:"column:ledger" json:"ledger"`
	Topic			pq.StringArray	`gorm:"column:topic;type:text[]" json:"topic"`
	Value			string			`gorm:"column:value" json:"value"`
	TransactionHash	string			`gorm:"column:transaction_hash" json:"transaction_hash"`
	ProcessTime		time.Time		`gorm:"column:process_time;type:timestamp" json:"process_time"`
	BlockDate		time.Time		`gorm:"column:block_date;type:timestamp" json:"block_date"`
}

type EventDao struct {
	sourceDB  *gorm.DB
	m         *Event
	schema    string
}

func NewEventDao(schema string, dbs ...*gorm.DB) *EventDao {
	dao := new(EventDao)
	switch len(dbs) {
	case 0:
		panic("database connection required")
	default:
		dao.sourceDB = dbs[0]
	}
	dao.schema = schema
	return dao
}

func (d *EventDao) List(blockNumber int64) []Event {
	var o []Event
	query := fmt.Sprintf(`
		SELECT id, type, contract_id, ledger, ledger_closed_at, topic, value, transaction_hash, block_date
        FROM %s.events
        WHERE ledger = %d order by id asc;
    `, d.schema, blockNumber)
	d.sourceDB.Raw(query).Scan(&o)
	return o
}
