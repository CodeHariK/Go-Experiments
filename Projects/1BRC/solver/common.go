package solver

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type WeatherStation struct {
	Total float32
	Min   float32
	Max   float32
	Len   int
}

func reduce(buffer []byte) map[string]WeatherStation {
	memFile := string(buffer)

	memSplit := strings.Split(memFile, "\n")

	weatherMap := make(map[string]WeatherStation)

	for i := 0; i < len(memSplit); i++ {
		updateMap(memSplit[i], weatherMap)
	}

	return weatherMap
}

func updateMap(memSplit string, weatherMap map[string]WeatherStation) {
	s := strings.Split(memSplit, ";")

	if len(s) < 2 {
		return
	}

	v, err := (strconv.ParseFloat(s[1], 32))
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}

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
}

func readFile() (*os.File, int64) {
	file, err := os.Open("measurements.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filesize := fileinfo.Size()
	return file, filesize
}

func printWeatherMap(weatherMap map[string]WeatherStation) {
	for k, v := range weatherMap {
		fmt.Printf("%s  avg:%f  min:%f  max:%f \n", k, v.Total/float32(v.Len), v.Min, v.Max)
	}
}
