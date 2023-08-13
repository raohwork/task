// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package future

// Resolve creates an already resolved future.
func Resolve[T any](v T) *Future[T] {
	ret, res, _ := New[T]()
	res(v)
	return ret
}

// Reject creates an already rejected future.
func Reject[T any](e error) *Future[T] {
	ret, _, rej := New[T]()
	rej(e)
	return ret
}
