// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package forge defines [Generator], a (maybe) cancellable funtion that generates
// one value on each run.
//
// Like in task package, [Tiny] means the underlying function does not receives a
// context, and [Micro] is never-fail [Tiny].
package forge
