// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/raohwork/task"
	"github.com/raohwork/task/forge"
)

func newServer() task.Task {
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", `attachment; filename="ret.txt"`)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(1024*1024))

		buf := make([]byte, 1024*1024)
		for idx := range buf {
			buf[idx] = '0'
		}
		time.Sleep(100 * time.Millisecond)
		http.ServeContent(w, r, "ret.txt", time.Now(), bytes.NewReader(buf))
	})
	mux.HandleFunc("/fast", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fast"))
	})
	return Server(&http.Server{
		Addr:    "127.0.0.1:9487",
		Handler: mux,
	}, func(ctx context.Context) (context.Context, func()) {
		ret, stop := context.WithCancel(ctx)
		stop()
		return ret, func() {}
	})
}

func ExampleReqGen_Do_incorrect() {
	globalCtx, stop := context.WithCancel(context.TODO())
	defer stop()

	// Start a http server at http://127.0.0.1:9487 for test
	//
	// It provides an endpoint "/slow" which will write response body after
	// waiting 100ms.
	done := newServer().Go(globalCtx)
	time.Sleep(10 * time.Millisecond)

	// this is identical to
	//
	//     func a(ctx context.Context) (*http.Response, error) {
	//         ctx, cancel := context.WithTimeout(time.Minute)
	//         defer cancel()
	//         req, err := http.NewRequestWithContext(
	//             ctx,
	//             http.MethodGet,
	//             "http://127.0.0.1:9487/slow",
	//             nil,
	//         )
	//         if err != nil {
	//             return nil, err
	//         }
	//         return http.DefaultClient.Do(req)
	//     }
	//     resp, err := a(globalCtx)
	//
	// Context is canceled before you actually read the body.
	resp, err := NewRequest(
		http.MethodGet, "http://127.0.0.1:9487/slow",
	).Do().With(task.Timeout(time.Minute)).Run(globalCtx)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	defer resp.Body.Close()
	// simulates some processing work or network latency
	time.Sleep(200 * time.Millisecond)

	// now read body after context is canceled, causing error
	_, err = io.ReadAll(resp.Body)
	fmt.Println(errors.Is(err, context.Canceled))

	// stop http server
	stop()
	<-done

	//output: true
}

func ExampleReqGen_Do_correct() {
	globalCtx, stop := context.WithCancel(context.TODO())
	defer stop()

	// Start a http server at http://127.0.0.1:9487 for test
	//
	// It provides an endpoint "/slow" which will write response body after
	// waiting 100ms.
	done := newServer().Go(globalCtx)
	time.Sleep(10 * time.Millisecond)

	extractBody := func(r *http.Response) ([]byte, error) {
		defer r.Body.Close()
		return io.ReadAll(r.Body)
	}

	// It is identical to
	//
	//     func a(ctx context.Context) ([]byte, error) {
	//         ctx, cancel := context.WithTimeout(time.Minute)
	//         defer cancel()
	//         req, err := http.NewRequestWithContext(
	//             ctx,
	//             http.MethodGet,
	//             "http://127.0.0.1:9487/slow",
	//             nil,
	//         )
	//         if err != nil {
	//             return nil, err
	//         }
	//         resp, err := http.DefaultClient.Do(req)
	//         if err != nil {
	//             return nil, err
	//         }
	//         defer resp.Body.Close()
	//         return io.ReadAll(resp.Body)
	//     }
	//     buf, err := a(globalCtx)
	//
	// Context is canceled after body is consumed.
	buf, err := forge.Convert(
		NewRequest(http.MethodGet, "http://127.0.0.1:9487/slow").Do(),
		extractBody,
	).With(task.Timeout(time.Minute)).Run(globalCtx)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	// simulates some processing work or network latency
	time.Sleep(200 * time.Millisecond)

	fmt.Println(err)
	fmt.Println(len(buf))

	// stop http server
	stop()
	<-done

	//output: <nil>
	// 1048576
}
