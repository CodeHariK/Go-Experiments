package solver

import (
	"fmt"
	"sync"
)

type chunk struct {
	bufsize int
	offset  int64
}

func AsyncChunk03(buffersize int) {
	file, filesize := readFile()
	defer file.Close()

	concurrency := int(filesize) / buffersize
	chunksizes := make([]chunk, concurrency)
	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = buffersize
		chunksizes[i].offset = int64(buffersize * i)
	}

	if remainder := int(filesize) % buffersize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * buffersize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(chunksizes []chunk, i int) {
			defer wg.Done()

			chunk := chunksizes[i]
			buffer := make([]byte, chunk.bufsize)
			bytesread, err := file.ReadAt(buffer, chunk.offset)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("bytes read : ", bytesread)
			// fmt.Println("", string(buffer)) <- print this and time will increase
		}(chunksizes, i)
	}

	wg.Wait()
}
