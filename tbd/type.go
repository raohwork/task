// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tbd

import (
	"context"
	"sync"
)

// TBD represents a value might be computed some time later.
//
// It is a placeholder for the value. There must be a resolver who give it the
// result, no matter a value or an error. Once the result is provided, the channel
// returned by Resolved() is closed.
type TBD[T any] interface {
	// Indicates if it is resolved.
	Resolved() <-chan struct{}
	// Wait and get the value. Can be used as a generator.
	Get(context.Context) (T, error)
}

// New creates a TBD and provides 2 functions to resolve it.
func New[T any]() (TBD[T], func(T) error, func(error) error) {
	ret := newbase[T]()
	return ret, ret.resolve, ret.reject
}

// Create creates a TBD and provides a function to resolve it.
func Create[T any]() (TBD[T], func(T, error) error) {
	ret := newbase[T]()
	return ret, ret.determine
}

func newbase[T any]() *_basetbd[T] {
	return &_basetbd[T]{ch: make(chan struct{})}
}

type _basetbd[T any] struct {
	once sync.Once
	ch   chan struct{}
	data T
	err  error
}

func (f *_basetbd[T]) determine(v T, e error) error {
	f.once.Do(func() {
		f.data = v
		f.err = e
		close(f.ch)
	})
	return e
}

func (f *_basetbd[T]) resolve(v T) error {
	f.once.Do(func() {
		f.data = v
		close(f.ch)
	})
	return nil
}

func (f *_basetbd[T]) reject(e error) error {
	f.once.Do(func() {
		f.err = e
		close(f.ch)
	})
	return e
}

func (f *_basetbd[T]) Resolved() <-chan struct{} { return f.ch }
func (f *_basetbd[T]) Get(ctx context.Context) (t T, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case <-f.Resolved():
		return f.data, f.err
	}

}
