package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func toStringPtr(s string) *string { return &s }

func Test_ToLinesSettings_GetValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings        ToLinesSettings
		indent          string
		fieldPrefix     string
		lastFieldPrefix string
	}{
		"empty settings": {
			indent:          "    ",
			fieldPrefix:     "├── ",
			lastFieldPrefix: "└── ",
		},
		"indent set": {
			settings: ToLinesSettings{
				Indent: toStringPtr("a"),
			},
			indent:          "a",
			fieldPrefix:     "├── ",
			lastFieldPrefix: "└── ",
		},
		"all fields set": {
			settings: ToLinesSettings{
				Indent:          toStringPtr("a"),
				FieldPrefix:     toStringPtr("b"),
				LastFieldPrefix: toStringPtr("c"),
			},
			indent:          "a",
			fieldPrefix:     "b",
			lastFieldPrefix: "c",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			indent, fieldPrefix, lastFieldPrefix :=
				testCase.settings.GetValues()

			assert.Equal(t, testCase.indent, indent)
			assert.Equal(t, testCase.fieldPrefix, fieldPrefix)
			assert.Equal(t, testCase.lastFieldPrefix, lastFieldPrefix)
		})
	}
}
