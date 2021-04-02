package logging

// Level is the level of the logger.
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

func (level Level) String() (s string) {
	return string(level)
}
