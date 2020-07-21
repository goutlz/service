package service

import (
	"os"
	"os/signal"
)

type logicFunc func(stopping <-chan struct{}) struct{}

func Run(logic logicFunc) {
	stopChan := make(chan struct{})
	doneChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)

	defer func() {
		close(stopChan)
		close(doneChan)
		close(sigChan)
	}()

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	go func() {
		doneChan <- logic(stopChan)
	}()

	go func() {
		<-sigChan
		stopChan <- struct{}{}
	}()

	<-doneChan
}
