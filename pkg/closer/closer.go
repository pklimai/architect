package closer

import (
	"os"
	"os/signal"
	"sync"

	"gitlab.com/zigal0/architect/pkg/logger"
)

// Closer - entity for closing connections and etc.
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New creates new instance of Closer.
// If len(sig) != 0 Closer will automatically call CloseAll when one of them appeared from OS,
// no need to call CloseAll() once more.
func New(sigs ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}

	if len(sigs) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)

			signal.Notify(ch, sigs...)

			<-ch

			signal.Stop(ch)

			c.CloseAll()
		}()
	}

	return c
}

// Add appends close-func to pool.
func (c *Closer) Add(f func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f)
	c.mu.Unlock()
}

// Wait wiaits until chan is closed or something is got from.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll calls all close-functions.
// If you call it more than one time - panic.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))

		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				logger.Errorf("error from Closer: %v", err)
			}
		}
	})
}
