// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package future defines [Future], a value which is determined some time in future.
//
// Future is designed for writing long-running program which doing some works
// repeatedly like crawlers. It's good to be used as a broadcaster, or as a param
// to task generator.
//
//	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
//	defer stop()
//	crawlerTask, pageFuture := downloadWebPage()
//	log.Print(task.Skip(
//		task.Iter(
//			crawlerTask,
//			parseAndSavePage(pageFuture),
//			task.Skip(
//				notifyWebDash(pageFuture),
//				cleanUpOldData(pageFuture),
//			),
//		).Timed(time.Minute).RetryN(3).Loop(),
//		task.HTTPServer(myWebDash),
//	).Run(ctx))
package future
