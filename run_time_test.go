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

func ExampleHelper_Timed() {
	ctx := context.Background()
	begin := time.Now()
	quickTask := F(func(_ context.Context) error {
		// simulates a quick task like computing 1+1
		fmt.Printf("quick done at +%d socond\n", time.Since(begin)/time.Second)
		return nil
	}).Timed(time.Second)
	quickTask.Run(ctx)
	fmt.Printf("quick returns at +%d second\n", time.Since(begin)/time.Second)

	begin = time.Now()
	slowTask := F(func(_ context.Context) error {
		// simulates a slow task like calling web api
		time.Sleep(2 * time.Second)
		fmt.Printf("slow done at +%d socond\n", time.Since(begin)/time.Second)
		return nil
	}).Timed(time.Second)
	slowTask.Run(ctx)
	fmt.Printf("slow returns at +%d second\n", time.Since(begin)/time.Second)

	// output: quick done at +0 socond
	// quick returns at +1 second
	// slow done at +2 socond
	// slow returns at +2 second
}

func ExampleHelper_TimedDone() {
	ctx := context.Background()
	begin := time.Now()
	doneTask := F(func(_ context.Context) error {
		// a task which always success
		return nil
	}).TimedDone(time.Second)
	doneTask.Run(ctx)
	fmt.Printf("done returns at +%d second\n", time.Since(begin)/time.Second)

	begin = time.Now()
	failTask := F(func(_ context.Context) error {
		return errors.New("a task which always fail")
	}).TimedDone(time.Second)
	failTask.Run(ctx)
	fmt.Printf("fail returns at +%d second\n", time.Since(begin)/time.Second)

	// output: done returns at +1 second
	// fail returns at +0 second
}
