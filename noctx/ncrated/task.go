// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ncrated

import (
	"time"

	"github.com/raohwork/task/noctx/nctask"
	"golang.org/x/time/rate"
)

// New creates a [nctask.Task] that respects the rate limit.
//
// Say you have an empty task r with rate limit to once per second:
//
//	r.Run() // executed immediately
//	r.Run() // executed after a second
//
// Deprecated: use [Task] instead.
func New(l *rate.Limiter, t nctask.Task) (ret nctask.Task) {
	return Task(l, t)
}

// Task creates a [nctask.Task] that respects the rate limit.
//
// Say you have an empty task r with rate limit to once per second:
//
//	r.Run() // executed immediately
//	r.Run() // executed after a second
func Task(l *rate.Limiter, t nctask.Task) nctask.Task {
	return func() error {
		reserve := l.Reserve()
		time.Sleep(reserve.Delay())
		return t.Run()
	}
}
