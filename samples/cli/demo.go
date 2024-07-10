package main

import (
	"log"
	"net/http"

	"github.com/strotz/goplaycalls/samples/hello"
)

// Run the web server that prints request to the log.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
	})
	http.HandleFunc("/hello", hello.Handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
