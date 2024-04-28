// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"

	"github.com/raohwork/task"
)

func (r *Runner) runSomeAsync(names ...string) task.Task {
	ctrls := map[string]chan struct{}{}
	names = append(r.ListDeps(names...), names...)
	tasks := make([]task.Task, 0, len(names))

	for _, n := range names {
		name := n
		ctrls[name] = make(chan struct{})
		tasks = append(tasks, func(ctx context.Context) error {
			defer close(ctrls[name])
			for _, dep := range r.deps[name] {
				<-ctrls[dep]
				if err := r.tasks[dep].err; err != nil {
					return err
				}
			}

			return r.tasks[name].Task(name, r.pre, r.post)(ctx)
		})
	}

	return task.Skip(tasks...)
}
