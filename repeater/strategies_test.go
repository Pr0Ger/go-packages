package repeater

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoff(t *testing.T) {
	// tick   delay  total
	//    1   0.000  0.000
	//    2   0.100  0.100
	//    3   0.200  0.300
	//    4   0.400  0.700
	//    5   0.800  1.500
	//    6   1.600  3.100 — past the 2.000 deadline, never emitted
	synctest.Test(t, func(t *testing.T) {
		ch := ExponentialBackoff{
			InitialInterval:     100 * time.Millisecond,
			MaxElapsedTime:      2000 * time.Millisecond,
			Multiplier:          2,
			RandomizationFactor: 0,
		}.Start(context.Background())

		var delays []time.Duration
		now := time.Now()
		for range ch {
			delays = append(delays, time.Since(now))
			now = time.Now()
		}

		assert.Equal(t, []time.Duration{
			0,
			100 * time.Millisecond,
			200 * time.Millisecond,
			400 * time.Millisecond,
			800 * time.Millisecond,
		}, delays)
	})
}

func TestExponentialBackoffWithRandomization(t *testing.T) {
	// tick   nominal  randomized   earliest total
	//    1   0.000    0.000        0.000
	//    2   0.100    0.050-0.150  0.050
	//    3   0.200    0.100-0.300  0.150
	//    4   0.400    0.200-0.600  0.350 (may land past the 1.000 deadline)
	//    5   0.800    0.400-1.200  0.750 (may land past the 1.000 deadline)
	//    6   1.600    0.800-2.400  1.550 — always past the deadline
	synctest.Test(t, func(t *testing.T) {
		ch := ExponentialBackoff{
			InitialInterval:     100 * time.Millisecond,
			MaxElapsedTime:      1000 * time.Millisecond,
			Multiplier:          2,
			RandomizationFactor: 0.5,
		}.Start(context.Background())

		start := time.Now()
		now := start
		count := 0
		for range ch {
			assert.Less(t, time.Since(start), 1000*time.Millisecond, "tick %d emitted past MaxElapsedTime", count+1)
			if count != 0 {
				duration := time.Since(now)
				expectedDelay := (1 << (count - 1)) * 100 * time.Millisecond
				assert.GreaterOrEqual(t, duration, expectedDelay/2, "tick %d below randomization window", count+1)
				assert.LessOrEqual(t, duration, expectedDelay*3/2, "tick %d above randomization window", count+1)
			}
			now = time.Now()
			count++
		}

		assert.GreaterOrEqual(t, count, 3, "should have at least 3 ticks")
		assert.LessOrEqual(t, count, 5, "should have no more than 5 ticks")
	})
}

func TestExponentialBackoffCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		ch := DefaultExponentialBackoff().Start(ctx)

		count := 0
		for range ch {
			count++
		}
		assert.Equal(t, 1, count, "only the initial tick fits before cancellation")
	})
}

func TestFixedDelay(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := FixedDelay(5, 10*time.Millisecond).Start(context.Background())

		var delays []time.Duration
		now := time.Now()
		for range ch {
			delays = append(delays, time.Since(now))
			now = time.Now()
		}

		assert.Equal(t, []time.Duration{
			0,
			10 * time.Millisecond,
			10 * time.Millisecond,
			10 * time.Millisecond,
			10 * time.Millisecond,
		}, delays)
	})
}

func TestFixedDelayCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()

		ch := FixedDelay(10, 10*time.Millisecond).Start(ctx)

		count := 0
		for range ch {
			count++
		}
		assert.Equal(t, 3, count, "ticks at 0/10/20ms fit before the 25ms deadline")
	})
}

func TestOnce(t *testing.T) {
	ch := Once().Start(context.Background())

	count := 0
	for range ch {
		count++
	}

	assert.Equal(t, 1, count, "should emit one tick")
}
