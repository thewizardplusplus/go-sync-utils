# go-sync-utils

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-sync-utils?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-sync-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-sync-utils)](https://goreportcard.com/report/github.com/thewizardplusplus/go-sync-utils)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-sync-utils.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-sync-utils)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-sync-utils/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-sync-utils)

The library that provides utility entities for syncing.

## Features

- interface of the `sync.WaitGroup` type;
- operating with a set of such interfaces as a whole.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-sync-utils.git
$ cd go-sync-utils
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

## Examples

`syncutils.MultiWaitGroup`:

```go
package main

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

func main() {
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
```

`syncutils.UnboundedSend`:

```go
package main

import (
	"fmt"

	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

func main() {
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
```

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
