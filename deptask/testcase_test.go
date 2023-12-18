// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/raohwork/task"
)

var nop = task.Task(func(_ context.Context) error { return nil })

func checkArr(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, x := range a {
		if b[idx] != x {
			return false
		}
	}

	return true
}

func checkArr2(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	m := map[string]bool{}
	for _, x := range a {
		m[x] = true
	}

	for _, x := range b {
		if !m[x] {
			return false
		}
	}

	return true
}

type testCase struct {
	creator func() Runner
}

func (c testCase) TestAll(t *testing.T) {
	t.Run("check missing", c.Missing)
	t.Run("check cyclic", c.Cyclic)
	t.Run("run order", c.RunOrder)
	t.Run("run some order", c.RunSomeOrder)
}

func (c testCase) Missing(t *testing.T) {
	r := c.creator()
	r.Add("a", nop, "b")
	if err := r.Validate(); !errors.Is(err, ErrMissing("b")) {
		t.Fatalf("missing is not reported: %+v", err)
	}

	r = c.creator()
	r.Add("a", nop)
	r.Add("b", nop, "a", "c")
	if err := r.Validate(); !errors.Is(err, ErrMissing("c")) {
		t.Fatalf("missing is not reported: %+v", err)
	}

	r = c.creator()
	r.Add("a", nop)
	r.Add("b", nop, "a")
	if err := r.Validate(); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
}

func (c testCase) Cyclic(t *testing.T) {
	r := c.creator()
	r.Add("a", nop)
	r.Add("b", nop, "c")
	r.Add("c", nop, "d")
	r.Add("d", nop, "b", "e")
	r.Add("e", nop, "f")
	r.Add("f", nop, "d")
	if r.Validate() != ErrCyclic {
		t.Fatal("cyclic is not detected")
	}

	r = c.creator()
	r.Add("a", nop)
	r.Add("b", nop, "c")
	r.Add("c", nop, "d")
	r.Add("d", nop, "b")
	r.Add("e", nop, "f")
	r.Add("f", nop, "g")
	r.Add("g", nop, "e")
	if r.Validate() != ErrCyclic {
		t.Fatal("cyclic is not detected")
	}

	r = c.creator()
	r.Add("a", nop)
	r.Add("b", nop, "c")
	r.Add("c", nop, "d")
	r.Add("d", nop, "b")
	r.Add("e", nop)
	if r.Validate() != ErrCyclic {
		t.Fatal("cyclic is not detected")
	}

	r = c.creator()
	r.Add("a", nop, "b", "e")
	r.Add("b", nop, "c", "e")
	r.Add("c", nop, "d", "e")
	r.Add("d", nop, "e")
	r.Add("e", nop)
	if e := r.Validate(); e != nil {
		t.Fatal("cyclic is incorrectly detected")
	}
}

func checkDep(deps map[string][]string, output []string) bool {
	for n, tdeps := range deps {
		found := false
		for _, cur := range output {
			if !found {
				if cur == n {
					found = true
				}
				continue
			}

			for _, dep := range tdeps {
				if cur == dep {
					return false
				}
			}
		}
	}

	return true
}

type orderHelper struct {
	r      Runner
	output []string
	lock   sync.Mutex
}

func (h *orderHelper) body(s string) task.Task {
	return task.Task(func(_ context.Context) error {
		h.lock.Lock()
		defer h.lock.Unlock()
		h.output = append(h.output, s)
		return nil
	})
}

func newHelper(creator func() Runner, deps map[string][]string) *orderHelper {
	ret := &orderHelper{
		r: creator(),
	}
	for name, dep := range deps {
		ret.r.Add(name, ret.body(name), dep...)
	}
	return ret
}

func (c *testCase) run(deps map[string][]string, run func(Runner) error) func(*testing.T) {
	return func(t *testing.T) {
		h := newHelper(c.creator, deps)

		if err := run(h.r); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		t.Logf("deps: %+v", deps)
		t.Logf("output: %+v", h.output)
		if !checkDep(deps, h.output) {
			t.Fatal("unexpected order")
		}
	}
}

