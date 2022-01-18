package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := startServer()

	listenOSKillSignals()
	_ = server.Shutdown(context.Background())
}

func startServer() *http.Server {
	http.HandleFunc("/health/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Status string `json:"status"`
		}{"OK"})
	})

	srv := &http.Server{Addr: ":8000"}
	go func() {
		_ = srv.ListenAndServe()
	}()
	return srv
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
