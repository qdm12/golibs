package command

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StreamMerger(t *testing.T) {
	t.Parallel()
	const wait = 100 * time.Millisecond
	ctx := context.Background()
	collectLinesCtx, collectLinesCancel := context.WithCancel(ctx)
	mergeCtx, mergeCancel := context.WithCancel(ctx)
	defer mergeCancel()
	streamMerger := NewStreamMerger()
	streamA := ioutil.NopCloser(strings.NewReader("1\n2\n"))
	streamB := ioutil.NopCloser(strings.NewReader("3\n"))
	streamC := ioutil.NopCloser(strings.NewReader("4"))
	streamD := ioutil.NopCloser(strings.NewReader("5"))
	go streamMerger.Merge(mergeCtx, streamA, MergeName("A"))
	go streamMerger.Merge(mergeCtx, streamB, MergeName("B"))
	go streamMerger.Merge(mergeCtx, streamC, MergeName("C"))
	go streamMerger.Merge(mergeCtx, streamD, MergeName("D"))

	start := time.Now()
	time.AfterFunc(wait/2, mergeCancel)      // closes streams
	time.AfterFunc(wait, collectLinesCancel) // stops collecting lines
	lines := []string{}
	streamMerger.CollectLines(collectLinesCtx, func(line string) {
		lines = append(lines, line)
	}, func(err error) {
		require.NoError(t, err)
	})
	assert.ElementsMatch(t, []string{"A: 1", "A: 2", "B: 3", "C: 4", "D: 5"}, lines)
	delta := time.Since(start)
	if delta < wait || delta > wait+time.Millisecond {
		t.Errorf("test lasted %s instead of %s", delta, wait)
	}
}
