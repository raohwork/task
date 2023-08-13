// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"errors"
)

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
