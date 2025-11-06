package main

import (
	"fmt"
	"os"
	"strings"
)

func DetectSingleCharacterXor() {
	if file, err := os.ReadFile("./4.DetectSingleCharacterXor.txt"); err != nil {
		fmt.Println("Error:", err)
	} else {

		bestScore := 0.0
		var bestPlaintext string
		var bestKey byte

		for _, line := range strings.Split(string(file), "\n") {

			text, key, score := singleByteXorCipher(string(line), false)

			if score > bestScore {
				bestScore = score
				bestKey = key
				bestPlaintext = text
			}

		}

		fmt.Printf("\nKey: %c (%d) : %f\n", bestKey, bestKey, bestScore)
		fmt.Println("Decrypted message:", bestPlaintext)
	}
}
