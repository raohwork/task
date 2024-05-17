// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"

	"github.com/raohwork/task"
)

// Converter is a function to convert one type of data into another.
type Converter[I, O any] func(context.Context, I) (O, error)

// Get creates a Converter, mostly for type converting.
func Get[I, O any](f func(context.Context, I) (O, error)) Converter[I, O] { return f }

// NoCtxGet creates Converter from f by ignoring context.
func NoCtxGet[I, O any](f func(I) (O, error)) Converter[I, O] {
	return func(_ context.Context, i I) (O, error) {
		return f(i)
	}
}

// NoErrGet creates Converter from f by ignoring context and error.
func NoErrGet[I, O any](f func(I) O) Converter[I, O] {
	return func(_ context.Context, i I) (O, error) {
		return f(i), nil
	}
}

// From creates a [Data] of output type by feeding i to the converter.
func (c Converter[I, O]) From(i Data[I]) Data[O] {
	return func(ctx context.Context) (ret O, err error) {
		input, err := i(ctx)
		if err != nil {
			return
		}
		return c(ctx, input)
	}
}

// By is like From but uses raw value instead of [Data].
func (c Converter[I, O]) By(i I) Data[O] {
	return func(ctx context.Context) (ret O, err error) { return c(ctx, i) }
}

// With wraps c to modify the context before run it.
func (c Converter[I, O]) With(mod task.CtxMod) Converter[I, O] {
	return func(ctx context.Context, i I) (ret O, err error) {
		ctx, cancel := mod(ctx)
		defer cancel()
		return c(ctx, i)
	}
}

// Pre wraps c to run f before it.
func (c Converter[I, O]) Pre(f func(I)) Converter[I, O] {
	return func(ctx context.Context, i I) (ret O, err error) {
		f(i)
		return c(ctx, i)
	}
}

// Post wraps c to run f after it.
func (c Converter[I, O]) Post(f func(I, O, error)) Converter[I, O] {
	return func(ctx context.Context, i I) (ret O, err error) {
		ret, err = c(ctx, i)
		f(i, ret, err)
		return
	}
}

// Defer wraps c to run f after it.
func (c Converter[I, O]) Defer(f func()) Converter[I, O] {
	return func(ctx context.Context, i I) (ret O, err error) {
		ret, err = c(ctx, i)
		f()
		return
	}
}
