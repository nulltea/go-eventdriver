package eventdriver

import "context"

// SubscribeHandler subscribes new `handler` for given `event`.
func SubscribeHandler(event string, handler EventHandlerFunc) context.CancelFunc {
	driver.handlers[event] = append(driver.handlers[event], &handler)

	return func() {
		handler = nil
	}
}

// SubscribeChannel subscribes to given `event` and redirects its payload to returned channel
func SubscribeChannel(event string) <- chan interface{} {
	var (
		ch = make(chan interface{})
	)

	SubscribeHandler(event, func(_ context.Context, v interface{}) error {
		ch <- v
		return nil
	})

	return ch
}
