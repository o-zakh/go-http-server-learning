package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitsCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(fmt.Appendf([]byte{}, `
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, cfg.fileserverHits.Load()))
}

func (cfg *apiConfig) hitsReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func main() {

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsCounter)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsReset)
	mux.HandleFunc("POST /api/validate_chirp", validate_handler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}
