package signals

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit(cleanup func(signal string) (exitCode int)) {
	signalsCh := make(chan os.Signal)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		os.Interrupt,
	)
	signal := <-signalsCh
	signalStr := fmt.Sprintf("%s", signal)
	os.Exit(cleanup(signalStr))
}
