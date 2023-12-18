// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"io"
)

// Copy wraps [io.Copy] into a cancellable task. Cancelling context will close src.
func Copy(dst io.Writer, src io.ReadCloser) Task {
	return func(ctx context.Context) (err error) {
		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				src.Close()
			case <-done:
			}
		}()

		defer func() { close(done) }()
		_, err = io.Copy(dst, src)
		return
	}
}
