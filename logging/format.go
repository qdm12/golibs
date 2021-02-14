package logging

import (
	"github.com/qdm12/golibs/format"
)

func formatWithSettings(settings settings, args ...interface{}) (s string) {
	s = format.ArgsToString(args...)

	if settings.color != nil {
		s = settings.color.Sprintf(s)
	}

	if settings.postProcess != nil {
		s = settings.postProcess(s)
	}

	return s
}
