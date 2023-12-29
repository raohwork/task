// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func ExampleGenerator_TimedFail() {
	ctx := context.Background()
	begin := time.Now()
	doneTask := Generator[int](func(_ context.Context) (int, error) {
		// a task which always success
		return 1, nil
	}).TimedFail(time.Second)
	doneTask.Run(ctx)
	fmt.Printf("done returns at +%d second\n", time.Since(begin)/time.Second)

	begin = time.Now()
	failTask := Generator[int](func(_ context.Context) (int, error) {
		return 0, errors.New("a task which always fail")
	}).TimedFail(time.Second)
	failTask.Run(ctx)
	fmt.Printf("fail returns at +%d second\n", time.Since(begin)/time.Second)

	// output: done returns at +0 second
	// fail returns at +1 second
}
