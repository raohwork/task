// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"
	"errors"

	"github.com/raohwork/task"
)

type stack struct {
	data []string
	tail int
	size int
}

func (s *stack) push(name string) {
	if s.tail == s.size {
		s.data = append(s.data, name)
		s.size = len(s.data)
		s.tail = s.size
		return
	}
	s.data[s.tail] = name
	s.tail++
}

func (s *stack) pop() (ret string, ok bool) {
	if s.tail == 0 {
		return
	}

	s.tail--
	return s.data[s.tail], true
}

func (s *stack) len() int { return s.tail }

// ErrMissing indicates a dependency is missing.
type ErrMissing string

func (e ErrMissing) Error() string {
	return "missing dependency: " + string(e)
}

var (
	// ErrDup indicates task name is already used.
	ErrDup = errors.New("duplicated task name")
	// ErrCyclic indicates there're cyclic dependencies.
	ErrCyclic = errors.New("cyclic dependencies detected")
)

type state int

const (
	// pending means the task might be running or waiting to run
	pending state = iota
	// executed means that task has been executed
	executed
	// skipped means the task will not be executed, see [Runner.Skip].
	skipped
)

func (s state) String() string {
	switch s {
	case pending:
		return "pending"
	case executed:
		return "executed"
	case skipped:
		return "skipped"
	default:
		return "INVALID"
	}
}

type taskStat struct {
	task task.Task
	state
	err error
}

func (t *taskStat) Task(name string, pre func(string), post func(string, bool, error)) task.Task {
	return func(ctx context.Context) error {
		if t.state == executed {
			return t.err
		}

		if t.state == pending {
			if pre != nil {
				pre(name)
			}
			t.err = t.task.Run(ctx)
			t.state = executed
		}

		if post != nil {
			post(name, t.state == skipped, t.err)
		}
		return t.err
	}
}

// New is shortcut to WithHook(nil, nil)
func New() *Runner { return WithHook(nil, nil) }

// WithHook creates a Runner with two hooks.
//
// pre is called right before actually executing the task. post is called after the
// execution, or everytime a skipped task is detected.
//
// This was originally designed to show debug log.
func WithHook(pre func(string), post func(string, bool, error)) *Runner {
	return &Runner{
		deps:  map[string][]string{},
		tasks: map[string]*taskStat{},
		pre:   pre,
		post:  post,
	}
}

// Runner manages task dependencies and runs the tasks.
//
// Runner is not thread-safe, you MUST NOT share same instance among multiple
// goroutines.
//
// Runner remembers whether a task is executed or not. Take a look at examples for
// detail.
//
// Zero value denotes an empty Runner, internal data structures are initialized when
// [Runner.Add] or [Runner.Validate] is called.
type Runner struct {
	deps    map[string][]string
	tasks   map[string]*taskStat
	pre     func(string)
	post    func(string, bool, error)
	checked bool
	lastErr error
	groups  [][]string
}

func (r *Runner) init() {
	if r.deps == nil {
		r.deps = map[string][]string{}
	}
	if r.tasks == nil {
		r.tasks = map[string]*taskStat{}
	}
}
