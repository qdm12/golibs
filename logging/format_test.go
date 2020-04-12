package logging

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
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

func Test_Debug(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockZap := NewMockZap(mockCtrl)
	mockZap.EXPECT().Debug("prefix: I am 25").Times(1)
	l := &logger{
		prefix:    "prefix: ",
		zapLogger: mockZap,
	}
	l.Debug("I am %d", 25)
}

func Test_Info(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockZap := NewMockZap(mockCtrl)
	mockZap.EXPECT().Info("prefix: I am 25").Times(1)
	l := &logger{
		prefix:    "prefix: ",
		zapLogger: mockZap,
	}
	l.Info("I am %d", 25)
}

func Test_Warn(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockZap := NewMockZap(mockCtrl)
	mockZap.EXPECT().Warn("prefix: I am 25").Times(1)
	l := &logger{
		prefix:    "prefix: ",
		zapLogger: mockZap,
	}
	l.Warn("I am %d", 25)
}

func Test_Error(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockZap := NewMockZap(mockCtrl)
	mockZap.EXPECT().Error("prefix: I am 25").Times(1)
	l := &logger{
		prefix:    "prefix: ",
		zapLogger: mockZap,
	}
	l.Error("I am %d", 25)
}
