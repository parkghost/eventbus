package eventbus

import (
	"sync"
	"testing"
)

type SimpleEvent struct {
}

func (self SimpleEvent) Event() string {
	return ""
}

type SubtypeEvent string

func (self SubtypeEvent) Event() string {
	return string(self)
}

type EventSubscriber struct {
	expected func(evt Event)
}

func (self EventSubscriber) OnEvent(evt Event) {
	self.expected(evt)
}

func TestEventBus(t *testing.T) {

	eventbus := &EventBus{
		handlers: make(map[string]map[EventHandler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    false,
	}

	var simpleEvent = SimpleEvent{}
	var helloEvent = SubtypeEvent("hello")
	var hiEvent = SubtypeEvent("hi")

	expectedReceiveEvent := func(expected Event, t *testing.T, count int) func(Event) {
		received := 0
		return func(actual Event) {
			received += 1
			if actual != expected {
				t.Errorf("expected %s, get %s", expected, actual)
			}

			if received > count {
				t.Errorf("received events is over expected count %d, get event %T", count, actual)
			}
		}
	}

	var a, b, c, d = &EventSubscriber{expectedReceiveEvent(simpleEvent, t, 1)}, &EventSubscriber{expectedReceiveEvent(helloEvent, t, 1)}, &EventSubscriber{expectedReceiveEvent(helloEvent, t, 2)}, &EventSubscriber{expectedReceiveEvent(hiEvent, t, 2)}

	eventbus.Subscribe(simpleEvent, a)
	eventbus.Publish(simpleEvent)

	eventbus.Subscribe(helloEvent, b)
	eventbus.Subscribe(helloEvent, c)
	eventbus.Subscribe(hiEvent, d)
	eventbus.Publish(helloEvent)
	eventbus.Publish(hiEvent)

	eventbus.Unsubscribe(helloEvent, b)
	eventbus.Publish(helloEvent)
	eventbus.Publish(hiEvent)

}

func TestAsyncEventBus(t *testing.T) {

	eventbus := &EventBus{
		handlers: make(map[string]map[EventHandler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    true,
	}

	var simpleEvent = SimpleEvent{}
	var helloEvent = SubtypeEvent("hello")
	var hiEvent = SubtypeEvent("hi")
	wg := &sync.WaitGroup{}

	expectedReceiveEvent := func(expected Event, t *testing.T, count int, wg *sync.WaitGroup) func(Event) {
		received := 0
		wg.Add(count)
		return func(actual Event) {
			received += 1
			wg.Add(-1)
			if actual != expected {
				t.Errorf("expected %s, get %s", expected, actual)
			}

			if received > count {
				t.Errorf("received events is over expected count %d, get event %T", count, actual)
			}
		}
	}

	var a, b, c, d = &EventSubscriber{expectedReceiveEvent(simpleEvent, t, 1, wg)}, &EventSubscriber{expectedReceiveEvent(helloEvent, t, 1, wg)}, &EventSubscriber{expectedReceiveEvent(helloEvent, t, 2, wg)}, &EventSubscriber{expectedReceiveEvent(hiEvent, t, 2, wg)}

	eventbus.Subscribe(simpleEvent, a)
	eventbus.Publish(simpleEvent)

	eventbus.Subscribe(helloEvent, b)
	eventbus.Subscribe(helloEvent, c)
	eventbus.Subscribe(hiEvent, d)
	eventbus.Publish(helloEvent)
	eventbus.Publish(hiEvent)

	eventbus.Unsubscribe(helloEvent, b)
	eventbus.Publish(helloEvent)
	eventbus.Publish(hiEvent)
	wg.Wait()
}

func TestEventChannel(t *testing.T) {

	eventbus := &EventBus{
		handlers: make(map[string]map[EventHandler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    true,
	}

	var simpleEvent = SimpleEvent{}

	start := make(chan bool)
	worker := func() {
		ch := NewEventChannel()
		eventbus.Subscribe(simpleEvent, ch)
		close(start)
		actual := <-ch.C

		if actual != simpleEvent {
			t.Errorf("expected %s, get %s", simpleEvent, actual)
		}
	}
	go worker()
	<-start

	eventbus.Publish(simpleEvent)
}

func TestResolveType(t *testing.T) {
	var evt1 = SimpleEvent{}
	result1 := resolveType(evt1)
	if result1 != "eventbus.SimpleEvent" {
		t.Errorf("expected eventbus.SimpleEvent, get %s", result1)
	}

	var evt2 = SubtypeEvent("hello")
	result2 := resolveType(evt2)
	if result2 != "eventbus.SubtypeEvent.hello" {
		t.Errorf("expected eventbus.SubtypeEvent.hello, get %s", result2)
	}
}
