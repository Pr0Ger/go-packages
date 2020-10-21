package repeater_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.pr0ger.dev/x/repeater"
)

func TestRepeater_Do(t *testing.T) {
	e := errors.New("error")

	called := 0
	fun := func(ctx context.Context) error {
		called++
		return e
	}

	err := repeater.New(repeater.Once()).Do(context.Background(), fun)
	assert.EqualError(t, err, "error")
	assert.Equal(t, 1, called)

	e = nil
	called = 0
	err = repeater.New(repeater.Once()).Do(context.Background(), fun)
	assert.NoError(t, err)
	assert.Equal(t, 1, called)
}

func TestRepeater_DoContextCancellation(t *testing.T) {
	e := errors.New("error")

	called := 0
	fun := func(ctx context.Context) error {
		called++
		return e
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repeater.New(repeater.Once()).Do(ctx, fun)
	assert.EqualError(t, err, "context canceled")
	assert.Equal(t, 0, called)
}
