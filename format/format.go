package format

import (
	"fmt"
	"strings"
)

// ArgsToString converts arguments to a single string, using the first argument as a format if possible.
func ArgsToString(args ...interface{}) (s string) {
	switch {
	case len(args) == 0:
		return ""
	case len(args) == 1:
		return fmt.Sprintf("%v", args[0])
	default:
		arg0, ok := args[0].(string)
		if ok { // possible format string
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
	return base58Encode(bytes)
}

func base58Encode(data []byte) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	const radix = 58

	zcount := 0
	for zcount < len(data) && data[zcount] == 0 {
		zcount++
	}

	// integer simplification of ceil(log(256)/log(58))
	ceilLog256Div58 := (len(data)-zcount)*555/406 + 1 //nolint:gomnd
	size := zcount + ceilLog256Div58

	output := make([]byte, size)

	high := size - 1
	for _, b := range data {
		i := size - 1
		for carry := uint32(b); i > high || carry != 0; i-- {
			carry += 256 * uint32(output[i]) //nolint:gomnd
			output[i] = byte(carry % radix)
			carry /= radix
		}
		high = i
	}

	// Determine the additional "zero-gap" in the output buffer
	additionalZeroGapEnd := zcount
	for additionalZeroGapEnd < size && output[additionalZeroGapEnd] == 0 {
		additionalZeroGapEnd++
	}

	val := output[additionalZeroGapEnd-zcount:]
	size = len(val)
	for i := range val {
		output[i] = alphabet[val[i]]
	}

	return string(output[:size])
}
