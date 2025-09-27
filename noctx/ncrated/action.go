// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ncrated

import (
	"time"

	"github.com/raohwork/task/noctx/ncaction"
	"golang.org/x/time/rate"
)

// Data creates an [ncaction.Data] that respects the rate limit, quite like [Task].
//
// There's no rate limited [ncaction.Action] or [ncaction.Converter]. For actions, rate
// limit should be applied on resulted task. For converters, it should be applied on
// resulted data.
func Data[T any](l *rate.Limiter, d ncaction.Data[T]) ncaction.Data[T] {
	return func() (ret T, err error) {
		reserve := l.Reserve()
		time.Sleep(reserve.Delay())
		return d()
	}
}
