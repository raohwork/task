// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lossy

import "context"

type notifier struct {
	ch chan chan struct{}
}

// NewNotifier creates a new pair of [Notifier]/[Waiter].
//
// It provides a lossy wait/notify implementation.
//
// They shares a channel, which is closed by [Notifier.Notify]. After the channel
// is closed, a new channel is shared instead, so further notification is lost
// before you wait again.
func NewNotifier() (Notifier, Waiter) {
	ret := &notifier{
		ch: make(chan chan struct{}, 1),
	}
	ret.ch <- make(chan struct{})
	return ret, ret
}

// Notifier is a channel-based lossy notifier, coupled with [Waiter].
//
// It notifies current waiters, and prepares for next notification.
//
// See [Waiter.Wait] for detailed info.
type Notifier interface {
	// Notifies current waiting waiters.
	Notify()
	// Deprecated: Use Notify() which is identical.
	N()
}

// Waiter is a channel-bases lossy waiter, coupled with [Notifier].
//
// It receives only one notification once you call Wait().
type Waiter interface {
	// Waits for next notification. Returned channel is closed by next
	// Notifier.Notify call. You'll lose further notifications before
	// you wait again.
	Wait() <-chan struct{}
	// Deprecated: Use Wait() which is identical.
	W() <-chan struct{}

	// A cancellable wait.
	WaitCtx(context.Context) error
}

// Notify closes current channel, and creates new one for new waiter.
func (n *notifier) Notify() {
	ch := <-n.ch
	close(ch)
	n.ch <- make(chan struct{})
}

// Wait retrieves current channel, which will be closed by next Notify.
func (n *notifier) Wait() <-chan struct{} {
	ret := <-n.ch
	n.ch <- ret
	return ret
}

// N is identical to Notify.
//
// Deprecated: Use Notify() instead.
func (n *notifier) N() { n.Notify() }

// W is identical to Wait.
//
// Deprecated: Use Wait() instead.
func (n *notifier) W() <-chan struct{} { return n.Wait() }

// WaitCtx waits current channel until ctx is canceled.
func (n *notifier) WaitCtx(ctx context.Context) error {
	select {
	case <-n.Wait():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
