// This is a generated source file. DO NOT EDIT
// Source: enums/monitor_type__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidMonitorType = errors.New("invalid MonitorType type")

func ParseMonitorTypeFromString(s string) (MonitorType, error) {
	switch s {
	default:
		return MONITOR_TYPE_UNKNOWN, InvalidMonitorType
	case "":
		return MONITOR_TYPE_UNKNOWN, nil
	case "CONTRACT_LOG":
		return MONITOR_TYPE__CONTRACT_LOG, nil
	case "CHAIN_HEIGHT":
		return MONITOR_TYPE__CHAIN_HEIGHT, nil
	case "CHAIN_TX":
		return MONITOR_TYPE__CHAIN_TX, nil
	}
}

func ParseMonitorTypeFromLabel(s string) (MonitorType, error) {
	switch s {
	default:
		return MONITOR_TYPE_UNKNOWN, InvalidMonitorType
	case "":
		return MONITOR_TYPE_UNKNOWN, nil
	case "CONTRACT_LOG":
		return MONITOR_TYPE__CONTRACT_LOG, nil
	case "CHAIN_HEIGHT":
		return MONITOR_TYPE__CHAIN_HEIGHT, nil
	case "CHAIN_TX":
		return MONITOR_TYPE__CHAIN_TX, nil
	}
}

func (v MonitorType) Int() int {
	return int(v)
}

func (v MonitorType) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case MONITOR_TYPE_UNKNOWN:
		return ""
	case MONITOR_TYPE__CONTRACT_LOG:
		return "CONTRACT_LOG"
	case MONITOR_TYPE__CHAIN_HEIGHT:
		return "CHAIN_HEIGHT"
	case MONITOR_TYPE__CHAIN_TX:
		return "CHAIN_TX"
	}
}

func (v MonitorType) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case MONITOR_TYPE_UNKNOWN:
		return ""
	case MONITOR_TYPE__CONTRACT_LOG:
		return "CONTRACT_LOG"
	case MONITOR_TYPE__CHAIN_HEIGHT:
		return "CHAIN_HEIGHT"
	case MONITOR_TYPE__CHAIN_TX:
		return "CHAIN_TX"
	}
}

func (v MonitorType) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.MonitorType"
}

func (v MonitorType) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{MONITOR_TYPE__CONTRACT_LOG, MONITOR_TYPE__CHAIN_HEIGHT, MONITOR_TYPE__CHAIN_TX}
}

func (v MonitorType) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidMonitorType
	}
	return []byte(s), nil
}

func (v *MonitorType) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseMonitorTypeFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *MonitorType) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = MonitorType(i)
	return nil
}

func (v MonitorType) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
