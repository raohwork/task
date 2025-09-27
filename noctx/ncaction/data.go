// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ncaction

import (
	"sync"
	"sync/atomic"

	"github.com/raohwork/task/noctx/nctask"
)

// Data is a function which can generate some data.
//
// Though it is named "data", it's more like a factory or "promise in JS". It is
// used to generate value to be used with [Action] or [Converter]. Take a look at
// their document for more info.
type Data[T any] func() (T, error)

// Use converts f to Data, for type conversion purpose.
func Use[T any](f func() (T, error)) Data[T] { return f }

// NoErrUse converts f into Data and never fail.
func NoErrUse[T any](f func() T) Data[T] {
	return func() (T, error) { return f(), nil }
}

// UseValue creates a Data from fixed value.
func UseValue[T any](v T) Data[T] { return func() (T, error) { return v, nil } }

// UseError creates a Data that always fail.
func UseError[T any](err error) Data[T] {
	return func() (v T, e error) {
		e = err
		return
	}
}

// Get generates a value from Data.
func (d Data[T]) Get() (T, error) { return d() }

// NoErr converts Data into simple function that ignores error.
func (d Data[T]) NoErr() func() T {
	return func() T { v, _ := d(); return v }
}

// Then creates a new Data by converting d with c.
func (d Data[T]) Then(c Converter[T, T]) Data[T] {
	return c.From(d)
}

// Saved wraps d to cache its result only when success.
func (d Data[T]) Saved() Data[T] {
	var (
		lock sync.Mutex
		v    T
		done atomic.Uint32
	)
	return func() (ret T, err error) {
		if done.Load() == 1 {
			return v, nil
		}

		lock.Lock()
		defer lock.Unlock()
		if done.Load() == 1 {
			return v, nil
		}

		ret, err = d()
		if err == nil {
			defer done.Store(1)
			v = ret
		}
		return
	}
}

// Cached wraps d to cache its result, no matter success or failed.
func (d Data[T]) Cached() Data[T] {
	var (
		once sync.Once
		v    T
		err  error
	)
	return func() (T, error) {
		once.Do(func() {
			v, err = d()
		})
		return v, err
	}
}

// Do creates a [nctask.Task] by doing something with its value.
func (d Data[T]) Do(a Action[T]) nctask.Task { return a.Use(d) }
