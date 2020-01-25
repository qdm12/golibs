package logging

// Level is the level of the logger (info warn error)
type Level string

const (
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

type Encoding string

const (
	JSONEncoding    Encoding = "json"
	ConsoleEncoding Encoding = "console"
)
