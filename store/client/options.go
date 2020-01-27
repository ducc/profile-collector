package client

import "time"

// address of the data service
func WithAddress(address string) Option {
	return func(b *builder) {
		b.address = address
	}
}

// max number of request retries
func WithMaxRetries(retries uint) Option {
	return func(b *builder) {
		b.maxRetries = retries
	}
}

// backoff between connection retries
func WithBackoff(backoff time.Duration) Option {
	return func(b *builder) {
		b.backoff = backoff
	}
}

// enables opentracing - defaults to true but checks if you have a registered tracer first so this should only be set
// if you explicitly want to disable tracing
func WithTracing(tracing bool) Option {
	return func(b *builder) {
		b.useTracing = tracing
	}
}
