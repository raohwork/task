// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nctask

import (
	"sync"

	"github.com/raohwork/task"
)

// Task repeasents a routine that not cancellable.
type Task func() error

// Run runs the task, equals to t().
func (t Task) Run() error { return t() }

// Go runs t in separated goroutine and returns a channel to retrieve error.
//
// It's safe to ignore the channel if you don't need the result.
func (t Task) Go() <-chan error {
	ret := make(chan error, 1)
	t.GoWithChan(ret)
	return ret
}

// GoWithChan runs t in separated goroutine and sends returned error into ch.
func (t Task) GoWithChan(ch chan<- error) {
	go func() { ch <- t.Run() }()
}

// NoErr converts the task into a simple function that ignores returned error.
func (t Task) NoErr() func() {
	return func() { t() }
}

// Once creates a task that can be run only once, further attempt returns task.ErrOnce.
func (t Task) Once() Task {
	var once sync.Once
	return func() (err error) {
		err = task.ErrOnce
		once.Do(func() {
			err = t()
		})
		return
	}
}

// Cached wraps t to cache the result, and reuse it in later call.
func (t Task) Cached() Task {
	var (
		once sync.Once
		err  error
	)
	return func() error {
		once.Do(func() {
			err = t()
		})
		return err
	}
}

// Defer wraps t to run f after it.
func (t Task) Defer(f func()) Task {
	return func() (err error) {
		err = t()
		f()
		return err
	}
}

// Pre wraps t to run f before it.
func (t Task) Pre(f func()) Task {
	return func() (err error) {
		f()
		err = t()
		return
	}
}

// Post wraps t to run f after it.
func (t Task) Post(f func(error)) Task {
	return func() (err error) {
		err = t()
		f(err)
		return
	}
}

// AlterError wraps t to run f to alter the error before returning.
func (t Task) AlterError(f func(error) error) Task {
	return func() error {
		return f(t())
	}
}
