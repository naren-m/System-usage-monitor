package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func test(rw http.ResponseWriter, req *http.Request) {
	// fmt.Println("--> %s\n\n", formatRequest(req))
	decoder := json.NewDecoder(req.Body)

	var t test_struct
	err := decoder.Decode(&t)

	fmt.Println("t:", t)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer req.Body.Close()
	fmt.Println("t.Test", t.Test)
}

func main() {
	http.HandleFunc("/", test)
	http.ListenAndServe(":8080", nil)
}
