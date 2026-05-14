// Package streamer provides incremental, streaming comparison of .env files.
//
// Unlike the batch diff pipeline, streamer reads files line by line and emits
// diff.Result values through a channel as soon as each key is resolved. This
// makes it suitable for large files or pipelines where low memory footprint is
// preferred over random access.
//
// Basic usage:
//
//	ch := make(chan streamer.Event, 32)
//	streamer.Stream(".env.staging", ".env.production", nil, ch)
//	for e := range ch {
//		if e.Err != nil {
//			log.Fatal(e.Err)
//		}
//		fmt.Println(e.Result.Key, e.Result.Status)
//	}
package streamer
