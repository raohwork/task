// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
)

// Loop creates a task that repeatedly runs t with same context until it returns an
// error.
func (t Task) Loop() Task {
	return func(ctx context.Context) (err error) {
		for {
			err = t.Run(ctx)
			if err != nil {
				return
			}
		}
	}
}

// Retry creates a task thats repeatedly runs t with same context until it returns
// nil.
func (t Task) Retry() Task {
	return func(ctx context.Context) (err error) {
		for {
			err = t.Run(ctx)
			if err == nil {
				return
			}
		}
	}
}

// RetryN is like Retry, but retries no more than n times.
//
// In other words, RetryN(2) will run at most 3 times:
//
//   - first try
//   - first retry
//   - second retry
func (t Task) RetryN(n int) Task {
	if n < 0 {
		n = 0
	}
	n++
	return func(ctx context.Context) (err error) {
		for i := 0; i < n; i++ {
			err = t.Run(ctx)
			if err == nil {
				return
			}
		}
		return
	}
}
