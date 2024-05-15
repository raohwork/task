// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package task provides some helper to work with common routines so that it can be
// cancellable or repeatable.
package task

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Task repeasents a (maybe) cancellable routine.
type Task func(context.Context) error

// Exec runs the task with empty context ([context.Background]).
func (t Task) Exec() error {
	return t(context.Background())
}

// Run runs the task, equals to t(ctx).
func (t Task) Run(ctx context.Context) error {
	return t(ctx)
}

// NoCtx converts the task into a simple function by feeding empty context when run.
func (t Task) NoCtx() func() error {
	return func() error { return t(context.Background()) }
}

// NoErr converts the task into a simple function by feeding empty context when run.
func (t Task) NoErr() func() {
	return func() { t(context.Background()) }
}

// Go runs t in separated goroutine and returns a channel to retrieve error.
//
// It's safe to ignore the channel if you don't need the result.
func (t Task) Go(ctx context.Context) <-chan error {
	ret := make(chan error, 1)
	t.GoWithChan(ctx, ret)
	return ret
}

// GoWithChan runs t in separated goroutine and sends returned error into ch.
func (t Task) GoWithChan(ctx context.Context, ch chan<- error) {
	go func() { ch <- t.Run(ctx) }()
}

// CtxMod defines how you modify a context.
type CtxMod func(context.Context) (context.Context, func())

// Timeout creates a CtxMod which adds timeout info to a context.
func Timeout(dur time.Duration) CtxMod {
	return func(ctx context.Context) (context.Context, func()) {
		return context.WithTimeout(ctx, dur)
	}
}

// With creates a task that the context is derived using modder before running t.
func (t Task) With(modder CtxMod) Task {
	return func(ctx context.Context) error {
		x, c := modder(ctx)
		defer c()
		return t.Run(x)
	}
}

// HandleErr creates a task that handles specific error after running t.
// It could change the error returned by Run. f is called only if t.Run returns an
// error.
func (t Task) HandleErr(f func(error) error) Task {
	return func(ctx context.Context) error {
		err := t.Run(ctx)
		if err != nil {
			err = f(err)
		}

		return err
	}
}

// HandleErrWithContext is like HandleErr, but uses same context used in f.
func (t Task) HandleErrWithContext(f func(context.Context, error) error) Task {
	return func(ctx context.Context) error {
		err := t.Run(ctx)
		if err != nil {
			err = f(ctx, err)
		}

		return err
	}
}

// IgnoreErr ignores the error returned by t.Run if it is not context error.
//
// Context error means [context.Canceled] and [context.DeadlineExceeded].
func (t Task) IgnoreErr() Task {
	return t.OnlyErrs(context.Canceled, context.DeadlineExceeded)
}

// IgnoreErrs creates a function to be used in [Task.IgnoreErrs] that ignores
// specific errors.
//
// Errors are compared using [errors.Is].
func IgnoreErrs(errorList ...error) func(error) error {
	return func(err error) error {
		for _, e := range errorList {
			if errors.Is(err, e) {
				return nil
			}
		}

		return err
	}
}

// IgnoreErrs ignores specific error.
func (t Task) IgnoreErrs(errorList ...error) Task {
	return t.HandleErr(IgnoreErrs(errorList...))
}

// OnlyErrs creates a function to be used in [Task.OnlyErrs] that ignores all
// errors excepts specified in errorList.
//
// Errors are compared using [errors.Is].
func OnlyErrs(errorList ...error) func(error) error {
	return func(err error) error {
		for _, e := range errorList {
			if errors.Is(err, e) {
				return err
			}
		}

		return nil
	}
}

// OnlyErrs preserves only specific error.
func (t Task) OnlyErrs(errorList ...error) Task {
	return t.HandleErr(OnlyErrs(errorList...))
}

var ErrOnce = errors.New("the task can only be executed once.")

// Once creates a task that can be run only once, further attempt returns ErrOnce.
func (t Task) Once() Task {
	var once sync.Once
	return func(ctx context.Context) (err error) {
		err = ErrOnce
		once.Do(func() {
			err = t(ctx)
		})
		return
	}
}

// Cached wraps t to cache the result, and reuse it in later call.
func (t Task) Cached() Task {
	var (
		once sync.Once
		err  error
	)
	return func(ctx context.Context) error {
		once.Do(func() {
			err = t(ctx)
		})
		return err
	}
}

// Defer wraps t to run f after it.
func (t Task) Defer(f func()) Task {
	return func(ctx context.Context) (err error) {
		err = t(ctx)
		f()
		return err
	}
}

// Pre wraps t to run f before it.
func (t Task) Pre(f func()) Task {
	return func(ctx context.Context) (err error) {
		f()
		err = t(ctx)
		return
	}
}

// Post wraps t to run f after it.
func (t Task) Post(f func(error)) Task {
	return func(ctx context.Context) (err error) {
		err = t(ctx)
		f(err)
		return
	}
}
