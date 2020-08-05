package syncutils

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestMultiWaitGroup_Add(test *testing.T) {
	type args struct {
		delta int
	}

	for _, data := range []struct {
		name       string
		waitGroups MultiWaitGroup
		args       args
	}{
		{
			name:       "without wait groups",
			waitGroups: nil,
			args: args{
				delta: 23,
			},
		},
		{
			name: "with wait groups",
			waitGroups: MultiWaitGroup{
				func() WaitGroup {
					waitGroup := new(MockWaitGroup)
					waitGroup.On("Add", 23).Return()

					return waitGroup
				}(),
				func() WaitGroup {
					waitGroup := new(MockWaitGroup)
					waitGroup.On("Add", 23).Return()

					return waitGroup
				}(),
			},
			args: args{
				delta: 23,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.waitGroups.Add(data.args.delta)

			for _, waitGroup := range data.waitGroups {
				mock.AssertExpectationsForObjects(test, waitGroup)
			}
		})
	}
}
