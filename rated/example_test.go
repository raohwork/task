// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rated

import (
	"context"
	"fmt"
	"time"

	"github.com/raohwork/task"
)

func Example() {
	t := func(_ context.Context) error { return nil }
	timed := task.Func(t).Helper().Timed(time.Second)
	rl := Every(time.Second, task.Func(t))
	ctx := context.Background()

	begin := time.Now()
	timed.Run(ctx) // run t, wait a second
	timed.Run(ctx) // run t, wait a second
	timed.Run(ctx) // run t, wait a second
	fmt.Printf("timed task: elapsed %d seconds\n", time.Since(begin)/time.Second)

	begin = time.Now()
	rl.Run(ctx) // run t
	rl.Run(ctx) // wait a second, run t
	rl.Run(ctx) // wait a second, run t
	fmt.Printf("ratelimited task: elapsed %d seconds\n", time.Since(begin)/time.Second)

	// output:timed task: elapsed 3 seconds
	// ratelimited task: elapsed 2 seconds
}
