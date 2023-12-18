// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"errors"
	"fmt"
)

func ExampleTask_RetryN() {
	ctx := context.Background()
	n := 1
	errTask := func(_ context.Context) error {
		fmt.Println(n)
		n++
		return errors.New("")
	}

	retry := Task(errTask).RetryN(2)
	retry.Run(ctx)

	// output: 1
	// 2
	// 3
}
