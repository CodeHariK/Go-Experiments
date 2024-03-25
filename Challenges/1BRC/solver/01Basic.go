package solver

import (
	"fmt"
)

func Basic01() {
	file, filesize := readFile()
	defer file.Close()

	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("bytes read: ", bytesread)

	weatherMap := reduce(buffer)

	printWeatherMap(weatherMap)
}
