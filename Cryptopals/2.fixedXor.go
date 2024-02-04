package main

import (
	"encoding/hex"
	"fmt"
)

// https://cryptopals.com/sets/1/challenges/2

func xorStrings(str1, str2 string) (string, error) {
	// Decode hexadecimal strings
	bytes1, err := hex.DecodeString(str1)
	if err != nil {
		return "", err
	}

	bytes2, err := hex.DecodeString(str2)
	if err != nil {
		return "", err
	}

	// Check if the lengths are the same
	if len(bytes1) != len(bytes2) {
		return "", fmt.Errorf("hexadecimal strings must have the same length")
	}

	// XOR the corresponding bytes
	resultBytes := make([]byte, len(bytes1))
	for i := 0; i < len(bytes1); i++ {
		resultBytes[i] = bytes1[i] ^ bytes2[i]
	}

	// Encode the result as hexadecimal
	resultHex := hex.EncodeToString(resultBytes)
	return resultHex, nil
}

func fixedXorCheck(str1, str2, expected string) {
	// XOR the strings
	result, err := xorStrings(str1, str2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result
	fmt.Println("XOR Result:", result)
	if result == expected {
		fmt.Println("Ok")
	} else {
		fmt.Println("Error")
	}
}
