package main

import (
	"fmt"
	"net/http"
)

type test_struct struct {
	Test string
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	fmt.Printf("Here in main code, %s\n", request.Method)
	fmt.Println(request.Body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	fmt.Printf("Here in main code, %s\n", r.Method)
}

func main() {
	http.HandleFunc("/", parseGhPost)
	http.ListenAndServe(":8080", nil)
}
