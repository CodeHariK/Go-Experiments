package solver

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var lock sync.Mutex

func FullFanScanner05() {
	file, _ := readFile()
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	weatherMap := make(map[string]WeatherStation)

	read := scanner.Scan()

	var wg sync.WaitGroup

	for read {
		if read {

			wg.Add(1)

			go processData(scanner.Text(), &wg, weatherMap)

		}
		read = scanner.Scan()
	}

	wg.Wait()

	printWeatherMap(weatherMap)
}

func processData(memSplit string, wg *sync.WaitGroup, weatherMap map[string]WeatherStation) {
	defer wg.Done()

	s := strings.Split(memSplit, ";")

	if len(s) < 2 {
		return
	}

	v, err := (strconv.ParseFloat(s[1], 32))
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}

	lock.Lock()

	w, ok := weatherMap[s[0]]
	if ok == false {
		weatherMap[s[0]] = WeatherStation{
			Min:   float32(v),
			Total: float32(v),
			Max:   float32(v),
			Len:   1,
		}
	} else {
		weatherMap[s[0]] = WeatherStation{
			Min:   min(float32(v), w.Min),
			Total: w.Total + float32(v),
			Max:   max(float32(v), w.Max),
			Len:   w.Len + 1,
		}
	}

	lock.Unlock()
}
