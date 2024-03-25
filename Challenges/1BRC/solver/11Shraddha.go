package solver

import (
	"cmp"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"slices"
	"strings"
	"sync"
)

func Shraddha11() {
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

	fmt.Println(evaluate5())

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

func evaluate5() string {
	mapOfTemp, err := readFileLineByLineIntoAMap4("measurements.txt")
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
		go func(city string, temps []float64) {
			defer wg.Done()
			var min, max, avg float64
			min, max = math.MaxFloat64, math.MinInt64

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

			updateResult(city, fmt.Sprintf("%.1f/%.1f/%.1f", min, avg, max))
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
