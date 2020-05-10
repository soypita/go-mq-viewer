package main

import (
	"context"
	"encoding/json"
	"log"
	"mqViewer/internals/handler"
	"mqViewer/internals/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// HeartbeatHandler is healtmonitor response to check server status
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

// ContentTypeMiddleware basic wrapper for json response type
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	serviceHdl := handler.NewMQViewerHandler(services.NewDefaultMQService())
	router := mux.NewRouter()
	router.Use(ContentTypeMiddleware)
	// Provide healthcheck
	router.HandleFunc("/health", HeartbeatHandler).Methods("GET")

	// Basic workflow
	router.HandleFunc("/connect", serviceHdl.CreateNewConnectionWithParams).Methods("POST")
	router.HandleFunc("/getAll", serviceHdl.BrowseAllMessages).Methods("GET")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("Start server")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
	}

	<-done
	log.Println("Shutdown server success")
}
