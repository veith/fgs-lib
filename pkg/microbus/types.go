package microbus

import (
	"context"
	"sync"
)

// read more here https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997

type DataEvent struct {
	Ctx   context.Context
	Data  interface{}
	Topic string
}

// DataChannel is a channel which can accept an DataEvent
type DataChannel chan DataEvent

// DataChannelSlice is a slice of DataChannels
type DataChannelSlice []DataChannel

// EventBus stores the information about subscribers interested for // a particular topic
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}
