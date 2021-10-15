// Package future provides a type parameter based primitive for running
// function calls in the background.
package future

import (
	"context"
	"sync"
	"time"
)

// A Future contains information about a function running in the
// background and how to fetch its result.
type Future[T interface{}] struct {
	result T
	err    error
	cancel context.CancelFunc
	start  time.Time
	finish time.Time
	lock   sync.RWMutex
}

// New returns a Future from a function that takes a context.Context and
// returns any object and an error. The function is immediately started
// in a Go routine and a future is returned that can be used to get
// the result of the function.
func New[T interface{}](ctx context.Context, fn func(ctx context.Context) (T, error)) *Future[T] {
	ctx, cancel := context.WithCancel(ctx)

	f := Future[T]{
		start:  time.Now(),
		cancel: cancel,
	}

	f.lock.Lock()
	go func() {
		f.result, f.err = fn(ctx)
		f.finish = time.Now()
		cancel()
		f.lock.Unlock()
	}()

	return &f
}

// Wait returns the result of the future once it's done.
func (f *Future[T]) Wait() (T, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.result, f.err
}

// MustWait panics instead of returning an error if there is an error,
// otherwise returns the result from Wait.
func (f *Future[T]) MustWait() T {
	res, err := f.Wait()
	if err != nil {
		panic(err)
	}
	return res
}

// Elapsed returns the time since the start of the function's
// execution.
func (f *Future[T]) Elapsed() time.Duration {
	if !f.Done() {
		return time.Since(f.start)
	}
	return f.finish.Sub(f.start)
}

// Cancel sends a cancel signal to the function's context.Context.
func (f *Future[T]) Cancel() {
	f.cancel()
}

// Done returns true if the function has finished running and false
// otherwise.
func (f *Future[T]) Done() bool {
	return !f.finish.IsZero()
}
