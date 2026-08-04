package main

import (
	// "fmt";
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
	}
	s.ListenAndServe()
// ----
	// http.Handle("/foo", fooHandler)

	// http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// })

	// log.Fatal(http.ListenAndServe(":8080", nil))

	// s := &http.Server{
	// 	Addr:           ":8080",
	// 	Handler:        myHandler,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// log.Fatal(s.ListenAndServe())
}

