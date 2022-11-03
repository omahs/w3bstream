// This is a generated source file. DO NOT EDIT
// Source: enums/monitor_state__generated.go

package enums

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enum"
)

var InvalidMonitorState = errors.New("invalid MonitorState type")

func ParseMonitorStateFromString(s string) (MonitorState, error) {
	switch s {
	default:
		return MONITOR_STATE_UNKNOWN, InvalidMonitorState
	case "":
		return MONITOR_STATE_UNKNOWN, nil
	case "SYNCING":
		return MONITOR_STATE__SYNCING, nil
	case "SYNCED":
		return MONITOR_STATE__SYNCED, nil
	}
}

func ParseMonitorStateFromLabel(s string) (MonitorState, error) {
	switch s {
	default:
		return MONITOR_STATE_UNKNOWN, InvalidMonitorState
	case "":
		return MONITOR_STATE_UNKNOWN, nil
	case "SYNCING":
		return MONITOR_STATE__SYNCING, nil
	case "SYNCED":
		return MONITOR_STATE__SYNCED, nil
	}
}

func (v MonitorState) Int() int {
	return int(v)
}

func (v MonitorState) String() string {
	switch v {
	default:
		return "UNKNOWN"
	case MONITOR_STATE_UNKNOWN:
		return ""
	case MONITOR_STATE__SYNCING:
		return "SYNCING"
	case MONITOR_STATE__SYNCED:
		return "SYNCED"
	}
}

func (v MonitorState) Label() string {
	switch v {
	default:
		return "UNKNOWN"
	case MONITOR_STATE_UNKNOWN:
		return ""
	case MONITOR_STATE__SYNCING:
		return "SYNCING"
	case MONITOR_STATE__SYNCED:
		return "SYNCED"
	}
}

func (v MonitorState) TypeName() string {
	return "github.com/machinefi/w3bstream/pkg/enums.MonitorState"
}

func (v MonitorState) ConstValues() []enum.IntStringerEnum {
	return []enum.IntStringerEnum{MONITOR_STATE__SYNCING, MONITOR_STATE__SYNCED}
}

func (v MonitorState) MarshalText() ([]byte, error) {
	s := v.String()
	if s == "UNKNOWN" {
		return nil, InvalidMonitorState
	}
	return []byte(s), nil
}

func (v *MonitorState) UnmarshalText(data []byte) error {
	s := string(bytes.ToUpper(data))
	val, err := ParseMonitorStateFromString(s)
	if err != nil {
		return err
	}
	*(v) = val
	return nil
}

func (v *MonitorState) Scan(src interface{}) error {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	i, err := enum.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*(v) = MonitorState(i)
	return nil
}

func (v MonitorState) Value() (driver.Value, error) {
	offset := 0
	o, ok := interface{}(v).(enum.ValueOffset)
	if ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}
