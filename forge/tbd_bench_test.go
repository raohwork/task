// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"testing"

	"github.com/raohwork/task/tbd"
)

var nopCtx = context.Background()

func BenchmarkTBDBind(b *testing.B) {
	g := G(func(_ context.Context) (int, error) {
		return 1, nil
	})

	for i := 0; i < b.N; i++ {
		x, f := tbd.Create[int]()
		resolver := func(ctx context.Context) error {
			return f(g.Run(ctx))
		}

		y := tbd.Bind(x, resolver)
		_, _ = y.Get(nopCtx)
	}
}

func BenchmarkTBDGeneraotr(b *testing.B) {
	g := G(func(_ context.Context) (int, error) {
		return 1, nil
	})

	for i := 0; i < b.N; i++ {
		_, _ = g.TBD().Get(nopCtx)
	}
}
