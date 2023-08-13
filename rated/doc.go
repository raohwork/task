// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package rated controls rate of a task with [rate.Limiter].
//
// A rated task focus on "How long I should wait before I can run again". You might
// want to take a look at package example to compare it with [task.Task.Timed].
//
// It is putted in seperated package so you won't link to unused external
// dependencies if you're not using it.
//
// It's quite common to write following code:
//
//	err = Every(
//		time.Minute,
//		task.T(mytask).
//			TimedFail(3*time.Second).
//			RetryN(3).
//			HandleErr(logError).
//			IgnoreErr(),
//	).Loop().Run(ctx)
//
// The above program will repeatedly execute mytask, with a maximum of one
// successful execution per minute. If the execution is not successful, it will be
// retried at an interval of once every three seconds for a total of three attempts.
// If all three retry attempts fail, the final error will be logged, treating it as
// a successful attempt, and continues to next run.
package rated
