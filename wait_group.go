package syncutils

// WaitGroup ...
type WaitGroup interface {
	Add(delta int)
	Done()
	Wait()
}
