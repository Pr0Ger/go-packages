package repeater

import (
	"context"
	"math"
	"math/rand"
	"time"

	"go.pr0ger.dev/cryptosource"
)

type ExponentialBackoff struct {
	InitialInterval     time.Duration
	MaxElapsedTime      time.Duration
	Multiplier          float64
	RandomizationFactor float64
}

func (s ExponentialBackoff) Start(ctx context.Context) <-chan struct{} {
	timeout := time.After(s.MaxElapsedTime)

	ch := make(chan struct{})
	go func() {
		defer close(ch)

		for count := 0; ; count++ {
			select {
			case <-ctx.Done():
				return
			case <-timeout:
				return
			case ch <- struct{}{}:
			}
			random := 1 - (rand.New(cryptosource.NewSource()).Float64()*2*s.RandomizationFactor - s.RandomizationFactor)
			time.Sleep(s.InitialInterval * time.Duration(math.Pow(s.Multiplier, float64(count))*random))
		}
	}()
	return ch
}

type fixedDelay struct {
	Repeats int
	Delay   time.Duration
}

func (s fixedDelay) Start(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)

		for i := 0; i < s.Repeats; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- struct{}{}:
			}
			time.Sleep(s.Delay)
		}
	}()
	return ch
}

type once struct {
}

func (once) Start(context.Context) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		close(ch)
	}()
	return ch
}

func DefaultExponentialBackoff() Strategy {
	return &ExponentialBackoff{
		InitialInterval:     500 * time.Millisecond,
		MaxElapsedTime:      15 * time.Minute,
		Multiplier:          1.5,
		RandomizationFactor: 0.5,
	}
}

func FixedDelay(repeats int, delay time.Duration) Strategy {
	return &fixedDelay{
		Repeats: repeats,
		Delay:   delay,
	}
}

func Once() Strategy {
	return &once{}
}
