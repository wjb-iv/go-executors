package executors

import (
	"sync"
	"time"
)

// Future - represents the results of some work to perform in a worker
type Future struct {
	waitGrp   *sync.WaitGroup // internal use (to wait gracefully)
	value     interface{}     // the actual return value
	cancelled bool
}

// TimeoutError happens if the timeout expires on 'Get' method
type TimeoutError struct {
}

func (te TimeoutError) Error() string {
	return "Operation timed out."
}

// Get returns the value or times out after 'timeout' duration
func (f *Future) Get(timeout time.Duration) (interface{}, error) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		// If the worker also holding a ref to this Wait Group
		// never calls 'Done()' - this goroutine blocks forever,
		// but we get away with it because the only way this
		// happens is with a panic which is handled within
		// the executors worker by using 'recover()'
		f.waitGrp.Wait()
	}()
	select {
	case <-c:
		return f.value, nil
	case <-time.After(timeout):
		return nil, TimeoutError{}
	}
}

// Private setter called by worker to signal completion
func (f *Future) setVal(val interface{}) {
	if f.cancelled {
		return
	}
	f.value = val
	f.waitGrp.Done()
}

// Private cancel to clean up waitgroup is error happens
func (f *Future) cancel() {
	if !f.cancelled {
		f.cancelled = true
		f.waitGrp.Done()
	}
}
