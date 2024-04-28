// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

func (r *Runner) validate() error {
	r.init()

	if len(r.deps) == 0 {
		return nil
	}

	degs := map[string]int{}

	// compute degree and check if missing
	for name, deps := range r.deps {
		if _, ok := degs[name]; !ok {
			degs[name] = 0
		}
		for _, dep := range deps {
			if _, ok := r.deps[dep]; !ok {
				return ErrMissing(dep)
			}
			degs[dep]++
		}
	}

	// check cyclic deps and make cache
	arr := r.removeLeaf(degs)
	r.groups = nil
	if len(arr) > 0 {
		r.groups = append(r.groups, arr)
	}
	for len(arr) > 0 && len(degs) > 0 {
		arr = r.removeLeaf(degs)
		if len(arr) > 0 {
			r.groups = append(r.groups, arr)
		}
	}

	if len(degs) > 0 {
		return ErrCyclic
	}

	return nil
}

func (r *Runner) removeLeaf(degs map[string]int) (ret []string) {
	for name, deg := range degs {
		if deg == 0 {
			ret = append(ret, name)
		}
	}

	for _, name := range ret {
		delete(degs, name)
		for _, dep := range r.deps[name] {
			degs[dep]--
		}
	}

	return ret
}
