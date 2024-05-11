package solver

import (
	"fmt"
	"io"
)

func Chunk02(buffersize int) {
	file, _ := readFile()
	defer file.Close()

	buffer := make([]byte, buffersize)

	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}

		fmt.Println("bytes read: ", bytesread)
		// fmt.Println("\n", string(buffer[:bytesread])) <- print this and time will increase
	}
}
