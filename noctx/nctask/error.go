// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nctask

// HandleErr creates a task that handles specific error after running t.
// It could change the error returned by Run. f is called only if t.Run returns an
// error.
func (t Task) HandleErr(f func(error) error) Task {
	return func() error {
		err := t.Run()
		if err != nil {
			err = f(err)
		}

		return err
	}
}

// OnlyErrs preserves errors if f(error) is true.
func (t Task) OnlyErrs(f func(error) bool) Task {
	return t.HandleErr(func(err error) error {
		if f(err) {
			return err
		}
		return nil
	})
}

// IgnoreErrs ignores errors if f(error) is true.
func (t Task) IgnoreErrs(f func(error) bool) Task {
	return t.HandleErr(func(err error) error {
		if f(err) {
			return nil
		}
		return err
	})
}
