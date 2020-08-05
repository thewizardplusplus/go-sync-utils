package syncutils

//go:generate mockery -name=WaitGroup -inpkg -case=underscore -testonly

// WaitGroup ...
type WaitGroup interface {
	Add(delta int)
	Done()
	Wait()
}
