package syncutils

// MultiWaitGroup allows operating with a set of WaitGroup interfaces
// as a whole. It sequentially calls corresponding methods on each interface
// in the set in the same order in which interfaces are presented.
//
// It might be useful for the simultaneous use of the sync.WaitGroup object
// and its mock. Attention! In this case, the real object must go last to avoid
// data races.
//
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
