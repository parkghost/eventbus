package main

import (
	"fmt"
	"github.com/parkghost/eventbus"
)

//define custom event which implements eventbus.Event interface
type HelloEvent struct{}

func (self HelloEvent) String() string {
	return "Hello"
}
func (self HelloEvent) Event() {}

//define custom subscriber which implements eventbus.EventHandler interface
type EventSubscriber struct{}

func (self EventSubscriber) OnEvent(evt *eventbus.Event) {
	fmt.Printf("receive event: %s", *evt)
}

func main() {
	var helloEvent eventbus.Event = HelloEvent{}
	var handler eventbus.EventHandler = EventSubscriber{}

	//the handler subscribe to HelloEvent via eventbus
	eventbus.Subscribe(&helloEvent, &handler)

	//to publish HelloEvent to eventbus and then eventbus will notify who has subscribed to this event
	eventbus.Publish(&helloEvent)

	//to block main goroutine for print out async result from eventbus
	var input string
	fmt.Scanln(&input)
}
