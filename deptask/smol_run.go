// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"

	"github.com/raohwork/task"
)

type stack struct {
	data []string
	size int
}

func (s *stack) push(x string) {
	s.data[s.size] = x
	s.size++
}

func (s *stack) get() string { return s.data[s.size-1] }

func (s *stack) pop() { s.size-- }

func newStack(l int) *stack {
	return &stack{data: make([]string, l)}
}

func (r *smolRunner) listDeps(name string, visited map[string]bool) (ret []string) {
	ret = make([]string, 0, len(r.tasks))
	q := newStack(len(r.tasks))
	q.push(name)

	for q.size > 0 {
		cur := q.get()

		visited[cur] = true
		fulfilled := true

		for _, d := range r.deps[cur] {
			if !visited[d] {
				fulfilled = false
				q.push(d)
			}
		}

		if !fulfilled {
			continue
		}

		ret = append(ret, cur)
		q.pop()
	}

	return ret
}

func (r *smolRunner) RunSync(ctx context.Context) (err error) {
	if err = r.Validate(); err != nil {
		return
	}

	visited := map[string]bool{}
	for n := range r.tasks {
		if err = r.runSomeSync(ctx, n, visited); err != nil {
			return
		}
	}

	return
}

type taskNode struct {
	body *execStat
	deps []*taskNode
	done chan struct{}
}

func (n *taskNode) over(e error) error {
	close(n.done)
	return e
}

func (n *taskNode) Run(ctx context.Context) (err error) {
	for _, dep := range n.deps {
		select {
		case <-ctx.Done():
			return n.over(ctx.Err())
		case <-dep.done:
			if dep.body.err != nil {
				return n.over(dep.body.err)
			}
		}
	}

	return n.over(n.body.run(ctx))
}

func (r *smolRunner) Run(ctx context.Context) (err error) {
	if err = r.Validate(); err != nil {
		return
	}

	nodes := map[string]*taskNode{}
	for n, t := range r.tasks {
		nodes[n] = &taskNode{
			body: t,
			done: make(chan struct{}),
		}
	}

	for n, deps := range r.deps {
		for _, dep := range deps {
			nodes[n].deps = append(nodes[n].deps, nodes[dep])
		}
	}

	tasks := make([]task.Task, 0, len(nodes))
	for _, n := range nodes {
		tasks = append(tasks, n.Run)
	}

	return task.Skip(tasks...).Run(ctx)
}
