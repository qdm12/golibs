package verification

import (
	"context"
	"net"
)

// MXLookuper is an interface for looking up MX records on a
// domain name. It's usually a *net.Resolver.
type MXLookuper interface {
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
}
