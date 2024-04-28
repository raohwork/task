// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"
	"errors"
	"testing"

	"github.com/raohwork/task"
)

type asyncTestCase struct {
	run func(*Runner) func(...string) task.Task
}

func (c asyncTestCase) Run(t asyncTestCaseTool) error {
	return c.run(t.Runner)()(context.TODO())
}
func (c asyncTestCase) RunSome(t asyncTestCaseTool, names ...string) error {
	return c.run(t.Runner)(names...)(context.TODO())
}

// test cases

func (c asyncTestCase) Test(t *testing.T) {
	// test Run/RunAsync
	t.Run("empty", c.empty)
	t.Run("nodep", c.nodep)
	t.Run("simple_dep", c.simpleDep)
	t.Run("multiple_dep", c.multiDep)
	t.Run("deep_dep", c.deepDep)
	t.Run("complex_dep", c.complexDep)
	t.Run("very_complex_dep", c.complexDep2)
	t.Run("skipped", c.skipped)
	t.Run("same", c.same)
	t.Run("some", c.some)
	t.Run("error", c.error)
}

// test Run/RunAsync

func (c asyncTestCase) empty(t *testing.T) {
	r := newTool()
	r.Validate()
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func (c asyncTestCase) nodep(t *testing.T) {
	r := newTool()
	r.add("a")
	r.Validate()
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "a," {
		t.Fatal("unexpected result:", str)
	}
}

func (c asyncTestCase) simpleDep(t *testing.T) {
	r := newTool()
	r.add("a", "b")
	r.add("b")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "b,a," {
		t.Fatal("unexpected result:", str)
	}
}

func (c asyncTestCase) multiDep(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c")
	r.add("b")
	r.add("c")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "b,c,a," && str != "c,b,a," {
		t.Fatal("unexpected result:", str)
	}
}

func (c asyncTestCase) deepDep(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c")
	r.add("b", "d")
	r.add("c")
	r.add("d")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.checkOrder("b"); str != "" {
		t.Fatal(str)
	}
	if str := r.checkOrder("a"); str != "" {
		t.Fatal(str)
	}
}

func (c asyncTestCase) complexDep(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c")
	r.add("b", "c", "d")
	r.add("c")
	r.add("d")
	r.add("e", "d")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}

	defer t.Log("full result:", r.String())
	if str := r.checkOrder("b"); str != "" {
		t.Fatal(str)
	}
	if str := r.checkOrder("e"); str != "" {
		t.Fatal(str)
	}
	if str := r.checkOrder("a"); str != "" {
		t.Fatal(str)
	}
}

func (c asyncTestCase) complexDep2(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c", "d")
	r.add("b", "c", "e")
	r.add("c")
	r.add("d")
	r.add("e", "d")
	r.add("f", "a", "g")
	r.add("g", "e")
	r.add("z")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}

	defer t.Log("full result:", r.String())
	names := r.ListDeps()
	for _, name := range names {
		if str := r.checkOrder(name); str != "" {
			t.Fatal(name, ":", str)
		}
	}
}

func (c asyncTestCase) skipped(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c")
	r.add("b", "d")
	r.add("c")
	r.add("d")
	r.add("e", "d")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}
	r.Skip("c")
	if err := c.Run(r); err != nil {
		t.Fatal("unexpected error:", err)
	}

	defer t.Log("full result:", r.String())
	if str := r.checkOrder("b"); str != "" {
		t.Fatal(str)
	}
	if str := r.checkOrder("e"); str != "" {
		t.Fatal(str)
	}
	if str := r.checkOrder("a"); str != "" {
		t.Fatal(str)
	}
}

func (c asyncTestCase) same(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c")
	r.add("b", "d")
	r.add("c")
	r.add("d")
	r.add("e", "d")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}

	if err := c.RunSome(r, "b"); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "d,b," {
		t.Fatal("incorrect run #1, you should also check other test cases")
	}
	r.Buffer.Reset()

	if err := c.RunSome(r, "b"); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "" {
		t.Fatal("unexpected output:", str)
	}
}

func (c asyncTestCase) some(t *testing.T) {
	r := newTool()
	r.add("a", "b", "c")
	r.add("b", "d")
	r.add("c")
	r.add("d")
	r.add("e", "d")
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}

	if err := c.RunSome(r, "b"); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "d,b," {
		t.Log(str)
		t.Log(r.groups)
		t.Fatal("incorrect run #1, you should also check other test cases")
	}
	r.Buffer.Reset()

	if err := c.RunSome(r, "a"); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if str := r.String(); str != "c,a," {
		t.Fatal("unexpected output:", str)
	}
}

func (c asyncTestCase) error(t *testing.T) {
	expect := errors.New("1")
	r := newTool()
	r.add("a", "b")
	r.addErr("b", expect)
	if err := r.Validate(); err != nil {
		r.dumpRunner(t)
		t.Fatal("test case error:", err)
	}

	if err := c.Run(r); err != expect {
		t.Log("output:", r.String())
		t.Fatal("unexpected result:", err)
	}
	if str := r.String(); str != "" {
		t.Fatal("unexpected output:", str)
	}
}

// benchmarks

func (c asyncTestCase) benchRun(b *testing.B) {
	r := newTool()
	r.add("a", "b", "c", "d")
	r.add("b", "c", "e")
	r.add("c")
	r.add("d")
	r.add("e", "d")
	r.add("f", "a", "g")
	r.add("g", "e")
	r.add("z")
	r.Validate()
	ctx := context.TODO()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, t := range r.tasks {
			t.state = pending
			t.err = nil
		}
		c.run(r.Runner)()(ctx)
	}
}
