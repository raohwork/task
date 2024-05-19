// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/raohwork/task"
	"github.com/raohwork/task/action"
)

const Addr = ":32100"

func server() task.Task {
	srv := &http.Server{Addr: Addr}
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		size, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"size": size,
		})
	})

	return Server(srv)
}

type APIResp struct {
	Size int64 `json:"size"`
}

func apiResp(_ context.Context, res *http.Response) (ret *APIResp, err error) {
	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)
	var v APIResp
	if err = json.NewDecoder(res.Body).Decode(&v); err == nil {
		ret = &v
	}
	return
}

func Example() {
	ctx, cancel := context.WithCancel(context.TODO())
	done := server().Go(ctx)
	defer func() { <-done }()
	defer cancel()

	// this will open and close the file multiple times when retrying (file is
	// closed by http client). this is simpler and has nearly no impact on
	// performance comparing to reading and sending large file over network.
	req := Request(http.MethodPost, "http://127.0.0.1"+Addr).
		Then(UseBody(action.NoCtxGet(os.Open).By("testdata/super_large_file")))
	// this will cache the file, preventing from opening again, so you have to
	// close it after retrying.
	// file := action.NoCtxGet(os.Open).By("testdata/super_large_file").Cached()
	// req := Request(http.MethodPost, "http://127.0.0.1"+Addr).
	// 	Then(UseBodyReader(file))
	resp := action.Get(apiResp).
		From(GetResp().From(req)).
		With(task.Timeout(time.Second)).
		TimedFailF(task.FixedDur(100 * time.Millisecond)).
		RetryN(3).
		// if you have cached the file, close it here.
		// Defer(file.Do(action.CloseIt).NoErr()).
		Cached()
	actual, err := resp.Get(ctx)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	fmt.Println(actual.Size)
	//output: 7
}
