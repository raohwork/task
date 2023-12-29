// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rated

import (
	"context"
	"time"

	"github.com/raohwork/task"
	"golang.org/x/time/rate"
)

// New creates a [task.Task] that respects the rate limit.
//
// Say you have an empty task r with rate limit to once per second:
//
//	r.Run() // executed immediatly
//	r.Run() // executed after a second
func New(l *rate.Limiter, t task.Task) (ret task.Task) {
	return func(ctx context.Context) error {
		reserve := l.Reserve()
		if err := task.Sleep(reserve.Delay()).Run(ctx); err != nil {
			reserve.Cancel()
			return err
		}

		return t.Run(ctx)
	}
}

// Every is a wrapper of New.
func Every(dur time.Duration, f task.Task) task.Task {
	return New(rate.NewLimiter(rate.Every(dur), 1), f)
}
