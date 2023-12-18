// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lossy

import "github.com/raohwork/task/future"

type pubsub[T any] struct {
	*future.Future[T]
	res func(T)
}

// PubSub provides a lossy publish/subscribe implementation.
//
// It always holds an unresolved [future.Future], which is resolved (or rejected) by
// publishing, for subscripting. After the Future is resolved, you lost further
// publishing before subcribes again.
//
// Zero value is not usable, use [NewPubSub] to create one.
type PubSub[T any] struct {
	ch chan *pubsub[T]
}

// NewPubSub creates a new PubSub.
func NewPubSub[T any]() *PubSub[T] {
	return (&PubSub[T]{make(chan *pubsub[T], 1)}).addElement()
}

func (p *PubSub[T]) addElement() *PubSub[T] {
	fut, res, _ := future.New[T]()
	p.ch <- &pubsub[T]{fut, res}
	return p
}

// Pub publishes v by resolving current [future.Future].
//
// It also creates a new Future for later use.
func (p *PubSub[T]) Pub(v T) {
	el := <-p.ch
	el.res(v)
	p.addElement()
}

// Sub subscribes by retrieving current [future.Future].
//
// The Future will be resolved (or rejected) by next Pub. After that, further
// publishing will lost before you subscribe again.
func (p *PubSub[T]) Sub() *future.Future[T] {
	el := <-p.ch
	p.ch <- el
	return el.Future
}
