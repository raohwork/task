// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"errors"
	"fmt"

	"github.com/raohwork/task/tbd"
)

func ExampleGenerator_RetryN() {
	ctx := context.Background()
	n := 1
	errTask := func(_ context.Context) (int, error) {
		fmt.Println(n)
		n++
		return 0, errors.New("")
	}

	retry := G(errTask).RetryN(2)
	retry.Run(ctx)

	// output: 1
	// 2
	// 3
}

func ExampleGenerator_Retry() {
	ctx := context.Background()
	n := 1
	errTask := func(_ context.Context) (int, error) {
		fmt.Println(n)
		if n >= 3 {
			return 100, nil
		}

		n++
		return 0, errors.New("")
	}

	retry := G(errTask).Retry()
	fu := retry.Go(ctx)
	fmt.Printf("result: %d", tbd.Value(fu))

	// output: 1
	// 2
	// 3
	// result: 100
}
