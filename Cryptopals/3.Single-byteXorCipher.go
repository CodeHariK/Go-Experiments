package main

import (
	"encoding/hex"
	"fmt"
)

func frequencyAnalysis(text string) map[byte]float64 {
	frequencies := make(map[rune]int)
	totalChars := 0

	// Count the occurrences of each character
	for _, char := range text {
		frequencies[char]++
		totalChars++
	}

	// Calculate the frequency percentages
	frequencyPercentages := make(map[byte]float64)
	for char, count := range frequencies {
		frequencyPercentages[byte(char)] = float64(count) / float64(totalChars) * 100
	}

	return frequencyPercentages
}

var englishLetterFrequency = map[byte]float64{}

func scorePlaintext(plaintext string) float64 {
	var score float64
	for _, char := range plaintext {
		score += englishLetterFrequency[byte(char)]
	}
	return score
}

func xorWithKey(hexStr string, key byte) string {
	// Decode hexadecimal string
	bytes, _ := hex.DecodeString(hexStr)

	// XOR with the key
	resultBytes := make([]byte, len(bytes))
	for i := 0; i < len(bytes); i++ {
		resultBytes[i] = bytes[i] ^ key
	}

	// Convert the result to a string
	resultString := string(resultBytes)
	return resultString
}

func singleByteXorCipher(hexString string, print bool) (bestPlaintext string, bestKey byte, bestScore float64) {
	text := `The world is the totality of entities, the whole of reality, or everything that is.[1] The nature of the world has been conceptualized differently in different fields. Some conceptions see the world as unique while others talk of a "plurality of worlds". Some treat the world as one simple object while others analyze the world as a complex made up of parts. In scientific cosmology, the world or universe is commonly defined as "[t]he totality of all space and time; all that is, has been, and will be". Theories of modality talk of possible worlds as complete and consistent ways how things could have been. Phenomenology, starting from the horizon of co-given objects present in the periphery of every experience, defines the world as the biggest horizon or the "horizon of all horizons". In philosophy of mind, the world is contrasted with the mind as that which is represented by the mind. Theology conceptualizes the world in relation to God, for example, as God's creation, as identical to God or as the two being interdependent. In religions, there is a tendency to downgrade the material or sensory world in favor of a spiritual world to be sought through religious practice. A comprehensive representation of the world and our place in it, as is found in religions, is known as a worldview. Cosmogony is the field that studies the origin or creation of the world while eschatology refers to the science or doctrine of the last things or of the end of the world.`

	// Perform frequency analysis
	englishLetterFrequency = frequencyAnalysis(text)
	// fmt.Println(englishLetterFrequency)

	bestScore = 0.0

	// Try XORing with each possible key
	for key := byte(0); key < 255; key++ {
		decrypted := xorWithKey(hexString, key)

		// Calculate the score for the decrypted text
		score := scorePlaintext(decrypted)

		// fmt.Printf("Key: %c (%d) : %f\n", key, key, score)
		// fmt.Println("Decrypted message:", decrypted)

		// Update the best score and plaintext if this one is better
		if score > bestScore {
			bestScore = score
			bestPlaintext = decrypted
			bestKey = key
		}
	}

	if print {
		fmt.Printf("\nKey: %c (%d) : %f\n", bestKey, bestKey, bestScore)
		fmt.Println("Decrypted message:", bestPlaintext)
	}

	return bestPlaintext, bestKey, bestScore
}
