package syncutils

import (
	"context"
	"sync"
)

//go:generate mockery --name=Handler --inpackage --case=underscore --testonly

// Handler represents the interface of an abstract handler.
type Handler interface {
	Handle(ctx context.Context, data interface{})
}

// ConcurrentHandler wraps an abstract handler and allows to call it
// concurrently.
type ConcurrentHandler struct {
	dataChannel  chan interface{}
	innerHandler Handler

	startMode            *startModeHolder
	stoppingCtx          context.Context
	stoppingCtxCanceller context.CancelFunc
}

// NewConcurrentHandler creates a concurrent wrapper for the passed abstract
// handler. The buffer size specifies the capacity of the inner channel used
// for passing data to the inner handler.
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

// Handle sends the passed data to the inner handler via the inner channel.
// This method does not block the execution flow.
func (handler ConcurrentHandler) Handle(data interface{}) {
	handler.dataChannel <- data
}

// Start processes data from the inner channel by the inner handler. This method
// performs the processing directly in the caller goroutine and blocks
// the execution flow until the processing will be stopped.
func (handler ConcurrentHandler) Start(ctx context.Context) {
	handler.basicStart(started, func() {
		for data := range handler.dataChannel {
			handler.innerHandler.Handle(ctx, data)
		}
	})
}

// StartConcurrently processes data from the inner channel by the inner handler.
// This method performs the processing in a goroutine pool (the concurrency
// factor specifies goroutine count in the pool). Regardless, it blocks
// the execution flow anyway until the processing will be stopped.
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

// Stop interrupts the processing data from the inner channel by the inner
// handler. This method can be called after both the Start()
// and StartConcurrently() methods. This method blocks the execution flow
// until the interrupting will be completed.
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
