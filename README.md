# 1BRC

| Attempt Number | Approach | Execution Time | Diff | Commit |
|-----------------|---|---|---|--|
|1| Naive Implementation: Read temperatures into a map of cities. Iterate serially over each key (city) in map to find min, max and average temperatures.| 6:13.15 | || 
|2| Evaluate each city in map concurrently using goroutines.|4:32.80|-100.35||
|3|Remove sorting float64 slices. Calculate min, max and average by iterating.|4:25.59|-7.21||