// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import "context"

// Default wraps d to provide default value whenever it failed.
func (d Data[T]) Default(v T) Data[T] {
	return d.DefaultIf(func(_ error) bool { return true }, v)
}

// DefaultIf wraps d to return v instead if any error matched by errf occurred.
func (d Data[T]) DefaultIf(errf func(error) bool, v T) Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		ret, err = d(ctx)
		if err != nil && errf(err) {
			return v, nil
		}
		return
	}
}

// DefaultIfNot uses v if error occurred and is NOT matched by errf.
func (d Data[T]) DefaultIfNot(errf func(error) bool, v T) Data[T] {
	return d.DefaultIf(func(err error) bool {
		return !errf(err)
	}, v)
}

// Retry wraps d to run it repeatly until success.
func (d Data[T]) Retry() Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		for {
			ret, err = d(ctx)
			if err == nil {
				return
			}
		}
	}
}

// RetryN is like Retry, but no more than n times.
//
// RetryN(3) will run at most 4 times, first attempt is not considered as retrying.
func (d Data[T]) RetryN(n int) Data[T] {
	if n < 0 {
		n = 0
	}
	return func(ctx context.Context) (ret T, err error) {
		for x := 0; x <= n; x++ {
			ret, err = d(ctx)
			if err == nil {
				return
			}
		}
		return
	}
}

// RetryIf wraps d to run it repeatly until success or errf returns false.
//
// Error passed to errf will never be nil.
func (d Data[T]) RetryIf(errf func(error) bool) Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		for {
			ret, err = d(ctx)
			if err == nil || !errf(err) {
				return
			}
		}
	}
}

// RetryNIf is like RetryIf, but no more than n times.
//
// Error passed to errf will never be nil.
func (d Data[T]) RetryNIf(n int, errf func(error) bool) Data[T] {
	if n < 0 {
		n = 0
	}
	return func(ctx context.Context) (ret T, err error) {
		for x := 0; x <= n; x++ {
			ret, err = d(ctx)
			if err == nil || !errf(err) {
				return
			}
		}
		return
	}
}
