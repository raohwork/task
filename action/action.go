// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"

	"github.com/raohwork/task"
)

// Action is a specialized [task.Task] which accepts one param.
type Action[T any] func(context.Context, T) error

// Do wraps f into Action, mostly for for typing purpose.
func Do[T any](f func(context.Context, T) error) Action[T] { return f }

// NoCtxDo creates Action from non-cancellable function.
func NoCtxDo[T any](f func(T) error) Action[T] {
	return func(_ context.Context, v T) error { return f(v) }
}

// NoErrDo creates Action from non-cancellable, never-fail function.
func NoErrDo[T any](f func(T)) Action[T] {
	return func(_ context.Context, v T) error { f(v); return nil }
}

// Then creates an Action which runs next after act is finished successfully.
func (act Action[T]) Then(next Action[T]) Action[T] {
	return func(ctx context.Context, v T) error {
		if err := act(ctx, v); err != nil {
			return err
		}
		return next(ctx, v)
	}
}

// NoCtx converts Action into simple function that ignores the context.
func (a Action[T]) NoCtx() func(T) error {
	return func(v T) error { return a(context.TODO(), v) }
}

// NoErr converts Action into simple function that ignores context and error.
func (a Action[T]) NoErr() func(T) { return func(v T) { a(context.TODO(), v) } }

// Use creates a [task.Task] that executes the action with Data. The value of data
// is generated on-the-fly.
func (a Action[T]) Use(d Data[T]) task.Task {
	return func(ctx context.Context) error {
		v, err := d(ctx)
		if err != nil {
			return err
		}
		return a(ctx, v)
	}
}

// Apply creates a [task.Task] that executes the action with a value.
func (a Action[T]) Apply(v T) task.Task {
	return func(ctx context.Context) error { return a(ctx, v) }
}

// With wraps a to modify context before run it.
func (a Action[T]) With(mod task.CtxMod) Action[T] {
	return func(ctx context.Context, v T) error {
		ctx, cancel := mod(ctx)
		defer cancel()
		return a(ctx, v)
	}
}

// Pre wraps a to run f before it.
func (a Action[T]) Pre(f func(v T)) Action[T] {
	return func(ctx context.Context, v T) error {
		f(v)
		return a(ctx, v)
	}
}

// Post wraps a to run f after it.
func (a Action[T]) Post(f func(T, error)) Action[T] {
	return func(ctx context.Context, v T) error {
		err := a(ctx, v)
		f(v, err)
		return err
	}
}

// Defer wraps a to run f after it.
func (a Action[T]) Defer(f func()) Action[T] {
	return func(ctx context.Context, v T) error {
		err := a(ctx, v)
		f()
		return err
	}
}

// Task is a predefined Action used to convert Data into [task.Task] like
//
//	err := Get(saveToDB).By(dbConn).From(myData).To(Task).Run(ctx)
func Task[T any](_ context.Context, _ T) error { return nil }