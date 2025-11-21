package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var requestCount int

func helloHandler(w http.ResponseWriter, r *http.Request) {
	requestCount++

	// Print to stdout
	fmt.Printf("[%s] Request #%d received from %s\n",
		time.Now().Format("15:04:05"),
		requestCount,
		r.RemoteAddr)

	// Get query parameters
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	// Prepare response
	message := fmt.Sprintf("Hello, %s! This is request #%d", name, requestCount)

	// Print message being sent
	fmt.Printf("[%s] Sending response: %s\n",
		time.Now().Format("15:04:05"),
		message)

	// Send response
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", message)

	// Additional logging
	fmt.Printf("[%s] Response sent successfully\n", time.Now().Format("15:04:05"))
}

func main() {
	port := ":8080"
	pid := os.Getpid()

	fmt.Printf("Process PID: %d\n", pid)
	fmt.Printf("Starting Hello World web server on http://localhost%s\n", port)
	fmt.Println()
	fmt.Println("To attach Delve debugger to this process:")
	fmt.Printf("  dlv attach %d\n", pid)
	fmt.Println()
	fmt.Println("Try visiting:")
	fmt.Println("  http://localhost:8080/")
	fmt.Println("  http://localhost:8080/?name=Alice")
	fmt.Println("  http://localhost:8080/?name=Bob")
	fmt.Println()

	http.HandleFunc("/", helloHandler)

	fmt.Printf("Server is listening on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
