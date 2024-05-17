// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"errors"
)

// HandleErr creates a task that handles specific error after running t.
// It could change the error returned by Run. f is called only if t.Run returns an
// error.
func (t Task) HandleErr(f func(error) error) Task {
	return func(ctx context.Context) error {
		err := t.Run(ctx)
		if err != nil {
			err = f(err)
		}

		return err
	}
}

// HandleErrWithContext is like HandleErr, but uses same context used in f.
func (t Task) HandleErrWithContext(f func(context.Context, error) error) Task {
	return func(ctx context.Context) error {
		err := t.Run(ctx)
		if err != nil {
			err = f(ctx, err)
		}

		return err
	}
}

// ContextError detects if err is [context.Canceled] or [context.DeadlineExceeded].
func ContextError(err error) bool {
	return ErrorIs(context.Canceled, context.DeadlineExceeded)(err)
}

// ErrorIs creates a function to detect if the error is listed.
func ErrorIs(errs ...error) func(error) bool {
	return func(err error) bool {
		for _, e := range errs {
			if errors.Is(err, e) {
				return true
			}
		}
		return false
	}
}

// OnlyErrs preserves errors if f(error) is true.
func (t Task) OnlyErrs(f func(error) bool) Task {
	return t.HandleErr(func(err error) error {
		if f(err) {
			return err
		}
		return nil
	})
}

// IgnoreErrs ignores errors if f(error) is true.
func (t Task) IgnoreErrs(f func(error) bool) Task {
	return t.HandleErr(func(err error) error {
		if f(err) {
			return nil
		}
		return err
	})
}

// IgnoreErr ignores the error returned by t.Run if it is not context error.
//
// It is shortcut to t.OnlyErrs(ContextError).
func (t Task) IgnoreErr() Task {
	return t.OnlyErrs(ContextError)
}
