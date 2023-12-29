// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lossy

import "github.com/raohwork/task/future"

type pubsubData[T any] struct {
	*future.Future[T]
	res func(T)
	rej func(error)
}

// NewPubSub creates a new pair of [Pub]/[Sub].
//
// It provides a lossy publish/subscribe implementation.
//
// They shares an unresolved [future.Future], which is resolved (or rejected) by
// publishing. After the future is resolved, you lost further publishing before
// subcribe again.
func NewPubSub[T any]() (Pub[T], Sub[T]) {
	x := (&pubsub[T]{make(chan *pubsubData[T], 1)}).addElement()
	return x, x
}

// Pub is a [future.Future] based lossy publisher.
//
// Take a look at [Sub.Sub] for detailed info.
type Pub[T any] interface {
	// Publishs v to current subscribers.
	V(v T)
	// Publish an error to current subscribers.
	E(e error)
}

// Sub is a [future.Future] based lossy subscriber.
//
// It receives only one value once you called Sub().
type Sub[T any] interface {
	// Subscribes single value. Returned future is resolved by next Pub.V or
	// rejected by Pub.E. You'll lose further values before you subscribe again.
	Sub() *future.Future[T]
}

// Zero value is not usable, use [NewPubSub] to create one.
type pubsub[T any] struct {
	ch chan *pubsubData[T]
}

func (p *pubsub[T]) addElement() *pubsub[T] {
	fut, res, rej := future.New[T]()
	p.ch <- &pubsubData[T]{fut, res, rej}
	return p
}

// V publishes v by resolving current [future.Future].
//
// It also creates a new Future for later use.
func (p *pubsub[T]) V(v T) {
	el := <-p.ch
	el.res(v)
	p.addElement()
}

// E publishes e by rejecting current [future.Future].
//
// It also creates a new Future for later use.
func (p *pubsub[T]) E(e error) {
	el := <-p.ch
	el.rej(e)
	p.addElement()
}

// Sub subscribes by retrieving current [future.Future].
//
// The Future will be resolved (or rejected) by next [Pub.V] (or [Pub.E]). After
// that, further publishing will lost before you subscribe again.
func (p *pubsub[T]) Sub() *future.Future[T] {
	el := <-p.ch
	p.ch <- el
	return el.Future
}
