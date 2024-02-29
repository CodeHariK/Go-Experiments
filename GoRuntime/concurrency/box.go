package main

import (
	"fmt"
	"sync"
	"time"
)

var lock sync.Mutex

func process(data int) int {
	time.Sleep(time.Second * 2)
	return data * 2
}

func processData(wg *sync.WaitGroup, result *[]int, data int) {
	// lock.Lock()//--------------
	defer wg.Done()

	processedData := process(data)

	fmt.Println(data, " : ", processedData, " : ", *result)

	lock.Lock() //--------------

	*result = append(*result, processedData)
	fmt.Println(data, " : ", processedData, " : ", *result)

	lock.Unlock() //-------------
}

func processDataBox(wg *sync.WaitGroup, result *int, data int) {
	defer wg.Done()

	processedData := process(data)

	*result = processedData
}

func box() {
	start := time.Now()
	var wg sync.WaitGroup

	input := []int{1, 2, 3, 4, 5}
	result := []int{}

	for _, data := range input {
		wg.Add(1)
		go processData(&wg, &result, data)
	}

	wg.Wait()
	fmt.Println(result)
	fmt.Println(time.Since(start))

	resBox := make([]int, len(input))
	for i, data := range input {
		wg.Add(1)
		go processDataBox(&wg, &resBox[i], data)
	}
	wg.Wait()
	fmt.Println(resBox)
	fmt.Println(time.Since(start))
}
