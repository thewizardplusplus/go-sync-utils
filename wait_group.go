package syncutils

//go:generate mockery --name=WaitGroup --inpackage --case=underscore --testonly

// WaitGroup represents the interface of the sync.WaitGroup type.
// It might be useful for supporting the ability to mock the latter.
type WaitGroup interface {
	Add(delta int)
	Done()
	Wait()
}
