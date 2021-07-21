package connectivity

import (
	"context"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Checker,SingleChecker

type Checker interface {
	ParallelChecks(ctx context.Context, urls []string) []error
	SingleChecker
}

type SingleChecker interface {
	Check(ctx context.Context, url string) error
}
