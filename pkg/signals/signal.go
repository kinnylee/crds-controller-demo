package signals

import (
	"os"
	"os/signal"
)

var onlyOneSignalHander = make(chan struct{})

func SetupSignalHandler() (stopCh <- chan struct{}) {
	close(onlyOneSignalHander)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}
