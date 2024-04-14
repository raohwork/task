// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tbd

import "context"

// CreateImpl creates a GeneralImpl and initializes it.
//
// You *MUST NOT* use it directly. See documentation of [GeneralImpl] for detail.
func CreateImpl[T any]() *GeneralImpl[T] {
	return &GeneralImpl[T]{ch: make(chan struct{})}
}

// GeneralImpl is helper to write your own TBD implementation.
//
// You *MUST NOT* use it directly as resolving it with Resolve(), Reject() or
// Determine() on a resolved GeneralImpl leads to panic, due to closing a closed
// channel. You have to ensure that it is not resolved multiple times, likely by
// protecting it with [sync.Once]. Take a look at source code of [New], [Create] or
// forge.Generaotr.TBD for example.
type GeneralImpl[T any] struct {
	ch   chan struct{}
	data T
	err  error
}

// Resolve tries to resolves it without validating if it is resolved.
//
// Resolving a resolved GeneralImpl leads to panic.
func (i *GeneralImpl[T]) Resolve(v T) error {
	i.data = v
	close(i.ch)
	return nil
}

// Reject tries to resolves it without validating if it is resolved.
//
// Rejecting a resolved GeneralImpl leads to panic.
func (i *GeneralImpl[T]) Reject(e error) error {
	i.err = e
	close(i.ch)
	return e
}

// Determine tries to resolves it without validating if it is resolved.
//
// Determining a resolved GeneralImpl leads to panic.
func (i *GeneralImpl[T]) Determine(v T, e error) error {
	i.data, i.err = v, e
	close(i.ch)
	return e
}

// Resolved implements TBD.
func (i *GeneralImpl[T]) Resolved() <-chan struct{} {
	return i.ch
}

// Get implements TBD.
func (i *GeneralImpl[T]) Get(ctx context.Context) (ret T, err error) {
	select {
	case <-i.ch:
		ret, err = i.data, i.err
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}
