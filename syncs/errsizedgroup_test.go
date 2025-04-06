package syncs

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrSizedGroup(t *testing.T) {
	swg := NewErrSizedGroup(context.TODO(), 10)
	var c uint32

	for i := 0; i < 100; i++ {
		swg.Go(func(context.Context) error {
			time.Sleep(5 * time.Millisecond)
			atomic.AddUint32(&c, 1)
			return nil
		})
	}
	assert.Greaterf(t, runtime.NumGoroutine(), 50, "goroutines %d", runtime.NumGoroutine())
	require.NoError(t, swg.Wait())
	assert.EqualValues(t, 100, c, "%d, not all routines have been executed", c)
}

func TestErrSizedGroupMaxLimit(t *testing.T) {
	swg := NewErrSizedGroup(context.Background(), 10)
	var c int32

	for i := 0; i < 100; i++ {
		swg.Go(func(context.Context) error {
			atomic.AddInt32(&c, 1)
			defer atomic.AddInt32(&c, -1)

			time.Sleep(5 * time.Millisecond)

			assert.LessOrEqual(t, atomic.LoadInt32(&c), int32(10), "more than 10 goroutines are running")

			return nil
		})
	}
	require.NoError(t, swg.Wait())
}

func TestErrSizedGroup_Cancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	swg := NewErrSizedGroup(ctx, 10)
	var c int32

	for i := 0; i < 100; i++ {
		swg.Go(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Millisecond):
				atomic.AddInt32(&c, 1)
			}

			return nil
		})
	}
	require.NoError(t, swg.Wait())
	assert.Less(t, c, int32(100), "not all goroutines should be executed")
}

func TestErrSizedGroup_CancellationWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	swg := NewErrSizedGroup(ctx, 1)
	var c int32

	for i := 0; i < 100; i++ {
		swg.Go(func(context.Context) error {
			if atomic.LoadInt32(&c) == 10 {
				cancel()
			} else {
				atomic.AddInt32(&c, 1)
			}
			time.Sleep(5 * time.Millisecond)

			return nil
		})
	}
	require.NoError(t, swg.Wait())
	assert.EqualValues(t, 10, c, "only 10 goroutines should be executed")
}

func TestErrSizedGroup_WithErrors(t *testing.T) {
	swg := NewErrSizedGroup(context.TODO(), 10)

	for i := 0; i < 100; i++ {
		i := i
		swg.Go(func(context.Context) error {
			if i == 50 {
				return errors.New("error")
			}
			return nil
		})
	}
	err := swg.Wait()
	require.EqualError(t, err, "multierror: 1 errors: [0] error")
}

func TestErrSizedGroupWithPreLock(t *testing.T) {
	swg := NewErrSizedGroup(context.TODO(), 10, WithPreLock)
	var c uint32

	for i := 0; i < 100; i++ {
		swg.Go(func(context.Context) error {
			time.Sleep(5 * time.Millisecond)
			atomic.AddUint32(&c, 1)
			return nil
		})
	}
	assert.Less(t, runtime.NumGoroutine(), 15, "goroutines %d", runtime.NumGoroutine())
	require.NoError(t, swg.Wait())
	assert.EqualValues(t, 100, c, "%d, not all routines have been executed", c)
}
