// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func ExampleIter() {
	e := errors.New("err")
	a := func(_ context.Context) error { fmt.Println("a"); return nil }
	b := func(_ context.Context) error { fmt.Println("b"); return e }
	c := func(_ context.Context) error { fmt.Println("c"); return nil }

	err := Iter(a, b, c).Run(context.Background())
	if err != e {
		fmt.Println("unexpected error:", err)
		return
	}

	// output: a
	// b
}

func ExampleSkip() {
	e := errors.New("err")
	a := func(ctx context.Context) error {
		if err := Sleep(time.Minute).Run(ctx); err != nil {
			fmt.Println("context canceled")
			return err
		}
		fmt.Println("a")
		return nil
	}
	b := func(_ context.Context) error { fmt.Println("b"); return e }
	c := func(ctx context.Context) error {
		if err := Sleep(time.Minute).Run(ctx); err != nil {
			fmt.Println("context canceled")
			return err
		}
		fmt.Println("c")
		return nil
	}

	err := Skip(a, b, c).Run(context.Background())
	if err != e {
		fmt.Println("unexpected error:", err)
		return
	}

	// output: b
	// context canceled
	// context canceled
}
