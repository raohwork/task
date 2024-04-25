// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"context"
	"sync"

	"github.com/raohwork/task/tbd"
)

// TBD creates a binded [tbd.TBD] that resolved by g.
//
// It's semantically identical to following code, with better performance:
//
//	ret, resolve := tbd.Create[T]()
//	return tbd.Bind(ret, func(ctx context.Context) error {
//		return resolve(g.Run(ctx))
//	}
func (g Generator[T]) TBD() tbd.TBD[T] {
	return &asTBD[T]{g: g, ch: make(chan struct{})}
}

type asTBD[T any] struct {
	g    Generator[T]
	once sync.Once
	ch   chan struct{}
	data T
	err  error
}

func (t *asTBD[T]) Resolved() <-chan struct{} { return t.ch }
func (t *asTBD[T]) Get(ctx context.Context) (T, error) {
	t.once.Do(func() {
		t.data, t.err = t.g(ctx)
	})
	return t.data, t.err
}

// Cached wraps g to cache the result, and reuse it in later call without running g.
func Cached[T any](g Generator[T]) Generator[T] {
	var (
		once sync.Once
		data T
		err  error
	)
	return func(ctx context.Context) (T, error) {
		once.Do(func() {
			data, err = g.Run(ctx)
		})
		return data, err
	}
}
