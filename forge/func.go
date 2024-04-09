// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"sync"
)

// Fixed creates a micro generator that always return same value.
func Fixed[T any](v T) Generator[T] {
	return Micro(func() T { return v })
}

// StringReader creates a micro generator which generates reader from same string.
func StringReader(str string) Generator[io.Reader] {
	return Micro(func() io.Reader { return strings.NewReader(str) })
}

// BytesReader creates a micro generator which generates reader from same byte slice.
func BytesReader(str []byte) Generator[io.Reader] {
	return Micro(func() io.Reader { return bytes.NewReader(str) })
}

// OpenFile wraps [os.Open] into a tiny generator.
func OpenFile(name string) Generator[*os.File] {
	return Tiny(func() (*os.File, error) { return os.Open(name) })
}

// Cached wraps g to cache the result, and reuse it in later call without running g.
func Cached[T any](g Generator[T]) Generator[T] {
	var (
		val  T
		once sync.Once
		err  error
	)
	return func(ctx context.Context) (T, error) {
		once.Do(func() {
			val, err = g.Run(ctx)
		})
		return val, err
	}
}
