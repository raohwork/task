// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import "testing"

func TestRunner(t *testing.T) {
	t.Run("smol", (testCase{
		creator: func() Runner { return newSmolRunner() },
	}).TestAll)
}

func BenchmarkListDeps(b *testing.B) {
	r := newSmolRunner()
	r.Add("a", nop)
	r.Add("b", nop)
	r.Add("c", nop, "a", "b")
	r.Add("d", nop, "a")
	r.Add("e", nop, "b")
	r.Add("f", nop, "c", "e")
	r.Validate()

	b.ResetTimer()
	for x := 0; x < b.N; x++ {
		r.listDeps("f", map[string]bool{})
	}
}
