package runner

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SignalHandler struct {
	cancel  func()
	stop    func()
	once    sync.Once
	cleanup sync.Once
}

func setupSignalHandler(cancel func()) *SignalHandler {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	h := &SignalHandler{
		cancel: cancel,
	}

	h.stop = func() {
		signal.Stop(sigChan)
		close(sigChan)
	}

	go func() {
		<-sigChan
		h.cancel()
	}()

	return h
}

func (h *SignalHandler) Cleanup() {
	h.cleanup.Do(func() {
		h.stop()
	})
}