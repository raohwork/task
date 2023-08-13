// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import "context"

func f2t(f []func(context.Context) error) []Task {
	args := make([]Task, len(f))
	for idx, t := range f {
		args[idx] = Func(t)
	}
	return args
}

// Iter creates a task run tasks with same context and stops at first error.
func Iter(tasks ...Task) Helper {
	return Func(func(ctx context.Context) error {
		for _, t := range tasks {
			if err := t.Run(ctx); err != nil {
				return err
			}
		}

		return nil
	}).Helper()
}

// IterF is helper to run Iter.
func IterF(tasks ...func(context.Context) error) Helper {
	return Iter(f2t(tasks)...)
}

// Wait creates a task that runs all task concurrently, wait them get done, and
// return first non-nil error.
func Wait(tasks ...Task) Helper {
	return Func(func(ctx context.Context) (err error) {
		ch := make(chan error)
		for _, t := range tasks {
			Helper{t}.GoWithChan(ctx, ch)
		}

		for range tasks {
			e := <-ch
			if err == nil && e != nil {
				err = e
			}
		}

		return
	}).Helper()
}

// WaitF is helper to run Wait.
func WaitF(tasks ...func(context.Context) error) Helper {
	return Wait(f2t(tasks)...)
}

// Skip creates a task that runs tasks concurrently, cancel others if any error, and
// wait them done.
func Skip(tasks ...Task) Helper {
	return Func(func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ch := make(chan error)
		for _, t := range tasks {
			Helper{t}.GoWithChan(ctx, ch)
		}

		for range tasks {
			e := <-ch
			if err == nil && e != nil {
				err = e
				cancel()
			}
		}

		return
	}).Helper()
}

// SkipF is helper to run Skip.
func SkipF(tasks ...func(context.Context) error) Helper {
	return Skip(f2t(tasks)...)
}
