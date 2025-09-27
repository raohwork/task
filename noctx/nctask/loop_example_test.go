// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nctask

import (
	"errors"
	"fmt"
)

func ExampleTask_RetryN() {
	n := 1
	errTask := func() error {
		fmt.Println(n)
		n++
		return errors.New("")
	}

	retry := Task(errTask).RetryN(2)
	retry.Run()

	// output: 1
	// 2
	// 3
}
