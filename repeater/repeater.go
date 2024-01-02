package repeater

import "context"

type Strategy interface {
	Start(context.Context) <-chan struct{}
}

type Repeater struct {
	Strategy
}

func New(strategy Strategy) *Repeater {
	return &Repeater{
		Strategy: strategy,
	}
}

func (r Repeater) Do(ctx context.Context, fn func(context.Context) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := r.Start(ctx)
	var err error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err() //nolint:wrapcheck
		case _, ok := <-ch:
			if !ok {
				return err
			}
			if err = fn(ctx); err == nil {
				return nil
			}
		}
	}
}
