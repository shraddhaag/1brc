# 1BRC

1️⃣🐝🏎️ [The One Billion Row Challenge](https://github.com/gunnarmorling/1brc) -- A fun exploration of how quickly 1B rows from a text file can be aggregated. The challenge was primarily foces on Java but I decided to solve it in Golang! 

I wrote a detailed blog about my implementation approach, you can check it out [here](https://www.bytesizego.com/blog/one-billion-row-challenge-go). 
## Record of iterations

Final implementation approach looks like this: 

![final iteration visualised](/excalidraw/final-iteration.png)

Here is a more detailed record of each individual iteration:

| Attempt Number | Approach | Execution Time | Diff | Commit |
|-----------------|---|---|---|--|
|0| Naive Implementation: Read temperatures into a map of cities. Iterate serially over each key (city) in map to find min, max and average temperatures.| 6:13.15 | || 
|1| Evaluate each city in map concurrently using goroutines.|4:32.80|-100.35| [8bd5f43](https://github.com/shraddhaag/1brc/commit/8bd5f437e8cc231e3ee18348b83f4dc694137546)|
|2|Remove sorting float64 slices. Calculate min, max and average by iterating.|4:25.59|-7.21|[830e5df](https://github.com/shraddhaag/1brc/commit/830e5dfacff9fb7a41d12027e21399736bc34701)|
|3|Decouple reading and processing of file content. A buffered goroutine is used to communicate between the two processes.|5:22.83|+57.24|[2babf7d](https://github.com/shraddhaag/1brc/commit/2babf7dda72d92c72722b220b8b663e747075bd7)|
|4|Instead of sending each line to the channel, now sending 100 lines chunked together. Also, to minimise garbage collection, not freeing up memory when resetting a slice. |3:41.76|-161.07|[b7b1781](https://github.com/shraddhaag/1brc/commit/b7b1781f58fd258a06940bd6c05eb404c8a14af6)|
|5|Read file in chunks of 100 MB instead of reading line by line. |3:32.62|-9.14|[c26fea4](https://github.com/shraddhaag/1brc/commit/c26fea40019552a7e4fc1c864236f433b1b686f0)|
|6|Convert temperature from `string` to `int64`, process in `int64` and convert to `float64` at the end. |2:51.50|-41.14|[7812da4](https://github.com/shraddhaag/1brc/commit/7812da4d0be07dd4686d5f9b9df1e93b08cd0dd1)|
|7|In the city <> temperatures map, replaced the value for each key (city) to preprocessed min, max, count and sum of all temperatures instead of storing all recorded temperatures for the city.|1:39.81|-71.79|[e5213a8](https://github.com/shraddhaag/1brc/commit/e5213a836b17bec0a858474a11f07c902e724bba)|
|8|Use producer consumer pattern to read file in chunks and process the chunks in parallel.|1:43.82|+14.01|[067f2a4](https://github.com/shraddhaag/1brc/commit/067f2a44c0d6b3bb7cc073639364f733bce09e3e)|
|9|Reduce memory allocation by processing each read chunk into a map. Result channel now can collate the smaller processed chunk maps.|0:28.544|-75.286|[d4153ac](https://github.com/shraddhaag/1brc/commit/d4153ac7a841170a5ceee47d930e97738b5a19f6)|
|10|Avoid string concatenation overhead by not reading the decimal point when processing city temperature.|0:24.571|-3.973|[90f2fe1](https://github.com/shraddhaag/1brc/commit/90f2fe121f454f3f1b5cdaeaaebe639bb86d4578)|
|11|Convert byte slice to string directly instead of using a `strings.Builder`.|0:18.910|-5.761|[88bb6da](https://github.com/shraddhaag/1brc/commit/88bb6da8b85424d46a8c836f3c35a49466df1ea4)|
|12|Replace `strconv.ParseInt` with a custom `string` to `int` parser.|0:14.008|-4.902|[17d575f](https://github.com/shraddhaag/1brc/commit/17d575fd0f143aed18d285713d030a5b52b478df)|
|13|Reduce map access calls when constructing final result string.|0:12.017|-1.9991||