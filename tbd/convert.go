// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tbd

import (
	"context"
	"sync"
)

type convert[S, D any] struct {
	source TBD[S]
	f      func(S) (D, error)
	once   sync.Once
	val    D
	e      error
}

func (c *convert[S, D]) Resolved() <-chan struct{} { return c.source.Resolved() }
func (c *convert[S, D]) Get(ctx context.Context) (D, error) {
	c.once.Do(func() {
		s, e := c.source.Get(ctx)
		if e != nil {
			c.e = e
			return
		}
		c.val, c.e = c.f(s)
	})
	return c.val, c.e
}

func EzConvert[S, D any](s TBD[S], f func(S) D) TBD[D] {
	return Convert(s, func(s S) (D, error) {
		return f(s), nil
	})
}

// Convert creates a TBD that will be resolved by converting src to new value.
//
// Returned TBD will be resolved when you get the value after src is resolved.
//
// # Improtant Note
//
// In general, chainning TBD (like [Convert]) is a bad idea: Returned TBD is
// double-locked: it self is locked and the src is also locked. The performance and
// memory usage is much more than using forge.Chain. You may write following code
// instead to get better performance:
//
//	forge.Convert(src.Get, f).TBD()
//
// If f is simple and fast, you could use pure Generator:
//
//	forge.Convert(src.Get, f)
//
// ONLY USE THIS IF f IS SIMPLE AND FAST, because code above calls f every time
// you get value from returned Generator. To fix this, use cached generator:
//
//	forge.Cache(forge.Convert(src.Get, f))
func Convert[S, D any](src TBD[S], f func(S) (D, error)) TBD[D] {
	ret, res, rej := New[D]()
	return Bind(ret, func(ctx context.Context) error {
		s, err := src.Get(ctx)
		if err != nil {
			rej(err)
			return err
		}

		d, err := f(s)
		if err != nil {
			rej(err)
			return err
		}

		res(d)
		return nil
	})
}
