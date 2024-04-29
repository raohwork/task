// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"time"

	"github.com/raohwork/task"
)

func (g Generator[T]) timed(dur func(time.Duration) time.Duration, e func(error) bool) Generator[T] {
	return func(ctx context.Context) (T, error) {
		begin := time.Now()
		v, err := g.Run(ctx)
		wait := dur(time.Since(begin))
		if wait > 0 && e(err) {
			er := task.Sleep(wait).Run(ctx)
			if err == nil {
				err = er
			}
		}

		return v, err
	}
}

func delta(d time.Duration) func(time.Duration) time.Duration {
	return func(dur time.Duration) time.Duration { return d - dur }
}

// TimedFail is like [task.TimedFail], but applies to Generator.
func (g Generator[T]) TimedFail(dur time.Duration) Generator[T] {
	return g.TimedFailF(delta(dur))
}

// TimedFailF is like [task.TimedFailF], but applies to Generator.
func (g Generator[T]) TimedFailF(f func(time.Duration) time.Duration) Generator[T] {
	return g.timed(f, func(e error) bool { return e != nil })
}

// TimedDone is like [task.TimedDone], but applies to Generator.
func (g Generator[T]) TimedDone(dur time.Duration) Generator[T] {
	return g.TimedDoneF(delta(dur))
}

// TimedDoneF is like [task.TimedDoneF], but applies to Generator.
func (g Generator[T]) TimedDoneF(f func(time.Duration) time.Duration) Generator[T] {
	return g.timed(f, func(e error) bool { return e == nil })
}

// Timed is like [task.Timed], but applies to Generator.
func (g Generator[T]) Timed(dur time.Duration) Generator[T] {
	return g.TimedF(delta(dur))
}

// TimedF is like [task.TimedF], but applies to Generator.
func (g Generator[T]) TimedF(f func(time.Duration) time.Duration) Generator[T] {
	return g.timed(f, func(e error) bool { return true })
}
