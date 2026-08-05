package main

import (
	// "fmt";
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
	}
	s.ListenAndServe()
}
