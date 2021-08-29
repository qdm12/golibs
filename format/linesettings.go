package format

type ToLinesSettings struct {
	// Indent defaults to 4 spaces "    ".
	Indent *string
	// FieldPrefix defaults to "├── ".
	FieldPrefix *string
	// LastFieldPrefix defaults to "└── ".
	LastFieldPrefix *string
}

// GetValues returns the value of each field of the settings
// with sensible defaults if they are left unset.
func (settings *ToLinesSettings) GetValues() (
	indent, fieldPrefix, lastFieldPrefix string,
) {
	indent = "    "
	if settings.Indent != nil {
		indent = *settings.Indent
	}

	fieldPrefix = "├── "
	if settings.FieldPrefix != nil {
		fieldPrefix = *settings.FieldPrefix
	}

	lastFieldPrefix = "└── "
	if settings.LastFieldPrefix != nil {
		lastFieldPrefix = *settings.LastFieldPrefix
	}

	return indent, fieldPrefix, lastFieldPrefix
}
