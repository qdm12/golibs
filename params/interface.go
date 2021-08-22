package params

import (
	liburl "net/url"
	"time"

	"github.com/qdm12/golibs/logging"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Interface

type Interface interface {
	Get(key string, optionSetters ...OptionSetter) (value string, err error)
	Int(key string, optionSetters ...OptionSetter) (n int, err error)
	IntRange(key string, lower, upper int, optionSetters ...OptionSetter) (n int, err error)
	YesNo(key string, optionSetters ...OptionSetter) (yes bool, err error)
	OnOff(key string, optionSetters ...OptionSetter) (on bool, err error)
	Inside(key string, possibilities []string, optionSetters ...OptionSetter) (value string, err error)
	CSV(key string, optionSetters ...OptionSetter) (values []string, err error)
	CSVInside(key string, possibilities []string, optionSetters ...OptionSetter) (values []string, err error)
	Duration(key string, optionSetters ...OptionSetter) (duration time.Duration, err error)
	Port(key string, optionSetters ...OptionSetter) (port uint16, err error)
	ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error)
	ListeningAddress(key string, optionSetters ...OptionSetter) (address, warning string, err error)
	RootURL(key string, optionSetters ...OptionSetter) (rootURL string, err error)
	Path(key string, optionSetters ...OptionSetter) (path string, err error)
	LogCaller(key string, optionSetters ...OptionSetter) (caller logging.Caller, err error)
	LogLevel(key string, optionSetters ...OptionSetter) (level logging.Level, err error)
	URL(key string, optionSetters ...OptionSetter) (URL *liburl.URL, err error)
}
