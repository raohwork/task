// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/raohwork/task/forge"
	"github.com/raohwork/task/tbd"
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

// ContentType is a shortcut so set Content-Type header.
func (r ReqGen) ContentType(typ string) ReqGen {
	return r.SetHeader("Content-Type", typ)
}

// URLEncoded is a shortcut so set Content-Type header.
func (r ReqGen) URLEncoded() ReqGen {
	return r.ContentType("application/x-www-form-urlencoded")
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

// Location sets request url from string.
func (r ReqGen) Location(locGen tbd.TBD[string]) ReqGen {
	return r.URL(tbd.Convert(locGen, url.Parse))
}

// URL sets request url.
func (r ReqGen) URL(urlGen tbd.TBD[*url.URL]) ReqGen {
	return func(ctx context.Context) (ret *http.Request, err error) {
		ret, err = r(ctx)
		if err != nil {
			return
		}

		u, err := urlGen.Get(ctx)
		if err != nil {
			return
		}
		ret.URL = u
		return
	}
}

// Body sets the request body to request.
func (r ReqGen) Body(bodyGen tbd.TBD[io.ReadCloser]) ReqGen {
	return func(ctx context.Context) (req *http.Request, err error) {
		req, err = r(ctx)
		if err != nil {
			return
		}

		body, err := bodyGen.Get(ctx)
		if err != nil {
			return
		}

		req.Body = body
		return
	}
}

// GetBody sets the request body and [http.Request.GetBody] to request.
func (r ReqGen) GetBody(bodyGen forge.Generator[io.ReadCloser]) ReqGen {
	return func(ctx context.Context) (req *http.Request, err error) {
		req, err = r(ctx)
		if err != nil {
			return
		}

		body, err := bodyGen(ctx)
		if err != nil {
			return
		}

		req.Body = body
		req.GetBody = bodyGen.Tiny
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
//
// It generates a request with empty body by default, use [ReqGen.Body] or
// [ReqGen.GetBody] to set a body. You might want to take a look at
// [forge.StringReader], [forge.BytesReader], [forge.OpenFile] and [forge.FsFile]
// to save your life.
func NewRequest[T io.Reader](method, url string) ReqGen {
	return func(ctx context.Context) (ret *http.Request, err error) {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}
}
