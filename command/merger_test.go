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
	const wait = 30 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	streamMerger := NewStreamMerger(ctx)
	streamA := ioutil.NopCloser(strings.NewReader("1\n2\n"))
	streamB := ioutil.NopCloser(strings.NewReader("3\n"))
	streamC := ioutil.NopCloser(strings.NewReader("4"))
	streamD := ioutil.NopCloser(strings.NewReader("5"))
	go streamMerger.Merge("A", streamA)
	go streamMerger.Merge("B", streamB)
	go streamMerger.Merge("C", streamC)
	go streamMerger.Merge("D", streamD)

	start := time.Now()
	time.AfterFunc(wait, cancel)
	lines := []string{}
	err := streamMerger.CollectLines(func(line string) {
		lines = append(lines, line)
	})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"A: 1", "A: 2", "B: 3", "C: 4", "D: 5"}, lines)
	delta := time.Since(start)
	if delta < wait || delta > wait+time.Millisecond {
		t.Errorf("test lasted %s instead of %s", delta, wait)
	}
}
