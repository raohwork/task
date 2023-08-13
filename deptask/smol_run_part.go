// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"

	"github.com/raohwork/task"
)

func (r *smolRunner) runSomeSync(ctx context.Context, name string, visited map[string]bool) (err error) {
	for _, dep := range r.listDeps(name, visited) {
		err = r.tasks[dep].run(ctx)
		if err != nil {
			return
		}
	}
	return
}

func (r *smolRunner) RunSomeSync(ctx context.Context, names ...string) (err error) {
	if len(names) == 0 {
		return r.RunSync(ctx)
	}

	if err = r.Validate(); err != nil {
		return
	}

	visited := map[string]bool{}
	for _, name := range names {
		if err = r.runSomeSync(ctx, name, visited); err != nil {
			return
		}
	}
	return
}

func (r *smolRunner) RunSome(ctx context.Context, names ...string) (err error) {
	if len(names) == 0 {
		return r.Run(ctx)
	}

	if err = r.Validate(); err != nil {
		return
	}

	visited := map[string]bool{}
	nodes := map[string]*taskNode{}
	for _, n := range names {
		for _, dep := range r.listDeps(n, visited) {
			nodes[dep] = &taskNode{
				body: r.tasks[dep],
				done: make(chan struct{}),
			}
		}
	}

	for n, t := range nodes {
		deps := r.deps[n]
		l := len(deps)
		t.deps = make([]*taskNode, 0, l)
		for _, dep := range deps {
			t.deps = append(t.deps, nodes[dep])
		}
	}

	tasks := make([]task.Task, 0, len(nodes))
	for _, n := range nodes {
		tasks = append(tasks, n)
	}

	return task.Skip(tasks...).Run(ctx)
}
