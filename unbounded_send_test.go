package syncutils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnboundedSend(test *testing.T) {
	type args struct {
		channel interface{}
		data    interface{}
	}

	for _, data := range []struct {
		name     string
		args     args
		wantData interface{}
	}{
		{
			name: "buffered channel",
			args: args{
				channel: make(chan int, 1),
				data:    23,
			},
			wantData: 23,
		},
		{
			name: "unbuffered channel",
			args: args{
				channel: make(chan int),
				data:    23,
			},
			wantData: 23,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			UnboundedSend(data.args.channel, data.args.data)

			gotData, _ := reflect.ValueOf(data.args.channel).Recv()
			assert.Equal(test, data.wantData, gotData.Interface())
		})
	}
}
