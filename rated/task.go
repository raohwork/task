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
//	r.Run() // executed immediately
//	r.Run() // executed after a second
//
// Deprecated: use [Task] instead.
func New(l *rate.Limiter, t task.Task) (ret task.Task) {
	return Task(l, t)
}

// Task creates a [task.Task] that respects the rate limit.
//
// Say you have an empty task r with rate limit to once per second:
//
//	r.Run() // executed immediately
//	r.Run() // executed after a second
func Task(l *rate.Limiter, t task.Task) task.Task {
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
//
// Deprecated: create your own [rate.Limiter] and use [Task] instead.
func Every(dur time.Duration, f task.Task) task.Task {
	return Task(rate.NewLimiter(rate.Every(dur), 1), f)
}
