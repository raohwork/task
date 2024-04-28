// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"errors"
)

// Iter creates a task run tasks with same context and stops at first error.
func Iter(tasks ...Task) Task {
	return func(ctx context.Context) error {
		for _, t := range tasks {
			if err := t.Run(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}

var ErrOneHasDone = errors.New("another task has been done")

// First creates a task that runs tasks concurrently, return first result and cancel
// others. Other tasks receives ErrOneHasDone as cancel cause.
//
// Take care of [Tiny] tasks as it cannot be cancelled by context.
func First(tasks ...Task) Task {
	return func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancelCause(ctx)

		ch := make(chan error)
		for _, t := range tasks {
			t.GoWithChan(ctx, ch)
		}

		err = <-ch
		cancel(ErrOneHasDone)
		go func() {
			for i := 1; i < len(tasks); i++ {
				<-ch
			}
		}()
		return

	}
}

// Wait creates a task that runs all task concurrently, wait them get done, and
// return first non-nil error.
func Wait(tasks ...Task) Task {
	return func(ctx context.Context) (err error) {
		ch := make(chan error)
		for _, t := range tasks {
			t.GoWithChan(ctx, ch)
		}

		for range tasks {
			e := <-ch
			if err == nil && e != nil {
				err = e
			}
		}

		return
	}
}

type ErrOthers struct {
	cause error
}

func (e ErrOthers) Error() string {
	return "canceled by error from other task: " + e.cause.Error()
}

func (e ErrOthers) Unwrap() error         { return e.cause }
func (e ErrOthers) Is(err error) bool     { return errors.Is(err, e.cause) }
func (e ErrOthers) As(v interface{}) bool { return errors.As(e.cause, v) }

// Skip creates a task that runs tasks concurrently, cancel others if any error, and
// wait them done.
//
// Tasks canceled by this receieves ErrOthers which wraps the error as cancel cause.
//
// Take care of [Tiny] and [Micro] tasks as it cannot be cancelled by context.
func Skip(tasks ...Task) Task {
	return func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancelCause(ctx)
		defer cancel(nil)

		ch := make(chan error)
		for _, t := range tasks {
			t.GoWithChan(ctx, ch)
		}

		for range tasks {
			e := <-ch
			if err == nil && e != nil {
				err = e
				cancel(ErrOthers{e})
			}
		}

		return
	}
}
