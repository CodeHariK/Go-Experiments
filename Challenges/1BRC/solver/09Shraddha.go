package solver

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
)

func Shraddha09() {
	evaluate3()
}

func evaluate3() string {
	mapOfTemp, err := readFileLineByLineIntoAMap3("measurements.txt")
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

func readFileLineByLineIntoAMap3(filepath string) (map[string][]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	mapOfTemp := make(map[string][]float64)

	chanOwner := func() <-chan []string {
		resultStream := make(chan []string, 100)
		toSend := make([]string, 100)
		go func() {
			defer close(resultStream)
			scanner := bufio.NewScanner(file)
			var count int
			for scanner.Scan() {
				if count == 100 {
					localCopy := make([]string, 100)
					copy(localCopy, toSend)
					resultStream <- localCopy
					count = 0
				}
				toSend[count] = scanner.Text()
				count++
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
			temp := convertStringToFloat(text[index+1:])
			if _, ok := mapOfTemp[city]; ok {
				mapOfTemp[city] = append(mapOfTemp[city], temp)
			} else {
				mapOfTemp[city] = []float64{temp}
			}
		}
	}
	return mapOfTemp, nil
}

type cityTemp struct {
	city string
	temp float64
}
