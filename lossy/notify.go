// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lossy

// Notifier is a channel-based lossy Waiter/Notifier. Zero value is not usable.
type Notifier struct {
	ch chan chan struct{}
}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	ret := &Notifier{
		ch: make(chan chan struct{}, 1),
	}
	ret.ch <- make(chan struct{})
	return ret
}

// Notify closes current channel, and creates new one for new waiter.
func (n *Notifier) Notify() {
	ch := <-n.ch
	close(ch)
	n.ch <- make(chan struct{})
}

// Wait retrieves current channel, which will be closed by next Notify.
func (n *Notifier) Wait() <-chan struct{} {
	ret := <-n.ch
	n.ch <- ret
	return ret
}
