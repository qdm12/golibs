package logging

func formatWithSettings(settings Settings, s string) (formatted string) {
	if settings.Color != nil {
		s = settings.Color().Sprintf(s)
	}

	if settings.PreProcess != nil {
		s = settings.PreProcess(s)
	}

	return settings.Prefix + s
}
