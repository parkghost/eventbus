package main

import (
	"../../eventbus"
	"fmt"
)

//define custom event which implements eventbus.Event interface
type SimpleEvent struct{}

func (self SimpleEvent) Event() string { return "" }

//define custom subscriber which implements eventbus.EventHandler interface
type EventSubscriber struct{}

func (self EventSubscriber) OnEvent(evt eventbus.Event) {
	//OUTPUT: receive event: main.SimpleEvent
	fmt.Printf("receive event: %T\n", evt)
}

func main() {

	var simpleEvent = SimpleEvent{}
	var handler = &EventSubscriber{}

	//the handler subscribe to main.SimpleEvent via eventbus
	eventbus.Subscribe(simpleEvent, handler)

	//to publish main.SimpleEvent to eventbus and then eventbus will notify who has subscribed to this event
	eventbus.Publish(simpleEvent)

	//to block main goroutine for print out async result from eventbus
	var input string
	fmt.Scanln(&input)
}
