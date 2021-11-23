package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Query()["a"][0] == "b" {
		fmt.Fprintf(w, "Hello")
		return nil
	}
	//fmt.Fprintf(w, "Hello")
	//return nil
	return fmt.Errorf("aaa")
}

func main() {
	//监听
	h := InitHystrix()

	http.HandleFunc("/aaa", h.Wrap(handler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

