// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ncaction

// Action2 is an action that accepts two parameters.
type Action2[A, B any] func(A, B) error

// Do2 creates an Action2, mostly for type converting purpose.
func Do2[A, B any](f func(A, B) error) Action2[A, B] { return f }

// NoErrDo2 is like Do2, but the function is not cancellable and never fail.
func NoErrDo2[A, B any](f func(A, B)) Action2[A, B] {
	return func(a A, b B) error { f(a, b); return nil }
}

// Use creates an Action by currifying act with [Data].
func (act Action2[A, B]) Use(a Data[A]) Action[B] {
	return func(vb B) error {
		va, err := a()
		if err != nil {
			return err
		}

		return act(va, vb)
	}
}

// Apply creates an Action by currifying act with a raw value.
func (act Action2[A, B]) Apply(va A) Action[B] {
	return func(vb B) error {
		return act(va, vb)
	}
}

// Then creates an Action2 by running next after act if finished successfully.
func (act Action2[A, B]) Then(next Action2[A, B]) Action2[A, B] {
	return func(va A, vb B) error {
		if err := act(va, vb); err != nil {
			return err
		}
		return next(va, vb)
	}
}

// Action3 is like Action2, but accepts one more param.
type Action3[A, B, C any] func(A, B, C) error

// Do3 creates an Action3, mostly for type converting purpose.
func Do3[A, B, C any](f func(A, B, C) error) Action3[A, B, C] { return f }

// NoErrDo3 is like Do3, but the function is not cancellable and never fail.
func NoErrDo3[A, B, C any](f func(A, B, C)) Action3[A, B, C] {
	return func(a A, b B, c C) error { f(a, b, c); return nil }
}

// Use creates an Action2 by currifying act with [Data].
func (act Action3[A, B, C]) Use(a Data[A]) Action2[B, C] {
	return func(vb B, vc C) error {
		va, err := a()
		if err != nil {
			return err
		}

		return act(va, vb, vc)
	}
}

// Apply creates an Action2 by currifying act with a raw value.
func (act Action3[A, B, C]) Apply(va A) Action2[B, C] {
	return func(vb B, vc C) error {
		return act(va, vb, vc)
	}
}

// Then creates an Action3 by running next after act if finished successfully.
func (act Action3[A, B, C]) Then(next Action3[A, B, C]) Action3[A, B, C] {
	return func(va A, vb B, vc C) error {
		if err := act(va, vb, vc); err != nil {
			return err
		}
		return next(va, vb, vc)
	}
}

// Action4 is like Action3, but accepts one more param.
type Action4[A, B, C, D any] func(A, B, C, D) error

// Do4 creates an Action4, mostly for type converting purpose.
func Do4[A, B, C, D any](f func(A, B, C, D) error) Action4[A, B, C, D] { return f }

// NoErrDo4 is like Do4, but the function is not cancellable.
func NoErrDo4[A, B, C, D any](f func(A, B, C, D)) Action4[A, B, C, D] {
	return func(a A, b B, c C, d D) error {
		f(a, b, c, d)
		return nil
	}
}

// Use creates an Action3 by currifying act with [Data].
func (act Action4[A, B, C, D]) Use(a Data[A]) Action3[B, C, D] {
	return func(vb B, vc C, vd D) error {
		va, err := a()
		if err != nil {
			return err
		}

		return act(va, vb, vc, vd)
	}
}

// Apply creates an Action3 by currifying act with a raw value.
func (act Action4[A, B, C, D]) Apply(va A) Action3[B, C, D] {
	return func(vb B, vc C, vd D) error {
		return act(va, vb, vc, vd)
	}
}

// Then creates an Action4 by running next after act if finished successfully.
func (act Action4[A, B, C, D]) Then(next Action4[A, B, C, D]) Action4[A, B, C, D] {
	return func(va A, vb B, vc C, vd D) error {
		if err := act(va, vb, vc, vd); err != nil {
			return err
		}
		return next(va, vb, vc, vd)
	}
}
