package event

import (
	"fmt"
	"reflect"
	"sync"
)

// Bus for handlers and callbacks.
type Bus struct {
	sync.Mutex
	handlers map[string][]*handler
}

type handler struct {
	callback reflect.Value
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *Bus) Subscribe(topic string, fn interface{}) error {
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	bus.Lock()
	defer bus.Unlock()
	bus.handlers[topic] = append(bus.handlers[topic], &handler{reflect.ValueOf(fn)})
	return nil
}

// Publish executes callback defined for a topic.
// Any additional argument will be transferred to the callback.
func (bus *Bus) Publish(topic string, args ...interface{}) {
	bus.Lock()
	defer bus.Unlock()
	if handlers, ok := bus.handlers[topic]; ok && 0 < len(handlers) {
		for _, handler := range handlers {
			go bus.doPublishAsync(handler, topic, args...)
		}
	}
}

func (bus *Bus) doPublishAsync(handler *handler, topic string, args ...interface{}) {
	passedArguments := bus.setUpPublish(handler, args...)
	handler.callback.Call(passedArguments)
}

func (bus *Bus) setUpPublish(callback *handler, args ...interface{}) []reflect.Value {
	funcType := callback.callback.Type()
	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments
}

// New returns new Bus with empty handlers.
func New() *Bus {
	return &Bus{
		sync.Mutex{},
		make(map[string][]*handler),
	}
}
