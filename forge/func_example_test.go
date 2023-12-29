// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"errors"
	"fmt"
)

func ExampleFixed() {
	g := Fixed(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v, err := g.Run(ctx)
	if err != nil {
		fmt.Println("Unexpected error: ", err)
		return
	}

	fmt.Println(v)
	if v != 1 {
		fmt.Println("Unexpected value: ", v)
		return
	}

	cancel()

	v, err = g.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		fmt.Println("Unexpected error: ", err)
		return
	}

	fmt.Println("context canceled")

	// output: 1
	// context canceled
}
