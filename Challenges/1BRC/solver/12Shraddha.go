package solver

import (
	"cmp"
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
	"slices"
	"strconv"
	"strings"
	"sync"
)

func Shraddha12() {
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

	fmt.Println(evaluate6())

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

func evaluate6() string {
	mapOfTemp, err := readFileLineByLineIntoAMap6("measurements.txt")
	if err != nil {
		panic(err)
	}

	var resultArr []result
	var wg sync.WaitGroup
	var mx sync.Mutex

	updateResult := func(city, temp string) {
		mx.Lock()
		defer mx.Unlock()

		resultArr = append(resultArr, result{city, temp})
	}

	for city, temps := range mapOfTemp {
		wg.Add(1)
		go func(city string, temps []int64) {
			defer wg.Done()
			var min, max, avg int64
			min, max = math.MaxInt64, math.MinInt64

			for _, temp := range temps {
				if temp < min {
					min = temp
				}

				if temp > max {
					max = temp
				}
				avg += temp
			}

			updateResult(city, fmt.Sprintf("%.1f/%.1f/%.1f", round(float64(min)/10.0), round(float64(avg)/10.0/float64(len(temps))), round(float64(max)/10.0)))
		}(city, temps)
	}

	wg.Wait()
	slices.SortFunc(resultArr, func(i, j result) int {
		return cmp.Compare(i.city, j.city)
	})

	var stringsBuilder strings.Builder
	for _, i := range resultArr {
		stringsBuilder.WriteString(fmt.Sprintf("%s=%s, ", i.city, i.temp))
	}
	return stringsBuilder.String()[:stringsBuilder.Len()-2]
}

func readFileLineByLineIntoAMap6(filepath string) (map[string][]int64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	mapOfTemp := make(map[string][]int64)

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
			if _, ok := mapOfTemp[city]; ok {
				mapOfTemp[city] = append(mapOfTemp[city], temp)
			} else {
				mapOfTemp[city] = []int64{temp}
			}
		}
	}
	return mapOfTemp, nil
}

func convertStringToInt64(input string) int64 {
	input = input[:len(input)-2] + input[len(input)-1:]
	output, _ := strconv.ParseInt(input, 10, 64)
	return output
}

func round(x float64) float64 {
	rounded := math.Round(x * 10)
	if rounded == -0.0 {
		return 0.0
	}
	return rounded / 10
}
