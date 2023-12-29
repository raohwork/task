// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"io"
	"net/http"

	"github.com/raohwork/task/forge"
)

// ReqGen is a [forge.Generator] which generates HTTP request.
type ReqGen forge.Generator[*http.Request]

// AddCookie adds c to generated request.
func (r ReqGen) AddCookie(c *http.Cookie) ReqGen {
	return r.Update(func(i *http.Request) *http.Request {
		i.AddCookie(c)
		return i
	})
}

// SetHeader overwrites HTTP header to generated request.
func (r ReqGen) SetHeader(k, v string) ReqGen {
	return r.Update(func(i *http.Request) *http.Request {
		i.Header.Set(k, v)
		return i
	})
}

// AddHeader adds HTTP header to generated request.
func (r ReqGen) AddHeader(k, v string) ReqGen {
	return r.Update(func(i *http.Request) *http.Request {
		i.Header.Add(k, v)
		return i
	})
}

// Customize use a function to customize the generated request.
func (r ReqGen) Customize(f func(*http.Request) (*http.Request, error)) ReqGen {
	return func(ctx context.Context) (ret *http.Request, err error) {
		ret, err = r(ctx)
		if err != nil {
			return
		}

		if f != nil {
			ret, err = f(ret)
		}
		return
	}
}

// Update use a function to customize the generated request.
func (r ReqGen) Update(f func(*http.Request) *http.Request) ReqGen {
	return r.Customize(func(q *http.Request) (*http.Request, error) {
		return f(q), nil
	})
}

// Do is shortcut to DoWith(nil)
//
// You might want to use it in most case, since timeout info can be set without
// changing http client (using context). Following code set a 3 seconds timeout to
// request, send it to server, wait a second and retry for once if failed.
//
//	body := `{"a":1}`
//	resp, err := NewRequest(method, url, forge.StringReader(body)).SetHeader(
//		"Content-Type", "application/json",
//	).Do().With(
//		task.Timeout(3*time.Second),
//	).TimedFail(time.Second).RetryN(1).Run(ctx)
func (r ReqGen) Do() forge.Generator[*http.Response] {
	return r.DoWith(nil)
}

// DoWith creates a [forge.Generator] that generates [http.Response].
//
// It builds http request using r and send the request using cl to get response.
//
// Like idiom of http package, pass nil to cl will use [http.DefaultClient], or you
// might use [ReqGen.Do] for lesser key strokes.
func (r ReqGen) DoWith(cl *http.Client) forge.Generator[*http.Response] {
	if cl == nil {
		cl = http.DefaultClient
	}

	return func(ctx context.Context) (ret *http.Response, err error) {
		req, err := r(ctx)
		if err != nil {
			return
		}

		return cl.Do(req)
	}
}

// NewRequest wraps [http.NewRequestWithContext] into a [ReqGen].
func NewRequest[T io.Reader](method, url string, bodyGen forge.Generator[T]) ReqGen {
	return func(ctx context.Context) (ret *http.Request, err error) {
		body, err := bodyGen.Run(ctx)
		if err != nil {
			return
		}
		return http.NewRequestWithContext(ctx, method, url, body)
	}
}
