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

```
Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0