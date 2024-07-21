package main

import (
	"fmt"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", fs))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})
	srv := &http.Server {
		Addr: ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Server failed:", err)
	}
}
