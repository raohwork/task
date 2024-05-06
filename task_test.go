// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"errors"
	"fmt"
)

func ExampleTask_Once() {
	t := Tiny(func() error {
		fmt.Println("executed")
		return errors.New("error")
	})
	ctx := context.TODO()
	fmt.Println(t.Run(ctx)) // error

	once := t.Once()
	fmt.Println(once.Run(ctx))                     // error
	fmt.Println(errors.Is(once.Run(ctx), ErrOnce)) // ErrOnce

	// output: executed
	// error
	// executed
	// error
	// true
}

func ExampleTask_Cached() {
	t := Tiny(func() error {
		fmt.Println("executed")
		return errors.New("error")
	})
	ctx := context.TODO()
	fmt.Println(t.Run(ctx)) // error

	cached := t.Cached()
	fmt.Println(cached.Run(ctx)) // error
	fmt.Println(cached.Run(ctx)) // error

	// output: executed
	// error
	// executed
	// error
	// error
}
