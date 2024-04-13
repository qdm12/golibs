package connectivity

import (
	"context"
)

type Checker interface {
	ParallelChecks(ctx context.Context, urls []string) []error
	SingleChecker
}

type SingleChecker interface {
	Check(ctx context.Context, url string) error
}
