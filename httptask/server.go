// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"net/http"

	"github.com/raohwork/task"
)

// Server wraps s into a task so it can shutdown gracefully when canceled.
//
// s.Shutdown is called with a new context modified by shutdownMods.
func Server(s *http.Server, shutdownMods ...task.CtxMod) task.Task {
	return task.FromServer(
		s.ListenAndServe,
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			for _, mod := range shutdownMods {
				x, c := mod(ctx)
				defer c()
				ctx, _ = x, c
			}

			s.Shutdown(ctx)
		},
	)
}
