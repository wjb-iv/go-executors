package executors

import (
	"fmt"
	"log"
	"sync"
)

// ExecutorService - the obect we hold to perform invocations
type ExecutorService struct {
	input chan call
}

// The "job" that will be submitted internally
type call struct {
	callable Callable
	future   *Future
}

// Invoke actually enqueues the callable to be executed and returns a future
func (exec *ExecutorService) Invoke(c Callable) *Future {
	var wg sync.WaitGroup
	wg.Add(1)
	var ft = Future{waitGrp: &wg}
	var cl = call{c, &ft}
	exec.input <- cl

	return &ft
}

// Close will terminate the input channel and cause workers to exit
func (exec *ExecutorService) Close() {
	close(exec.input)
}

// New creates a new executor service
func New(poolName string, numWorkers int, queueSize int) ExecutorService {
	work := make(chan call, queueSize)
	// This starts up the workers, initially blocked
	// because there is no work yet.
	for w := 1; w <= int(numWorkers); w++ {
		go worker(fmt.Sprintf("%s-worker_%d", poolName, w), work)
	}
	return ExecutorService{work}
}

// InternalError is a catch-all exception for any unhandled error
// arising from execution of a callable's 'Call()' method
type InternalError struct {
	reason string
}

func (ie InternalError) Error() string {
	return fmt.Sprintf("Callable.Call() failed due to: %s", ie.reason)
}

// Internal worker function uses 'safelyCall' to trap panics that
// could arise from a bad 'Callable' implementation
func worker(id string, work <-chan call) {
	myid := id
	for c := range work {
		log.Printf("[%s] processing callable", myid)
		if val, err := safelyCall(c.callable); err == nil {
			c.future.setVal(val)
		} else {
			// Setting the futures value to any error will cause
			// the 'Get()' method to return that value in err,
			// thus relaying info about the unhandled error in
			// the passed Callable.
			c.future.setVal(err)
		}
	}
	log.Printf("[%s] Shutting down...", myid)
}

func safelyCall(callable Callable) (result interface{}, err error) {
	defer func() {
		if code := recover(); code != nil {
			err = InternalError{fmt.Sprintf("%v", code)}
		}
	}()
	result = callable.Call()
	return result, err
}
