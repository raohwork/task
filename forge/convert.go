// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import "context"

// Chain creates a Generator that generates new value from result of input using f.
//
// # Context
//
// Returned generator might use different context with input: The context passed to
// f is derived from ctx, and context passed to input is derived from the context
// passed to f.
//
//	g1 := myGenerator.With(task.Timeout(time.Second))
//	g2 := Chain(g1, f).With(task.Timeout(3 * time.Second))
//	g2.Run(ctx)
//
// is same as
//
//	func(ctx context.Context) (O, error) {
//		ctx2, can2 := context.WithTimeout(ctx, 3*time.Second)
//		defer can2()
//
//		v, err := func(ctx context.Context) (I, error) (
//			ctx1, can1 := context.WithTimeout(ctx, time.Second)
//			defer can1()
//			return g1(ctx1)
//		}(ctx2)
//		if err != nil {
//			return err
//		}
//
//		return f(ctx2, v)
//	}(ctx)
//
// Same rule applies to [Combine].
//
// # Retrying
//
// You should take extra care if you want to use [Generator.RetryN]. Each call to
// the new generator also calls input to generate a value. Same rule applies to
// [Combine].
func Chain[I, O any](input Generator[I], f func(context.Context, I) (O, error)) Generator[O] {
	return func(ctx context.Context) (o O, err error) {
		i, err := input.Run(ctx)
		if err != nil {
			return
		}
		return f(ctx, i)
	}
}

// ChainTiny is "tiny" version of Chain.
func ChainTiny[I, O any](input Generator[I], f func(I) (O, error)) Generator[O] {
	return Chain(input, func(_ context.Context, v I) (O, error) {
		return f(v)
	})
}

// ChainMicro is "micro" version of Chain.
//
// It will return error only if input returns an error.
func ChainMicro[I, O any](input Generator[I], f func(I) O) Generator[O] {
	return Chain(input, func(_ context.Context, v I) (O, error) {
		return f(v), nil
	})
}

// Convert is non-cancellable version of Chain.
//
// Deprecated: use [ChainTiny] instead, for better nameing convention.
func Convert[I, O any](input Generator[I], f func(I) (O, error)) Generator[O] {
	return func(ctx context.Context) (o O, err error) {
		i, err := input.Run(ctx)
		if err != nil {
			return
		}
		return f(i)
	}
}

// Combine creates a Generator that uese result of i1 and i2 to generate new value.
//
// You should take extra care about context and retrying. Take a look at [Convert]
// for explaination.
func Combine[I1, I2, O any](
	i1 Generator[I1],
	i2 Generator[I2],
	f func(context.Context, I1, I2) (O, error),
) Generator[O] {
	return func(ctx context.Context) (ret O, err error) {
		a, err := i1(ctx)
		if err != nil {
			return
		}

		b, err := i2(ctx)
		if err != nil {
			return
		}

		return f(ctx, a, b)
	}
}

// CombineTiny is "tiny" version of Combine.
func CombineTiny[I1, I2, O any](
	i1 Generator[I1],
	i2 Generator[I2],
	f func(I1, I2) (O, error),
) Generator[O] {
	return Combine(i1, i2, func(_ context.Context, a I1, b I2) (O, error) {
		return f(a, b)
	})
}

// CombineMicro is "micro" version of Combine.
//
// It will return error only if i1 or i2 returns an error.
func CombineMicro[I1, I2, O any](
	i1 Generator[I1],
	i2 Generator[I2],
	f func(I1, I2) O,
) Generator[O] {
	return Combine(i1, i2, func(_ context.Context, a I1, b I2) (O, error) {
		return f(a, b), nil
	})
}
