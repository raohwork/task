// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"

	"github.com/raohwork/task"
	"github.com/raohwork/task/tbd"
)

// Tiny wraps a non-cancellable function into generator.
func Tiny[T any](f func() (T, error)) Generator[T] {
	return func(ctx context.Context) (v T, e error) {
		select {
		case <-ctx.Done():
			return v, ctx.Err()
		default:
			return f()
		}
	}
}

// Micro wraps a never-fail, non-cancellable function into generator.
func Micro[T any](f func() T) Generator[T] {
	return func(ctx context.Context) (v T, e error) {
		select {
		case <-ctx.Done():
			return v, ctx.Err()
		default:
			return f(), nil
		}
	}
}

// G is shortcut to create a generator.
func G[T any](f func(context.Context) (T, error)) Generator[T] { return f }

// Generator much likes [task.Task] but it generates a new value each time.
//
// A generator should generate n values if it is executed n times successfully. So
// some helpers like [task.Task.Loop] or [task.Iter] cannot applies to it.
//
// Take a look at httptask package to see how generator makes your code more clear.
type Generator[T any] func(context.Context) (T, error)

// Run runs the Generator.
func (g Generator[T]) Run(ctx context.Context) (T, error) { return g(ctx) }

// Tiny transforms g to be a non-cancellable generator.
func (g Generator[T]) Tiny() (T, error) {
	return g.Run(context.Background())
}

// Go runs g in separated goroutine and returns a TBD to retrieve result.
func (g Generator[T]) Go(ctx context.Context) tbd.TBD[T] {
	ret, d := tbd.Create[T]()
	task.Task(func(ctx context.Context) error {
		v, err := g(ctx)
		d(v, err)
		return err
	}).Go(ctx)
	return ret
}

// With creates a Generator that the context is derived using modder before running.
func (g Generator[T]) With(modder task.CtxMod) Generator[T] {
	return func(ctx context.Context) (T, error) {
		x, c := modder(ctx)
		defer c()
		return g.Run(x)
	}
}

// HandleErr creates a Generator that handles specific error after running g.
// It will change the error returned by Run. f is called only if g.Run returns an
// error.
func (g Generator[T]) HandleErr(f func(error) error) Generator[T] {
	return func(ctx context.Context) (T, error) {
		v, err := g.Run(ctx)
		if err != nil {
			err = f(err)
		}

		return v, err
	}
}

// IgnoreErr ignores the error returned by g.Run if it is not context error.
//
// Context error means [context.Canceled] and [context.DeadlineExceeded].
func (g Generator[T]) IgnoreErr() Generator[T] {
	return g.OnlyErrs(context.Canceled, context.DeadlineExceeded)
}

// IgnoreErrs ignores specific error.
func (g Generator[T]) IgnoreErrs(errorList ...error) Generator[T] {
	return g.HandleErr(task.IgnoreErrs(errorList...))
}

// OnlyErrs preserves only specific error.
func (g Generator[T]) OnlyErrs(errorList ...error) Generator[T] {
	return g.HandleErr(task.OnlyErrs(errorList...))
}

// Defer wraps g to run f after it.
func (g Generator[T]) Defer(f func()) Generator[T] {
	return func(ctx context.Context) (ret T, err error) {
		ret, err = g(ctx)
		f()
		return
	}
}

// Pre wraps g to run f before it.
func (g Generator[T]) Pre(f func()) Generator[T] {
	return func(ctx context.Context) (ret T, err error) {
		f()
		ret, err = g(ctx)
		return
	}
}

// Post wraps g to run f after it.
func (g Generator[T]) Post(f func(T, error)) Generator[T] {
	return func(ctx context.Context) (ret T, err error) {
		ret, err = g(ctx)
		f(ret, err)
		return
	}
}

// Chain creates a new Generator that generates value from input using f.
//
// Returned generator might use different context with input:
//
//	g1 := myGenerator().With(task.Timeout(time.Second))
//	g2 := Chain(g1, func).With(task.Timeout(3 * time.Second))
//	g2.Run(ctx)
//
// Running output.Run() will be a max 3 second timeout: The context passed to f is
// derived from ctx, and context passed to input is derived from the context passed
// to f.
//
// Returned generator should take extra care if you want to use [Generator.RetryN].
// If the error is returned from f, retrying will feed a new value to f, not the one
// leads to error. If it's not what you want, you should implement Chain-like code
// to fit your own need.
func Chain[I, O any](input Generator[I], f func(context.Context, I) (O, error)) Generator[O] {
	return func(ctx context.Context) (o O, err error) {
		i, err := input.Run(ctx)
		if err != nil {
			return
		}
		return f(ctx, i)
	}
}

// Convert is non-cancellable version of Chain.
func Convert[I, O any](input Generator[I], f func(I) (O, error)) Generator[O] {
	return func(ctx context.Context) (o O, err error) {
		i, err := input.Run(ctx)
		if err != nil {
			return
		}
		return f(i)
	}
}
