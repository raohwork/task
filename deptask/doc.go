// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package deptask provides a tool, [Runner], to run tasks in order according to its
// dependency.
//
// [Runner] is designed to handle complex initializing process of large project. All
// method that executes the tasks follows fail-fast principle.
package deptask
