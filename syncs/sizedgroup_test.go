package syncs

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSizedGroup(t *testing.T) {
	swg := NewSizedGroup(context.TODO(), 10)
	var c uint32

	for i := 0; i < 100; i++ {
		swg.Go(func(ctx context.Context) {
			time.Sleep(5 * time.Millisecond)
			atomic.AddUint32(&c, 1)
		})
	}
	assert.True(t, runtime.NumGoroutine() > 50, "goroutines %d", runtime.NumGoroutine())
	swg.Wait()
	assert.Equal(t, c, uint32(100), fmt.Sprintf("%d, not all routines have been executed", c))
}

func TestSizedGroupMaxLimit(t *testing.T) {
	swg := NewSizedGroup(context.Background(), 10)
	var c int32

	for i := 0; i < 100; i++ {
		swg.Go(func(ctx context.Context) {
			atomic.AddInt32(&c, 1)
			defer atomic.AddInt32(&c, -1)

			time.Sleep(5 * time.Millisecond)

			assert.LessOrEqual(t, atomic.LoadInt32(&c), int32(10), "more than 10 goroutines are running")
		})
	}
	swg.Wait()
}

func TestSizedGroup_Cancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	swg := NewSizedGroup(ctx, 10)
	var c int32

	for i := 0; i < 100; i++ {
		swg.Go(func(ctx context.Context) {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Millisecond):
				atomic.AddInt32(&c, 1)
			}
		})
	}
	swg.Wait()
	assert.Less(t, c, int32(100), "not all goroutines should be executed")
}

func TestSizedGroup_CancellationWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	swg := NewSizedGroup(ctx, 1)
	var c int32

	for i := 0; i < 100; i++ {
		swg.Go(func(ctx context.Context) {
			if atomic.LoadInt32(&c) == 10 {
				cancel()
			} else {
				atomic.AddInt32(&c, 1)
			}
			time.Sleep(5 * time.Millisecond)
		})
	}
	swg.Wait()
	assert.Equal(t, c, int32(10), "only 10 goroutines should be executed")
}
