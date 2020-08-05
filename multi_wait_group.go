package syncutils

// MultiWaitGroup ...
type MultiWaitGroup []WaitGroup

// Add ...
func (waitGroups MultiWaitGroup) Add(delta int) {
	for _, waitGroup := range waitGroups {
		waitGroup.Add(delta)
	}
}

// Done ...
func (waitGroups MultiWaitGroup) Done() {
	for _, waitGroup := range waitGroups {
		waitGroup.Done()
	}
}

// Wait ...
func (waitGroups MultiWaitGroup) Wait() {
	for _, waitGroup := range waitGroups {
		waitGroup.Wait()
	}
}
