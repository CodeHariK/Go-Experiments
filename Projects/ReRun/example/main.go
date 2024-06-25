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
		fmt.Fprintln(w, "Html")

		fmt.Fprint(os.Stderr, `Html`)
	})

	fmt.Println("Server:8080")
	http.ListenAndServe(":8080", nil)
}
