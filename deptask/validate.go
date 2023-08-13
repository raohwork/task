// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

func (r *smolRunner) Validate() error {
	if r.checked {
		return r.lastCheck
	}
	r.checked = true

	if err := r.checkMissing(); err != nil {
		r.lastCheck = err
		return err
	}

	r.lastCheck = r.checkCyclic()
	return r.lastCheck
}

func (r *smolRunner) checkMissing() error {
	for t := range r.tasks {
		for _, d := range r.deps[t] {
			if _, ok := r.tasks[d]; !ok {
				return ErrMissing(d)
			}
		}
	}
	return nil
}

func (r *smolRunner) checkCyclic() error {
	indeg := r.inDeg()
	arr := r.removeLeaf(indeg)
	for len(arr) > 0 && len(indeg) > 0 {
		arr = r.removeLeaf(indeg)
	}

	if len(indeg) > 0 {
		return ErrCyclic
	}

	return nil
}

func (r *smolRunner) inDeg() map[string]int {
	ret := map[string]int{}
	for n, deps := range r.deps {
		if _, ok := ret[n]; !ok {
			ret[n] = 0
		}

		for _, dep := range deps {
			ret[dep]++
		}
	}
	return ret
}

func (r *smolRunner) removeLeaf(deg map[string]int) (ret []string) {
	for name, d := range deg {
		if d == 0 {
			ret = append(ret, name)
		}
	}

	for _, n := range ret {
		delete(deg, n)
		for _, dep := range r.deps[n] {
			deg[dep]--
		}
	}

	return
}
