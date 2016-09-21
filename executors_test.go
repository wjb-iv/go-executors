package executors_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	executors "github.com/wjb-iv/go-executors"
)

type TestCallable struct {
	name  string
	delay int
}

func (tc TestCallable) Call() interface{} {
	time.Sleep(time.Duration(tc.delay) * time.Millisecond)
	return tc.name + " Results"
}

func TestBasic(t *testing.T) {
	callable1 := TestCallable{"Job1", 1000}
	callable2 := TestCallable{"Job2", 1000}
	callable3 := TestCallable{"Job3", 1000}

	exec := executors.New("my-go-pool", 2, 50)
	defer func() {
		exec.Close()
		time.Sleep(time.Millisecond * 100)
	}()

	f1 := exec.Invoke(callable1)
	f2 := exec.Invoke(callable2)
	f3 := exec.Invoke(callable3)

	fmt.Println(time.Now())
	if _, err := f1.Get(time.Millisecond * 100); err == nil {
		t.Error("Timeout error should have been received.")
	} else {
		log.Printf("[TestBasic] %v ", err)
	}
	if res, err := f2.Get(time.Second * 5); err == nil {
		log.Printf("[TestBasic] %s - %s ", time.Now(), res)
	} else {
		t.Error(err)
	}
	if res, err := f3.Get(time.Second * 5); err == nil {
		log.Printf("[TestBasic] %s - %s ", time.Now(), res)
	} else {
		t.Error(err)
	}
}

type TestObject struct {
	Val1 int
	Val2 int
}

type TestCallableType struct {
	testObject *TestObject
}

func (tc TestCallableType) Call() interface{} {
	return tc.testObject.Val1 + tc.testObject.Val2
}

func TestTypes(t *testing.T) {
	to := TestObject{3, 5}
	callable1 := TestCallableType{&to}

	exec := executors.New("my_go_pool", 2, 50)
	defer func() {
		exec.Close()
		time.Sleep(time.Millisecond * 100)
	}()

	f1 := exec.Invoke(callable1)

	if res, err := f1.Get(time.Millisecond * 100); err == nil {

		if resultVal, ok := res.(int); ok {
			log.Printf("[TestTypes] %s - %d ", time.Now(), resultVal)
			if resultVal != 8 {
				t.Error("Expected result is '8'")
			}
		} else {
			t.Error("Expected type is 'int'")
		}

	} else {
		log.Printf("[TestTypes] %v ", err)
		t.Error(err)
	}
}

func TestMultiplePools(t *testing.T) {
	exec1 := executors.New("pool-one", 2, 50)
	exec2 := executors.New("pool-two", 2, 50)
	defer func() {
		exec1.Close()
		exec2.Close()
		time.Sleep(time.Millisecond * 100)
	}()

	to1 := TestObject{3, 5}
	callable1 := TestCallableType{&to1}
	to2 := TestObject{6, 5}
	callable2 := TestCallableType{&to2}

	f1 := exec1.Invoke(callable1)
	f2 := exec2.Invoke(callable2)

	if res, err := f1.Get(time.Millisecond * 100); err == nil {
		log.Printf("[TestMultiplePools] %s - %d ", time.Now(), res)
	}

	if res, err := f2.Get(time.Millisecond * 100); err == nil {
		log.Printf("[TestMultiplePools] %s - %d ", time.Now(), res)
	}
}

type CallablePanic struct {
}

func (cp CallablePanic) Call() interface{} {
	log.Println("About to panic...")
	time.Sleep(time.Millisecond * 200)
	panic("Test error")
}

func TestCallablePanic(t *testing.T) {
	exec := executors.New("pool-of-uno", 1, 10)
	defer func() {
		exec.Close()
		time.Sleep(time.Millisecond * 100)
	}()

	callable := CallablePanic{}
	f := exec.Invoke(callable)
	if res, err := f.Get(time.Second); err == nil {
		log.Printf("[TestCallablePanic] %s - %d ", time.Now(), res)
		if res != nil {
			t.Error("Result is expected to be 'nil' if callable panics!")
		}
	}

	callable2 := TestCallable{"Job1", 100}
	f2 := exec.Invoke(callable2)
	if res, err := f2.Get(time.Second); err == nil {
		log.Printf("[TestCallablePanic] %s - %d ", time.Now(), res)
		if res != "Job1 Results" {
			t.Error("Unexpected return value")
		}
	}
}
