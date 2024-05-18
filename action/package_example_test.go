// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"os/signal"
)

// read value from remote api, can be canceled by context.
func readData(ctx context.Context) (int, error) {
	// should add logging and error processing here
	return 1, nil
}

// create db connection.
func connectDB(dsn string) (*sql.DB, error) {
	// should add logging and error processing here
	_ = dsn
	return nil, errors.New("fake")
}

// close db
func closeDB(db *sql.DB) error {
	// should add logging and error processing here
	return db.Close()
}

// functions accepting [Data] is not recommended as it introduces extra complexity,
// but performance is barely improved.
func fasterCloseDB(db Data[*sql.DB]) func() {
	return func() {
		if conn, err := db(context.TODO()); err == nil {
			conn.Close()
		}
	}
}

// save val to db, can be canceled by context.
func saveToDB(ctx context.Context, db *sql.DB, val int) error {
	// should add logging and error processing here
	_, _, _ = ctx, db, val
	return errors.New("fake")
}

func Example() {
	const dsn = "my db dsn"
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer stop()

	// Traditional approach
	func(ctx context.Context) {
		db, err := connectDB(dsn)
		if err != nil {
			return
		}
		defer db.Close() // should add error processing

		data, err := readData(ctx)
		if err != nil {
			return
		}

		err = saveToDB(ctx, db, data)
		if err != nil {
			return
		}

		// what if you want to retry db connecting if failed?
	}(ctx)

	// Using this package, shared func should be renamed to be more descriptive
	//
	// Pros:
	//     - less if err != nil block
	//     - descriptive
	//     - easier to write testable code
	//     - useful helper to do common jobs like retrying
	//
	// Cons:
	//     - hard to trace
	//     - slower (can be ignored when doing IO task)
	//     - use more resource sometimes
	//     - need to cache the result sometimes, which uses extra sync.Once
	func(ctx context.Context) {
		db := NoCtxGet(connectDB).By(dsn).Cached()
		// to retry connecting, use this line instead
		// db := NoCtxGet(connectDB).By(dsn).RetryN(3).Cached()

		err := Do2(saveToDB).
			Use(db).
			Use(readData).
			Defer(NoCtxDo(closeDB).Use(db).NoErr()).
			// can be a little bit faster, but not worthing it
			// Defer(fasterCloseDB(db)).
			Run(ctx)
		if err != nil {
			return
		}
	}(ctx)
}
