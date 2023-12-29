// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"time"
)

func (t Task) timed(dur func(time.Duration) time.Duration, e func(error) bool) Task {
	return func(ctx context.Context) error {
		begin := time.Now()
		err := t.Run(ctx)
		wait := dur(time.Since(begin))
		if wait > 0 && e(err) {
			er := Sleep(wait).Run(ctx)
			if err == nil {
				err = er
			}
		}

		return err
	}
}

func delta(d time.Duration) func(time.Duration) time.Duration {
	return func(dur time.Duration) time.Duration { return d - dur }
}

// Timed wraps t into a task ensures that it is not returned before dur passed.
//
// It focuses on "How long I should wait before returning". Take a look at example
// for how it works.
//
// If you're looking for rate limiting solution, you should take a look at "rated"
// subdirectory.
func (t Task) Timed(dur time.Duration) Task {
	return t.TimedF(delta(dur))
}

// TimedF is like Timed, but use function instead.
func (t Task) TimedF(f func(time.Duration) time.Duration) Task {
	return t.timed(f, func(_ error) bool { return true })
}

// TimedDone is like Timed, but limits only successful run.
//
// If you're looking for rate limiting solution, you should take a look at "rated"
// subdirectory.
func (t Task) TimedDone(dur time.Duration) Task {
	return t.TimedDoneF(delta(dur))
}

// TimedDoneF is like TimedDone, but use function instead.
func (t Task) TimedDoneF(f func(time.Duration) time.Duration) Task {
	return t.timed(f, func(e error) bool { return e == nil })
}

// TimedFail is like Timed, but limits only failed run.
//
// If you're looking for rate limiting solution, you should take a look at "rated"
// subdirectory.
func (t Task) TimedFail(dur time.Duration) Task {
	return t.TimedFailF(delta(dur))
}

// TimedFailF is like TimedFail, but use function instead.
func (t Task) TimedFailF(f func(time.Duration) time.Duration) Task {
	return t.timed(f, func(e error) bool { return e != nil })
}
