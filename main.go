package main

import (
	// "fmt";
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))
	// mux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
	}
	s.ListenAndServe()
}
