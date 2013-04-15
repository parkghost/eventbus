package eventbus

import (
	"fmt"
)

type SimpleHandler struct {
	quit chan bool
}

func (s *SimpleHandler) OnEvent(event Event) {
	fmt.Printf("receives event %T\n", event)
	close(s.quit)
}

func ExampleHandler() {
	quit := make(chan bool)

	handler := &SimpleHandler{quit}
	Subscribe(&SimpleEvent{}, handler)

	Publish(&SimpleEvent{})

	<-quit
	Unsubscribe(&SimpleEvent{}, handler)
	// Output: receives event *eventbus.SimpleEvent
}

func ExampleCallback() {
	quit := make(chan bool)

	handler := SubscribeWithCallback(&SimpleEvent{}, func(evt Event) {
		fmt.Printf("receives event %T\n", evt)
		close(quit)
	})

	Publish(&SimpleEvent{})

	<-quit
	Unsubscribe(&SimpleEvent{}, handler)
	// Output: receives event *eventbus.SimpleEvent
}

func ExampleChannel() {
	quit := make(chan bool)

	handler := NewChannel()
	Subscribe(&SimpleEvent{}, handler)

	go func() {
		fmt.Printf("receives event %T\n", <-handler.C)
		close(quit)
	}()

	Publish(&SimpleEvent{})

	<-quit
	Unsubscribe(&SimpleEvent{}, handler)
	// Output: receives event *eventbus.SimpleEvent
}
