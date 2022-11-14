package syncs

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// SizedGroup has the same role as WaitingGroup but adds a limit of the amount of goroutines started concurrently.
// Uses similar Go() scheduling as errgroup.Group, thread safe.
type SizedGroup struct {
	options

	ctx context.Context

	sem *semaphore.Weighted
	wg  sync.WaitGroup
}

// NewSizedGroup makes wait group with limited size of alive goroutines
func NewSizedGroup(ctx context.Context, size int64, opts ...GroupOption) *SizedGroup {
	options := options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &SizedGroup{
		options: options,
		ctx:     ctx,
		sem:     semaphore.NewWeighted(size),
	}
}

// Go calls the given function in a new goroutine.
// Every call will be unblocked, but some goroutines may wait
func (g *SizedGroup) Go(fn func(ctx context.Context)) {
	g.wg.Add(1)

	if g.preLock {
		if err := g.sem.Acquire(g.ctx, 1); err != nil {
			return
		}
	}

	go func() {
		defer g.wg.Done()

		if !g.preLock {
			if err := g.sem.Acquire(g.ctx, 1); err != nil {
				return
			}
		}
		defer g.sem.Release(1)

		fn(g.ctx)
	}()
}

// Wait blocks until all function calls from the Go method have returned
func (g *SizedGroup) Wait() {
	g.wg.Wait()
}
