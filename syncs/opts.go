package syncs

type options struct {
	preLock bool
}

// GroupOption functional option type
type GroupOption func(o *options)

// WithPreLock option will prevent spawning goroutines if group capacity is exceeded
func WithPreLock(opts *options) {
	opts.preLock = true
}
