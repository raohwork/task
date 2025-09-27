// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nctask

// Next creates a task that runs next after t finished successfully.
func (t Task) Then(next Task) Task {
	return func() error {
		if err := t.Run(); err != nil {
			return err
		}

		return next.Run()
	}
}

// Iter creates a task run tasks with same context and stops at first error.
func Iter(tasks ...Task) Task {
	return func() error {
		for _, t := range tasks {
			if err := t.Run(); err != nil {
				return err
			}
		}

		return nil
	}
}

// Wait creates a task that runs all task concurrently, wait them get done, and
// return first non-nil error.
func Wait(tasks ...Task) Task {
	return func() (err error) {
		ch := make(chan error)
		for _, t := range tasks {
			t.GoWithChan(ch)
		}

		for range tasks {
			e := <-ch
			if err == nil && e != nil {
				err = e
			}
		}

		return
	}
}
