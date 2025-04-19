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
	defaultServerAddress = ":8888" // Default port if PORT env var is not set
	readTimeout          = 5 * time.Second
	writeTimeout         = 10 * time.Second
	idleTimeout          = 120 * time.Second
	portEnvVar           = "PORT"
)

func main() {
	// Determine the server address from the environment variable or use the default
	port := os.Getenv(portEnvVar)
	if port == "" {
		port = defaultServerAddress
		fmt.Printf("Using default port: %s\n", defaultServerAddress)
	} else {
		port = ":" + port
		fmt.Printf("Using port from environment variable %s: %s\n", portEnvVar, port)
	}

	// Create a new ServeMux to handle different routes if needed in the future
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlePage)

	// Configure the HTTP server with timeouts for better resource management
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout, // Recommended for long-lived connections
		ErrorLog:     log.New(os.Stderr, "http: ", log.LstdFlags), // Custom error log
	}

	// Start the server in a goroutine so it doesn't block the main function
	go func() {
		fmt.Printf("Server listening on %s\n", server.Addr)
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

