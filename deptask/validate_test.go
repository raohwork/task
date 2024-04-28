// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"errors"
	"testing"

	"github.com/raohwork/task"
)

func TestCyclicOf2(t *testing.T) {
	r := New()
	r.MustAdd("a", task.Micro(func() {}), "b")
	r.MustAdd("b", task.Micro(func() {}), "a")
	if err := r.Validate(); err != ErrCyclic {
		t.Fatal("unexpected result:", err)
	}
}

func TestCyclicOf3(t *testing.T) {
	r := New()
	r.MustAdd("a", task.Micro(func() {}), "b")
	r.MustAdd("b", task.Micro(func() {}), "c")
	r.MustAdd("c", task.Micro(func() {}), "d")
	r.MustAdd("d", task.Micro(func() {}), "b")
	if err := r.Validate(); err != ErrCyclic {
		t.Fatal("unexpected result:", err)
	}
}

func TestMissing(t *testing.T) {
	r := New()
	r.MustAdd("a", task.Micro(func() {}), "b")
	err := r.Validate()
	var m ErrMissing
	if !errors.As(err, &m) {
		t.Fatal("unexpected result:", err)
	}
}

func BenchmarkValidate(b *testing.B) {
	r := New()
	r.MustAdd("a", task.Micro(func() {}), "b", "c", "d")
	r.MustAdd("b", task.Micro(func() {}), "c", "e")
	r.MustAdd("c", task.Micro(func() {}))
	r.MustAdd("d", task.Micro(func() {}))
	r.MustAdd("e", task.Micro(func() {}), "d")
	r.MustAdd("f", task.Micro(func() {}), "a", "g")
	r.MustAdd("g", task.Micro(func() {}), "e")
	r.MustAdd("z", task.Micro(func() {}))

	for i := 0; i < b.N; i++ {
		r.checked = false
		r.Validate()
	}
}
