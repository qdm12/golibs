package logging

// Caller is the configuration to show the caller of the log.
type Caller uint8

const (
	CallerHidden Caller = iota
	CallerShort
)
