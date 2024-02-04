package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// https://cryptopals.com/sets/1/challenges/1

func hexToBase64(hexString string) (string, error) {
	// Decode hexadecimal string
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}

	// Encode as base64
	base64String := base64.StdEncoding.EncodeToString(decoded)
	return base64String, nil
}

func hex2base64Check(hexString, expected string) {
	// Convert hex to base64
	base64String, err := hexToBase64(hexString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result
	fmt.Println("Base64:", base64String)
	if base64String == expected {
		fmt.Println("Ok")
	} else {
		fmt.Println("Error")
	}
}
