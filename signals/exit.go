package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForExit(cleanup func(signal string) (exitCode int)) {
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	signal := <-signalsCh
	os.Exit(cleanup(signal.String()))
}
