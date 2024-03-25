package solver

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strings"
)

func Shraddha13() {
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

	fmt.Println(evaluate7())

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

func evaluate7() string {
	mapOfTemp, err := readFileLineByLineIntoAMap7("measurements.txt")
	if err != nil {
		panic(err)
	}

	resultArr := make([]string, len(mapOfTemp))
	var count int
	for city := range mapOfTemp {
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

func readFileLineByLineIntoAMap7(filepath string) (map[string]cityTemperatureInfo, error) {
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
