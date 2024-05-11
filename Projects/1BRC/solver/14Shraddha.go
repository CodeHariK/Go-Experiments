package solver

import (
	"bytes"
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
	"sync"
)

func Shraddha14() {
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

	evaluate8()

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

func evaluate8() string {
	mapOfTemp, err := readFileLineByLineIntoAMap8("measurements.txt")
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

func readFileLineByLineIntoAMap8(filepath string) (map[string]cityTemperatureInfo, error) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mapOfTemp := make(map[string]cityTemperatureInfo)
	resultStream := make(chan []string, 100)
	chunkStream := make(chan []byte, 15)
	chunkSize := 64 * 1024 * 1024
	var wg sync.WaitGroup

	// spawn workers to consume (process) file chunks read
	for i := 0; i < runtime.NumCPU()-1; i++ {
		wg.Add(1)
		go func() {
			for chunk := range chunkStream {
				processReadChunk8(chunk, resultStream)
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

	return mapOfTemp, nil
}

func processReadChunk8(buf []byte, resultStream chan<- []string) {
	var count int
	var stringsBuilder strings.Builder
	toSend := make([]string, 100)

	for _, char := range buf {
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
	if count != 0 {
		resultStream <- toSend[:count]
	}
}
