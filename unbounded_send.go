package syncutils

import (
	"reflect"
)

// UnboundedSend sends the provided data to the provided channel
// without blocking even if the channel is busy. This function works
// with any types, but the data should have the type assignable to the type
// of the channel item. The order of receiving data may not correspond
// to the order of sending one.
func UnboundedSend(channel interface{}, data interface{}) {
	channelReflection := reflect.ValueOf(channel)
	dataReflection := reflect.ValueOf(data)

	chosenCase, _, _ := reflect.Select([]reflect.SelectCase{
		{
			Dir:  reflect.SelectSend,
			Chan: channelReflection,
			Send: dataReflection,
		},
		{
			Dir: reflect.SelectDefault,
		},
	})
	if chosenCase == 0 {
		return
	}

	go channelReflection.Send(dataReflection)
}
