package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"1brc/solver"
)

func main() {
	start := time.Now()

	if len(os.Args) == 2 {
		size := 0
		var err error
		if size, err = strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println("Not a integer")
			os.Exit(1)
		}
		create_measurements(size)
	} else {
		// solver.Basic01()

		/*
			// We need line seperated scanner, fixed byte buffer scanner won't work.
			// solver.Chunk02(1024)
			// solver.Chunk02(10240)
			// solver.Chunk02(102400)

			// solver.AsyncChunk03(1024)
			// solver.AsyncChunk03(10240)
			// solver.AsyncChunk03(102400)
		*/

		// solver.SeqScanner04()
		// solver.FullFanScanner05()
		// solver.FixedFanScanner06()

		// solver.ShraddhaBasic1()
		// solver.ShraddhaBasic2()
		// solver.ShraddhaBasic3()
		// solver.ShraddhaBasic4()
		// solver.Shraddha5()
		// solver.Shraddha6()
		solver.Shraddha13()
	}

	fmt.Printf("Time : %d ms\n", time.Since(start).Milliseconds())
}
