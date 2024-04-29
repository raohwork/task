// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"

	"github.com/raohwork/task"
)

// ListDeps lists the dependencies of given tasks, or all tasks if no task is given.
// Order is unspecified.
func (r *Runner) ListDeps(names ...string) (ret []string) {
	if len(names) == 0 {
		ret = make([]string, 0, len(r.deps))
		for name := range r.deps {
			ret = append(ret, name)
		}
		return
	}

	done := map[string]bool{}
	for _, spec := range names {
		done[spec] = true
	}

	for _, spec := range names {
		stack := &stack{data: make([]string, 0, len(r.deps))}
		for _, name := range r.deps[spec] {
			stack.push(name)
		}

		name, ok := stack.pop()
		for ok {
			if done[name] {
				name, ok = stack.pop()
				continue
			}
			done[name] = true

			ret = append(ret, name)
			for _, dep := range r.deps[name] {
				stack.push(dep)
			}
			name, ok = stack.pop()
		}
	}

	return
}

// Add adds a task to Runner, will return ErrDup if name has been used.
func (r *Runner) Add(name string, f task.Task, deps ...string) error {
	r.init()

	if _, ok := r.tasks[name]; ok {
		return ErrDup
	}

	r.checked = false
	r.tasks[name] = &taskStat{task: f, state: pending}
	r.deps[name] = deps
	return nil
}

// MustAdd is like Add, but panics instead of returning error.
func (r *Runner) MustAdd(name string, f task.Task, deps ...string) {
	if err := r.Add(name, f, deps...); err != nil {
		panic(err)
	}
}

// Mark some tasks to be skipped. Skipped task will not be executed, just pretends
// that it has finished. Other tasks in dependency tree are unaffected.
//
// Unknown task is ignored silently.
func (r *Runner) Skip(names ...string) {
	for _, name := range names {
		stat, ok := r.tasks[name]
		if !ok {
			continue
		}
		if stat.state == pending {
			stat.state = skipped
		}
	}
}

// Skipped check if named task is marked as skip. If the task does not exist,
// ErrMissing is returned.
func (r *Runner) Skipped(name string) (bool, error) {
	stat, ok := r.tasks[name]
	if !ok {
		return false, ErrMissing(name)
	}

	return stat.state == skipped, nil
}

// CopyTo copies specified tasks and their deps to dst using dst.Add. Useful when
// testing.
//
// Non-exist tasks are ignored silently. Say you have a Runner contains four tasks:
// a, b, c (depends b) and d. Calling with a, c, f will add task a, b and c into
// dst.
//
// ErrDup returned by dst.Add is silently ignored.
//
// It will call Runner.Validate on r (and returns error if any) before actually
// coping tasks.
//
// State is not copied! Use with caution!
func (r *Runner) CopyTo(dst *Runner, names ...string) error {
	if err := r.Validate(); err != nil {
		return err
	}
	deps := append(r.ListDeps(names...), names...)

	for _, name := range deps {
		dst.Add(name, r.tasks[name].task, r.deps[name]...)
	}
	return nil
}

// RunSync is shortcut to "RunSomesync(ctx)"
func (r *Runner) RunSync(ctx context.Context) error {
	return r.runSomeSync()(ctx)
}

// Run is shortcut to "RunSome(ctx)"
func (r *Runner) Run(ctx context.Context) error {
	return r.runSomeAsync()(ctx)
}

// RunSomeSync validates dependencies and runs some tasks (and deps) synchronously.
// The order is unspecified, only dependencies are ensured.
//
// It returns immediately after first error.
func (r *Runner) RunSomeSync(ctx context.Context, names ...string) error {
	if err := r.Validate(); err != nil {
		return err
	}
	return r.runSomeSync(names...)(ctx)
}

// RunSome validates dependencies and runs some tasks (and deps) concurrently.
//
// When an error occurred, other tasks are canceled, prevents further execution.
func (r *Runner) RunSome(ctx context.Context, names ...string) error {
	if err := r.Validate(); err != nil {
		return err
	}
	return r.runSomeAsync(names...)(ctx)
}

// Validate validates the Runner, reports following errors if any:
//
//   - ErrMissing: one or more dependecy is missing.
//   - ErrCyclic: there are cyclic dependencies.
//
// It caches the result until you call Runner.Add. Feel free to run it
// multiple times.
func (r *Runner) Validate() error {
	if !r.checked {
		r.lastErr = r.validate()
		r.checked = true
	}

	return r.lastErr
}
