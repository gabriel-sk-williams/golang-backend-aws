package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello World! append meeee reeee For example, use %s/Mary to say hello to Mary.", r.Host)
	} else {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	}
}

func main() {
	fmt.Println("test server activated")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5000", nil)
}
