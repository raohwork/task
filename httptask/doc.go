// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package httptask provides some helper to wrap http server and some client job
// into task.
//
// For client job, timeout is controled by context. the timeout info should be
// applied to whole task (send request + get response + read body). For example:
//
//	resp := GetResp().
//		From(Request("GET", "http://example.com")).
//		RetryN(3).
//		Cached()
//	buf, err := ReadBody().
//		From(resp).
//		With(task.Timeout(10 * time.Second)).
//		Defer(Consume(GetBody().From(resp))).
//		Get(ctx)
package httptask
