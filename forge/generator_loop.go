// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import "context"

// Retry creates a Generator thats repeatedly runs g with same context until it
// generates a value without error.
//
// For generators created using [Chain], you should take extra care.
func (g Generator[T]) Retry() Generator[T] {
	return func(ctx context.Context) (T, error) {
		for {
			v, err := g.Run(ctx)
			if err == nil {
				return v, err
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
//
// For generators created using [Chain], you should take extra care.
func (g Generator[T]) RetryN(n int) Generator[T] {
	if n < 0 {
		n = 0
	}
	n++
	return func(ctx context.Context) (ret T, err error) {
		for i := 0; i < n; i++ {
			ret, err = g.Run(ctx)
			if err == nil {
				return
			}
		}
		return
	}
}
