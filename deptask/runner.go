// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"

	"github.com/raohwork/task"
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

func newStat(t task.Task) *execStat {
	return &execStat{
		task:  t,
		state: pending,
	}
}

// execStat represents execution state of a task.
//
// It's designed for debug purpose, use it in production is kind of bad smell.
type execStat struct {
	task  task.Task
	state state
	err   error
}

func (n *execStat) run(ctx context.Context) error {
	switch n.state {
	case executed:
		return n.err
	case skipped:
		return nil
	}

	n.err = n.task.Run(ctx)
	n.state = executed
	return n.err
}

// Runner manages task dependencies and runs the tasks.
//
// Runner is not thread-safe, you MUST NOT share same instance amoung multiple
// goroutines.
//
// Runner remembers whether a task is executed or not. Take a look at examples for
// detail.
type Runner interface {
	// Add adds a task to Runner, will return ErrDup if name has been used.
	Add(name string, taskBody task.Task, deps ...string) error
	// MustAdd is like Add, but panics instead of returning error.
	MustAdd(name string, taskBody task.Task, deps ...string)
	// Mark some tasks to be skipped. Skipped task will not be executed, just
	// pretends that it has finished successfully.
	Skip(name ...string)
	// CopyTo copies specified tasks and their deps to dst using dst.Add.
	//
	// Non-exist tasks are ignored silently. Say you have a Runner contains four
	// tasks: a, b, c (depends b) and d. Calling with a, c, f will add task a,
	// b and c into dst.
	//
	// ErrDup returned by dst.Add is silently ignored.
	//
	// It will call Runner.Validate (and returns error if any) before actually
	// coping tasks.
	CopyTo(dst smolRunner, name ...string) error
	// RunSync validates dependencies and runs all tasks synchronously.
	//
	// The order is unspecified, only dependencies are ensured.
	RunSync(ctx context.Context) error
	// Run validates dependencies and runs all tasks concurrently.
	// You have to take care of race conditions.
	Run(ctx context.Context) error
	// RunSomeSync validates dependencies and runs some tasks (and deps)
	// synchronously.
	RunSomeSync(ctx context.Context, names ...string) error
	// RunSome validates dependencies and runs some tasks (and deps)
	// concurrently.
	RunSome(ctx context.Context, names ...string) error
	// Validate validates the Runner, reports following errors if any:
	//
	//   - ErrMissing: one or more dependecy is missing.
	//   - ErrCyclic: there are cyclic dependencies.
	//
	// It caches the result until you call Runner.Add. Feel free to run it
	// multiple times.
	Validate() error
}

// New creates a Runner.
//
// You're suggested to use New to create runner, as it ensures internal map is
// created brefore you use it. However, Runner.Add can initialize the Runner too.
func New() Runner { return newSmolRunner() }
