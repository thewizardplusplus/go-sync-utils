package syncutils

import (
	"reflect"
)

// UnboundedSend ...
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
