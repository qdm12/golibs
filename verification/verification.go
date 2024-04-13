package verification

import "net"

type Verifier struct {
	mxLookup func(name string) ([]*net.MX, error)
	*Regex
}

func NewVerifier() *Verifier {
	return &Verifier{
		mxLookup: net.LookupMX,
		Regex:    NewRegex(),
	}
}
