// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package forge

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"strings"
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

// FsFile wraps [fs.FS.Open] into a tiny generator.
func FsFile(f fs.FS, name string) Generator[fs.File] {
	return Tiny(func() (fs.File, error) { return f.Open(name) })
}
