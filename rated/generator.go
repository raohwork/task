// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rated

import (
	"context"
	"time"

	"github.com/raohwork/task"
	"github.com/raohwork/task/forge"
	"golang.org/x/time/rate"
)

// New creates a [forge.Generator] that respects the rate limit.
//
// Say you have an fixed generator r with rate limit to once per second:
//
//	r.Run() // executed immediatly
//	r.Run() // executed after a second
func NewG[T any](l *rate.Limiter, g forge.Generator[T]) forge.Generator[T] {
	return func(ctx context.Context) (ret T, err error) {
		reserve := l.Reserve()
		if err = task.Sleep(reserve.Delay()).Run(ctx); err != nil {
			reserve.Cancel()
			return
		}
		return g.Run(ctx)
	}
}

// EveryG is a wrapper of NewG.
func EveryG[T any](dur time.Duration, g forge.Generator[T]) forge.Generator[T] {
	return NewG(rate.NewLimiter(rate.Every(dur), 1), g)
}
