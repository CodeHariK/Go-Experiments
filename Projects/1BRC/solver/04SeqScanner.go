package solver

import (
	"bufio"
)

func SeqScanner04() {
	file, _ := readFile()
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	weatherMap := make(map[string]WeatherStation)

	read := scanner.Scan()

	for read {
		if read {
			updateMap(scanner.Text(), weatherMap)
			// fmt.Println("read byte array: ", scanner.Bytes())
			// fmt.Println(scanner.Text())
		}
		read = scanner.Scan()
	}

	printWeatherMap(weatherMap)
}
