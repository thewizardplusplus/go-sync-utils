package syncutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConcurrentHandler(test *testing.T) {
	innerHandler := new(MockHandler)
	handler := NewConcurrentHandler(1000, innerHandler)

	mock.AssertExpectationsForObjects(test, innerHandler)
	for _, field := range []interface{}{
		handler.dataChannel,
		handler.stoppingCtx,
		handler.stoppingCtxCanceller,
	} {
		assert.NotNil(test, field)
	}
	assert.Len(test, handler.dataChannel, 0)
	assert.Equal(test, 1000, cap(handler.dataChannel))
	assert.Equal(test, innerHandler, handler.innerHandler)
	assert.Equal(test, &startModeHolder{}, handler.startMode)
}

func TestConcurrentHandler_Handle(test *testing.T) {
	dataChannel := make(chan interface{}, 1)
	handler := ConcurrentHandler{dataChannel: dataChannel}
	handler.Handle("data")

	gotData := <-handler.dataChannel
	assert.Equal(test, "data", gotData)
}

func TestConcurrentHandler_starting(test *testing.T) {
	type fields struct {
		dataGroup []interface{}

		startMode            *startModeHolder
		stoppingCtxCanceller ContextCancellerInterface
	}

	for _, data := range []struct {
		name         string
		fields       fields
		startHandler func(ctx context.Context, handler ConcurrentHandler)
	}{
		{
			name: "with the Start() method",
			fields: fields{
				dataGroup: []interface{}{"one", "two"},

				startMode: &startModeHolder{},
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return().Times(1)

					return stoppingCtxCanceller
				}(),
			},
			startHandler: func(ctx context.Context, handler ConcurrentHandler) {
				handler.Start(ctx)
			},
		},
		{
			name: "with the StartConcurrently() method",
			fields: fields{
				dataGroup: []interface{}{"one", "two"},

				startMode: &startModeHolder{},
				stoppingCtxCanceller: func() ContextCancellerInterface {
					stoppingCtxCanceller := new(MockContextCancellerInterface)
					stoppingCtxCanceller.On("CancelContext").Return().Times(1)

					return stoppingCtxCanceller
				}(),
			},
			startHandler: func(ctx context.Context, handler ConcurrentHandler) {
				handler.StartConcurrently(ctx, 10)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			innerHandler := new(MockHandler)
			for _, data := range data.fields.dataGroup {
				innerHandler.On("Handle", context.Background(), data).Return().Times(1)
			}

			dataChannel := make(chan interface{}, len(data.fields.dataGroup))
			for _, data := range data.fields.dataGroup {
				dataChannel <- data
			}
			close(dataChannel)

			handler := ConcurrentHandler{
				dataChannel:  dataChannel,
				innerHandler: innerHandler,

				startMode:            data.fields.startMode,
				stoppingCtxCanceller: data.fields.stoppingCtxCanceller.CancelContext,
			}
			data.startHandler(context.Background(), handler)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.stoppingCtxCanceller,
				innerHandler,
			)
		})
	}
}

func TestConcurrentHandler_Stop(test *testing.T) {
	stoppingCtx, stoppingCtxCanceller := context.WithCancel(context.Background())
	stoppingCtxCanceller()

	dataChannel := make(chan interface{})
	handler := ConcurrentHandler{
		dataChannel: dataChannel,

		stoppingCtx: stoppingCtx,
	}
	handler.Stop()

	isNotClosed := true
	select {
	case _, isNotClosed = <-handler.dataChannel:
	default: // to prevent blocking
	}

	assert.False(test, isNotClosed)
}
