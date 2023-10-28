package main

import (
	"context"
	"fetch/api"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// set up logging
	file, _ := os.Create("log")
	log.SetOutput(file)

	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", api.ProcessReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", api.GetPointsHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server in a separate Goroutine
	go func() {
		fmt.Println("Server is running on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			log.Fatal(err)
		}
	}()

	// Create a channel to receive signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the interrupt signal
	<-stop
	fmt.Println("Shutting down the server...")

	// wrap up logging
	file.Close()

	// Create a context with a timeout to force shutdown after a given duration (e.g., 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
