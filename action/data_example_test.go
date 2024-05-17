// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
)

func ExampleData_Cached() {
	a := NoCtxUse(func() (int32, error) {
		fmt.Println("executed")
		ret := rand.Int31n(100)
		if ret < 50 {
			return 0, errors.New("50/50")
		}
		return ret, nil
	}).Cached()

	v1, e1 := a(context.TODO()) // print "executed"
	v2, e2 := a(context.TODO()) // nothing printed

	if v1 != v2 {
		fmt.Println("v1 != v2")
	}
	if e1 != e2 {
		fmt.Println("e1 != e2")
	}

	// output: executed
}

func ExampleData_Saved() {
	err := errors.New("error")
	cnt := 0
	a := NoCtxUse(func() (int, error) {
		fmt.Println("executed")
		cnt++
		if cnt < 2 {
			return cnt, err
		}
		return cnt, nil
	}).Saved()

	fmt.Println(a(context.TODO()))
	fmt.Println(a(context.TODO()))
	fmt.Println(a(context.TODO()))

	// output: executed
	// 1 error
	// executed
	// 2 <nil>
	// 2 <nil>
}
