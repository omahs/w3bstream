package enums

//go:generate toolkit gen enum MonitorType
type MonitorType uint8

const (
	MONITOR_TYPE_UNKNOWN MonitorType = iota
	MONITOR_TYPE__CONTRACT_LOG
	MONITOR_TYPE__CHAIN_HEIGHT
	MONITOR_TYPE__CHAIN_TX
)
