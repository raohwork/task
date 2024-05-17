// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import "context"

// Converter2 is an Converter that accepts two input.
type Converter2[A, B, O any] func(context.Context, A, B) (O, error)

// Get2 creates an Converter2, mostly for type converting purpose.
func Get2[A, B, O any](f func(context.Context, A, B) (O, error)) Converter2[A, B, O] { return f }

// From creates a Converter by currifying c with a [Data].
func (c Converter2[A, B, O]) From(a Data[A]) Converter[B, O] {
	return func(ctx context.Context, vb B) (ret O, err error) {
		va, err := a(ctx)
		if err != nil {
			return
		}

		return c(ctx, va, vb)
	}
}

// By creates a Converter by currifying c with a value.
func (c Converter2[A, B, O]) By(va A) Converter[B, O] {
	return func(ctx context.Context, vb B) (ret O, err error) {
		return c(ctx, va, vb)
	}
}

// Converter3 is an Converter2 with additional input.
type Converter3[A, B, C, O any] func(context.Context, A, B, C) (O, error)

// Get3 creates an Converter3, mostly for type converting purpose.
func Get3[A, B, C, O any](f func(context.Context, A, B, C) (O, error)) Converter3[A, B, C, O] {
	return f
}

// From creates a Converter2 by currifying c with a [Data].
func (c Converter3[A, B, C, O]) From(a Data[A]) Converter2[B, C, O] {
	return func(ctx context.Context, vb B, vc C) (ret O, err error) {
		va, err := a(ctx)
		if err != nil {
			return
		}

		return c(ctx, va, vb, vc)
	}
}

// By creates a Converter2 by currifying c with a value.
func (c Converter3[A, B, C, O]) By(va A) Converter2[B, C, O] {
	return func(ctx context.Context, vb B, vc C) (ret O, err error) {
		return c(ctx, va, vb, vc)
	}
}

// Converter4 is an Converter3 with additional input.
type Converter4[A, B, C, D, O any] func(context.Context, A, B, C, D) (O, error)

// Get4 creates an Converter4, mostly for type converting purpose.
func Get4[A, B, C, D, O any](f func(context.Context, A, B, C, D) (O, error)) Converter4[A, B, C, D, O] {
	return f
}

// From creates a Converter3 by currifying c with a [Data].
func (c Converter4[A, B, C, D, O]) From(a Data[A]) Converter3[B, C, D, O] {
	return func(ctx context.Context, vb B, vc C, vd D) (ret O, err error) {
		va, err := a(ctx)
		if err != nil {
			return
		}

		return c(ctx, va, vb, vc, vd)
	}
}

// By creates a Converter3 by currifying c with a value.
func (c Converter4[A, B, C, D, O]) By(va A) Converter3[B, C, D, O] {
	return func(ctx context.Context, vb B, vc C, vd D) (ret O, err error) {
		return c(ctx, va, vb, vc, vd)
	}
}
