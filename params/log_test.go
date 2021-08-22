package params

import (
	"errors"
	"testing"

	"github.com/qdm12/golibs/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LogCaller(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		caller        logging.Caller
		err           error
	}{
		"hidden": {
			envValue: "hidden",
			caller:   logging.CallerHidden,
		},
		"short": {
			envValue: "short",
			caller:   logging.CallerShort,
		},
		"get error": {
			optionSetters: []OptionSetter{Compulsory()},
			caller:        logging.CallerHidden,
			err:           ErrNoValue,
		},
		"invalid value": {
			envValue: "bla",
			caller:   logging.CallerHidden,
			err:      errors.New("unknown log caller: bla: can be one of: hidden, short"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv: map[string]string{"LOG_CALLER": tc.envValue},
			}
			caller, err := e.LogCaller("LOG_CALLER", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.caller, caller)
			}
		})
	}
}

func Test_LogLevel(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		level         logging.Level
		err           error
	}{
		"debug": {
			envValue: "debug",
			level:    logging.LevelDebug,
		},
		"info": {
			envValue: "info",
			level:    logging.LevelInfo,
		},
		"warning": {
			envValue: "warning",
			level:    logging.LevelWarn,
		},
		"error": {
			envValue: "error",
			level:    logging.LevelError,
		},
		"get error": {
			optionSetters: []OptionSetter{Compulsory()},
			level:         logging.LevelInfo,
			err:           ErrNoValue,
		},
		"invalid value": {
			envValue: "bla",
			level:    logging.LevelInfo,
			err:      errors.New("unknown log level: bla: can be one of: debug, info, warning, error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv: map[string]string{"LOG_LEVEL": tc.envValue},
			}
			level, err := e.LogLevel("LOG_LEVEL", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.level, level)
			}
		})
	}
}
