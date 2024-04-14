// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tbd

import (
	"context"
	"sync"
)

var nopCtx = context.Background()

// TBD represents a value might be computed some time later.
//
// It is a placeholder for the value. There must be a resolver who give it the
// result, no matter a value or an error. Once the result is provided, the channel
// returned by Resolved() is closed.
type TBD[T any] interface {
	// Indicates if it is resolved.
	Resolved() <-chan struct{}
	// Wait and get the value. Can be used as a generator.
	Get(context.Context) (T, error)
}

// New creates a TBD and provides 2 functions to resolve it.
func New[T any]() (TBD[T], func(T) error, func(error) error) {
	var once sync.Once
	ret := CreateImpl[T]()
	return ret,
		func(v T) error {
			once.Do(func() { ret.Resolve(v) })
			return nil
		},
		func(e error) error {
			once.Do(func() { ret.Reject(e) })
			return e
		}
}

// Create creates a TBD and provides a function to resolve it.
func Create[T any]() (TBD[T], func(T, error) error) {
	var once sync.Once
	ret := CreateImpl[T]()
	return ret,
		func(v T, e error) error {
			once.Do(func() { ret.Determine(v, e) })
			return e
		}
}
