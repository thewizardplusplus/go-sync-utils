package syncutils

import (
	"context"
	"sync"
)

//go:generate mockery --name=Handler --inpackage --case=underscore --testonly

// Handler ...
type Handler interface {
	Handle(ctx context.Context, data interface{})
}

// ConcurrentHandler ...
type ConcurrentHandler struct {
	dataChannel  chan interface{}
	innerHandler Handler

	startMode            *startModeHolder
	stoppingCtx          context.Context
	stoppingCtxCanceller context.CancelFunc
}

// NewConcurrentHandler ...
func NewConcurrentHandler(
	bufferSize int,
	innerHandler Handler,
) ConcurrentHandler {
	stoppingCtx, stoppingCtxCanceller := context.WithCancel(context.Background())
	return ConcurrentHandler{
		dataChannel:  make(chan interface{}, bufferSize),
		innerHandler: innerHandler,

		startMode:            &startModeHolder{},
		stoppingCtx:          stoppingCtx,
		stoppingCtxCanceller: stoppingCtxCanceller,
	}
}

// Handle ...
func (handler ConcurrentHandler) Handle(data interface{}) {
	handler.dataChannel <- data
}

// Start ...
func (handler ConcurrentHandler) Start(ctx context.Context) {
	handler.basicStart(started, func() {
		for data := range handler.dataChannel {
			handler.innerHandler.Handle(ctx, data)
		}
	})
}

// StartConcurrently ...
func (handler ConcurrentHandler) StartConcurrently(
	ctx context.Context,
	concurrencyFactor int,
) {
	handler.basicStart(startedConcurrently, func() {
		var waiter sync.WaitGroup
		waiter.Add(concurrencyFactor)

		for threadID := 0; threadID < concurrencyFactor; threadID++ {
			go func() {
				defer waiter.Done()

				handler.Start(ctx)
			}()
		}

		waiter.Wait()
	})
}

// Stop ...
func (handler ConcurrentHandler) Stop() {
	close(handler.dataChannel)
	<-handler.stoppingCtx.Done()
}

func (handler ConcurrentHandler) basicStart(
	mode startMode,
	startHandler func(),
) {
	handler.startMode.setStartModeOnce(mode)
	defer func() {
		if handler.startMode.getStartMode() == mode {
			handler.stoppingCtxCanceller()
		}
	}()

	startHandler()
}
