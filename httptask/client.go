// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/raohwork/task/action"
)

// Consume creates a function to consume the body of a response.
//
// See example in DecodeWith.
func Consume(body action.Data[io.ReadCloser]) func() {
	return func() {
		b, e := body.Get(context.TODO())
		if e == nil {
			io.Copy(io.Discard, b)
			b.Close()
		}
	}
}

// Reader creates an [action.Data] to help type checking.
//
// See example in DecodeWith.
func Reader[T io.Reader](i action.Data[T]) action.Data[io.Reader] {
	return action.NoCtxGet(func(r T) (io.Reader, error) {
		return r, nil
	}).From(i)
}

// ParseWith creates an [action.Converter] to parse response body.
func ParseWith[T any](f func([]byte, any) error) action.Converter[[]byte, T] {
	return action.NoCtxGet(func(b []byte) (ret T, err error) {
		err = f(b, &ret)
		return
	})
}

// Decoder abstracts some common decoder like json.Decoder.
type Decoder interface {
	Decode(any) error
}

// DecodeWith creates an [action.Converter] to decode response body.
func DecodeWith[T any, D Decoder](f func(io.Reader) D) action.Converter[io.Reader, T] {
	return action.NoCtxGet(func(r io.Reader) (ret T, err error) {
		var dec = f(r)
		err = dec.Decode(&ret)
		return
	})
}

// BodyGenFrom creates an [action.Converter] to generate request body.
//
// f is a function to generate request body from input parameter, like os.Open.
func BodyGenFrom[P any, T io.Reader](f func(P) (T, error)) action.Converter[P, io.Reader] {
	return func(_ context.Context, p P) (io.Reader, error) {
		return f(p)
	}
}

// BodyGenBy creates an [action.Data] to generate request body.
//
// f is a function to generate request body, typically a custom function.
func BodyGenBy[T io.Reader](f func() (T, error)) action.Data[io.Reader] {
	return func(_ context.Context) (io.Reader, error) {
		return f()
	}
}

// Req creates an [action.Converter] to build http client request.
//
// Request context is set by [ReadResp] or [GetResp].
func Req(method, url string) action.Converter[io.Reader, *http.Request] {
	return func(_ context.Context, body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, url, body)
	}
}

// Request is like [Req], but without body.
func Request(method, uri string) action.Data[*http.Request] {
	return Req(method, uri).By(nil)
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

// ReadBody creates an [action.Converter] to read the whole body of a response.
func ReadBody() action.Converter[*http.Response, []byte] {
	return func(_ context.Context, resp *http.Response) ([]byte, error) {
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}
}

// GetBody creates an [action.Converter] to get the body of a response.
//
// Extracted body is guaranteed to be non-nil.
func GetBody() action.Converter[*http.Response, io.ReadCloser] {
	return action.NoCtxGet(func(resp *http.Response) (io.ReadCloser, error) {
		if resp == nil || resp.Body == nil {
			return nil, errors.New("GetBody: response body does not exist")
		}
		return resp.Body, nil
	})
}
