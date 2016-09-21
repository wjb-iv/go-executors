package executors

import (
	"sync"
	"time"
)

// Future - represents the results of some work to perform in a worker
type Future struct {
	waitGrp *sync.WaitGroup // internal use (to wait gracefully)
	value   interface{}     // the actual return value
}

// TimeoutError happens if the timeout expires on 'Get' method
type TimeoutError struct {
}

func (te TimeoutError) Error() string {
	return "Operation timed out."
}

// Get returns the value or times out after 'timeout' duration
func (f *Future) Get(timeout time.Duration) (val interface{}, err error) {
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
		if ex, ok := f.value.(error); ok {
			val = nil
			err = ex
		} else {
			val = f.value
			err = nil
		}
	case <-time.After(timeout):
		val = nil
		err = TimeoutError{}
	}
	return val, err
}

// Private setter called by worker to signal completion
func (f *Future) setVal(val interface{}) {
	f.value = val
	f.waitGrp.Done()
}
