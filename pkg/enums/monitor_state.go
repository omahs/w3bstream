package enums

//go:generate toolkit gen enum MonitorState
type MonitorState uint8

const (
	MONITOR_STATE_UNKNOWN MonitorState = iota
	MONITOR_STATE__SYNCING
	MONITOR_STATE__SYNCED
	MONITOR_STATE__FAILED_UNKNOWN
	MONITOR_STATE__FAILED_UNIQ
)
