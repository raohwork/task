// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
)

// Pre wraps d to run f before it.
func (d Data[T]) Pre(f func()) Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		f()
		return d(ctx)
	}
}

// Post wraps d to run f after it.
func (d Data[T]) Post(f func(T, error)) Data[T] {
	return func(ctx context.Context) (T, error) {
		ret, err := d(ctx)
		f(ret, err)
		return ret, err
	}
}

// Defer wraps d to run f after it.
func (d Data[T]) Defer(f func()) Data[T] {
	return func(ctx context.Context) (T, error) {
		ret, err := d(ctx)
		f()
		return ret, err
	}
}
