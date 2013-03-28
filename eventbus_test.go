package eventbus

import (
	"sync"
	"testing"
)

type HelloEvent struct{}

func (self HelloEvent) Event() {}

type HiEvent struct{}

func (self HiEvent) Event() {}

type EventSubscriber struct {
	expected func(evt *Event)
}

func (self EventSubscriber) OnEvent(evt *Event) {
	self.expected(evt)
}

func TestEventBus(t *testing.T) {

	eventbus := &EventBus{
		locks:    NewSegmentedRWLock(32),
		handlers: make(map[string]map[*EventHandler]None),
		async:    false,
	}

	var helloEvent Event = HelloEvent{}
	var hiEvent Event = HiEvent{}

	expectedReceiveEvent := func(expected *Event, t *testing.T, count int) func(*Event) {
		received := 0
		return func(actual *Event) {
			received += 1
			if actual != expected {
				t.Errorf("expected %s, get %s", *expected, *actual)
			}

			if received > count {
				t.Errorf("received events is over expected count %d, get event %s", count, *actual)
			}
		}
	}

	var a, b, c EventHandler = EventSubscriber{expectedReceiveEvent(&helloEvent, t, 1)}, EventSubscriber{expectedReceiveEvent(&helloEvent, t, 2)}, EventSubscriber{expectedReceiveEvent(&hiEvent, t, 2)}

	eventbus.Subscribe(&helloEvent, &a)
	eventbus.Subscribe(&helloEvent, &b)
	eventbus.Subscribe(&hiEvent, &c)
	eventbus.Publish(&helloEvent)
	eventbus.Publish(&hiEvent)

	eventbus.Unsubscribe(&helloEvent, &a)
	eventbus.Publish(&helloEvent)
	eventbus.Publish(&hiEvent)

}

func TestAsyncEventBus(t *testing.T) {

	eventbus := &EventBus{
		locks:    NewSegmentedRWLock(32),
		handlers: make(map[string]map[*EventHandler]None),
		async:    true,
	}

	var helloEvent Event = HelloEvent{}
	var hiEvent Event = HiEvent{}
	wg := &sync.WaitGroup{}

	expectedReceiveEvent := func(expected *Event, t *testing.T, count int, wg *sync.WaitGroup) func(*Event) {
		received := 0
		wg.Add(count)
		return func(actual *Event) {
			received += 1
			wg.Add(-1)
			if actual != expected {
				t.Errorf("expected %s, get %s", *expected, *actual)
			}

			if received > count {
				t.Errorf("received events is over expected count %d, get event %s", count, *actual)
			}
		}
	}

	var a, b, c EventHandler = EventSubscriber{expectedReceiveEvent(&helloEvent, t, 1, wg)}, EventSubscriber{expectedReceiveEvent(&helloEvent, t, 2, wg)}, EventSubscriber{expectedReceiveEvent(&hiEvent, t, 2, wg)}

	eventbus.Subscribe(&helloEvent, &a)
	eventbus.Subscribe(&helloEvent, &b)
	eventbus.Subscribe(&hiEvent, &c)
	eventbus.Publish(&helloEvent)
	eventbus.Publish(&hiEvent)

	eventbus.Unsubscribe(&helloEvent, &a)
	eventbus.Publish(&helloEvent)
	eventbus.Publish(&hiEvent)
	wg.Wait()
}

func TestResolveType(t *testing.T) {
	var evt Event = HelloEvent{}
	result := resolveType(&evt)
	if result != "eventbus.HelloEvent" {
		t.Errorf("expected eventbus.HelloEvent, get %s", result)
	}
}
