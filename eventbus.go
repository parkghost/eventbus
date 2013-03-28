/*
The EventBus allows publish-subscribe-style communication between
components without requiring the components to explicitly register with one
another 
*/
package eventbus

import (
	"fmt"
	"reflect"
)

type EventBus struct {
	locks    *SegmentedRWLock
	handlers map[string]map[*EventHandler]None
	async    bool
}

type EventHandler interface {
	OnEvent(evt *Event)
}

type Event interface {
	Event()
}

type None struct{}

func (self *EventBus) Unsubscribe(evt *Event, handler *EventHandler) {
	eventType := resolveType(evt)
	locker := self.locks.locker(eventType)
	locker.Lock()
	defer locker.Unlock()

	handlers, ok := self.handlers[eventType]
	if ok {
		delete(handlers, handler)
	} else {
		panic(fmt.Sprint("the event '%s' doesn't exist", *evt))
	}
}

func (self *EventBus) Subscribe(evt *Event, handler *EventHandler) {
	eventType := resolveType(evt)
	locker := self.locks.locker(eventType)
	locker.Lock()
	defer locker.Unlock()

	handlers, ok := self.handlers[eventType]
	if ok {
		handlers[handler] = None{}
	} else {
		handlers := make(map[*EventHandler]None)
		handlers[handler] = None{}
		self.handlers[eventType] = handlers
	}

}

func (self *EventBus) Publish(evt *Event) {
	eventType := resolveType(evt)
	locker := self.locks.locker(eventType)
	locker.RLock()
	defer locker.RUnlock()

	handlers, ok := self.handlers[eventType]
	if ok {
		for handler, _ := range handlers {
			self.dispatch(evt, handler)
		}
	}

}

func (self *EventBus) dispatch(evt *Event, handler *EventHandler) {
	if self.async {
		go EventHandler(*handler).OnEvent(evt)
	} else {
		EventHandler(*handler).OnEvent(evt)
	}
}

func resolveType(evt *Event) string {
	return reflect.TypeOf(*evt).String()
}

var DefaultEventBus = &EventBus{
	locks:    NewSegmentedRWLock(32),
	handlers: make(map[string]map[*EventHandler]None),
	async:    true,
}

func Unsubscribe(evt *Event, handler *EventHandler) {
	DefaultEventBus.Unsubscribe(evt, handler)
}

func Subscribe(evt *Event, handler *EventHandler) {
	DefaultEventBus.Subscribe(evt, handler)
}

func Publish(evt *Event) {
	DefaultEventBus.Publish(evt)
}
