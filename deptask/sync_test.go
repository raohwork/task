// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"testing"

	"github.com/raohwork/task"
)

func TestSync(t *testing.T) {
	c := asyncTestCase{func(r *Runner) func(...string) task.Task {
		return r.runSomeSync
	}}
	c.Test(t)
}

func BenchmarkSync(b *testing.B) {
	c := asyncTestCase{func(r *Runner) func(...string) task.Task {
		return r.runSomeSync
	}}
	c.benchRun(b)
}
