// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/raohwork/task"
)

// Data is a function which can generate some data.
type Data[T any] func(context.Context) (T, error)

// Use converts f to Data, for type conversion purpose.
func Use[T any](f func(context.Context) (T, error)) Data[T] { return f }

// NoCtxUse converts f into Data, so context is ignored.
func NoCtxUse[T any](f func() (T, error)) Data[T] {
	return func(_ context.Context) (T, error) { return f() }
}

// NoErrUse converts f into Data, so context is ignored and never fail.
func NoErrUse[T any](f func() T) Data[T] {
	return func(_ context.Context) (T, error) { return f(), nil }
}

// UseValue creates a Data from fixed value.
func UseValue[T any](v T) Data[T] { return NoErrUse(func() T { return v }) }

// UseError creates a Data that always fail.
func UseError[T any](err error) Data[T] {
	return NoCtxUse(func() (v T, e error) {
		e = err
		return
	})
}

// Get generates a value from Data.
func (d Data[T]) Get(ctx context.Context) (T, error) { return d(ctx) }

// NoCtx converts Data into simple function that ignores context.
func (d Data[T]) NoCtx() func() (T, error) {
	return func() (T, error) { return d(context.TODO()) }
}

// NoErr converts Data into simple function that ignores context and error.
func (d Data[T]) NoErr() func() T {
	return func() T { v, _ := d(context.TODO()); return v }
}

// With wraps d to adjust its context.
func (d Data[T]) With(mod task.CtxMod) Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		ctx, cancel := mod(ctx)
		defer cancel()
		return d(ctx)
	}
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
	return func(ctx context.Context) (ret T, err error) {
		if done.Load() == 1 {
			return v, nil
		}

		lock.Lock()
		defer lock.Unlock()
		if done.Load() == 1 {
			return v, nil
		}

		ret, err = d(ctx)
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
	return func(ctx context.Context) (T, error) {
		once.Do(func() {
			v, err = d(ctx)
		})
		return v, err
	}
}

// To creates a [task.Task] by doing something with its value.
func (d Data[T]) To(a Action[T]) task.Task { return a.Use(d) }
