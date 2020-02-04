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
	streamA := ioutil.NopCloser(strings.NewReader("A1\nA2\n"))
	streamB := ioutil.NopCloser(strings.NewReader("B1\nB2\n"))
	streamC := ioutil.NopCloser(strings.NewReader("C1"))
	streamD := ioutil.NopCloser(strings.NewReader("D1\nD2\n"))
	go streamMerger.Merge("streamA", streamA)
	go streamMerger.Merge("streamB", streamB)
	go streamMerger.Merge("streamC", streamC)
	go streamMerger.Merge("streamD", streamD)

	start := time.Now()
	time.AfterFunc(wait, cancel)
	lines := []string{}
	err := streamMerger.CollectLines(func(line string) {
		lines = append(lines, line)
	})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"B1", "B2", "A1", "A2", "C1", "D1", "D2"}, lines)
	delta := time.Since(start)
	if delta < wait || delta > wait+time.Millisecond {
		t.Errorf("test lasted %s instead of %s", delta, wait)
	}
}
