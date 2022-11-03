package models

import (
	"database/sql/driver"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
	"github.com/machinefi/w3bstream/pkg/enums"
)

// Monitor project monitor info
// @def primary                       ID
// @def unique_index UI_monitor_id   MonitorID
//
//go:generate toolkit gen model Monitor --database DB
type Monitor struct {
	datatypes.PrimaryID
	RelMonitor
	RelProject
	MonitorData
	datatypes.OperationTimesWithDeleted
}

type RelMonitor struct {
	MonitorID types.SFID `db:"f_monitor_id" json:"monitorID"`
}

type MonitorData struct {
	State enums.MonitorState `db:"f_state,default='0'"     json:"-"`
	Data  MonitorInfo        `db:"f_data"                  json:"-"`
}

type MonitorInfo struct {
	Type        enums.MonitorType `json:"type,omitempty"`
	ChainTx     *ChaintxInfo      `json:"chainTx,omitempty"`
	ContractLog *ContractlogInfo  `json:"contractLog,omitempty"`
	ChainHeight *ChainHeightInfo  `json:"chainHeight,omitempty"`
}

func (MonitorInfo) DataType(drv string) string { return "text" }

func (v MonitorInfo) Value() (driver.Value, error) {
	return datatypes.JSONValue(v)
}

func (v *MonitorInfo) Scan(src interface{}) error {
	return datatypes.JSONScan(src, v)
}
