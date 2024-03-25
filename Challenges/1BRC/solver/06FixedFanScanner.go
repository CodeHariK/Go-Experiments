package solver

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func FixedFanScanner06() {
	file, _ := readFile()
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	weatherMap := make(map[string]WeatherStation)

	sss := make(chan string, 100)
	for i := 0; i < 4; i++ {
		go processData2(sss, weatherMap)
	}

	for scanner.Scan() {
		sss <- scanner.Text()
	}

	printWeatherMap(weatherMap)
}

func processData2(sss chan string, weatherMap map[string]WeatherStation) {
	for {
		s := strings.Split(<-sss, ";")

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
}
