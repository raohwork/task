// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/raohwork/task/action"
)

// Req creates an [action.Data] of empty http client request. The context is ignored.
func Req() action.Data[*http.Request] {
	return func(_ context.Context) (*http.Request, error) {
		return http.NewRequest("", "", nil)
	}
}

// Request is like [Req], but with method and url prefilled. The context is ignored.
func Request(method, uri string) action.Data[*http.Request] {
	return func(_ context.Context) (*http.Request, error) {
		return http.NewRequest(method, uri, nil)
	}
}

// SetMethod creates an [action.Converter2] to set request method of a request.
//
// Common usage: request = request.Then(SetMethod().By(http.MethodPost))
func SetMethod() action.Converter2[string, *http.Request, *http.Request] {
	return func(_ context.Context, v string, r *http.Request) (*http.Request, error) {
		r.Method = v
		return r, nil
	}
}

// SetURL creates an [action.Converter2] to set url of a request.  It's suggested to
// write tour own converter to setup request at once.
//
// Common usage: request = request.Then(SetURL().By(myurl))
func SetURL() action.Converter2[*url.URL, *http.Request, *http.Request] {
	return func(_ context.Context, v *url.URL, r *http.Request) (*http.Request, error) {
		r.URL = v
		return r, nil
	}
}

// SetHeader creates an [action.Converter2] to set a header to the request. It's
// suggested to write tour own converter to setup request at once.
//
// Common usage: request = request.Then(SetHeader("Content-Type").By("text/json"))
func SetHeader(k string) action.Converter2[string, *http.Request, *http.Request] {
	return func(_ context.Context, v string, r *http.Request) (*http.Request, error) {
		r.Header.Set(k, v)
		return r, nil
	}
}

// AddHeader creates an [action.Converter2] to add a header to the request. It's
// suggested to write tour own converter to setup request at once.
//
// Common usage: request = request.Then(AddHeader("Content-Type").By("text/json"))
func AddHeader(k string) action.Converter2[string, *http.Request, *http.Request] {
	return func(_ context.Context, v string, r *http.Request) (*http.Request, error) {
		r.Header.Add(k, v)
		return r, nil
	}
}

// AddCookie creates an [action.Converter2] to add a cookie to the request. It's
// suggested to write tour own converter to setup request at once.
func AddCookie() action.Converter2[*http.Cookie, *http.Request, *http.Request] {
	return func(_ context.Context, c *http.Cookie, r *http.Request) (*http.Request, error) {
		r.AddCookie(c)
		return r, nil
	}
}

// UseBody is shortcut to SetBody[some_type]().From(body)
func UseBody[T io.ReadCloser](body action.Data[T]) action.Converter[*http.Request, *http.Request] {
	return SetBody[T]().From(body)
}

// ApplyBody is shortcut to SetBody[some_type]().By(body)
func ApplyBody[T io.ReadCloser](body T) action.Converter[*http.Request, *http.Request] {
	return SetBody[T]().By(body)
}

// UseBodyReader is shortcut to SetBodyReader[some_type]().From(body)
func UseBodyReader[T io.Reader](body action.Data[T]) action.Converter[*http.Request, *http.Request] {
	return SetBodyReader[T]().From(body)
}

// ApplyBodyReader is shortcut to SetBodyReader[some_type]().By(body)
func ApplyBodyReader[T io.Reader](body T) action.Converter[*http.Request, *http.Request] {
	return SetBodyReader[T]().By(body)
}

// SetBodyReader creates a converter to set a reader as request body, preventing
// http client from closing it. Useful when using os.File as body.
func SetBodyReader[T io.Reader]() action.Converter2[T, *http.Request, *http.Request] {
	return func(_ context.Context, body T, r *http.Request) (*http.Request, error) {
		r.Body = io.NopCloser(body)
		return r, nil
	}
}

// SetBody creates a converter to set the `Body` of a request.
func SetBody[T io.ReadCloser]() action.Converter2[T, *http.Request, *http.Request] {
	return func(_ context.Context, body T, r *http.Request) (*http.Request, error) {
		r.Body = body
		return r, nil
	}
}

// SetContentLength creates a converter to update the request. It's suggested to
// write tour own converter to setup request at once.
//
// Common usage: request = request.Then(SetContentLength().By(len(body)))
func SetContentLength() action.Converter2[int64, *http.Request, *http.Request] {
	return func(_ context.Context, l int64, r *http.Request) (*http.Request, error) {
		r.ContentLength = l
		return r, nil
	}
}

// ReadResp creates an [action.Converter2] to convert client request into response
// by sending it using an http client.
//
// The context is applied to the request before sending.
//
// Common usage: resp, err := ReadResp().By(client).From(request).Get(ctx)
func ReadResp() action.Converter2[*http.Client, *http.Request, *http.Response] {
	return func(ctx context.Context, cl *http.Client, req *http.Request) (*http.Response, error) {
		return cl.Do(req.WithContext(ctx))
	}
}

// GetResp creates an [action.Converter] to convert client request into response by
// sending it using [http.DefaultClient].
//
// The context is applied to the request before sending.
//
// Common usage: resp, err := GetResp().From(request).Get(ctx)
func GetResp() action.Converter[*http.Request, *http.Response] {
	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(req.WithContext(ctx))
	}
}

// ReadBody creates an [action.Converter] to read the body of a response.
func ReadBody() action.Converter[*http.Response, []byte] {
	return func(_ context.Context, resp *http.Response) ([]byte, error) {
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}
}
