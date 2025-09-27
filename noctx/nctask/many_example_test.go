// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nctask

import (
	"errors"
	"fmt"
)

func ExampleIter() {
	e := errors.New("err")
	a := func() error { fmt.Println("a"); return nil }
	b := func() error { fmt.Println("b"); return e }
	c := func() error { fmt.Println("c"); return nil }

	err := Iter(a, b, c).Run()
	if err != e {
		fmt.Println("unexpected error:", err)
		return
	}

	// output: a
	// b
}
