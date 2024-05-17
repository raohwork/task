// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"sync"
)

// Future creates a cached Data, whose result is determined by a function. Getting
// value of the Data will be blocked until the result is determined.
func Future[T any]() (ret Data[T], determine func(T, error)) {
	var (
		once sync.Once
		v    T
		err  error
	)
	done := make(chan struct{})
	return func(ctx context.Context) (ret T, x error) {
			select {
			case <-done:
				return v, err
			case <-ctx.Done():
				return ret, ctx.Err()
			}
		}, func(data T, e error) {
			once.Do(func() { v, err = data, e; close(done) })
		}
}

// TBD is identical to [Future] but provides different type of function.
func TBD[T any]() (ret Data[T], resolve func(T), reject func(error)) {
	ret, f := Future[T]()
	return ret,
		func(v T) { f(v, nil) },
		func(e error) { var v T; f(v, e) }
}
