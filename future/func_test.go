// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package future

import (
	"context"
	"errors"
	"testing"
)

func TestResolve(t *testing.T) {
	fu := Resolve(1)
	v, err := fu.Await(context.Background())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if v != 1 {
		t.Fatal("unexpected value:", v)
	}
}

func TestReject(t *testing.T) {
	err := errors.New("")
	fu := Reject[int](err)
	v, err := fu.Await(context.Background())
	if err != err {
		t.Error("unexpected error:", err)
	}
	if v != 0 {
		t.Fatal("unexpected value:", v)
	}
}
