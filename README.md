# 1BRC

| Attempt Number | Approach | Execution Time | Diff | Commit |
|-----------------|---|---|---|--|
|0| Naive Implementation: Read temperatures into a map of cities. Iterate serially over each key (city) in map to find min, max and average temperatures.| 6:13.15 | || 
|1| Evaluate each city in map concurrently using goroutines.|4:32.80|-100.35| [8bd5f43](https://github.com/shraddhaag/1brc/commit/8bd5f437e8cc231e3ee18348b83f4dc694137546)|
|2|Remove sorting float64 slices. Calculate min, max and average by iterating.|4:25.59|-7.21|[830e5df](https://github.com/shraddhaag/1brc/commit/830e5dfacff9fb7a41d12027e21399736bc34701)|
|3|Decouple reading and processing of file content. A buffered goroutine is used to communicate between the two processes.|5:22.83|+57.24|[2babf7d](https://github.com/shraddhaag/1brc/commit/2babf7dda72d92c72722b220b8b663e747075bd7)|
|4|Instead of sending each line to the channel, now sending 100 lines chunked together. Also, to minimise garbage collection, not freeing up memory when resetting a slice. |3:41.76|-161.07|[b7b1781](https://github.com/shraddhaag/1brc/commit/b7b1781f58fd258a06940bd6c05eb404c8a14af6)|
|5|Read file in chunks of 100 MB instead of reading line by line. |3:32.62|-9.14|[c26fea4](https://github.com/shraddhaag/1brc/commit/c26fea40019552a7e4fc1c864236f433b1b686f0)|
|6|Convert temperature from `string` to `int64`, process in `int64` and convert to `float64` at the end. |2:51.50|-41.14|[7812da4](https://github.com/shraddhaag/1brc/commit/7812da4d0be07dd4686d5f9b9df1e93b08cd0dd1)|
|7|In the city <> temperatures map, replaced the value for each key (city) to preprocessed min, max, count and sum of all temperatures instead of storing all recorded temperatures for the city.|1:39.81|-71.79||