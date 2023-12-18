// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package task provides some helper to work with simple routine.
//
// Functions and methods which requires a task accepts raw [Task] type for best
// compatibility. Those create task returns [Helper] instead so it's easier to use.
package task

import (
	"context"
	"errors"
	"time"
)

// Task repeasents a (maybe) cancellable routine.
type Task func(context.Context) error

// Run runs the task, equals to t(ctx).
func (t Task) Run(ctx context.Context) error {
	return t(ctx)
}

// Go runs t in separated goroutine and returns a channel to retrieve error.
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

// WithTimeout creates a CtxMod which adds timeout info to a context.
func WithTimeout(dur time.Duration) CtxMod {
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
// It will change the error returned by Run. f is called only if t.Run returns an
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

// HandleErrWithContext is like HandleErr, but uses same context in f.
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
	return t.HandleErr(func(err error) error {
		if errors.Is(err, context.Canceled) {
			return err
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		return nil
	})
}
