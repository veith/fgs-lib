package microbus

// start with a bus and then add the topics an subscribers

func NewMicrobus() *EventBus {
	return &EventBus{
		subscribers: map[string]DataChannelSlice{},
	}
}
