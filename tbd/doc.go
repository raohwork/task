// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package tbd provides [TBD], a value might be computed some time later.
//
// A TBD *MUST* be "resolved" by a "resolver" or whoever use it might run into
// deadlock.
package tbd
