// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"
	"fmt"
	"slices"

	"github.com/raohwork/task"
)

func (r *Runner) runSomeSync(names ...string) task.Task {
	if len(names) > 0 {
		names = append(r.ListDeps(names...), names...)
	} else {
		names = r.ListDeps()
	}

	m := map[string]bool{}
	for _, name := range names {
		m[name] = true
	}

	var tasks []task.Task
	for _, group := range r.groups {
		for _, name := range group {
			if !m[name] {
				continue
			}
			name := name

			tasks = append(tasks, func(ctx context.Context) error {
				// check if any of depepdencies failed
				for _, dep := range r.deps[name] {
					if err := r.tasks[dep].err; err != nil {
						// skip
						fmt.Println()
						return err
					}
				}
				return r.tasks[name].Task(name, r.pre, r.post)(ctx)
			})
		}
	}

	slices.Reverse(tasks)
	return task.Iter(tasks...)
}
