// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package action is designed to write code like this:
//
//	func fetchAPI(context.Context) (*http.Response, error)
//	func parseResult(context.Context, *http.Response) (MyDataType, error) {}
//	func saveToDB(context.Context, sql.DB, MyDataType) error {}
//	func generateReport(context.Context, MyDataType) error {}
//	someData := action.Get(parseResult).
//		From(fetchAPI).
//		RetryN(3).
//		Cached()
//	err := action.Do2(saveToDB).
//		Apply(dbConn).
//		Then(generateReport).
//		Use(someData).
//		Run(ctx)
//
// This approach eliminates most error processing code in current scope.
package action