func (c *testCase) runOrder(run func(r Runner) error) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("1 dep 1", c.run(map[string][]string{
			"a": {},
			"b": {"a"},
			"c": {"b"},
		}, run))

		t.Run("1 dep many", c.run(map[string][]string{
			"a": {},
			"b": {},
			"c": {"b", "a"},
		}, run))

		t.Run("many dep 1", c.run(map[string][]string{
			"a": {},
			"b": {"a"},
			"c": {"a"},
		}, run))

		t.Run("no dep", c.run(map[string][]string{
			"a": {},
			"b": {},
			"c": {},
		}, run))

		t.Run("many dep many", c.run(map[string][]string{
			"a": {},
			"b": {},
			"c": {},
			"d": {"a"},
			"e": {"b"},
			"f": {"d", "b"},
			"g": {"e", "c"},
			"h": {"f", "g"},
			"i": {"c"},
		}, run))
	}
}

func (c testCase) RunOrder(t *testing.T) {
	ctx := context.Background()
	t.Run("sync", c.runOrder(func(r Runner) error { return r.RunSync(ctx) }))
	t.Run("async", c.runOrder(func(r Runner) error { return r.Run(ctx) }))
}

func (c testCase) RunSomeOrder(t *testing.T) {
	ctx := context.Background()
	t.Run("sync", c.runSomeOrder(func(r Runner, name ...string) error {
		return r.RunSomeSync(ctx, name...)
	}))
	t.Run("async", c.runSomeOrder(func(r Runner, name ...string) error {
		return r.RunSome(ctx, name...)
	}))
}

func (c testCase) runSomeOrder(run func(Runner, ...string) error) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("1 dep 1", c.runSome1Dep1(run))
		t.Run("1 dep many", c.runSome1DepMany(run))
	}
}

func (c testCase) runSome1Dep1(run func(Runner, ...string) error) func(*testing.T) {
	return func(t *testing.T) {
		h := newHelper(c.creator, map[string][]string{
			"a": {},
			"b": {"a"},
			"c": {"b"},
			"d": {"c"},
			"e": {"d"},
		})

		if err := run(h.r, "c"); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !checkArr(h.output, []string{"a", "b", "c"}) {
			t.Fatalf("first run: unexpected result: %+v", h.output)
		}

		h.output = nil
		if err := run(h.r, "c"); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !checkArr(h.output, []string{}) {
			t.Fatalf("re-run: unexpected result: %+v", h.output)
		}

		h.output = nil
		if err := run(h.r, "e"); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !checkArr(h.output, []string{"d", "e"}) {
			t.Fatalf("last run: unexpected result: %+v", h.output)
		}
	}
}

func (c testCase) runSome1DepMany(run func(Runner, ...string) error) func(*testing.T) {
	return func(t *testing.T) {
		h := newHelper(c.creator, map[string][]string{
			"a": {},
			"b": {},
			"c": {},
			"d": {},
			"e": {"a", "b"},
			"f": {"c", "d"},
		})

		if err := run(h.r, "e"); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !checkArr2(h.output, []string{"a", "b", "e"}) {
			t.Fatalf("first run: unexpected result: %+v", h.output)
		}
		if !checkDep(map[string][]string{
			"a": {},
			"b": {},
			"e": {"a", "b"},
		}, h.output) {
			t.Fatalf("first run: unexpected order: %+v", h.output)
		}

		h.output = nil
		if err := run(h.r, "e"); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !checkArr2(h.output, []string{}) {
			t.Fatalf("re-run: unexpected result: %+v", h.output)
		}

		h.output = nil
		if err := run(h.r, "f"); err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !checkArr2(h.output, []string{"c", "d", "f"}) {
			t.Fatalf("last run: unexpected result: %+v", h.output)
		}
		if !checkDep(map[string][]string{
			"c": {},
			"d": {},
			"f": {"c", "d"},
		}, h.output) {
			t.Fatalf("last run: unexpected order: %+v", h.output)
		}
	}
}
