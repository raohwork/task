// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rated

import (
	"context"

	"github.com/raohwork/task"
	"github.com/raohwork/task/action"
	"golang.org/x/time/rate"
)

// Data creates an [action.Data] that respects the rate limit, quite like [Task].
//
// There's no rate limited [action.Action] or [action.Converter]. For actions, rate
// limit should be applied on resulted task. For converters, it should be applied on
// resulted data.
func Data[T any](l *rate.Limiter, d action.Data[T]) action.Data[T] {
	return func(ctx context.Context) (ret T, err error) {
		reserve := l.Reserve()
		if err = task.Sleep(reserve.Delay()).Run(ctx); err != nil {
			reserve.Cancel()
			return
		}

		return d(ctx)
	}
}
