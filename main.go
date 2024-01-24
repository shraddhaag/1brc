package main

import (
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
	"strconv"
	"strings"
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

	fmt.Println(evaluate(*input))

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

type result struct {
	city string
	temp string
}

func evaluate(input string) string {
	mapOfTemp, err := readFileLineByLineIntoAMap(input)
	if err != nil {
		panic(err)
	}

	resultArr := make([]string, len(mapOfTemp))
	var count int
	for city, _ := range mapOfTemp {
		resultArr[count] = city
		count++
	}

	sort.Strings(resultArr)

	var stringsBuilder strings.Builder
	for _, i := range resultArr {
		stringsBuilder.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", i,
			round(float64(mapOfTemp[i].min)/10.0),
			round(float64(mapOfTemp[i].sum)/10.0/float64(mapOfTemp[i].count)),
			round(float64(mapOfTemp[i].max)/10.0)))
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

	mapOfTemp := make(map[string]cityTemperatureInfo)

	chanOwner := func() <-chan []string {
		resultStream := make(chan []string, 100)
		toSend := make([]string, 100)
		//  reading 100MB per request
		chunkSize := 100 * 1024 * 1024
		buf := make([]byte, chunkSize)
		var stringsBuilder strings.Builder
		stringsBuilder.Grow(500)
		var count int
		go func() {
			defer close(resultStream)
			for {
				readTotal, err := file.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						count = processReadChunk(buf, readTotal, count, &stringsBuilder, toSend, resultStream)
						break
					}
					panic(err)
				}
				count = processReadChunk(buf, readTotal, count, &stringsBuilder, toSend, resultStream)
			}
			if count != 0 {
				resultStream <- toSend[:count]
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for t := range resultStream {
		for _, text := range t {
			index := strings.Index(text, ";")
			if index == -1 {
				continue
			}
			city := text[:index]
			temp := convertStringToInt64(text[index+1:])
			if val, ok := mapOfTemp[city]; ok {
				val.count++
				val.sum += temp
				if temp < val.min {
					val.min = temp
				}

				if temp > val.max {
					val.max = temp
				}
				mapOfTemp[city] = val
			} else {
				mapOfTemp[city] = cityTemperatureInfo{
					count: 1,
					min:   temp,
					max:   temp,
					sum:   temp,
				}
			}
		}
	}
	// fmt.Println(mapOfTemp)
	return mapOfTemp, nil
}

func convertStringToInt64(input string) int64 {
	input = input[:len(input)-2] + input[len(input)-1:]
	output, _ := strconv.ParseInt(input, 10, 64)
	return output
}

func processReadChunk(buf []byte, readTotal, count int, stringsBuilder *strings.Builder, toSend []string, resultStream chan<- []string) int {
	for _, char := range buf[:readTotal] {
		if char == '\n' {
			if stringsBuilder.Len() != 0 {
				toSend[count] = stringsBuilder.String()
				stringsBuilder.Reset()
				count++

				if count == 100 {
					count = 0
					localCopy := make([]string, 100)
					copy(localCopy, toSend)
					resultStream <- localCopy
				}
			}
		} else {
			stringsBuilder.WriteByte(char)
		}
	}

	return count
}

func round(x float64) float64 {
	rounded := math.Round(x * 10)
	if rounded == -0.0 {
		return 0.0
	}
	return rounded / 10
}
