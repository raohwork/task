// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"testing"

	"github.com/raohwork/task"
)

func TestAsync(t *testing.T) {
	c := asyncTestCase{
		run: func(r *Runner) func(...string) task.Task {
			return r.runSomeAsync
		},
	}

	c.Test(t)
}

func BenchmarkAsync(b *testing.B) {
	c := asyncTestCase{func(r *Runner) func(...string) task.Task {
		return r.runSomeAsync
	}}
	c.benchRun(b)
}
