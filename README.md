EventBus
=======

*a concurrent eventbus library written in go*

### Example:

```go
package main

import (
	"fmt"
	"github.com/parkghost/eventbus"
)


//Subscribers and publisher run concurrently so here use WaitGroup to synchronize those 
var start = &sync.WaitGroup{}
var end = &sync.WaitGroup{}

//Define a event which implements eventbus.Event interface
type SimpleEvent struct{}

func (self SimpleEvent) Event() string { return "" }

//Define a event handler which implements eventbus.Handler interface
type Subscriber struct{}

func (self Subscriber) OnEvent(evt eventbus.Event) {
	fmt.Printf("SubscriberOne receives event %T\n", evt)
	end.Done()
}

//SubscriberOne use Handler for subscribe to main.SimpleEvent 
func SubscriberOne() {
	eventbus.Subscribe(SimpleEvent{}, &Subscriber{})
	start.Done()
}

//SubscriberTwo use Channel for subscribe to main.SimpleEvent 
func SubscriberTwo() {
	ch := eventbus.NewChannel()
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

```

*output*
	
	$ go run example/eventbus_example.go
	SubscriberOne receives event main.SimpleEvent
	SubscriberTwo receives event main.SimpleEvent

Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0