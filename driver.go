package eventdriver

import (
	"context"

	"github.com/pkg/errors"
)

var (
	driver *EventDriver

	ErrIncorrectPayload = errors.New("incorrect payload for event")
)

// EventDriver provides a tool for handling inner-service communication and responsibility segregation in eventual way.
type EventDriver struct {
	pipe     chan *EventMessage
	cancel   context.CancelFunc // Event loop goroutine cancels
	handlers map[string][]EventHandlerFunc

	logger     Logger
	bufferSize int
}

// EventMessage defines event message.
type EventMessage struct {
	event   string
	ctx     context.Context
	payload interface{}
}

// EventHandlerFunc defines signature func for handling event.
type EventHandlerFunc func(context.Context, interface{}) error

// Init performs initialisation of the EventDriver.
func Init(options ...Option) {
	var (
		ctx = context.Background()
		ed  = &EventDriver{
			handlers: map[string][]EventHandlerFunc{},

			bufferSize: 100,
			logger:     NopLogger{},
		}
	)

	for i := range options {
		options[i].Apply(ed)
	}

	ed.pipe = make(chan *EventMessage, ed.bufferSize)

	ctx, ed.cancel = context.WithCancel(ctx)

	go ed.eventLoop(ctx)

	driver = ed
}

// SubscribeHandler starts event loop and subscribes new handlers afterwords.
func SubscribeHandler(event string, handler EventHandlerFunc) context.CancelFunc {
	driver.handlers[event] = append(driver.handlers[event], handler)

	return func() {
		handler = nil
	}
}

// EmitEvent emits event for concurrent handlers.
func EmitEvent(ctx context.Context, event string, payload interface{}) {
	driver.pipe <- &EventMessage{
		event:   event,
		ctx:     ctx,
		payload: payload,
	}
}

// Close stops event loop and free resources.
func Close() {
	driver.cancel()
}

func (ed *EventDriver) eventLoop(ctx context.Context) {
	for {
		select {
		case message := <-ed.pipe:
			if handlers, ok := ed.handlers[message.event]; ok || len(handlers) == 0 {
				go ed.executeHandlers(handlers, message)

				continue
			}

			ed.logger.Warningf("no handlers subscribed for %s event", message)
		case <-ctx.Done():
			return
		}
	}
}

func (ed *EventDriver) executeHandlers(handlers []EventHandlerFunc, msg *EventMessage) {
	for i := range handlers {
		if handlers[i] == nil {
			continue
		} // handler has been canceled

		if err := handlers[i](msg.ctx, msg.payload); err != nil {
			if err == ErrIncorrectPayload {
				err = errors.Wrap(err, msg.event)
				continue
			}

			ed.logger.Error(errors.Wrapf(err, "failed to handle event %s"))
		}
	}
}
