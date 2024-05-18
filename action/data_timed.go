// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"time"

	"github.com/raohwork/task"
)

func (d Data[T]) timed(dur func(time.Duration) time.Duration, e func(error) bool) Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		begin := time.Now()
		ret, err = d(ctx)
		wait := dur(time.Since(begin))
		if wait > 0 && e(err) {
			er := task.Sleep(wait)(ctx)
			if err == nil {
				err = er
			}
		}

		return
	}
}

func delta(d time.Duration) func(time.Duration) time.Duration {
	return func(dur time.Duration) time.Duration { return d - dur }
}

// Timed wraps d into a Data and ensures that it is not returned before dur passed.
//
// It focuses on "How long I should wait before returning". Take a look at example
// of [task.Task.Timed] for how it works.
func (d Data[T]) Timed(dur time.Duration) Data[T] {
	return d.TimedF(delta(dur))
}

// TimedF is like Timed, but use function instead.
//
// The function accepts actual execution time, and returns how long it should wait.
func (d Data[T]) TimedF(f func(time.Duration) time.Duration) Data[T] {
	return d.timed(f, func(_ error) bool { return true })
}

// TimedDone is like Timed, but only successful run is limited.
func (d Data[T]) TimedDone(dur time.Duration) Data[T] {
	return d.TimedDoneF(delta(dur))
}

// TimedDoneF is like TimedDone, but use function instead.
//
// The function accepts actual execution time, and returns how long it should wait.
func (d Data[T]) TimedDoneF(f func(time.Duration) time.Duration) Data[T] {
	return d.timed(f, func(e error) bool { return e == nil })
}

// TimedFail is like Timed, but only failed run is limited.
func (d Data[T]) TimedFail(dur time.Duration) Data[T] {
	return d.TimedFailF(delta(dur))
}

// TimedFailF is like TimedFail, but use function instead.
//
// The function accepts actual execution time, and returns how long it should wait.
func (d Data[T]) TimedFailF(f func(time.Duration) time.Duration) Data[T] {
	return d.timed(f, func(e error) bool { return e != nil })
}
