package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)

		fmt.Println("GET docs")
		fmt.Fprint(os.Stderr, `Html`)
		fmt.Fprintln(w, `<body style="background:black;color:white;text-align: center;align-content: center;font: 30px monospace;"><span>Hello</span></body>`)
	})

	fmt.Println("Server:8080")
	http.ListenAndServe(":8080", nil)
}
