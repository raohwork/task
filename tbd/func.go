// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tbd

// Resolve creates an already successfully resolved TBD.
func Resolve[T any](v T) TBD[T] {
	ret, res, _ := New[T]()
	res(v)
	return ret
}

// Reject creates an already failed to be resoved TBD.
func Reject[T any](e error) TBD[T] {
	ret, _, rej := New[T]()
	rej(e)
	return ret
}

// Value waits and gets the value from t.
func Value[T any](t TBD[T]) T {
	v, _ := t.Get(nopCtx)
	return v
}

// Value waits and gets the error from t.
func Err[T any](t TBD[T]) error {
	_, e := t.Get(nopCtx)
	return e
}
