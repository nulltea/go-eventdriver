package eventdriver

import "github.com/go-kit/kit/log"

// An Option configures a EventDriver client.
type Option interface {
	Apply(*EventDriver)
}

// OptionFunc is a function that configures a EventDriver client.
type OptionFunc func(device *EventDriver)

// Apply calls OptionFunc on EventDriver client instance.
func (f OptionFunc) Apply(device *EventDriver) {
	f(device)
}

// WithBufferSize can be used to specify logger implementation for driver client.
// Default is 100.
func WithBufferSize(bufferSize int) Option {
	return OptionFunc(func(ed *EventDriver) {
		ed.bufferSize = bufferSize
	})
}

// WithLogger can be used to specify logger implementation for driver client.
// Default is NopLogger: the one that do anything.
func WithLogger(logger log.Logger) Option {
	return OptionFunc(func(ed *EventDriver) {
		ed.logger = logger
	})
}
