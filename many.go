// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import "context"

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

// First creates a task that runs tasks concurrently, return first result and cancel
// others.
func First(tasks ...Task) Task {
	return func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancel(ctx)

		ch := make(chan error)
		for _, t := range tasks {
			t.GoWithChan(ctx, ch)
		}

		err = <-ch
		cancel()
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

// Skip creates a task that runs tasks concurrently, cancel others if any error, and
// wait them done.
func Skip(tasks ...Task) Task {
	return func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ch := make(chan error)
		for _, t := range tasks {
			t.GoWithChan(ctx, ch)
		}

		for range tasks {
			e := <-ch
			if err == nil && e != nil {
				err = e
				cancel()
			}
		}

		return
	}
}
