// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package ncrated

import (
	"fmt"
	"time"

	"github.com/raohwork/task/noctx/nctask"
	"golang.org/x/time/rate"
)

func Example() {
	t := func() error { return nil }
	timed := nctask.Task(t).Timed(time.Second)
	rl := Task(rate.NewLimiter(rate.Every(time.Second), 1), t)

	begin := time.Now()
	timed.Run() // run t, wait a second
	timed.Run() // run t, wait a second
	timed.Run() // run t, wait a second
	fmt.Printf("timed task: elapsed %d seconds\n", time.Since(begin)/time.Second)

	begin = time.Now()
	rl.Run() // run t
	rl.Run() // wait a second, run t
	rl.Run() // wait a second, run t
	fmt.Printf("ratelimited task: elapsed %d seconds\n", time.Since(begin)/time.Second)

	// output:timed task: elapsed 3 seconds
	// ratelimited task: elapsed 2 seconds
}
