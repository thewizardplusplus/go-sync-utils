package syncutils_test

import (
	"fmt"
	"math/rand"
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
