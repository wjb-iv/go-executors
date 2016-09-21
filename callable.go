package executors

// Callable - the interface to implement to do actual work
type Callable interface {
	Call() interface{}
}
