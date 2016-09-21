# Project: go-executors
Another learning project in Go, this one to help me grasp Go's various features for concurrency and synchronization.

# Why Executors?
Since I'm coming from Java, I was wondering how one might implement a standing thread pool, send it work as an implementation of a 'Callable' interface, and collect results with something like a 'Future'.

Of course, Go probably needs Java-like executors like it's creators need a hole in the head, but it turns out that trying to emulate the pattern taught me a lot about Go's concurrency mechanisms. In addition, I was inspired by this blog post:

http://www.jtolds.com/writing/2016/03/go-channels-are-bad-and-you-should-feel-bad/

To try and capture this behavior while minimizing the use of channels. It turns out that I still used several, it's hard to ignore one of Go's most useful features! In the end, it is probably easier and more idiomatic to perform such tasks using Go's own features, but as a learning example, this project was quite useful.

## Usage

Create an executor service with 3 workers and a 25 item queue. Ensure that once your done with the service you close it (to close channels and shut down worker go routines):
```go
exec := executors.New("my-go-pool", 3, 25)
defer exec.Close()
```
Implement the Callable interface:
```go
type TestCallable struct {
	name  string
	delay int
}

func (tc TestCallable) Call() interface{} {
	time.Sleep(time.Duration(tc.delay) * time.Millisecond)
	return tc.name + " Results"
}
```
Create an instance of your callable, invoke it, collect your future, and get the result:
```go
callable := TestCallable{"Job1", 1000}
future := exec.Invoke(callable)
if res, err := future.Get(time.Second); err == nil {
    // Oh Go, can we not have generics someday?
    if resStr, ok := res.(string); ok {
        fmt.Println(resStr)
    }
}
```
