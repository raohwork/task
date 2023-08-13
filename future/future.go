// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package future

import (
	"context"
	"sync"
	"sync/atomic"
)

// New creates a future fut, which can be resolved by res or rejected by rej.
func New[T any]() (fut *Future[T], res func(T), rej func(error)) {
	ret := &Future[T]{}
	return ret, ret.resolve, ret.reject
}

// Future represents a value will be resolved or rejected some time in future.
// The term "resolve" indicates the value is computed successfully.
//
// Future implements [task.Task] so it's easier to wait for multiple futures.
//
//	err = task.Skip(future1, future2, future3).Run(ctx)
//	if err != nil {
//		// one or more futures are rejected
//	}
type Future[T any] struct {
	l      sync.Mutex
	ch     atomic.Value
	closed bool
	data   T
	err    error
}

func (f *Future[T]) prepareDone(locked bool) chan struct{} {
	if !locked {
		f.l.Lock()
		defer f.l.Unlock()
	}

	done, ok := f.ch.Load().(chan struct{})
	if !ok || done == nil {
		done = make(chan struct{})
		f.ch.Store(done)
	}

	return done
}

func (f *Future[T]) done() chan struct{} {
	done, ok := f.ch.Load().(chan struct{})
	if !ok || done == nil {
		return f.prepareDone(false)
	}
	return done
}

func (f *Future[T]) resolve(v T) {
	f.l.Lock()
	if !f.closed {
		f.data = v
		close(f.prepareDone(true))
		f.closed = true
	}
	f.l.Unlock()
}

func (f *Future[T]) reject(e error) {
	f.l.Lock()
	if !f.closed {
		f.err = e
		close(f.prepareDone(true))
		f.closed = true
	}
	f.l.Unlock()
}

// Done returns a channel which is closed when value is determined.
func (f *Future[T]) Done() <-chan struct{} { return f.done() }

// Get waits the value and retrieves it.
func (f *Future[T]) Get() T { <-f.Done(); return f.data }

// Err waits the value and retrieves the error.
func (f *Future[T]) Err() error { <-f.Done(); return f.err }

// Run waits the value to be resolved or rejected. It also implements [task.Task].
func (f *Future[T]) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.Done():
		return f.Err()
	}
}

// Await is a helper which waits the value and retrieves it along with the error.
func (f *Future[T]) Await(ctx context.Context) (T, error) {
	f.Run(ctx)
	return f.data, f.err
}
