package main

import (
	"../../eventbus"
	"fmt"
	"sync"
)

//Subscribers and publisher run concurrently so here use WaitGroup to synchronize those 
var start = &sync.WaitGroup{}
var end = &sync.WaitGroup{}

//Define a event which implements eventbus.Event interface
type SimpleEvent struct{}

func (self SimpleEvent) Event() string { return "" }

//Define a event handler which implements eventbus.EventHandler interface
type EventSubscriber struct{}

func (self EventSubscriber) OnEvent(evt eventbus.Event) {
	fmt.Printf("SubscriberOne receives event %T\n", evt)
	end.Done()
}

//SubscriberOne use EventHandler for subscribe to main.SimpleEvent 
func SubscriberOne() {
	eventbus.Subscribe(SimpleEvent{}, &EventSubscriber{})
	start.Done()
}

//SubscriberTwo use EventChannel for subscribe to main.SimpleEvent 
func SubscriberTwo() {
	ch := eventbus.NewEventChannel()
	eventbus.Subscribe(SimpleEvent{}, ch)
	start.Done()
	fmt.Printf("SubscriberTwo receives event %T\n", <-ch.C)
	end.Done()
}

func main() {
	start.Add(2)
	end.Add(2)

	go SubscriberOne()
	go SubscriberTwo()

	start.Wait()

	//To publish main.SimpleEvent to eventbus and then eventbus will notify who has subscribed to this event
	eventbus.Publish(SimpleEvent{})

	end.Wait()
}
