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
type Task interface {
	// Run runs the task.
	Run(context.Context) error
}

// Func converts specific function into task.
type Func func(context.Context) error

// Run runs the task, equals to t(ctx).
func (t Func) Run(ctx context.Context) error {
	return t(ctx)
}

// Helper creates a Helper.
func (t Func) Helper() Helper { return Helper{t} }

// F is shortcut to Func(f).Helper()
func F(f func(context.Context) error) Helper { return Func(f).Helper() }

// T wraps t into Helper.
func T(t Task) Helper { return Helper{t} }

// Helper provides some useful tools.
type Helper struct{ Task }

// Go runs t in separated goroutine and returns a channel to retrieve error.
func (t Helper) Go(ctx context.Context) <-chan error {
	ret := make(chan error, 1)
	t.GoWithChan(ctx, ret)
	return ret
}

// GoWithChan runs t in separated goroutine and sends returned error into ch.
func (t Helper) GoWithChan(ctx context.Context, ch chan<- error) {
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
func (t Helper) With(modder CtxMod) Helper {
	return Func(func(ctx context.Context) error {
		x, c := modder(ctx)
		defer c()
		return t.Run(x)
	}).Helper()
}

// HandleErr creates a task that handles specific error after running t.
// It will change the error returned by Run. f is called only if t.Run returns an
// error.
func (t Helper) HandleErr(f func(error) error) Helper {
	return F(func(ctx context.Context) error {
		err := t.Run(ctx)
		if err != nil {
			err = f(err)
		}

		return err
	})
}

// IgnoreErr ignores the error returned by t.Run if it is not context error.
//
// Context error means [context.Canceled] and [context.DeadlineExceeded].
func (t Helper) IgnoreErr() Helper {
	return F(func(ctx context.Context) error {
		err := t.Run(ctx)
		if errors.Is(err, context.Canceled) {
			return err
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		return nil
	})
}
