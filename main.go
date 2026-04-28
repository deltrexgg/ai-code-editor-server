package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deltrexgg/ai-code-editor-server/internals/ai"
	"github.com/deltrexgg/ai-code-editor-server/internals/config"
	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/joho/godotenv"
)

// LoggingMiddleware logs the time taken for each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Pass control to the next handler
		next.ServeHTTP(w, r)
		
		// This runs after the handler finishes
		fmt.Printf("Method: %s | Path: %s | Duration: %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	_ = godotenv.Load()

	//Load ENV
	cred := config.LoadConfig()

	//minio load
	infra.InitMinio(cred.Minio)

	//database
	infra.InitDB(cred.Postgres.DSN())

	// Create a standard handler
	mux := http.NewServeMux()
	
	mux.HandleFunc("/genfiles", func(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    defer r.Body.Close()

    type RequestBody struct {
        Content string `json:"content"`
    }

    var reqBody RequestBody

    err := json.NewDecoder(r.Body).Decode(&reqBody)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if reqBody.Content == "" {
        http.Error(w, "content is required", http.StatusBadRequest)
        return
    }

    result, err := ai.FileStructure(reqBody.Content, cred.AI.IP)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(result))
	})

	// Wrap the entire mux (router) with the middleware
	wrappedMux := LoggingMiddleware(mux)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", wrappedMux)
}

