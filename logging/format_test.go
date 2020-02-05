package logging

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Run with go test -race
func Test_Concurrency(t *testing.T) {
	l, err := NewEmptyLogger()
	require.NoError(t, err)
	logThousandTimes := func(log func(args ...interface{}), done chan struct{}) {
		defer func() {
			done <- struct{}{}
		}()
		now := time.Now()
		for time.Since(now) < 5*time.Millisecond {
			log("test")
		}
	}
	done := make(chan struct{})
	go logThousandTimes(l.Debug, done)
	go logThousandTimes(l.Debug, done)
	go logThousandTimes(l.Info, done)
	go logThousandTimes(l.Warn, done)
	go logThousandTimes(l.Error, done)
	for i := 0; i < 4; i++ {
		<-done
	}
}
