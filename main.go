package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("custom middleware")
	switch req.URL.Path {
		case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		case "/hello":
			for k, v := range req.Header {
				fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
			}
		default:
			fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}


func main() {
	fmt.Println("hello world")
	// http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/hello", helloHandler)

	// log.Fatal(http.ListenAndServe(":8080", nil))
	// custom middleware 
	var e Engine
	log.Fatal(http.ListenAndServe(":8080", &e))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("the path is %v\n", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		ret := fmt.Sprintf("Header[%vl] = %v\n", k, v)
		w.Write([]byte(ret))
	}
}
