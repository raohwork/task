// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tbd

import (
	"context"
	"sync"

	"github.com/raohwork/task"
)

type bindedTBD[T any] struct {
	TBD[T]
	task task.Task
	once sync.Once
}

func (v *bindedTBD[T]) Get(ctx context.Context) (T, error) {
	v.once.Do(func() {
		v.task.Run(ctx)
	})
	return v.TBD.Get(ctx)
}

// Bind creates a TBD that will be resolved by getting value.
func Bind[T any](v TBD[T], resolver task.Task) TBD[T] {
	return &bindedTBD[T]{TBD: v, task: resolver}
}
