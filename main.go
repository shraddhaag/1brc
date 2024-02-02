package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strings"
	"sync"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var executionprofile = flag.String("execprofile", "", "write tarce execution to `file`")
var input = flag.String("input", "", "path to the input file to evaluate")

func main() {

	flag.Parse()

	if *executionprofile != "" {
		f, err := os.Create("./profiles/" + *executionprofile)
		if err != nil {
			log.Fatal("could not create trace execution profile: ", err)
		}
		defer f.Close()
		trace.Start(f)
		defer trace.Stop()
	}

	if *cpuprofile != "" {
		f, err := os.Create("./profiles/" + *cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	evaluate(*input)

	if *memprofile != "" {
		f, err := os.Create("./profiles/" + *memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

type computedResult struct {
	city          string
	min, max, avg float64
}

func evaluate(input string) string {
	mapOfTemp, err := readFileLineByLineIntoAMap(input)
	if err != nil {
		panic(err)
	}

	resultArr := make([]computedResult, len(mapOfTemp))
	var count int
	for city, calculated := range mapOfTemp {
		resultArr[count] = computedResult{
			city: city,
			min:  round(float64(calculated.min) / 10.0),
			max:  round(float64(calculated.max) / 10.0),
			avg:  round(float64(calculated.sum) / 10.0 / float64(calculated.count)),
		}
		count++
	}

	sort.Slice(resultArr, func(i, j int) bool {
		return resultArr[i].city < resultArr[j].city
	})

	var stringsBuilder strings.Builder
	for _, i := range resultArr {
		stringsBuilder.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", i.city, i.min, i.avg, i.max))
	}
	return stringsBuilder.String()[:stringsBuilder.Len()-2]
}

type cityTemperatureInfo struct {
	count int64
	min   int64
	max   int64
	sum   int64
}

func readFileLineByLineIntoAMap(filepath string) (map[string]cityTemperatureInfo, error) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mapOfTemp := make(map[string]cityTemperatureInfo)
	resultStream := make(chan map[string]cityTemperatureInfo, 10)
	chunkStream := make(chan []byte, 15)
	chunkSize := 64 * 1024 * 1024
	var wg sync.WaitGroup

	// spawn workers to consume (process) file chunks read
	for i := 0; i < runtime.NumCPU()-1; i++ {
		wg.Add(1)
		go func() {
			for chunk := range chunkStream {
				processReadChunk(chunk, resultStream)
			}
			wg.Done()
		}()
	}

	// spawn a goroutine to read file in chunks and send it to the chunk channel for further processing
	go func() {
		buf := make([]byte, chunkSize)
		leftover := make([]byte, 0, chunkSize)
		for {
			readTotal, err := file.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				panic(err)
			}
			buf = buf[:readTotal]

			toSend := make([]byte, readTotal)
			copy(toSend, buf)

			lastNewLineIndex := bytes.LastIndex(buf, []byte{'\n'})

			toSend = append(leftover, buf[:lastNewLineIndex+1]...)
			leftover = make([]byte, len(buf[lastNewLineIndex+1:]))
			copy(leftover, buf[lastNewLineIndex+1:])

			chunkStream <- toSend

		}
		close(chunkStream)

		// wait for all chunks to be proccessed before closing the result stream
		wg.Wait()
		close(resultStream)
	}()

	// process all city temperatures derived after processing the file chunks
	for t := range resultStream {
		for city, tempInfo := range t {
			if val, ok := mapOfTemp[city]; ok {
				val.count += tempInfo.count
				val.sum += tempInfo.sum
				if tempInfo.min < val.min {
					val.min = tempInfo.min
				}

				if tempInfo.max > val.max {
					val.max = tempInfo.max
				}
				mapOfTemp[city] = val
			} else {
				mapOfTemp[city] = tempInfo
			}
		}
	}

	return mapOfTemp, nil
}

func processReadChunk(buf []byte, resultStream chan<- map[string]cityTemperatureInfo) {
	toSend := make(map[string]cityTemperatureInfo)
	var start int
	var city string

	stringBuf := string(buf)
	for index, char := range stringBuf {
		switch char {
		case ';':
			city = stringBuf[start:index]
			start = index + 1
		case '\n':
			if (index-start) > 1 && len(city) != 0 {
				temp := customStringToIntParser(stringBuf[start:index])
				start = index + 1

				if val, ok := toSend[city]; ok {
					val.count++
					val.sum += temp
					if temp < val.min {
						val.min = temp
					}

					if temp > val.max {
						val.max = temp
					}
					toSend[city] = val
				} else {
					toSend[city] = cityTemperatureInfo{
						count: 1,
						min:   temp,
						max:   temp,
						sum:   temp,
					}
				}

				city = ""
			}
		}
	}
	resultStream <- toSend
}

func round(x float64) float64 {
	rounded := math.Round(x * 10)
	if rounded == -0.0 {
		return 0.0
	}
	return rounded / 10
}

// input: string containing signed number in the range [-99.9, 99.9]
// output: signed int in the range [-999, 999]
func customStringToIntParser(input string) (output int64) {
	var isNegativeNumber bool
	if input[0] == '-' {
		isNegativeNumber = true
		input = input[1:]
	}

	switch len(input) {
	case 3:
		output = int64(input[0])*10 + int64(input[2]) - int64('0')*11
	case 4:
		output = int64(input[0])*100 + int64(input[1])*10 + int64(input[3]) - (int64('0') * 111)
	}

	if isNegativeNumber {
		return -output
	}
	return
}
