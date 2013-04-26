package eventbus

import (
	"fmt"
)

type SimpleHandler struct {
	quit chan bool
}

func (s *SimpleHandler) OnEvent(event Event) {
	fmt.Printf("receive event %T\n", event)
	close(s.quit)
}

func ExampleHandler() {
	quit := make(chan bool)

	handler := &SimpleHandler{quit}
	Subscribe(&SimpleEvent{}, handler)

	Publish(&SimpleEvent{})

	<-quit
	Unsubscribe(&SimpleEvent{}, handler)
	// Output: receive event *eventbus.SimpleEvent
}

func ExampleHandlerFunc() {
	quit := make(chan bool)

	handler := SubscribeFunc(&SimpleEvent{}, func(evt Event) {
		fmt.Printf("receive event %T\n", evt)
		close(quit)
	})

	Publish(&SimpleEvent{})

	<-quit
	Unsubscribe(&SimpleEvent{}, handler)
	// Output: receive event *eventbus.SimpleEvent
}

func ExampleChannel() {
	quit := make(chan bool)

	handler := NewChannel()
	Subscribe(&SimpleEvent{}, handler)

	go func() {
		fmt.Printf("receive event %T\n", <-handler.C)
		close(quit)
	}()

	Publish(&SimpleEvent{})

	<-quit
	Unsubscribe(&SimpleEvent{}, handler)
	// Output: receive event *eventbus.SimpleEvent
}
