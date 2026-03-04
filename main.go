package main

import (
	"log"
	"net/http"

	"github.com/WillSmithTE/ai-recipes/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard", handlers.DashboardHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
