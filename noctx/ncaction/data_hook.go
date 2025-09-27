// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ncaction

// Pre wraps d to run f before it.
func (d Data[T]) Pre(f func()) Data[T] {
	return func() (ret T, err error) {
		f()
		return d()
	}
}

// Post wraps d to run f after it.
func (d Data[T]) Post(f func(T, error)) Data[T] {
	return func() (T, error) {
		ret, err := d()
		f(ret, err)
		return ret, err
	}
}

// AlterOutput wraps d to run f on its output.
func (d Data[T]) AlterOutput(f func(T, error) (T, error)) Data[T] {
	return func() (T, error) {
		return f(d())
	}
}

// AlterError wraps d to run f on its error.
func (d Data[T]) AlterError(f func(error) error) Data[T] {
	return func() (T, error) {
		ret, err := d()
		return ret, f(err)
	}
}

// Defer wraps d to run f after it.
func (d Data[T]) Defer(f func()) Data[T] {
	return func() (T, error) {
		ret, err := d()
		f()
		return ret, err
	}
}
