package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create("./profiles/" + *cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	fmt.Println(evaluate())

	if *memprofile != "" {
		f, err := os.Create("./profiles/" + *memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func evaluate() string {
	mapOfTemp, err := readFileLineByLineIntoAMap("./test_cases/measurements-1.txt")
	if err != nil {
		panic(err)
	}

	var result []string
	var wg sync.WaitGroup
	var mx sync.Mutex

	updateResult := func(input string) {
		mx.Lock()
		defer mx.Unlock()

		result = append(result, input)
	}

	for city, temps := range mapOfTemp {
		wg.Add(1)
		go func(city string, temps []float64) {
			defer wg.Done()
			sort.Float64s(temps)

			var avg float64

			for _, temp := range temps {
				avg += temp
			}

			// fmt.Println(avg, len(temps))
			avg = avg / float64(len(temps))
			// fmt.Println(avg)
			avg = math.Ceil(avg*10) / 10
			// fmt.Println(avg)

			updateResult(fmt.Sprintf("%s=%.1f/%.1f/%.1f", city, temps[0], avg, temps[len(temps)-1]))

		}(city, temps)
	}

	wg.Wait()
	sort.Strings(result)
	return strings.Join(result, ", ")
}

func readFileLineByLineIntoAMap(filepath string) (map[string][]float64, error) {
	mapOfTemp := make(map[string][]float64)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()

		index := strings.Index(text, ";")
		city := text[:index]
		temp := convertStringToFloat(text[index+1:])

		if _, ok := mapOfTemp[city]; ok {
			mapOfTemp[city] = append(mapOfTemp[city], temp)
		} else {
			mapOfTemp[city] = []float64{temp}
		}
	}

	return mapOfTemp, nil
}

func convertStringToFloat(input string) float64 {
	output, _ := strconv.ParseFloat(input, 64)
	return output
}
