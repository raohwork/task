// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"sync"

	"github.com/raohwork/task/tbd"
)

type _bindTBD[T any] struct {
	g    Generator[T]
	once sync.Once
	ch   chan struct{}
	data T
	err  error
}

func (b *_bindTBD[T]) Resolved() <-chan struct{} { return b.ch }
func (b *_bindTBD[T]) Get(ctx context.Context) (T, error) {
	b.once.Do(func() {
		b.data, b.err = b.g.Run(ctx)
		close(b.ch)
	})
	return b.data, b.err
}

// TBD creates a binded [tbd.TBD] that resolved by g.
//
// It's optimized for performance and semantically identical to following code:
//
//	ret, resolve := tbd.Create[T]()
//	return tbd.Bind(ret, func(ctx context.Context) error {
//		return resolve(g.Run(ctx))
//	}
func (g Generator[T]) TBD() tbd.TBD[T] {
	return &_bindTBD[T]{g: g, ch: make(chan struct{})}
}

// Cached wraps g to cache the result, and reuse it in later call without running g.
func Cached[T any](g Generator[T]) Generator[T] {
	return g.TBD().Get
}
