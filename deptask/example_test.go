// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package deptask

import (
	"context"
	"fmt"

	"github.com/raohwork/task"
)

func ExampleRunner() {
	r := New()
	body := func(n string) task.Task {
		return task.Task(func(_ context.Context) error {
			fmt.Println(n)
			return nil
		})
	}
	r.MustAdd("a", body("a"))
	r.MustAdd("b", body("b"), "a")
	r.MustAdd("c", body("c"), "b")
	r.MustAdd("d", body("d"), "c")
	r.MustAdd("e", body("e"), "d")

	fmt.Println("RunSomeSync(c):")
	r.RunSomeSync(context.TODO(), "c")

	fmt.Println("RunSomeSync(e):")
	r.RunSomeSync(context.TODO(), "e")

	// output: RunSomeSync(c):
	// a
	// b
	// c
	// RunSomeSync(e):
	// d
	// e
}

func ExampleRunner_Skip() {
	r := New()
	body := func(n string) task.Task {
		return task.Task(func(_ context.Context) error {
			fmt.Println(n)
			return nil
		})
	}
	r.MustAdd("a", body("a"))
	r.MustAdd("b", body("b"), "a")
	r.MustAdd("c", body("c"), "b")
	r.MustAdd("d", body("d"), "c")
	r.MustAdd("e", body("e"), "d")
	r.Skip("s", "e", "a")

	fmt.Println("RunSomeSync(e):")
	r.RunSomeSync(context.TODO(), "e")

	// output: RunSomeSync(e):
	// b
	// c
	// d
}
