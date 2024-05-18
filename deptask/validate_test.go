// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/raohwork/task"
)

func TestListDepsSome(t *testing.T) {
	r := New()
	f := task.NoErr(func() {})
	r.MustAdd("a", f, "b", "c")
	r.MustAdd("b", f)
	r.MustAdd("c", f)
	r.MustAdd("d", f)

	ret := r.ListDeps("a")
	slices.Sort(ret)
	expect := []string{"b", "c"}

	if !reflect.DeepEqual(ret, expect) {
		t.Fatal("unexpected result:", ret)
	}

}

func TestCyclicOf2(t *testing.T) {
	r := New()
	r.MustAdd("a", task.NoErr(func() {}), "b")
	r.MustAdd("b", task.NoErr(func() {}), "a")
	if err := r.Validate(); err != ErrCyclic {
		t.Fatal("unexpected result:", err)
	}
}

func TestCyclicOf3(t *testing.T) {
	r := New()
	r.MustAdd("a", task.NoErr(func() {}), "b")
	r.MustAdd("b", task.NoErr(func() {}), "c")
	r.MustAdd("c", task.NoErr(func() {}), "d")
	r.MustAdd("d", task.NoErr(func() {}), "b")
	if err := r.Validate(); err != ErrCyclic {
		t.Fatal("unexpected result:", err)
	}
}

func TestMissing(t *testing.T) {
	r := New()
	r.MustAdd("a", task.NoErr(func() {}), "b")
	err := r.Validate()
	var m ErrMissing
	if !errors.As(err, &m) {
		t.Fatal("unexpected result:", err)
	}
}

func BenchmarkValidate(b *testing.B) {
	r := New()
	r.MustAdd("a", task.NoErr(func() {}), "b", "c", "d")
	r.MustAdd("b", task.NoErr(func() {}), "c", "e")
	r.MustAdd("c", task.NoErr(func() {}))
	r.MustAdd("d", task.NoErr(func() {}))
	r.MustAdd("e", task.NoErr(func() {}), "d")
	r.MustAdd("f", task.NoErr(func() {}), "a", "g")
	r.MustAdd("g", task.NoErr(func() {}), "e")
	r.MustAdd("z", task.NoErr(func() {}))

	for i := 0; i < b.N; i++ {
		r.checked = false
		r.Validate()
	}
}
