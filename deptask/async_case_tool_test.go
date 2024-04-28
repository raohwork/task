// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/raohwork/task"
)

type asyncTestCaseTool struct {
	*Runner
	// buffer with lock
	*sync.Mutex
	*bytes.Buffer
}

func newTool() asyncTestCaseTool {
	return asyncTestCaseTool{
		Runner: New(),
		Mutex:  &sync.Mutex{},
		Buffer: &bytes.Buffer{},
	}
}
func (c asyncTestCaseTool) add(name string, deps ...string) {
	c.Runner.Add(name, task.Micro(func() {
		c.Lock()
		defer c.Unlock()
		c.Write([]byte(name + ","))
	}), deps...)
}
func (c asyncTestCaseTool) addErr(name string, err error, deps ...string) {
	c.Runner.Add(name, task.Tiny(func() error { return err }), deps...)
}
func (c asyncTestCaseTool) dumpRunner(t *testing.T) {
	t.Log("DUMP DEPENDENCY:")
	for name, deps := range c.deps {
		t.Logf("    %s: %v", name, deps)
	}
	t.Log("DUMP STATE:")
	for name, state := range c.tasks {
		t.Logf("    %s: state=%s err=%v", name, state.state.String(), state.err)
	}
}

func (c asyncTestCaseTool) has(names ...string) string {
	if len(names) == 0 {
		return "unexpected test code: it does not provide name to has()"
	}

	arr := strings.Split(c.String(), ",")
	pos := map[string]int{}
	for idx, name := range arr {
		if _, ok := pos[name]; ok {
			return name + " is ran multiple times"
		}
		pos[name] = idx
	}

	lastPos, ok := pos[names[0]]
	if !ok {
		return names[0] + " does not exist"
	}
	for _, name := range names[1:] {
		me, ok := pos[name]
		if !ok {
			return name + " does not exist"
		}
		if me > lastPos {
			return name + " has incorrect order"
		}
		lastPos = me
	}

	return ""
}

func (c asyncTestCaseTool) checkOrder(names ...string) string {
	return checkOrder(c.Runner, c.String(), names...)
}

// check if it is ran after all deps fulfills.
// str in form of "a,b,c,d," *REMEMBER LAST COMMA*
func checkOrder(r *Runner, str string, names ...string) string {
	if len(names) == 0 {
		names = r.ListDeps()
	}
	arr := strings.Split(str, ",")
	pos := map[string]int{}
	for idx, name := range arr {
		if _, ok := pos[name]; ok {
			return name + " is ran multiple times"
		}
		pos[name] = idx
	}

	for _, name := range names {
		me, ok := pos[name]
		if !ok && r.tasks[name].state != skipped {
			return "missing output of " + name
		}
		if ok && r.tasks[name].state == skipped {
			return fmt.Sprintf("skipped %s has been ran", name)
		}
		for _, dep := range r.deps[name] {
			it, ok := pos[dep]
			if !ok && r.tasks[dep].state != skipped {
				return fmt.Sprintf("missing output or dep %s of %s", dep, name)
			}
			if ok && r.tasks[dep].state == skipped {
				return fmt.Sprintf("skipped %s has been ran", dep)
			}
			if it >= me {
				return fmt.Sprintf("%s and dep %s has incorrect order", name, dep)
			}
		}
	}

	return ""
}
