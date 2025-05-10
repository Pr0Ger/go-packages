package repeater

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const maxTries = 100

type errorEntry struct {
	format string
	args   []interface{}
}

type TestStrategySuite struct {
	suite.Suite

	errorLog []errorEntry
}

func (s *TestStrategySuite) Errorf(format string, args ...interface{}) {
	s.errorLog = append(s.errorLog, errorEntry{
		format: format,
		args:   args,
	})
}

func (s *TestStrategySuite) TestExponentialBackoff() {
	// tick   delay  total
	//    1   0.000  0.000
	//    2   0.100  0.100
	//    3   0.200  0.300
	//    4   0.400  0.700
	//    5   0.800  1.500
	//    6   1.600  end
	s.repeatRunner(func(t assert.TestingT) bool {
		result := true

		ch := ExponentialBackoff{
			InitialInterval:     100 * time.Millisecond,
			MaxElapsedTime:      2000 * time.Millisecond,
			Multiplier:          2,
			RandomizationFactor: 0,
		}.Start(context.Background())

		now := time.Now()
		count := 0
		for range ch {
			if count != 0 {
				duration := time.Since(now)
				expectedDelay := (1 << (count - 1)) * 100 * time.Millisecond

				result = result &&
					assert.Greater(t, duration, expectedDelay-time.Millisecond, "too fast") &&
					assert.Less(t, duration, expectedDelay+5*time.Millisecond, "too slow")
			}
			now = time.Now()

			count++
		}
		result = result &&
			assert.Equal(t, 5, count, "should emit 5 events")

		return result
	})
}

func (s *TestStrategySuite) TestExponentialBackoffWithRandomization() {
	// tick		delay 	randomized delay	total
	//    1   	0.000   0.000               0.000
	//    2  	0.100   0.050-0.150         >=0.050, <= 0.150
	//    3   	0.200   0.100-0.300         >=0.150, <= 0.450
	//    4   	0.400   0.200-0.600         >=0.350, <= 1.050 (end)
	//    5   	0.800   0.400-1.200         >=0.950, <= 2.250 (end)
	//    6   	1.600   0.800-2.400         >=1.550 (end)
	s.repeatRunner(func(t assert.TestingT) bool {
		result := true

		ch := ExponentialBackoff{
			InitialInterval:     100 * time.Millisecond,
			MaxElapsedTime:      1000 * time.Millisecond,
			Multiplier:          2,
			RandomizationFactor: 0.5,
		}.Start(context.Background())

		now := time.Now()
		count := 0
		for range ch {
			if count != 0 {
				duration := time.Since(now)
				expectedDelay := (1 << (count - 1)) * 100 * time.Millisecond

				result = result &&
					assert.Greater(t, duration, expectedDelay/2-time.Millisecond, "too fast") &&
					assert.Less(t, duration, expectedDelay*3/2+5*time.Millisecond, "too slow")
			}
			now = time.Now()

			count++
		}

		result = result &&
			assert.GreaterOrEqual(t, count, 4, "should have at least 4 tries") &&
			assert.LessOrEqual(t, count, 6, "should have no more than 6 tries")

		return result
	})
}

func (s *TestStrategySuite) TestExponentialBackoffCancellation() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	ch := DefaultExponentialBackoff().Start(ctx)

	count := 0
	for range ch {
		count++
	}
	s.Less(count, 10, "channel should be closed before all elements emitted")
}

func (s *TestStrategySuite) TestFixedDelay() {
	s.repeatRunner(func(t assert.TestingT) bool {
		result := true
		var now time.Time

		ch := FixedDelay(5, 10*time.Millisecond).Start(context.Background())

		count := 0
		for range ch {
			count++

			if count != 1 {
				result = result &&
					assert.Greater(t, int64(time.Since(now)), int64(9*time.Millisecond), "too fast") &&
					assert.Less(t, int64(time.Since(now)), int64(13*time.Millisecond), "too slow")
			}
			now = time.Now()
		}

		result = result && assert.Equal(t, 5, count, "should be called 5 times")

		return result
	})
}

func (s *TestStrategySuite) TestFixedDelayCancellation() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	ch := FixedDelay(10, 10*time.Millisecond).Start(ctx)

	count := 0
	for range ch {
		count++
	}
	s.Less(count, 10, "channel should be closed before all elements emitted")
}

func (s *TestStrategySuite) TestOnce() {
	ch := Once().Start(context.Background())

	count := 0
	for range ch {
		count++
	}

	s.Equal(1, count, "should emit one tick")
}

func (s *TestStrategySuite) repeatRunner(fn func(t assert.TestingT) bool) {
	s.T().Helper()

	success := false
	for i := 0; i < maxTries && !success; i++ {
		s.Run(fmt.Sprintf("attempt %d", i), func() {
			s.errorLog = nil
			result := fn(s)
			if result {
				success = true
			} else {
				s.T().Skipf("failed attempt (assertions failed: %d)", len(s.errorLog))
			}
		})
	}
	if !success {
		// if all attempts failed print log from the last one
		for _, logEntry := range s.errorLog {
			s.T().Errorf(logEntry.format, logEntry.args...)
		}
	}
}

func TestStrategies(t *testing.T) {
	suite.Run(t, new(TestStrategySuite))
}
