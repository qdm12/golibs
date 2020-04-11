package format

import (
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
)

// ArgsToString converts arguments to a single string, using the first argument as a format if possible
func ArgsToString(args ...interface{}) (s string) {
	switch {
	case len(args) == 0:
		return ""
	case len(args) == 1:
		return fmt.Sprintf("%v", args[0])
	default:
		if arg0, ok := args[0].(string); ok { // possible format string
			if strings.Count(arg0, "%") == len(args[1:]) { // treat as format string
				return fmt.Sprintf(arg0, args[1:]...)
			}
		}
		var words []string
		for _, arg := range args {
			words = append(words, fmt.Sprintf("%v", arg))
		}
		return strings.Join(words, " ")
	}
}

func ReadableBytes(bytes []byte) string {
	return base58.Encode(bytes)
}
