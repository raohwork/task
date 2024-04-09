// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"github.com/raohwork/task"
)

func newSmolRunner() *smolRunner {
	return &smolRunner{
		deps:  map[string][]string{},
		tasks: map[string]*execStat{},
	}
}

type smolRunner struct {
	deps      map[string][]string
	tasks     map[string]*execStat
	checked   bool
	lastCheck error
}

func (r *smolRunner) MustAdd(name string, taskBody task.Task, deps ...string) {
	if err := r.Add(name, taskBody, deps...); err != nil {
		panic(err)
	}
}
func (r *smolRunner) Add(name string, taskBody task.Task, deps ...string) error {
	if r.deps == nil {
		r.deps = map[string][]string{}
	}
	if r.tasks == nil {
		r.tasks = map[string]*execStat{}
	}

	if _, ok := r.tasks[name]; ok {
		return ErrDup
	}

	r.checked = false
	r.tasks[name] = newStat(taskBody)
	r.deps[name] = deps
	return nil
}

func (r *smolRunner) Skip(name ...string) {
	for _, n := range name {
		if _, ok := r.tasks[n]; !ok {
			return
		}
		r.tasks[n].state = skipped
	}
}

func (r *smolRunner) CopyTo(dst Runner, name ...string) error {
	if err := r.Validate(); err != nil {
		return err
	}

	deps := map[string]bool{}
	for _, n := range name {
		if _, ok := r.tasks[n]; !ok {
			// no such task, skip
			continue
		}

		r.listDeps(n, deps)
	}

	for n := range deps {
		dst.Add(n, r.tasks[n].task, r.deps[n]...)
	}

	return nil
}
