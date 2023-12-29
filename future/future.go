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
//
// Created future will only have at most one non-empty value. If it is resolved by
// res, fut.Err() will be nil; fut.Get() returns empty value if rejected.
func New[T any]() (fut *Future[T], res func(T), rej func(error)) {
	ret := &Future[T]{}
	return ret, ret.resolve, ret.reject
}

// Create creates a future fut, which can be resolved by determine.
//
// Created future might have both the value and the error with non-empty value.
func Create[T any]() (fut *Future[T], determine func(T, error)) {
	ret := &Future[T]{}
	return ret, ret.determine
}

// Future represents a value which is determined some time in future.
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

func (f *Future[T]) determine(v T, e error) {
	f.l.Lock()
	if !f.closed {
		f.data = v
		f.err = e
		close(f.prepareDone(true))
		f.closed = true
	}
	f.l.Unlock()
}

func (f *Future[T]) resolve(v T) {
	f.determine(v, nil)
}

func (f *Future[T]) reject(e error) {
	var v T
	f.determine(v, e)
}

// Done returns a channel which is closed when value is determined.
func (f *Future[T]) Done() <-chan struct{} { return f.done() }

// Get waits the value and retrieves it.
func (f *Future[T]) Get() T { <-f.Done(); return f.data }

// Err waits the value and retrieves the error.
func (f *Future[T]) Err() error { <-f.Done(); return f.err }

// Await is a helper which waits the value and retrieves it along with the error.
func (f *Future[T]) Await(ctx context.Context) (t T, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case <-f.Done():
		return f.data, f.err
	}

}
