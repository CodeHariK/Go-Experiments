package error

import (
	"fmt"
	"log"
)

func Error(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v", message, err)
	}
}

func Fatal(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
