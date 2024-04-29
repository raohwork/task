// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"
	"slices"
	"sync"
	"testing"

	"github.com/raohwork/task"
)

// This can be faster since the order has been ensured in Validate(). But it will
// change dependencies so I'm not using it.
//
// # Changes of dependencies
//
// Since tasks are splitted into some groups according to their degree, it becomes
// group A depends on group B.
func (r *Runner) runSomeAsync3(names ...string) task.Task {
	has := map[string]bool{}
	names = append(r.ListDeps(names...), names...)
	for _, name := range names {
		has[name] = true
	}
	tasks := make([]task.Task, 0, len(r.groups))
	for _, group := range r.groups {
		local := make([]task.Task, 0, len(group))
		for _, name := range group {
			if !has[name] {
				continue
			}
			local = append(local, r.tasks[name].Task(name, r.pre, r.post))
		}
		if len(local) > 0 {
			tasks = append(tasks, task.Skip(local...))
		}
	}

	slices.Reverse(tasks)
	return task.Iter(tasks...)
}

// This is here to test if using lock is better than channel.
//
// Using lock is not Go idiom since we're passing "state" to another goroutine,
// exactly what channel designed for. But I've found some articles explaining why
// lock is faster than channel in general. When in doubt, write test and benchmark.
func (r *Runner) runSomeAsync2(names ...string) task.Task {
	names = append(r.ListDeps(names...), names...)
	tasks := make([]task.Task, len(names))
	locks := map[string]*sync.RWMutex{}

	for _, name := range names {
		locks[name] = &sync.RWMutex{}
		locks[name].Lock()
		// it is unlocked when task has been processed
	}

	for idx, name := range names {
		name := name
		tasks[idx] = func(ctx context.Context) error {
			defer locks[name].Unlock()

			for _, dep := range r.deps[name] {
				// lock again to ensure dep has been processed
				locks[dep].RLock()
				locks[dep].RUnlock()

				if err := r.tasks[dep].err; err != nil {
					return err
				}
			}

			return r.tasks[name].Task(name, r.pre, r.post)(ctx)
		}
	}

	return task.Skip(tasks...)
}

func TestAsync2(t *testing.T) {
	c := asyncTestCase{
		run: func(r *Runner) func(...string) task.Task {
			return r.runSomeAsync2
		},
	}

	c.Test(t)
}

func TestAsync3(t *testing.T) {
	c := asyncTestCase{
		run: func(r *Runner) func(...string) task.Task {
			return r.runSomeAsync3
		},
	}

	c.Test(t)
}

func BenchmarkAsync2(b *testing.B) {
	c := asyncTestCase{func(r *Runner) func(...string) task.Task {
		return r.runSomeAsync2
	}}
	c.benchRun(b)
}

func BenchmarkAsync3(b *testing.B) {
	c := asyncTestCase{func(r *Runner) func(...string) task.Task {
		return r.runSomeAsync3
	}}
	c.benchRun(b)
}
