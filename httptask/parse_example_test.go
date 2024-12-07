// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package httptask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func ExampleParseWith() {
	// for content of server(), see package example
	ctx, cancel := context.WithCancel(context.TODO())
	done := server().Go(ctx)
	defer func() { <-done }()
	defer cancel()

	req := Request(http.MethodPost, "http://localhost"+Addr+"/get")
	resp := GetResp().From(req)
	body := ReadBody().From(resp)
	data, err := ParseWith[map[string]interface{}](json.Unmarshal).
		From(body).
		Get(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
	// Output: map[size:9]
}

func ExampleDecodeWith() {
	// for content of server(), see package example
	ctx, cancel := context.WithCancel(context.TODO())
	done := server().Go(ctx)
	defer func() { <-done }()
	defer cancel()

	req := Request(http.MethodPost, "http://localhost"+Addr+"/get")
	resp := GetResp().From(req)
	body := GetBody().From(resp)
	data, err := DecodeWith[map[string]interface{}](json.NewDecoder).
		From(Reader(body)).
		Defer(Consume(body)).
		Get(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
	// Output: map[size:9]
}
