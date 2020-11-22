package microbus

import (
	"fmt"
	"github.com/spf13/viper"
)

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

func logEvent(data DataEvent) {
	if viper.GetBool("messaging.log") {
		fmt.Printf("Topic: %s;\n", data.Topic)
	}
}

func (bus *EventBus)RegisterSubscriptionsOnBus(subscriptionList map[string]func(event DataEvent)) {
	for eventName, subscriberFunc := range subscriptionList {
		ch1 := make(chan DataEvent)
		bus.Subscribe(eventName, ch1)
		subscriberFunc := subscriberFunc
		go func() {
			for {
				select {
				case e := <-ch1:
					logEvent(e)
					go func() {
						subscriberFunc(e)

					}()
				}
}
		}()
	}
}


func RegisterSubscriptionsOnBus(subscriptionList map[string]func(event DataEvent), bus *EventBus) {
	for eventName, subscriberFunc := range subscriptionList {
		ch1 := make(chan DataEvent)
		bus.Subscribe(eventName, ch1)
		subscriberFunc := subscriberFunc
		go func() {
			for {
				select {
				case e := <-ch1:
					logEvent(e)
					go func() {
						subscriberFunc(e)

					}()
				}
			}
		}()
	}
}
