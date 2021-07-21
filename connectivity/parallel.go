package connectivity

import (
	"context"
)

type indexedError struct {
	i   int
	err error
}

func parallelChecks(ctx context.Context, checker SingleChecker,
	urls []string) (errs []error) {
	indexToURL := make(map[int]string, len(urls))
	for i, url := range urls {
		indexToURL[i] = url
	}

	indexedErrorsCh := make(chan indexedError)
	for i, u := range indexToURL {
		go func(i int, url string) {
			err := checker.Check(ctx, url)
			indexedErrorsCh <- indexedError{i: i, err: err}
		}(i, u)
	}

	for range urls {
		checkErr := <-indexedErrorsCh
		errs[checkErr.i] = checkErr.err
	}
	close(indexedErrorsCh)

	return errs
}
