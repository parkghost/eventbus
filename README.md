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
type HelloEvent struct{}

func (self HelloEvent) Event() {}

//define custom subscriber which implements eventbus.EventHandler interface
type EventSubscriber struct{}

func (self EventSubscriber) OnEvent(evt *eventbus.Event) {
	fmt.Printf("receive event: %T\n", *evt)
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

```
Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0