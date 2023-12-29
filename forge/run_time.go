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

// TimedFail creates a Generator and ensures the run time is longer than dur if it failed.
//
// It focuses on "How long I should wait before returning". Take a look at example
// for how it works.
func (g Generator[T]) TimedFail(dur time.Duration) Generator[T] {
	return g.TimedFailF(delta(dur))
}

// TimedFailF is like TimedFail, but use function instead.
func (g Generator[T]) TimedFailF(f func(time.Duration) time.Duration) Generator[T] {
	return g.timed(f, func(e error) bool { return e != nil })
}
