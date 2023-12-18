// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"time"
)

// Tiny wraps a non-cancellable function into task.
func Tiny(f func() error) Task {
	return func(_ context.Context) error { return f() }
}

// Micro wraps a never-fail, non-cancellable function into task.
func Micro(f func()) Task {
	return func(_ context.Context) error { f(); return nil }
}

// Sleep is a cancellable [time.Sleep] in task form.
func Sleep(timeout time.Duration) Task {
	return func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(timeout):
			return nil
		}
	}
}

// FromServer creates a task from something can be started or stopped. Running
// the task calls start, and cancelling context calls stop.
func FromServer(start func() error, stop func()) Task {
	return func(ctx context.Context) error {
		done := make(chan struct{})
		defer close(done)
		go func() {
			select {
			case <-done:
				return
			case <-ctx.Done():
			}

			stop()
		}()

		return start()
	}
}
