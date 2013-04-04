/*
The EventBus allows publish-subscribe-style communication between
components without requiring the components to explicitly register with one
another 
*/
package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

type Event interface {
	Event() string
}

type Handler interface {
	OnEvent(evt Event)
}

type Callback struct {
	fn func(evt Event)
}

func (self *Callback) OnEvent(evt Event) {
	self.fn(evt)
}

type Channel struct {
	C chan Event
}

func (self *Channel) OnEvent(evt Event) {
	self.C <- evt
}

type EventBus struct {
	handlers map[string]map[Handler]None
	Locks    *SegmentedRWLock
	Async    bool
}

type None struct{}

func (self *EventBus) Unsubscribe(evt Event, handler Handler) {
	eventType := resolveType(evt)
	locker := self.Locks.locker(eventType)
	locker.Lock()
	defer locker.Unlock()

	handlers, ok := self.handlers[eventType]
	if ok {
		delete(handlers, handler)
	} else {
		panic(fmt.Sprint("the event '%s' doesn't exist", eventType))
	}
}

func (self *EventBus) Subscribe(evt Event, handler Handler) {
	eventType := resolveType(evt)
	locker := self.Locks.locker(eventType)
	locker.Lock()
	defer locker.Unlock()

	handlers, ok := self.handlers[eventType]
	if ok {
		handlers[handler] = None{}
	} else {
		handlers := make(map[Handler]None)
		handlers[handler] = None{}
		self.handlers[eventType] = handlers
	}
}

func (self *EventBus) Publish(evt Event) {
	eventType := resolveType(evt)
	locker := self.Locks.locker(eventType)
	locker.RLock()
	defer locker.RUnlock()

	handlers, ok := self.handlers[eventType]
	if ok {
		for handler, _ := range handlers {
			self.dispatch(evt, handler)
		}
	}
}

func (self *EventBus) dispatch(evt Event, handler Handler) {
	if self.Async {
		go handler.OnEvent(evt)
	} else {
		handler.OnEvent(evt)
	}
}

func resolveType(evt Event) string {
	subEventType := evt.Event()
	if subEventType != "" {
		return reflect.TypeOf(evt).String() + "." + subEventType
	}
	return reflect.TypeOf(evt).String()
}

type SegmentedRWLock struct {
	segments int
	locks    []sync.RWMutex
}

func NewSegmentedRWLock(segments int) *SegmentedRWLock {
	return &SegmentedRWLock{
		segments: segments,
		locks:    make([]sync.RWMutex, segments),
	}
}

func (self *SegmentedRWLock) locker(key string) *sync.RWMutex {
	hash := abs(hash([]byte(key)))
	return &self.locks[hash%self.segments]
}

func abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func hash(bytes []byte) int {
	var h int
	for i := 0; i < len(bytes); i++ {
		h = 31*h ^ int(bytes[i])
	}
	return h
}

var DefaultEventBus = &EventBus{
	handlers: make(map[string]map[Handler]None),
	Locks:    NewSegmentedRWLock(32),
	Async:    true,
}

func Unsubscribe(evt Event, handler Handler) {
	DefaultEventBus.Unsubscribe(evt, handler)
}

func Subscribe(evt Event, handler Handler) {
	DefaultEventBus.Subscribe(evt, handler)
}

func SubscribeWithCallback(evt Event, fn func(evt Event)) {
	Subscribe(evt, &Callback{fn})
}

func Publish(evt Event) {
	DefaultEventBus.Publish(evt)
}

func NewChannel() *Channel {
	return &Channel{
		C: make(chan Event),
	}
}

func NewBufferedChannel(size int) *Channel {
	return &Channel{
		C: make(chan Event, size),
	}
}
