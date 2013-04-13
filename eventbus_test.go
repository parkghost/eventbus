package eventbus

import (
	"sync"
	"sync/atomic"
	"testing"
)

type SimpleEvent struct {
}

func (self *SimpleEvent) Event() string {
	return ""
}

type SubtypeEvent struct {
	subtype string
}

func (self *SubtypeEvent) Event() string {
	return string(self.subtype)
}

type Subscriber struct {
	expected func(evt Event)
}

func (self *Subscriber) OnEvent(evt Event) {
	self.expected(evt)
}

func TestEventBus(t *testing.T) {

	eventbus := &EventBus{
		handlers: make(map[string]map[Handler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    false,
	}

	var simpleEvent = &SimpleEvent{}
	var helloEvent = &SubtypeEvent{"hello"}
	var hiEvent = &SubtypeEvent{"hi"}

	expectedReceiveEvent := func(expected Event, t *testing.T, count int) func(Event) {
		var received int32 = 0
		return func(actual Event) {
			atomic.AddInt32(&received, 1)
			if actual != expected {
				t.Errorf("expected %s, get %s", expected, actual)
			}

			if atomic.LoadInt32(&received) > int32(count) {
				t.Errorf("received events is over expected count %d, get event %T", count, actual)
			}
		}
	}

	var a, b, c, d = &Subscriber{expectedReceiveEvent(simpleEvent, t, 1)}, &Subscriber{expectedReceiveEvent(helloEvent, t, 1)}, &Subscriber{expectedReceiveEvent(helloEvent, t, 2)}, &Subscriber{expectedReceiveEvent(hiEvent, t, 2)}

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
		handlers: make(map[string]map[Handler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    true,
	}

	var simpleEvent = &SimpleEvent{}
	var helloEvent = &SubtypeEvent{"hello"}
	var hiEvent = &SubtypeEvent{"hi"}
	wg := &sync.WaitGroup{}

	expectedReceiveEvent := func(expected Event, t *testing.T, count int, wg *sync.WaitGroup) func(Event) {
		var received int32 = 0
		wg.Add(count)
		return func(actual Event) {
			atomic.AddInt32(&received, 1)
			wg.Add(-1)
			if actual != expected {
				t.Errorf("expected %s, get %s", expected, actual)
			}

			if atomic.LoadInt32(&received) > int32(count) {
				t.Errorf("received events is over expected count %d, get event %T", count, actual)
			}
		}
	}

	var a, b, c, d = &Subscriber{expectedReceiveEvent(simpleEvent, t, 1, wg)}, &Subscriber{expectedReceiveEvent(helloEvent, t, 1, wg)}, &Subscriber{expectedReceiveEvent(helloEvent, t, 2, wg)}, &Subscriber{expectedReceiveEvent(hiEvent, t, 2, wg)}

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

func TestCallback(t *testing.T) {
	eventbus := &EventBus{
		handlers: make(map[string]map[Handler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    true,
	}

	var simpleEvent = &SimpleEvent{}

	start := make(chan bool)
	worker := func() {
		eventbus.Subscribe(simpleEvent, &Callback{func(actual Event) {
			if actual != simpleEvent {
				t.Errorf("expected %s, get %s", simpleEvent, actual)
			}
		}})
		close(start)
	}
	go worker()
	<-start

	eventbus.Publish(simpleEvent)
}

func TestChannel(t *testing.T) {
	eventbus := &EventBus{
		handlers: make(map[string]map[Handler]None),
		Locks:    NewSegmentedRWLock(32),
		Async:    true,
	}

	var simpleEvent = &SimpleEvent{}

	start := make(chan bool)
	worker := func() {
		ch := NewChannel()
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
	var evt1 = &SimpleEvent{}
	result1 := resolveType(evt1)
	if result1 != "*eventbus.SimpleEvent" {
		t.Errorf("expected *eventbus.SimpleEvent, get %s", result1)
	}

	var evt2 = &SubtypeEvent{"hello"}
	result2 := resolveType(evt2)
	if result2 != "*eventbus.SubtypeEvent.hello" {
		t.Errorf("expected *eventbus.SubtypeEvent.hello, get %s", result2)
	}
}

// TODO: concurrentcy testing and benchmark
