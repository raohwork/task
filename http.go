// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package task

import (
	"context"
	"net/http"
)

// HTTPServer wraps s into a task so it can shutdown gracefully when canceled.
func HTTPServer(s *http.Server, shutdown ...CtxMod) Helper {
	return FromServer(
		s.ListenAndServe,
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			for _, mod := range shutdown {
				x, c := mod(ctx)
				ctx, _ = x, c
			}

			s.Shutdown(ctx)
		},
	)
}
