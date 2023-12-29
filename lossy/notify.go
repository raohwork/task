// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lossy

type notifier struct {
	ch chan chan struct{}
}

// NewNotifier creates a new pair of [Notifier]/[Waiter].
//
// It provides a lossy wait/notify implementation.
//
// They shares a channel, which is closed by [Notifier.N]. After the channel is
// closed, a new channel is shared instead, so further notification is lost before
// you wait again.
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
// See [Waiter.W] for detailed info.
type Notifier interface {
	// Notifies current waiting waiters.
	N()
}

// Waiter is a channel-bases lossy waiter, coupled with [Notifier].
//
// It receives only one notification once you call W().
type Waiter interface {
	// Waits for next notification. Returned channel is closed by next
	// Notifier.N call. You'll lose further notifications before you
	// wait again.
	W() <-chan struct{}
}

// N closes current channel, and creates new one for new waiter.
func (n *notifier) N() {
	ch := <-n.ch
	close(ch)
	n.ch <- make(chan struct{})
}

// W retrieves current channel, which will be closed by next Notify.
func (n *notifier) W() <-chan struct{} {
	ret := <-n.ch
	n.ch <- ret
	return ret
}
