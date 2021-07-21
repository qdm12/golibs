package connectivity

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var _ Checker = new(MixedChecker)
var _ SingleChecker = new(MixedChecker)

var ErrOneOrMoreChecksFailed = errors.New("one or more checks failed")

// NewMixedChecker creates a new mixed checker using the given checkers.
func NewMixedChecker(checkers []Checker) *MixedChecker {
	return &MixedChecker{
		checkers: checkers,
	}
}

type MixedChecker struct {
	checkers []Checker
}

// ParallelChecks verifies the connectivity to each of the urls
// using each of the checkers specified in the MixedChecker.
// It returns a slice of errors with the same indexing and order as the
// urls, meaning that some errors might be nil or not. You should ensure
// to iterate over the errors and check each of them.
func (c *MixedChecker) ParallelChecks(ctx context.Context, urls []string) (errs []error) {
	return parallelChecks(ctx, c, urls)
}

// Check verifies the connectivity to the url using each of the checkers
// specified in the MixedChecker.
func (c *MixedChecker) Check(ctx context.Context, url string) (err error) {
	errorsCh := make(chan error)

	for _, checker := range c.checkers {
		go func(checker Checker) {
			err := checker.Check(ctx, url)
			errorsCh <- err
		}(checker)
	}

	errStrings := make([]string, 0, len(c.checkers))
	for range c.checkers {
		if err := <-errorsCh; err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}
	close(errorsCh)

	if len(errStrings) == 0 {
		return nil
	}

	return fmt.Errorf("%w: for URL %s: %s",
		ErrOneOrMoreChecksFailed, url, strings.Join(errStrings, "; "))
}
