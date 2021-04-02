package logging

import (
	"github.com/qdm12/golibs/format"
)

func formatWithSettings(settings Settings, args ...interface{}) (s string) {
	s = format.ArgsToString(args...)

	if settings.Color != nil {
		s = settings.Color().Sprintf(s)
	}

	if settings.PreProcess != nil {
		s = settings.PreProcess(s)
	}

	return settings.Prefix + s
}
