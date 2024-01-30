package main

import (
	"Kjasn/Kin/kin"
	"fmt"
	"net/http"
)

type Engine struct {
}



func main() {
	router := kin.New()

	router.GET("/", indexHandler)
	router.GET("/hello", helloHandler)
	err := router.Run(":80")
	if err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}


func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}