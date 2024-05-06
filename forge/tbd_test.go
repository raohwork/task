// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"errors"
	"fmt"

	"github.com/raohwork/task"
)

func ExampleGenerator_Once() {
	g := Tiny(func() (int, error) {
		fmt.Println("executed")
		return 1, nil
	}).Once()
	ctx := context.TODO()

	v, err := g.Run(ctx)
	fmt.Printf("val=%d ErrOnce=%v\n", v, errors.Is(err, task.ErrOnce))
	v, err = g.Run(ctx)
	fmt.Printf("val=%d ErrOnce=%v\n", v, errors.Is(err, task.ErrOnce))

	// output: executed
	// val=1 ErrOnce=false
	// val=0 ErrOnce=true
}

func ExampleGenerator_Cached() {
	g := Tiny(func() (int, error) {
		fmt.Println("executed")
		return 1, errors.New("error")
	}).Cached()
	ctx := context.TODO()

	v, err := g.Run(ctx)
	fmt.Printf("val=%d err=%v\n", v, err)
	v, err = g.Run(ctx)
	fmt.Printf("val=%d err=%v\n", v, err)

	// output: executed
	// val=1 err=error
	// val=1 err=error
}

func ExampleGenerator_Saved() {
	cnt := 0
	g := Tiny(func() (int, error) {
		fmt.Println("executed")
		cnt++
		var err error
		if cnt < 2 {
			err = errors.New("error")
		}
		return cnt, err
	}).Saved()
	ctx := context.TODO()

	v, err := g.Run(ctx)
	fmt.Printf("val=%d err=%v\n", v, err)
	v, err = g.Run(ctx)
	fmt.Printf("val=%d err=%v\n", v, err)
	v, err = g.Run(ctx)
	fmt.Printf("val=%d err=%v\n", v, err)

	// output: executed
	// val=1 err=error
	// executed
	// val=2 err=<nil>
	// val=2 err=<nil>
}
