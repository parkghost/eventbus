EventBus
=======

*a concurrent eventbus library written in go*

### Example:

```go
import (
	"fmt"
	"github.com/parkghost/eventbus"
	"sync"
)

var start, end sync.WaitGroup

//Define a event that implements eventbus.Event interface
type SimpleEvent struct{}

func (self *SimpleEvent) Event() string { return "" }

//SubscriberOne use Callback to subscribe *main.SimpleEvent 
func SubscriberOne() {
	eventbus.SubscribeWithCallback(&SimpleEvent{}, func(evt eventbus.Event) {
		fmt.Printf("SubscriberOne receives event %T\n", evt)
		end.Done()
	})
	start.Done()
}

//SubscriberTwo use Channel to subscribe *main.SimpleEvent 
func SubscriberTwo() {
	ch := eventbus.NewChannel()
	eventbus.Subscribe(&SimpleEvent{}, ch)
	start.Done()
	fmt.Printf("SubscriberTwo receives event %T\n", <-ch.C)
	end.Done()
}

func main() {
	start.Add(2)
	end.Add(2)

	//Two subscribers are running concurrently
	go SubscriberOne()
	go SubscriberTwo()

	start.Wait()

	//To publish *main.SimpleEvent to eventbus
	eventbus.Publish(&SimpleEvent{})

	end.Wait()
}
```

*output*
	
	$ go run example/eventbus_example.go
	SubscriberOne receives event *main.SimpleEvent
	SubscriberTwo receives event *main.SimpleEvent

Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

This project is licensed under the MIT license