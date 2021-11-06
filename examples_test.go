package syncutils_test

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"

	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

type Call struct {
	Method    string
	Arguments []interface{}
}

func (call Call) String() string {
	arguments := strings.Trim(fmt.Sprint(call.Arguments), "[]")
	return fmt.Sprintf("%s(%s)", call.Method, arguments)
}

type MockWaitGroup struct {
	sync.Mutex

	Calls []Call
}

func (mock *MockWaitGroup) Add(delta int) {
	mock.Lock()
	defer mock.Unlock()

	mock.Calls = append(mock.Calls, Call{"Add", []interface{}{delta}})
}

func (mock *MockWaitGroup) Done() {
	mock.Lock()
	defer mock.Unlock()

	mock.Calls = append(mock.Calls, Call{"Done", []interface{}{}})
}

func (mock *MockWaitGroup) Wait() {
	mock.Lock()
	defer mock.Unlock()

	mock.Calls = append(mock.Calls, Call{"Wait", []interface{}{}})
}

type Handler struct {
	Locker    sync.Mutex
	DataGroup []interface{}
}

func (handler *Handler) Handle(ctx context.Context, data interface{}) {
	handler.Locker.Lock()
	defer handler.Locker.Unlock()

	handler.DataGroup = append(handler.DataGroup, data)
}

func ExampleMultiWaitGroup() {
	waitGroupMock := new(MockWaitGroup)
	waitGroups := syncutils.MultiWaitGroup{waitGroupMock, new(sync.WaitGroup)}
	for _, duration := range []time.Duration{
		time.Duration(rand.Intn(100)) * time.Millisecond,
		time.Duration(rand.Intn(100)) * time.Millisecond,
	} {
		waitGroups.Add(1)

		go func(duration time.Duration) {
			defer waitGroups.Done()

			time.Sleep(duration)
		}(duration)
	}

	waitGroups.Wait()

	for _, call := range waitGroupMock.Calls {
		fmt.Println(call)
	}

	// Unordered output:
	// Add(1)
	// Add(1)
	// Done()
	// Done()
	// Wait()
}

func ExampleUnboundedSend() {
	numbers := make(chan int, 2)
	for number := 0; number < 10; number++ {
		syncutils.UnboundedSend(numbers, number)
	}

	for index := 0; index < 10; index++ {
		number := <-numbers
		fmt.Println(number)
	}

	// Unordered output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

func ExampleConcurrentHandler() {
	// start the data handling
	var innerHandler Handler
	concurrentHandler := syncutils.NewConcurrentHandler(1000, &innerHandler)
	go concurrentHandler.StartConcurrently(context.Background(), runtime.NumCPU())

	// handle the data
	for index := 0; index < 10; index++ {
		data := fmt.Sprintf("data #%d", index)
		concurrentHandler.Handle(data)
	}
	concurrentHandler.Stop()

	// print the handled data
	for _, data := range innerHandler.DataGroup {
		fmt.Println(data)
	}

	// Unordered output:
	// data #0
	// data #1
	// data #2
	// data #3
	// data #4
	// data #5
	// data #6
	// data #7
	// data #8
	// data #9
}
