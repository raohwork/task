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
