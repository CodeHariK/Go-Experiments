package solver

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
	"strings"
	"sync"
)

func Shraddha08() {
	// trace.Start(os.Stderr)
	// defer trace.Stop()

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

	fmt.Println(evaluate2())

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

func evaluate2() string {
	mapOfTemp, err := readFileLineByLineIntoAMap2("measurements.txt")
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
			var min, max, avg float64
			min, max = math.MaxFloat64, 0

			for _, temp := range temps {
				if temp < min {
					min = temp
				}

				if temp > max {
					max = temp
				}
				avg += temp
			}

			avg = avg / float64(len(temps))
			avg = math.Ceil(avg*10) / 10

			updateResult(fmt.Sprintf("%s=%.1f/%.1f/%.1f", city, min, avg, max))
		}(city, temps)
	}

	wg.Wait()
	sort.Strings(result)
	return strings.Join(result, ", ")
}

func readFileLineByLineIntoAMap2(filepath string) (map[string][]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	mapOfTemp := make(map[string][]float64)

	chanOwner := func() <-chan string {
		resultStream := make(chan string, 100)
		go func() {
			defer close(resultStream)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				text := scanner.Text()
				resultStream <- text
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()

	for text := range resultStream {
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
