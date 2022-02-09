package utils

import (
	"fmt"
	"sync"
)

// Any is a substitute for interface{}
type Any = interface{}

type Promise struct {
	pending bool

	// A function that is passed with the arguments resolve and reject.
	// The executor function is executed immediately by the Promise implementation,
	// passing resolve and reject functions (the executor is called
	// before the Promise constructor even returns the created object).
	// The resolve and reject functions, when called, resolve or reject
	// the promise, respectively. The executor normally initiates some
	// asynchronous work, and then, once that completes, either calls the
	// resolve function to resolve the promise or else rejects it if
	// an error or panic occurred.
	executor func(resolve func(Any), reject func(error))

	// Stores the result passed to resolve()
	result Any

	// Stores the error passed to reject()
	err error

	// Mutex protects against data race conditions.
	mutex sync.Mutex
	// WaitGroup allows to block until all callbacks are executed.
	wg sync.WaitGroup
}

// New instantiates and returns a pointer to a new Promise.
func New(executor func(resolve func(Any), reject func(error))) *Promise {
	var promise = &Promise{
		pending:  true,
		executor: executor,
		result:   nil,
		err:      nil,
		mutex:    sync.Mutex{},
		wg:       sync.WaitGroup{},
	}

	promise.wg.Add(1)

	go func() {
		defer promise.handlePanic()
		promise.executor(promise.resolve, promise.reject)
	}()

	return promise
}

// Resolve returns a Promise that has been resolved with a given value.
func Resolve(resolution Any) *Promise {
	return New(func(resolve func(Any), reject func(error)) {
		resolve(resolution)
	})
}

type resolutionHelper struct {
	index int
	data  Any
}

// Then appends fulfillment and rejection handlers to the promise,
// and returns a new promise resolving to the return value of the called handler.
func (promise *Promise) Then(fulfillment func(data Any) Any) *Promise {
	return New(func(resolve func(Any), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(err)
			return
		}
		resolve(fulfillment(result))
	})
}

// Catch Appends a rejection handler to the promise,
// and returns a new promise resolving to the return value of the handler.
func (promise *Promise) Catch(rejection func(err error) error) *Promise {
	return New(func(resolve func(Any), reject func(error)) {
		result, err := promise.Await()
		if err != nil {
			reject(rejection(err))
			return
		}
		resolve(result)
	})
}

// Await is a blocking function that waits for all callbacks to be executed.
func (promise *Promise) Await() (Any, error) {
	promise.wg.Wait()
	return promise.result, promise.err
}

// All waits for all Promises to be resolved, or for any to be rejected.
func All(promises ...*Promise) *Promise {
	psLen := len(promises)
	if psLen == 0 {
		return Resolve(make([]Any, 0))
	}

	return New(func(resolve func(Any), reject func(error)) {
		resolutionsChan := make(chan resolutionHelper, psLen)
		errorChan := make(chan error, psLen)

		for index, promise := range promises {
			func(i int) {
				promise.Then(func(data Any) Any {
					resolutionsChan <- resolutionHelper{i, data}
					return data
				}).Catch(func(err error) error {
					errorChan <- err
					return err
				})
			}(index)
		}

		resolutions := make([]Any, psLen)
		for x := 0; x < psLen; x++ {
			select {
			case resolution := <-resolutionsChan:
				resolutions[resolution.index] = resolution.data

			case err := <-errorChan:
				reject(err)
				return
			}
		}
		resolve(resolutions)
	})
}

// Reject returns a Promise that has been rejected with a given error.
func Reject(err error) *Promise {
	return New(func(resolve func(Any), reject func(error)) {
		reject(err)
	})
}

func (promise *Promise) resolve(resolution Any) {
	promise.mutex.Lock()

	if !promise.pending {
		promise.mutex.Unlock()
		return
	}

	switch result := resolution.(type) {
	case *Promise:
		flattenedResult, err := result.Await()
		if err != nil {
			promise.mutex.Unlock()
			promise.reject(err)
			return
		}
		promise.result = flattenedResult
	default:
		promise.result = result
	}
	promise.pending = false

	promise.wg.Done()
	promise.mutex.Unlock()
}

func (promise *Promise) reject(err error) {
	promise.mutex.Lock()
	defer promise.mutex.Unlock()

	if !promise.pending {
		return
	}

	promise.err = err
	promise.pending = false

	promise.wg.Done()
}

func (promise *Promise) handlePanic() {
	e := recover()
	if e != nil {
		switch err := e.(type) {
		case nil:
			promise.reject(fmt.Errorf("panic recovery with nil error"))
		case error:
			promise.reject(fmt.Errorf("panic recovery with error: %s", err.Error()))
		default:
			promise.reject(fmt.Errorf("panic recovery with unknown error: %s", fmt.Sprint(err)))
		}
	}
}
