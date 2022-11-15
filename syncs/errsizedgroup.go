package syncs

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// ErrSizedGroup has the same role as WaitingGroup but adds a limit of the amount of goroutines started concurrently.
// Uses similar Go() scheduling as errgroup.Group, thread safe.
type ErrSizedGroup struct {
	options

	ctx context.Context

	err *multierror
	sem *semaphore.Weighted
	wg  sync.WaitGroup
}

// NewErrSizedGroup makes wait group with limited size of alive goroutines
func NewErrSizedGroup(ctx context.Context, size int64, opts ...GroupOption) *ErrSizedGroup {
	options := options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &ErrSizedGroup{
		options: options,
		ctx:     ctx,
		err:     &multierror{},
		sem:     semaphore.NewWeighted(size),
	}
}

// Go calls the given function in a new goroutine.
// Every call will be unblocked, but some goroutines may wait
func (g *ErrSizedGroup) Go(fn func(ctx context.Context) error) {
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

		if err := fn(g.ctx); err != nil {
			g.err.Add(err)
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned
func (g *ErrSizedGroup) Wait() error {
	g.wg.Wait()

	return g.err.errorOrNil()
}
