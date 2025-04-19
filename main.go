package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serverAddress = ":8080"
	readTimeout   = 5 * time.Second
	writeTimeout  = 10 * time.Second
	idleTimeout   = 120 * time.Second
)

func main() {
	// Create a new ServeMux to handle different routes if needed in the future
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlePage)

	// Configure the HTTP server with timeouts for better resource management
	server := &http.Server{
		Addr:         serverAddress,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout, // Recommended for long-lived connections
		ErrorLog:     log.New(os.Stderr, "http: ", log.LstdFlags), // Custom error log
	}

	// Start the server in a goroutine so it doesn't block the main function
	go func() {
		fmt.Printf("Server listening on %s\n", serverAddress)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		fmt.Println("Server stopped.")
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK) // Use http.StatusOK constant
	page := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hello from Go!</title>
</head>
<body>
    <p> Hello from Docker! I'm a Go server. </p>
</body>
</html>
`
	_, err := w.Write([]byte(page))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

