package params

import (
	"github.com/qdm12/golibs/verification"

	"fmt"
	"strconv"
)

func verifyListeningPort(listeningPort string, uid int) (warning string, err error) {
	if err := verification.VerifyPort(listeningPort); err != nil {
		return "", fmt.Errorf("listening port: %w", err)
	}
	value, _ := strconv.Atoi(listeningPort)
	const maxPrivilegedPort = 1023
	const minDynamicPort = 49151
	switch {
	case value <= maxPrivilegedPort:
		switch uid {
		case 0:
			return fmt.Sprintf("listening port %s allowed to be in the reserved system ports range as you are running as root", listeningPort), nil
		case -1:
			return fmt.Sprintf("listening port %s allowed to be in the reserved system ports range as you are running in Windows", listeningPort), nil
		default:
			return "", fmt.Errorf("listening port %s cannot be in the reserved system ports range (1 to 1023) when running without root", listeningPort)
		}
	case value > minDynamicPort:
		// dynamic and/or private ports.
		return fmt.Sprintf("listening port %s is in the dynamic/private ports range (above 49151)", listeningPort), nil
	default:
		return "", nil
	}
}

func verifyRootURL(rootURL string) error {
	if !verification.MatchRootURL(rootURL) {
		return fmt.Errorf("root URL %q is invalid", rootURL)
	}
	return nil
}

func verifyHostname(hostname string) error {
	if verification.MatchHostname(hostname) {
		return nil
	}
	return fmt.Errorf("hostname %q is not valid", hostname)
}
