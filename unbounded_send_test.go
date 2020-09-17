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
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			UnboundedSend(data.args.channel, data.args.data)

			gotData, _ := reflect.ValueOf(data.args.channel).Recv()
			assert.Equal(test, data.wantData, gotData.Interface())
		})
	}
}
