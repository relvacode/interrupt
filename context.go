package interrupt

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type Error struct {
	Signal os.Signal
}

func (e Error) Error() string {
	return fmt.Sprintf("signal: %s", e.Signal)
}

type signalContext struct {
	context.Context

	done chan struct{}
	mu   sync.Mutex
	err  error
}

func (c *signalContext) String() string {
	return "interrupt.Context"
}

func (c *signalContext) Done() <-chan struct{} {
	return c.done
}

func (c *signalContext) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()

	return err
}

var defaultNotify = []os.Signal{os.Interrupt}

// Context returns a context.Context that is cancelled on any of the given os.Signal.
// If no signals are provided then os.Interrupt is used in its place.
// The returned Context will, on cancellation via a signal contain the Err() type of interrupt.Error.
func Context(ctx context.Context, signals ...os.Signal) context.Context {
	notify := make(chan os.Signal, 1)

	if len(signals) == 0 {
		signals = defaultNotify
	}

	signal.Notify(notify, signals...)

	c := &signalContext{
		Context: ctx,
		done:    make(chan struct{}),
	}

	go func() {
		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case sig := <-notify:
			err = Error{Signal: sig}
		}

		signal.Stop(notify)

		c.mu.Lock()
		c.err = err
		c.mu.Unlock()

		close(c.done)
	}()

	return c
}
